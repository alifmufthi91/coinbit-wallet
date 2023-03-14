package processor

import (
	"coinbit-wallet/config"
	"coinbit-wallet/generated/model"
	"coinbit-wallet/util"
	"coinbit-wallet/util/logger"
	"context"
	"fmt"

	"github.com/lovoo/goka"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type AboveThresholdProcessor struct {
	Processor *goka.Processor
}

func NewAboveThresholdProcessor(brokers []string, options ...goka.ProcessorOption) (*AboveThresholdProcessor, error) {
	logger.Info("Create new above threshold processor..")
	aboveThresholdGroup := goka.DefineGroup(config.AboveThresholdGroup,
		goka.Input(config.TopicDeposit, new(util.DepositCodec), processAboveThreshold),
		goka.Persist(new(util.AboveThresholdCodec)),
	)
	processor, err := goka.NewProcessor(brokers,
		aboveThresholdGroup,
		options...,
	)
	if err != nil {
		logger.Error("error creating processor: %v", err)
		return nil, err
	}
	return &AboveThresholdProcessor{
		Processor: processor,
	}, nil
}

func (atp *AboveThresholdProcessor) Run(ctx context.Context) error {
	logger.Info("Running above threshold processor..")
	return atp.Processor.Run(ctx)
}

func processAboveThreshold(ctx goka.Context, msg interface{}) {
	logger.Info("process above threshold, data = %v", msg)
	var aboveThreshold *model.AboveThreshold
	var ok bool

	deposit, ok := msg.(*model.Deposit)
	if !ok {
		ctx.Fail(fmt.Errorf("processing failed due to casting failure"))
	}

	if val := ctx.Value(); val != nil {
		aboveThreshold, ok = val.(*model.AboveThreshold)
		if !ok {
			ctx.Fail(fmt.Errorf("processing failed due to casting failure"))
		}
	} else {
		aboveThreshold = &model.AboveThreshold{
			WalletId:    deposit.GetWalletId(),
			Status:      false,
			StartPeriod: deposit.DepositedAt,
		}
	}

	// if deposit already past two minutes period
	if !util.IsWithinTwoMins(aboveThreshold.StartPeriod, deposit.DepositedAt) {
		resetAboveThreshold(aboveThreshold, deposit.DepositedAt)
	}
	aboveThreshold.AmountWithinTwoMins += deposit.GetAmount()
	if aboveThreshold.AmountWithinTwoMins > 10000 {
		aboveThreshold.Status = true
	}
	ctx.SetValue(aboveThreshold)
	logger.Info("wallet above threshold is procesed")
}

func resetAboveThreshold(aboveThreshold *model.AboveThreshold, currentTime *timestamppb.Timestamp) {
	aboveThreshold.AmountWithinTwoMins = 0
	aboveThreshold.Status = false
	aboveThreshold.StartPeriod = currentTime
}
