package processor

import (
	"coinbit-wallet/config"
	"coinbit-wallet/generated/model"
	"coinbit-wallet/util"
	"coinbit-wallet/util/logger"
	"context"

	"github.com/lovoo/goka"
)

var (
	AboveThresholdGroup goka.Group = "aboveThreshold"
	AboveThresholdTable goka.Table = goka.GroupTable(AboveThresholdGroup)
)

func RunAboveThresholdProcessor(ctx context.Context, brokers []string) func() error {
	return func() error {
		logger.Info("Running above threshold processor..")

		aboveThresholdGroup := goka.DefineGroup(AboveThresholdGroup,
			goka.Input(config.TopicDeposit, new(util.DepositCodec), processAboveThreshold),
			goka.Persist(new(util.AboveThresholdMapCodec)),
		)
		aboveThresholdProcessor, err := goka.NewProcessor(brokers,
			aboveThresholdGroup,
			goka.WithTopicManagerBuilder(goka.TopicManagerBuilderWithTopicManagerConfig(config.TMC)),
			goka.WithConsumerGroupBuilder(goka.DefaultConsumerGroupBuilder),
		)
		if err != nil {
			logger.Error("error creating processor: %v", err)
			panic(err)
		}

		return aboveThresholdProcessor.Run(ctx)
	}
}

func processAboveThreshold(ctx goka.Context, msg interface{}) {
	logger.Info("process above threshold, data = %v", msg)
	var aboveThresholdMap *model.AboveThresholdMap

	if val := ctx.Value(); val != nil {
		aboveThresholdMap = val.(*model.AboveThresholdMap)
	}

	deposit := msg.(*model.Deposit)

	if aboveThresholdMap.Items == nil {
		aboveThresholdMap.Items = make(map[string]*model.AboveThreshold)
	}
	aboveThreshold, ok := aboveThresholdMap.Items[deposit.GetWalletId()]
	if !ok {
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
	aboveThresholdMap.Items[aboveThreshold.GetWalletId()] = aboveThreshold
	ctx.SetValue(aboveThresholdMap)
	logger.Info("wallet above threshold is procesed")
}
