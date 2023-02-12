package view

import (
	"coinbit-wallet/config"
	"coinbit-wallet/processor"
	"coinbit-wallet/util"
	"coinbit-wallet/util/logger"
	"context"

	"github.com/lovoo/goka"
)

func CreateAboveThresholdView(brokers []string) *goka.View {
	aboveThresholdView, err := goka.NewView(
		config.Brokers,
		processor.AboveThresholdTable,
		new(util.AboveThresholdMapCodec),
	)
	if err != nil {
		panic(err)
	}
	return aboveThresholdView
}

func RunAboveThresholdView(view *goka.View, ctx context.Context) func() error {
	return func() error {
		logger.Info("Running Above Threshold View..")
		return view.Run(ctx)
	}
}
