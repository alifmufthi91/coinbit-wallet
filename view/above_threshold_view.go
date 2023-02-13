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
		err := view.Run(ctx)
		if err != nil {
			logger.Error("Error running aboveThresholdView: %v", err)
			panic(err)
		}
		return err
	}
}
