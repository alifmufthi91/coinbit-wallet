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
	balanceProcessor *goka.Processor
)

func RunBalanceProcessor() {
	logger.Info("Running balance processor..")
	var err error
	brokers := config.GetEnv().KafkaBrokers
	balanceGroup := goka.DefineGroup(config.BalanceGroup,
		goka.Input(config.TopicDeposit, new(util.DepositCodec), processBalance),
		goka.Persist(new(util.BalanceMapCodec)),
	)
	balanceProcessor, err = goka.NewProcessor(brokers,
		balanceGroup,
		goka.WithTopicManagerBuilder(goka.TopicManagerBuilderWithTopicManagerConfig(config.TMC)),
		goka.WithConsumerGroupBuilder(goka.DefaultConsumerGroupBuilder),
	)
	if err != nil {
		logger.Error("error creating processor: %v", err)
		panic(err)
	}

	err = balanceProcessor.Run(context.Background())
	if err != nil {
		logger.Error("error running processor: %v", err)
		panic(err)
	}
}

func processBalance(ctx goka.Context, msg interface{}) {
	logger.Info("process balance, data = %v", msg)
	var balanceMap *model.BalanceMap

	if val := ctx.Value(); val != nil {
		balanceMap = val.(*model.BalanceMap)
	}

	deposit := msg.(*model.Deposit)

	if balanceMap.Items == nil {
		balanceMap.Items = make(map[string]*model.Balance)
	}
	balance, ok := balanceMap.Items[deposit.GetWalletId()]
	if !ok {
		balance = &model.Balance{
			WalletId: deposit.GetWalletId(),
		}
	}
	balance.Balance += deposit.GetAmount()
	balanceMap.Items[balance.GetWalletId()] = balance
	ctx.SetValue(balanceMap)
	logger.Info("wallet balance is procesed")
}
