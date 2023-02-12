package util

import (
	"coinbit-wallet/generated/model"
	"errors"

	"google.golang.org/protobuf/proto"
)

type AboveThresholdMapCodec struct{}

func (c *AboveThresholdMapCodec) Encode(value interface{}) ([]byte, error) {
	if aboveThresholdMap, ok := value.(*model.AboveThresholdMap); ok {
		return proto.Marshal(aboveThresholdMap)
	}
	return nil, errors.New("fail to cast model")
}

func (c *AboveThresholdMapCodec) Decode(data []byte) (interface{}, error) {
	var m model.AboveThresholdMap
	return &m, proto.Unmarshal(data, &m)
}
