package util

import (
	"coinbit-wallet/generated/model"
	"errors"

	"google.golang.org/protobuf/proto"
)

type BalanceCodec struct{}

func (c *BalanceCodec) Encode(value interface{}) ([]byte, error) {
	if balance, ok := value.(*model.Balance); ok {
		return proto.Marshal(balance)
	}
	return nil, errors.New("fail to cast model")
}

func (c *BalanceCodec) Decode(data []byte) (interface{}, error) {
	var m model.Balance
	return &m, proto.Unmarshal(data, &m)
}
