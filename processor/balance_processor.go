package processor

import (
	"coinbit-wallet/config"
	"coinbit-wallet/generated/model"
	"coinbit-wallet/util"
	"coinbit-wallet/util/logger"
	"context"
	"fmt"

	"github.com/lovoo/goka"
)

type BalanceProcessor struct {
	group *goka.GroupGraph
}

func NewBalanceProcessor() BalanceProcessor {
	balanceGroup := goka.DefineGroup(config.BalanceGroup,
		goka.Input(config.TopicDeposit, new(util.DepositCodec), processBalance),
		goka.Persist(new(util.BalanceCodec)),
	)
	bp := BalanceProcessor{
		group: balanceGroup,
	}
	return bp
}

func RunBalanceProcessor(ctx context.Context, brokers []string) error {
	logger.Info("Running balance processor..")

	balanceProcessor := NewBalanceProcessor()
	processor, err := goka.NewProcessor(brokers,
		balanceProcessor.group,
		goka.WithTopicManagerBuilder(goka.TopicManagerBuilderWithTopicManagerConfig(config.TMC)),
		goka.WithConsumerGroupBuilder(goka.DefaultConsumerGroupBuilder),
	)
	if err != nil {
		logger.Error("error creating processor: %v", err)
		panic(err)
	}
	err = processor.Run(ctx)
	if err != nil {
		logger.Error("Error running balanceProcessor: %v", err)
	}
	return err
}

func processBalance(ctx goka.Context, msg interface{}) {
	logger.Info("process balance, data = %v", msg)
	var balance *model.Balance
	var ok bool

	deposit := msg.(*model.Deposit)

	if val := ctx.Value(); val != nil {
		balance, ok = val.(*model.Balance)
		if !ok {
			ctx.Fail(fmt.Errorf("processing failed due to casting failure"))
		}
	} else {
		balance = &model.Balance{}
	}

	balance.Balance += deposit.GetAmount()
	ctx.SetValue(balance)
	logger.Info("wallet balance is procesed")
}
