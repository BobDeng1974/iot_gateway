package downlink

import (
	"github.com/prometheus/common/log"
	"open/backend/gateway"
	pb "open/backend/proto"
	"open/config"
)

// DeviceQueueSchedulerLoop starts an infinit loop calling the scheduler loop for Class-B
// and Class-C sheduling.
var downPack  chan pb.DownlinkFrame

func deviceQueueSchedulerLoop() {

	for downPack := range downPack {
		// send the packet to the gateway
		if err := gateway.Backend().SendTXPacket(downPack); err != nil {
			log.Error("send downlink frame error",err)
		}
	}
}
// 面向应用层的下行
func AddDownlinkFrame(dw pb.DownlinkFrame){
	downPack <- dw
}

func Setup(conf config.Config) error {

	downPack = make( chan pb.DownlinkFrame ,conf.General.DownChannelSize)
	go deviceQueueSchedulerLoop()

	return nil
}