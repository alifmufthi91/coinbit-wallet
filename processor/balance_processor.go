package processor

import (
	"coinbit-wallet/config"
	"coinbit-wallet/generated/model"
	"coinbit-wallet/util"
	"coinbit-wallet/util/logger"
	"context"
	"sync"

	"github.com/lovoo/goka"
)

var (
	balanceMtx sync.Mutex
)

func RunBalanceProcessor(ctx context.Context, brokers []string) func() error {
	return func() error {
		logger.Info("Running balance processor..")

		balanceGroup := goka.DefineGroup(config.BalanceGroup,
			goka.Input(config.TopicDeposit, new(util.DepositCodec), processBalance),
			goka.Persist(new(util.BalanceMapCodec)),
		)
		balanceProcessor, err := goka.NewProcessor(brokers,
			balanceGroup,
			goka.WithTopicManagerBuilder(goka.TopicManagerBuilderWithTopicManagerConfig(config.TMC)),
			goka.WithConsumerGroupBuilder(goka.DefaultConsumerGroupBuilder),
		)
		if err != nil {
			logger.Error("error creating processor: %v", err)
			panic(err)
		}

		err = balanceProcessor.Run(ctx)
		if err != nil {
			logger.Error("Error running balanceProcessor: %v", err)
			panic(err)
		}
		return err
	}
}

func processBalance(ctx goka.Context, msg interface{}) {
	logger.Info("process balance, data = %v", msg)
	balanceMtx.Lock()
	var balanceMap *model.BalanceMap

	if val := ctx.Value(); val != nil {
		balanceMap = val.(*model.BalanceMap)
	} else {
		balanceMap = &model.BalanceMap{}
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
	balanceMtx.Unlock()
	logger.Info("wallet balance is procesed")
}
