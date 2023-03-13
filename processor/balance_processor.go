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
	Processor *goka.Processor
}

func NewBalanceProcessor(brokers []string, options ...goka.ProcessorOption) (*BalanceProcessor, error) {
	logger.Info("Create new balance processor..")
	balanceGroup := goka.DefineGroup(config.BalanceGroup,
		goka.Input(config.TopicDeposit, new(util.DepositCodec), processBalance),
		goka.Persist(new(util.BalanceCodec)),
	)
	processor, err := goka.NewProcessor(brokers,
		balanceGroup,
		options...,
	)
	if err != nil {
		logger.Error("error creating processor: %v", err)
		return nil, err
	}
	return &BalanceProcessor{
		Processor: processor,
	}, nil
}

func (bp *BalanceProcessor) Run(ctx context.Context) error {
	logger.Info("Running balance processor..")
	return bp.Processor.Run(ctx)
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
