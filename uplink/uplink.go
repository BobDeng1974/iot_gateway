package uplink

import (
	"context"
	"fmt"
	"open/config"
	"open/uplink/data"
	"open/uplink/proprietary"
	"sync"
	"time"

	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	pb "open/backend/proto"
	"open/backend/gateway"
	"open/logging"
)

var (
	deduplicationDelay time.Duration
)

// Setup configures the package.
func Setup(conf config.Config) error {
	if err := data.Setup(conf); err != nil {
		return errors.Wrap(err, "configure uplink/data error")
	}
	return nil
}

// Server represents a server listening for uplink packets.
type Server struct {
	wg sync.WaitGroup
}

// NewServer creates a new server.
func NewServer() *Server {
	return &Server{}
}

// Start starts the server.
func (s *Server) Start() error {
	go func() {
		s.wg.Add(1)
		defer s.wg.Done()
		HandleUplinkFrames(&s.wg) //以协成的方式启动接收，这里阻塞在接收管道上。
	}()

	//go func() {
	//	s.wg.Add(1)
	//	defer s.wg.Done()
	//	HandleDownlinkTXAcks(&s.wg)
	//}()
	return nil
}

// Stop closes the gateway backend and waits for the server to complete the
// pending packets.
func (s *Server) Stop() error {
	if err := gateway.Backend().Close(); err != nil {
		return fmt.Errorf("close gateway backend error: %s", err)
	}
	log.Info("uplink: waiting for pending actions to complete")
	s.wg.Wait()
	return nil
}

// HandleUplinkFrames consumes received packets by the gateway and handles them
// in a separate go-routine. Errors are logged.
func HandleUplinkFrames(wg *sync.WaitGroup) {
	for uplinkFrame := range gateway.Backend().RXPacketChan() { //遍历这个管道。可是这个管道是无缓冲阻塞的。
		go func(uplinkFrame pb.UplinkFrame) { //每接收到一条上行消息，就开一个协成去处理
			wg.Add(1) //控制并发数量，来自更上一层的
			defer wg.Done()
			// The ctxID will be available as context value "ctx_id" so that
			// this can be used when writing logs. This makes it easier to
			// group multiple log-lines to the same context.
			ctxID, err := uuid.NewV4()
			if err != nil {
				log.WithError(err).Error("uplink: get new uuid error")
			}

			ctx := context.Background()
			ctx = context.WithValue(ctx, logging.ContextIDKey, ctxID)

			if err := HandleUplinkFrame(ctx, uplinkFrame); err != nil {
				log.WithFields(log.Fields{
					"ctx_id": ctxID,
				}).WithError(err).Error("uplink: processing uplink frame error")
			}
		}(uplinkFrame)
	}
}

// HandleUplinkFrame handles a single uplink frame.
func HandleUplinkFrame(ctx context.Context, uplinkFrame pb.UplinkFrame) error {

	switch uplinkFrame.FrameType { // 开始解析应用层payload数据
	case gateway.UnconfirmedDataUp,gateway.ConfirmedDataUp:
		return data.Handle(ctx, uplinkFrame)
	case gateway.Proprietary:
		return proprietary.Handle(ctx,uplinkFrame)
	default:
		log.Info("[HandleUplinkFrame]unkown frame type")
		return nil
	}
	return nil
}

// HandleDownlinkTXAcks consumes received downlink tx acknowledgements from
// the gateway.
//func HandleDownlinkTXAcks(wg *sync.WaitGroup) {
//	for downlinkTXAck := range gateway.Backend().DownlinkTXAckChan() {
//		go func(downlinkTXAck gw.DownlinkTXAck) {
//			wg.Add(1)
//			defer wg.Done()
//
//			// The ctxID will be available as context value "ctx_id" so that
//			// this can be used when writing logs. This makes it easier to
//			// group multiple log-lines to the same context.
//			var ctxID uuid.UUID
//			if downlinkTXAck.DownlinkId != nil {
//				copy(ctxID[:], downlinkTXAck.DownlinkId)
//			}
//
//			ctx := context.Background()
//			ctx = context.WithValue(ctx, logging.ContextIDKey, ctxID)
//
//			if err := ack.HandleDownlinkTXAck(ctx, downlinkTXAck); err != nil {
//				log.WithFields(log.Fields{
//					"gateway_id": hex.EncodeToString(downlinkTXAck.GatewayId),
//					"token":      downlinkTXAck.Token,
//					"ctx_id":     ctxID,
//				}).WithError(err).Error("uplink: handle downlink tx ack error")
//			}
//
//		}(downlinkTXAck)
//	}
//}
// lora 协议上行帧处理
//func collectUplinkFrames(ctx context.Context, uplinkFrame gw.UplinkFrame) error {
//	// collectAndCallOnce 第三个参数是回调函数
//	return collectAndCallOnce(storage.RedisPool(), uplinkFrame, func(rxPacket models.RXPacket) error {
//		var uplinkIDs []uuid.UUID
//		for _, p := range rxPacket.RXInfoSet { // 这个set是按信号排序完毕的结构体对象
//			uplinkIDs = append(uplinkIDs, helpers.GetUplinkID(p))
//		}
//
//		log.WithFields(log.Fields{
//			"uplink_ids": uplinkIDs,
//			"mtype":      rxPacket.PHYPayload.MHDR.MType,
//			"ctx_id":     ctx.Value(logging.ContextIDKey),
//		}).Info("uplink: frame(s) collected")
//
//		// update the gateway meta-data
//		if err := gateway.UpdateMetaDataInRxInfoSet(ctx, storage.DB(), storage.RedisPool(), rxPacket.RXInfoSet); err != nil {
//			log.WithError(err).Error("uplink: update gateway meta-data in rx-info set error")
//		}
//
//		// log the frame for each receiving gatewa
//		if err := framelog.LogUplinkFrameForGateways(ctx, storage.RedisPool(), gw.UplinkFrameSet{
//			PhyPayload: uplinkFrame.PhyPayload,
//			TxInfo:     rxPacket.TXInfo,
//			RxInfo:     rxPacket.RXInfoSet,
//		}); err != nil {
//			log.WithFields(log.Fields{
//				"ctx_id": ctx.Value(logging.ContextIDKey),
//			}).WithError(err).Error("uplink: log uplink frames for gateways error")
//		}
//
//		// handle the frame based on message-type
//		switch rxPacket.PHYPayload.MHDR.MType { // 开始解析应用层payload数据
//		case lorawan.JoinRequest:
//			return join.Handle(ctx, rxPacket)
//		case lorawan.RejoinRequest:
//			return rejoin.Handle(ctx, rxPacket)
//		case lorawan.UnconfirmedDataUp, lorawan.ConfirmedDataUp:
//			return data.Handle(ctx, rxPacket)
//		case lorawan.Proprietary:
//			return proprietary.Handle(ctx, rxPacket)
//		default:
//			return nil
//		}
//	})
//}
