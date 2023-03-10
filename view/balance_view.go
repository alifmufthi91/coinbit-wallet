package view

import (
	"coinbit-wallet/config"
	"coinbit-wallet/generated/model"
	"coinbit-wallet/util"
	"coinbit-wallet/util/logger"
	"context"
	"errors"

	"github.com/lovoo/goka"
)

type BalanceView struct {
	view *goka.View
}

func NewBalanceView(brokers []string) (*BalanceView, error) {
	balanceView, err := goka.NewView(
		brokers,
		config.BalanceTable,
		new(util.BalanceCodec),
	)
	if err != nil {
		return nil, err
	}
	return &BalanceView{
		view: balanceView,
	}, nil
}

func (bv *BalanceView) Run(ctx context.Context) error {
	logger.Info("Running balance View..")
	return bv.view.Run(ctx)
}

func (bv *BalanceView) GetByKey(key string) (*model.Balance, error) {
	logger.Info("Get balance by key, key = %s", key)
	val, err := bv.view.Get(key)
	if err != nil {
		return nil, err
	}
	if val == nil {
		return nil, errors.New("balance is not found")
	}
	var balance *model.Balance
	var ok bool
	if balance, ok = val.(*model.Balance); !ok {
		return nil, errors.New("failed to cast type")
	}
	return balance, nil
}
