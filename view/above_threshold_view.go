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
	AboveThresholdView *goka.View
)

func RunAboveThresholdView() {
	logger.Info("Running Above Threshold View..")
	var err error
	AboveThresholdView, err = goka.NewView(
		config.Brokers,
		config.AboveThresholdTable,
		new(util.AboveThresholdMapCodec),
	)
	if err != nil {
		panic(err)
	}
	err = AboveThresholdView.Run(context.Background())
	if err != nil {
		logger.Error("error running view: %v", err)
		panic(err)
	}
}

func GetAboveThresholdView() *model.AboveThresholdMap {
	val, err := AboveThresholdView.Get(string(config.TopicDeposit))
	if err != nil {
		panic(err)
	}
	if val == nil {
		panic("view is not found")
	}
	return val.(*model.AboveThresholdMap)
}
