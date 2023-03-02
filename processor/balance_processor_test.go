package processor_test

import (
	"coinbit-wallet/config"
	"coinbit-wallet/generated/model"
	"coinbit-wallet/processor"
	"context"
	"testing"

	"github.com/lovoo/goka"
	"github.com/lovoo/goka/tester"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func Test_BalanceProcessor(t *testing.T) {
	gkt := tester.New(t)

	balanceProcessor := processor.NewBalanceProcessor()
	// create a new processor, registering the tester
	proc, _ := goka.NewProcessor([]string{}, balanceProcessor.Group,
		goka.WithTester(gkt),
	)

	go proc.Run(context.Background())

	deposit := model.Deposit{
		WalletId:    "111-222",
		Amount:      2000,
		DepositedAt: timestamppb.Now(),
	}

	balance := model.Balance{
		WalletId: deposit.WalletId,
		Balance:  float32(1000),
	}

	gkt.SetTableValue(config.BalanceTable, balance.WalletId, &balance)

	gkt.Consume(string(config.TopicDeposit), deposit.WalletId, &deposit)

	value := gkt.TableValue(config.BalanceTable, deposit.WalletId)
	received := value.(*model.Balance)

	// check wallet id
	require.Equal(t, balance.WalletId, received.WalletId)

	// check if balance is already added by 2000
	balance.Balance += 2000
	require.Equal(t, balance.Balance, received.Balance)

}
