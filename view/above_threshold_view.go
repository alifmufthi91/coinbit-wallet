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

type AboveThresholdView struct {
	view *goka.View
}

func NewAboveThresholdView(brokers []string) (*AboveThresholdView, error) {
	aboveThresholdView, err := goka.NewView(
		brokers,
		config.AboveThresholdTable,
		new(util.AboveThresholdCodec),
	)
	if err != nil {
		return nil, err
	}
	return &AboveThresholdView{
		view: aboveThresholdView,
	}, nil
}

func (at *AboveThresholdView) Run(ctx context.Context) error {
	logger.Info("Running above threshold View..")
	return at.view.Run(ctx)
}

func (abt *AboveThresholdView) GetByKey(key string) (*model.AboveThreshold, error) {
	logger.Info("Get aboveThreshold by key, key = %s", key)
	val, err := abt.view.Get(key)
	if err != nil {
		return nil, err
	}
	if val == nil {
		return nil, errors.New("aboveThreshold is not found")
	}
	var aboveThreshold *model.AboveThreshold
	var ok bool
	if aboveThreshold, ok = val.(*model.AboveThreshold); !ok {
		return nil, errors.New("failed to cast type")
	}
	return aboveThreshold, nil
}
