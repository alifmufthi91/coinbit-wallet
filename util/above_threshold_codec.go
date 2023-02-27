package util

import (
	"coinbit-wallet/generated/model"
	"errors"

	"google.golang.org/protobuf/proto"
)

type AboveThresholdCodec struct{}

func (c *AboveThresholdCodec) Encode(value interface{}) ([]byte, error) {
	if aboveThreshold, ok := value.(*model.AboveThreshold); ok {
		return proto.Marshal(aboveThreshold)
	}
	return nil, errors.New("fail to cast model")
}

func (c *AboveThresholdCodec) Decode(data []byte) (interface{}, error) {
	var m model.AboveThreshold
	return &m, proto.Unmarshal(data, &m)
}
