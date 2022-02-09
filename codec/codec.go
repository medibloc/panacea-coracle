package codec

import (
	"bytes"
	"github.com/gogo/protobuf/jsonpb"
	"github.com/gogo/protobuf/proto"
	"strings"
)

func ProtoMarshalJSON(msg proto.Message) ([]byte, error) {
	jm := &jsonpb.Marshaler{}

	buf := new(bytes.Buffer)

	if err := jm.Marshal(buf, msg); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func ProtoUnmarshalJSON(bz []byte, ptr proto.Message) error {
	unmarshaler := jsonpb.Unmarshaler{}
	err := unmarshaler.Unmarshal(strings.NewReader(string(bz)), ptr)
	if err != nil {
		return err
	}

	return nil
}