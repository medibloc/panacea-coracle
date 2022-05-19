package codec

import (
	"encoding/json"
	"github.com/gogo/protobuf/jsonpb"
	"github.com/gogo/protobuf/proto"
	"strings"
)

func ProtoMarshalJSON(msg proto.Message) ([]byte, error) {
	marshaledResp, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}

	return marshaledResp, nil
}

func ProtoUnmarshalJSON(bz []byte, ptr proto.Message) error {
	unmarshaler := jsonpb.Unmarshaler{}
	err := unmarshaler.Unmarshal(strings.NewReader(string(bz)), ptr)
	if err != nil {
		return err
	}

	return nil
}