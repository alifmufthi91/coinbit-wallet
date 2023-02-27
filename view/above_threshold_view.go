package view

import (
	"coinbit-wallet/config"
	"coinbit-wallet/util"
	"coinbit-wallet/util/logger"
	"context"

	"github.com/lovoo/goka"
)

func CreateAboveThresholdView(brokers []string) *goka.View {
	aboveThresholdView, err := goka.NewView(
		brokers,
		config.AboveThresholdTable,
		new(util.AboveThresholdCodec),
	)
	if err != nil {
		panic(err)
	}
	return aboveThresholdView
}

func RunAboveThresholdView(view *goka.View, ctx context.Context) error {
	logger.Info("Running Above Threshold View..")
	return view.Run(ctx)
}
