package mqtt

import (
	"bytes"
	"fmt"
	"github.com/brocaar/lorawan"
	paho "github.com/eclipse/paho.mqtt.golang"
	"open/backend/gateway"
	pb "open/backend/proto"
	"open/config"
	"open/helpers"
	"open/marshaler"
	"strings"
	"sync"
	"text/template"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/panjf2000/ants"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

const deduplicationLockTTL = time.Millisecond * 500
const (
	marshalerV2JSON = iota
	marshalerProtobuf
	marshalerJSON
)
type Frame struct {
	DeviceType string
	Msg paho.Message
}



// Backend implements a MQTT pub-sub backend.
type Backend struct {
	sync.RWMutex

	wg sync.WaitGroup
	HandlePool              *ants.PoolWithFunc

	rxPacketChan      chan *pb.UplinkFrame //上行有管道，下行不需要管道
	//statsPacketChan   chan gw.GatewayStats
	//downlinkTXAckChan chan gw.DownlinkTXAck

	conn                 paho.Client
	//redisPool            *redis.Pool
	eventTopic           string
	commandTopicTemplate *template.Template
	qos                  uint8

	gatewayMarshaler map[lorawan.EUI64]marshaler.Type
}

// NewBackend creates a new Backend.
func NewBackend( c config.Config) (gateway.Gateway, error) {
	conf := c.Mqtt
	var err error

	b := Backend{
		rxPacketChan:      make(chan *pb.UplinkFrame), // 无缓冲，意味着，并发量高的时候，将无法接收新的帧
		//statsPacketChan:   make(chan gw.GatewayStats),
		//downlinkTXAckChan: make(chan gw.DownlinkTXAck),
		//gatewayMarshaler:  make(map[lorawan.EUI64]marshaler.Type),
		eventTopic:        conf.EventTopic,
		qos:               uint8(conf.Qos),
	}

	b.commandTopicTemplate, err = template.New("command").Parse(conf.CommandTopicTemplate) // 构建下行命令通用topic格式
	if err != nil {
		return nil, errors.Wrap(err, "gateway/mqtt: parse command topic template error")
	}

	b.HandlePool ,_ = ants.NewPoolWithFunc(conf.MaxHandleGoroutine, func(i interface{}) {
		//.MessageCheck(i)
		b.FrameUnpack(i)
	})

	opts := paho.NewClientOptions()
	opts.AddBroker(conf.Server) //在root.go中
	opts.SetUsername(conf.Username)
	opts.SetPassword(conf.Password)
	opts.SetCleanSession(conf.CleanSession)
	opts.SetClientID(conf.ClientID)
	opts.SetOnConnectHandler(b.onConnected)
	opts.SetConnectionLostHandler(b.onConnectionLost)
	opts.SetMaxReconnectInterval(time.Duration(conf.MaxReconnectInterval)*time.Second)


	log.WithField("server", conf.Server).Info("gateway/mqtt: connecting to mqtt broker")
	b.conn = paho.NewClient(opts)
	for {
		if token := b.conn.Connect(); token.Wait() && token.Error() != nil {
			log.Errorf("gateway/mqtt: connecting to mqtt broker failed, will retry in 2s: %s", token.Error())
			time.Sleep(2 * time.Second)
		} else {
			break
		}
	}

	return &b, nil
}

// Close closes the backend.
// Note that this closes the backend one-way (gateway to backend).
// This makes it possible to perform a graceful shutdown (e.g. when there are
// still packets to send back to the gateway).
func (b *Backend) Close() error {
	log.Info("gateway/mqtt: closing backend")

	log.WithField("topic", b.eventTopic).Info("gateway/mqtt: unsubscribing from event topic")
	if token := b.conn.Unsubscribe(b.eventTopic); token.Wait() && token.Error() != nil {
		return fmt.Errorf("gateway/mqtt: unsubscribe from %s error: %s", b.eventTopic, token.Error())
	}

	log.Info("backend/gateway: handling last messages")
	b.wg.Wait()
	b.HandlePool.Release()
	close(b.rxPacketChan)
	//close(b.statsPacketChan)
	//close(b.downlinkTXAckChan)
	return nil
}

// RXPacketChan returns the uplink-frame channel.
// uplink.go 调用这个函数来获取 数据接收管道
func (b *Backend) RXPacketChan() chan *pb.UplinkFrame {
	return b.rxPacketChan
}

// StatsPacketChan returns the gateway stats channel.
//func (b *Backend) StatsPacketChan() chan gw.GatewayStats {
//	return b.statsPacketChan
//}

// DownlinkTXAckChan returns the downlink tx ack channel.
//func (b *Backend) DownlinkTXAckChan() chan gw.DownlinkTXAck {
//	return b.downlinkTXAckChan
//}

// SendTXPacket sends the given downlink-frame to the gateway.
func (b *Backend) SendTXPacket(txPacket pb.DownlinkFrame) error {
	//if txPacket.TxInfo == nil {
	//	return errors.New("tx_info must not be nil")
	//}

	//gatewayID := helpers.GetGatewayID(txPacket.TxInfo)
	downID := helpers.GetDownlinkID(&txPacket)
	//downID是上下文uuid，当初的GatewayID当做，DevAddr
	log.Debug("[SendTXPacket]downID=",downID)

	//TODO:这里应该去掉devName和devID，这里是网络通信不是lora通信
	return b.publishCommand(log.Fields{
		"downlink_id": downID,
	}, txPacket.DevName,txPacket.DevId,"down", &txPacket)
}

// SendGatewayConfigPacket sends the given GatewayConfigPacket to the gateway.
//func (b *Backend) SendGatewayConfigPacket(configPacket gw.GatewayConfiguration) error {
//	gatewayID := helpers.GetGatewayID(&configPacket)
//
//	return b.publishCommand(log.Fields{}, gatewayID, "config", &configPacket)
//}
// 下行命令. devAddr当做GatewayID
// command 是这里用做命令类型，即up 、 down
func (b *Backend) publishCommand(fields log.Fields, DevType string, DevId string, command string, msg proto.Message) error {
	//t := b.getGatewayMarshaler(gatewayID) //根据上行的格式
	bb, err := marshaler.MarshalCommand(marshaler.Protobuf, msg)
	if err != nil {
		return errors.Wrap(err, "gateway/mqtt: marshal gateway command error")
	}
	//TODO:DevAddr和设备类型应该再做一下区分。这里暂时写死SER_NAME 作为topic
	templateCtx := struct {
		DevType   string
		DevId string
		CommandType string
	}{DevType, DevId,command}
	topic := bytes.NewBuffer(nil) //把上面两个参数按照模板格式写入到topic这个buffer中。"gateway/{{ .DevType }}/{{ .DevId }}/command/{{ .CommandType }}"
	if err := b.commandTopicTemplate.Execute(topic, templateCtx); err != nil {
		return errors.Wrap(err, "execute command topic template error")
	}

	fields["DevId"] = DevId
	fields["command"] = command
	fields["qos"] = b.qos
	fields["topic"] = topic.String() // 下行的topic
	fields["proto_body"] = fmt.Sprintf("%02X",bb)

	log.WithFields(fields).Info("gateway/mqtt: publishing gateway command")

	mqttCommandCounter(command).Inc()

	if token := b.conn.Publish(topic.String(), b.qos, false, bb); token.Wait() && token.Error() != nil {
		return errors.Wrap(err, "gateway/mqtt: publish gateway command error")
	}

	return nil
}
// 订阅时候的回调函数
func (b *Backend) eventHandler(c paho.Client, msg paho.Message) {
	log.Info("[eventHandler]",msg.Topic())
	b.wg.Add(1)
	defer b.wg.Done()

	if strings.HasSuffix(msg.Topic(), "up") { //后缀是up，对来自mqtt消息的第一步处理
		mqttEventCounter("up").Inc()
		b.rxPacketHandler(c, msg)
	} else if strings.HasSuffix(msg.Topic(), "ack") {
		mqttEventCounter("ack").Inc()
		//b.ackPacketHandler(c, msg)
		log.Info("[eventHandler]ldm delete ack Handler")
	} else if strings.HasSuffix(msg.Topic(), "stats") {
		mqttEventCounter("stats").Inc()
		//b.statsPacketHandler(c, msg)
		log.Info("[eventHandler]ldm delete stats Handler")

	}
}
// 处理上行帧
func (b *Backend) rxPacketHandler(c paho.Client, msg paho.Message) {
	b.wg.Add(1)
	defer b.wg.Done()

	// 建立一个协成池，然后并发rpc请求相应的设备驱动来处理。根据topic区分不同的设备
	// topic 格式暂时定为，设备上行： up/deviceType/deviceAddr/xxxx
	topic := strings.Split( msg.Topic(), "/")
	if len(topic ) < 3 {
		log.Warn("[prxPacketHandler]get a error topic.",msg.Topic())
		return
	}
	frame := Frame {
		DeviceType:topic[1],
		Msg:msg,
	}
	_=b.HandlePool.Invoke(&frame)

	//var uplinkFrame pb.UplinkFrame
	//_, err := marshaler.UnmarshalUplinkFrame(msg.Payload(), &uplinkFrame)
	//if err != nil {
	//	log.Debugf("mqtt payload hex= %02x,str=%s\n",msg.Payload(),string(msg.Payload()))
	//	log.WithFields(log.Fields{
	//	}).WithError(err).Error("gateway/mqtt: unmarshal uplink frame error")
	//	return
	//}

	//if uplinkFrame.TxInfo == nil {
	//	log.WithFields(log.Fields{
	//		"data_base64": base64.StdEncoding.EncodeToString(msg.Payload()),
	//	}).Error("gateway/mqtt: tx_info must not be nil")
	//	return
	//}


	//log.WithFields(log.Fields{
	//	"uplink_id": uplinkFrame.UplinkId ,
	//	"DevAddr":uplinkFrame.DevAddr ,
	//}).Info("gateway/mqtt: uplink frame received")
	//

	//b.rxPacketChan <- uplinkFrame //写入接收管道，不会重复的消息。这个管道无缓冲，阻塞
}


func (b *Backend)FrameUnpack(frame interface{}){
	MsgInfo :=frame.(*Frame)
	// DeviceType 和 注册服务的id应该一一对应
		log.Debug("type=",MsgInfo.DeviceType )
		var uplinkFrame pb.UplinkFrame
		_, err := marshaler.UnmarshalUplinkFrame(MsgInfo.Msg.Payload(), &uplinkFrame)
		if err != nil {
			log.Debugf("mqtt payload hex= %02x,str=%s\n",MsgInfo.Msg.Payload(),string(MsgInfo.Msg.Payload()))
			log.WithFields(log.Fields{
			}).WithError(err).Error("gateway/mqtt: unmarshal uplink frame error")
			return
		}
		//r := upParserService(MsgInfo.DeviceType,uplinkFrame.DevName,uplinkFrame.PhyPayload)
		//if r == nil {
		//	log.Debug("[FrameUnpack]upParserService result is nil")
		//	return
		//}
		//log.Debug("respond: ",r.ServiceName,r.Payload)
		//uplinkFrame.PhyPayload = r.Payload
		b.sendResult(&uplinkFrame)
	return


}
func (b *Backend) sendResult(uplinkFrame *pb.UplinkFrame){
	b.rxPacketChan <- uplinkFrame //写入接收管道，不会重复的消息。这个管道无缓冲，阻塞
}

//
//func (b *Backend) statsPacketHandler(c paho.Client, msg paho.Message) {
//	b.wg.Add(1)
//	defer b.wg.Done()
//
//	var gatewayStats gw.GatewayStats
//	t, err := marshaler.UnmarshalGatewayStats(msg.Payload(), &gatewayStats)
//	if err != nil {
//		log.WithFields(log.Fields{
//			"data_base64": base64.StdEncoding.EncodeToString(msg.Payload()),
//		}).WithError(err).Error("gateway/mqtt: unmarshal gateway stats error")
//		return
//	}
//
//	gatewayID := helpers.GetGatewayID(&gatewayStats)
//	statsID := helpers.GetStatsID(&gatewayStats)
//	b.setGatewayMarshaler(gatewayID, t)
//
//	// Since with MQTT all subscribers will receive the stats messages sent
//	// by all the gateways, the first instance receiving the message must lock it,
//	// so that other instances can ignore the same message (from the same gw).
//	// As an unique id, the gw mac is used.
//	key := fmt.Sprintf("lora:ns:stats:lock:%s", gatewayID)
//	if locked, err := b.isLocked(key); err != nil || locked {
//		if err != nil {
//			log.WithError(err).WithFields(log.Fields{
//				"key":      key,
//				"stats_id": statsID,
//			}).Error("gateway/mqtt: acquire lock error")
//		}
//
//		return
//	}
//
//	log.WithFields(log.Fields{
//		"gateway_id": gatewayID,
//		"stats_id":   statsID,
//	}).Info("gateway/mqtt: gateway stats packet received")
//	b.statsPacketChan <- gatewayStats
//}
//
//func (b *Backend) ackPacketHandler(c paho.Client, msg paho.Message) {
//	b.wg.Add(1)
//	defer b.wg.Done()
//
//	var ack gw.DownlinkTXAck
//	t, err := marshaler.UnmarshalDownlinkTXAck(msg.Payload(), &ack)
//	if err != nil {
//		log.WithFields(log.Fields{
//			"data_base64": base64.StdEncoding.EncodeToString(msg.Payload()),
//		}).WithError(err).Error("backend/gateway: unmarshal downlink tx ack error")
//	}
//
//	gatewayID := helpers.GetGatewayID(&ack)
//	downlinkID := helpers.GetDownlinkID(&ack)
//	b.setGatewayMarshaler(gatewayID, t)
//
//	// Since with MQTT all subscribers will receive the ack messages sent
//	// by all the gateways, the first instance receiving the message must lock it,
//	// so that other instances can ignore the same message (from the same gw).
//	// As an unique id, the gw mac is used.
//	key := fmt.Sprintf("lora:ns:ack:lock:%s", gatewayID)
//	if locked, err := b.isLocked(key); err != nil || locked {
//		if err != nil {
//			log.WithError(err).WithFields(log.Fields{
//				"key":         key,
//				"downlink_id": downlinkID,
//			}).Error("gateway/mqtt: acquire lock error")
//		}
//
//		return
//	}
//
//	log.WithFields(log.Fields{
//		"gateway_id":  gatewayID,
//		"downlink_id": downlinkID,
//	}).Info("backend/gateway: downlink tx acknowledgement received")
//	b.downlinkTXAckChan <- ack // 传递给下行通道
//}

func (b *Backend) onConnected(c paho.Client) {
	log.Info("backend/gateway: connected to mqtt server")

	mqttConnectCounter().Inc()

	for {
		log.WithFields(log.Fields{
			"topic": b.eventTopic,
			"qos":   b.qos,
		}).Info("gateway/mqtt: subscribing to gateway event topic")
		//参数三是回调函数
		if token := b.conn.Subscribe(b.eventTopic, b.qos, b.eventHandler); token.Wait() && token.Error() != nil {
			log.WithError(token.Error()).WithFields(log.Fields{
				"topic": b.eventTopic,
				"qos":   b.qos,
			}).Errorf("gateway/mqtt: subscribe error")
			time.Sleep(time.Second)
			continue
		}
		break
	}
}

func (b *Backend) onConnectionLost(c paho.Client, reason error) {
	log.Errorf("gateway/mqtt: mqtt connection error: %s", reason)
	mqttDisconnectCounter().Inc()
}

func (b *Backend) setGatewayMarshaler(gatewayID lorawan.EUI64, t marshaler.Type) {
	b.Lock()
	defer b.Unlock()

	b.gatewayMarshaler[gatewayID] = t
}

func (b *Backend) getGatewayMarshaler(gatewayID lorawan.EUI64) marshaler.Type {
	b.RLock()
	defer b.RUnlock()

	return b.gatewayMarshaler[gatewayID] // 在数据上行的时候，会对这个map进行写操作
}

// isLocked returns if a lock exists for the given key, if false a lock is
// acquired.
//func (b *Backend) isLocked(key string) (bool, error) {
//	c := b.redisPool.Get()
//	defer c.Close()
//	// PX 是毫秒
//	_, err := redis.String(c.Do("SET", key, "lock", "PX", int64(deduplicationLockTTL/time.Millisecond), "NX"))
//	if err != nil {
//		if err == redis.ErrNil {
//			// the payload is already being processed by an other instance
//			return true, nil
//		}
//
//		return false, err
//	}
//
//	return false, nil
//}
//
//func newTLSConfig(cafile, certFile, certKeyFile string) (*tls.Config, error) {
//	if cafile == "" && certFile == "" && certKeyFile == "" {
//		return nil, nil
//	}
//
//	tlsConfig := &tls.Config{}
//
//	// Import trusted certificates from CAfile.pem.
//	if cafile != "" {
//		cacert, err := ioutil.ReadFile(cafile)
//		if err != nil {
//			log.WithError(err).Error("gateway/mqtt: could not load ca certificate")
//			return nil, err
//		}
//		certpool := x509.NewCertPool()
//		certpool.AppendCertsFromPEM(cacert)
//
//		tlsConfig.RootCAs = certpool // RootCAs = certs used to verify server cert.
//	}
//
//	// Import certificate and the key
//	if certFile != "" && certKeyFile != "" {
//		kp, err := tls.LoadX509KeyPair(certFile, certKeyFile)
//		if err != nil {
//			log.WithError(err).Error("gateway/mqtt: could not load mqtt tls key-pair")
//			return nil, err
//		}
//		tlsConfig.Certificates = []tls.Certificate{kp}
//	}
//
//	return tlsConfig, nil
//}
