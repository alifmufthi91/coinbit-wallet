package view

import (
	"coinbit-wallet/config"
	"coinbit-wallet/generated/model"
	"coinbit-wallet/util"
	"coinbit-wallet/util/logger"
	"context"

	"github.com/lovoo/goka"
)

var (
	BalanceView *goka.View
)

func RunBalanceView() {
	logger.Info("Running balance View..")
	var err error
	BalanceView, err = goka.NewView(
		config.Brokers,
		config.BalanceTable,
		new(util.BalanceMapCodec),
	)
	if err != nil {
		panic(err)
	}
	err = BalanceView.Run(context.Background())
	if err != nil {
		logger.Error("error running view: %v", err)
		panic(err)
	}
}

func GetBalanceView() *model.BalanceMap {
	val, err := BalanceView.Get(string(config.TopicDeposit))
	if err != nil {
		panic(err)
	}
	if val == nil {
		panic("view is not found")
	}
	return val.(*model.BalanceMap)
}
