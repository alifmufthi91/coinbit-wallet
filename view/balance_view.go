package view

import (
	"coinbit-wallet/processor"
	"coinbit-wallet/util"
	"coinbit-wallet/util/logger"
	"context"

	"github.com/lovoo/goka"
)

func CreateBalanceView(brokers []string) *goka.View {
	balanceView, err := goka.NewView(
		brokers,
		processor.BalanceTable,
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
		return view.Run(ctx)
	}
}
