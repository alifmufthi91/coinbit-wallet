package util

import (
	"coinbit-wallet/generated/model"
	"errors"

	"google.golang.org/protobuf/proto"
)

type BalanceMapCodec struct{}

func (c *BalanceMapCodec) Encode(value interface{}) ([]byte, error) {
	if balanceMap, ok := value.(*model.BalanceMap); ok {
		return proto.Marshal(balanceMap)
	}
	return nil, errors.New("fail to cast model")
}

func (c *BalanceMapCodec) Decode(data []byte) (interface{}, error) {
	var m model.BalanceMap
	return &m, proto.Unmarshal(data, &m)
}
