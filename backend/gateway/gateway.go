package gateway

import (
	pb "iot_gateway/backend/proto"
)


const (
	UnconfirmedDataUp  = iota
	ConfirmedDataUp
	Proprietary
	UnconfirmedDataDown
	ConfirmedDataDown
)
var backend Gateway

type Gateway interface {
	SendTXPacket(pb.DownlinkFrame) error                   // send the given packet to the gateway
	//SendGatewayConfigPacket(gw.GatewayConfiguration) error // SendGatewayConfigPacket sends the given GatewayConfigPacket to the gateway.
	RXPacketChan() chan *pb.UplinkFrame                     // channel containing the received packets
	//StatsPacketChan() chan gw.GatewayStats                 // channel containing the received gateway stats
	//DownlinkTXAckChan() chan gw.DownlinkTXAck              // channel containing the downlink tx acknowledgements
	Close() error                                    // close the gateway backend.
}

func Backend() Gateway {
	return backend
}

func SetBackend(b Gateway) {
	backend = b
}

