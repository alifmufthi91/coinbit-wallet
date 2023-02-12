package util

import (
	"coinbit-wallet/generated/model"
	"errors"

	"google.golang.org/protobuf/proto"
)

type DepositCodec struct{}

func (c *DepositCodec) Encode(value interface{}) ([]byte, error) {
	if deposit, ok := value.(*model.Deposit); ok {
		return proto.Marshal(deposit)
	}
	return nil, errors.New("fail to cast model")
}

func (c *DepositCodec) Decode(data []byte) (interface{}, error) {
	var m model.Deposit
	return &m, proto.Unmarshal(data, &m)
}
