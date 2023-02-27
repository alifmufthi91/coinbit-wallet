package processor

import (
	"coinbit-wallet/config"
	"coinbit-wallet/generated/model"
	"coinbit-wallet/util"
	"coinbit-wallet/util/logger"
	"context"

	"github.com/lovoo/goka"
)

type AboveThresholdProcessor struct {
	group *goka.GroupGraph
}

func NewAboveThresholdProcessor() AboveThresholdProcessor {
	aboveThresholdGroup := goka.DefineGroup(config.AboveThresholdGroup,
		goka.Input(config.TopicDeposit, new(util.DepositCodec), processAboveThreshold),
		goka.Persist(new(util.AboveThresholdCodec)),
	)
	atp := AboveThresholdProcessor{
		group: aboveThresholdGroup,
	}
	return atp
}

func RunAboveThresholdProcessor(ctx context.Context, brokers []string) error {
	logger.Info("Running above threshold processor..")
	aboveThresholdProcessor := NewAboveThresholdProcessor()
	processor, err := goka.NewProcessor(brokers,
		aboveThresholdProcessor.group,
		goka.WithTopicManagerBuilder(goka.TopicManagerBuilderWithTopicManagerConfig(config.TMC)),
		goka.WithConsumerGroupBuilder(goka.DefaultConsumerGroupBuilder),
	)
	if err != nil {
		logger.Error("error creating processor: %v", err)
		panic(err)
	}
	err = processor.Run(ctx)
	if err != nil {
		logger.Error("Error running aboveThresholdProcessor: %v", err)
		panic(err)
	}
	return err
}

func processAboveThreshold(ctx goka.Context, msg interface{}) {
	logger.Info("process above threshold, data = %v", msg)
	var aboveThreshold *model.AboveThreshold

	deposit := msg.(*model.Deposit)

	if val := ctx.Value(); val != nil {
		aboveThreshold = val.(*model.AboveThreshold)
	} else {
		aboveThreshold = &model.AboveThreshold{
			WalletId:    deposit.GetWalletId(),
			Status:      false,
			StartPeriod: deposit.DepositedAt,
		}
	}

	// if deposit already past two minutes period
	if deposit.DepositedAt.Seconds > aboveThreshold.StartPeriod.Seconds+120 {
		aboveThreshold.AmountWithinTwoMins = 0
		aboveThreshold.Status = false
		aboveThreshold.StartPeriod = deposit.DepositedAt
	}
	aboveThreshold.AmountWithinTwoMins += deposit.GetAmount()
	if aboveThreshold.AmountWithinTwoMins > 10000 {
		aboveThreshold.Status = true
	}
	ctx.SetValue(aboveThreshold)
	logger.Info("wallet above threshold is procesed")
}
