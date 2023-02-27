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
		new(util.BalanceCodec),
	)
	if err != nil {
		panic(err)
	}
	return balanceView
}

func RunBalanceView(view *goka.View, ctx context.Context) error {
	logger.Info("Running balance View..")
	return view.Run(ctx)
}
