package marshaler

import (
	"bytes"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	pb "iot_gateway/backend/proto"
)

// UnmarshalUplinkFrame unmarshals an UplinkFrame.
func UnmarshalUplinkFrame(b []byte, uf *pb.UplinkFrame) (Type, error) {
	var t Type

	t = Protobuf
	switch t {
	case Protobuf:
		return t, proto.Unmarshal(b, uf)
	case JSON:
		m := jsonpb.Unmarshaler{
			AllowUnknownFields: true,
		}
		return t, m.Unmarshal(bytes.NewReader(b), uf)
	}

	return t, nil
}
