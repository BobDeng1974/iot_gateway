package gateway

import (
	pb "open/backend/proto"
)


const (
	UnconfirmedDataUp  = iota
	ConfirmedDataUp
	Proprietary
	UnconfirmedDataDown
	ConfirmedDataDown
)

//
//type UplinkFrame struct {
//	FrameType int8 `json:"frame_type"`
//	DevAddr [4]byte `json:"dev_addr"` //设备唯一地址或是id
//	FrameNum uint32 `json:"frame_num"` //帧号，用来统计有没有丢包。
//	Port uint32 `json:"port"` //端口
//	PhyPayload []byte `json:"phy_payload"` // proto打包的应用层数据。aes加密数据
//
//	UplinkId uint16 `json:"uplink_id"`
//	// PHYPayload.
//	//PhyPayload []byte `protobuf:"bytes,1,opt,name=phy_payload,json=phyPayload,proto3" json:"phy_payload,omitempty"`
//	//// TX meta-data.
//	//TxInfo *UplinkTXInfo `protobuf:"bytes,2,opt,name=tx_info,json=txInfo,proto3" json:"tx_info,omitempty"`
//	//// RX meta-data.
//	//RxInfo               *UplinkRXInfo `protobuf:"bytes,3,opt,name=rx_info,json=rxInfo,proto3" json:"rx_info,omitempty"`
//}
//func (m *UplinkFrame) Reset()         { *m = UplinkFrame{} }
//func (m *UplinkFrame) String() string { return proto.CompactTextString(m) }
//func (*UplinkFrame) ProtoMessage()    {}
//type DownlinkFrame struct {
//	FrameType int8 `json:"frame_type"`
//	DevAddr [4]byte `json:"dev_addr"` //设备唯一地址或是id
//	FrameNum uint32 `json:"frame_num"` //帧号，用来统计有没有丢包。
//	Port uint32 `json:"port"` //端口
//	PhyPayload []byte `json:"phy_payload"` // proto打包的应用层数据。aes加密数据
//	// 非下行数据
//	DownlinkId []byte
//}
//
//func (m *DownlinkFrame) Reset()         { *m = DownlinkFrame{} }
//func (m *DownlinkFrame) String() string { return proto.CompactTextString(m) }
//func (*DownlinkFrame) ProtoMessage()    {}
// Gateway is the interface of a gateway backend.
// A gateway backend is responsible for the communication with the gateway.
type Gateway interface {
	SendTXPacket(pb.DownlinkFrame) error                   // send the given packet to the gateway
	//SendGatewayConfigPacket(gw.GatewayConfiguration) error // SendGatewayConfigPacket sends the given GatewayConfigPacket to the gateway.
	RXPacketChan() chan pb.UplinkFrame                     // channel containing the received packets
	//StatsPacketChan() chan gw.GatewayStats                 // channel containing the received gateway stats
	//DownlinkTXAckChan() chan gw.DownlinkTXAck              // channel containing the downlink tx acknowledgements
	Close() error                                          // close the gateway backend.
}


var backend Gateway


// Backend returns the gateway backend.
func Backend() Gateway {
	return backend
}

// SetBackend sets the given gateway backend.
//
func SetBackend(b Gateway) {
	backend = b
}

