package uplink

import (
	"context"
	"fmt"
	"iot_gateway/config"
	"iot_gateway/uplink/data"
	"iot_gateway/uplink/proprietary"
	"sync"
	"time"

	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"iot_gateway/backend/gateway"
	pb "iot_gateway/backend/proto"
	"iot_gateway/lib/logging"
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

// 把解析器结果，附加上id，发送给上层应用服务。
func HandleUplinkFrames(wg *sync.WaitGroup) {
	for uplinkFrame := range gateway.Backend().RXPacketChan() { //遍历这个管道。可是这个管道是无缓冲阻塞的。
		go func(uplinkFrame *pb.UplinkFrame) { //每接收到一条上行消息，就开一个协成去处理
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

			if err := distributeUplinkFrame(ctx, uplinkFrame); err != nil {
				log.WithFields(log.Fields{
					"ctx_id": ctxID,
				}).WithError(err).Error("uplink: processing uplink frame error")
			}
		}(uplinkFrame)
	}
}

// 根据帧头类型，分发帧。在lora中网络层ack帧，可以不给应用层
func distributeUplinkFrame(ctx context.Context, uplinkFrame *pb.UplinkFrame) error {

	switch uplinkFrame.FrameType {
	case gateway.UnconfirmedDataUp,gateway.ConfirmedDataUp:
		return data.Handle(ctx, uplinkFrame)
	case gateway.Proprietary:
		return proprietary.Handle(ctx,uplinkFrame)
	default:
		log.Info("[HandleUplinkFrame]unknown frame type")
		return nil
	}
	return nil
}
