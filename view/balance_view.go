package view

import (
	"coinbit-wallet/config"
	"coinbit-wallet/util"
	"coinbit-wallet/util/logger"
	"context"

	"github.com/lovoo/goka"
)

func CreateBalanceView(brokers []string) *goka.View {
	balanceView, err := goka.NewView(
		brokers,
		config.BalanceTable,
		new(util.BalanceMapCodec),
	)
	if err != nil {
		panic(err)
	}
	return balanceView
}

func RunBalanceView(view *goka.View, ctx context.Context) func() error {
	return func() error {
		logger.Info("Running balance View..")
		err := view.Run(ctx)
		if err != nil {
			logger.Error("Error running balanceView: %v", err)
			panic(err)
		}
		return err
	}
}
