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
	"github.com/stretchr/testify/suite"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type BalanceProcessorSuite struct {
	suite.Suite
	balanceProcessor processor.BalanceProcessor
	gkt              *tester.Tester
	proc             *goka.Processor
}

func TestBalanceProcessorSuite(t *testing.T) {
	suite.Run(t, new(BalanceProcessorSuite))
}

func (bp *BalanceProcessorSuite) SetupSuite() {
	bp.balanceProcessor = processor.NewBalanceProcessor()
	bp.gkt = tester.New(bp.T())

	proc, err := goka.NewProcessor([]string{}, bp.balanceProcessor.Group,
		goka.WithTester(bp.gkt),
	)

	require.Nil(bp.T(), err)

	go proc.Run(context.Background())
}

func (bp *BalanceProcessorSuite) BeforeTest(_, _ string) {
	bp.gkt.ClearValues()
}

func (bp *BalanceProcessorSuite) TestBalanceProcessor_Process() {

	deposit := model.Deposit{
		WalletId:    "111-222",
		Amount:      2000,
		DepositedAt: timestamppb.Now(),
	}

	balance := model.Balance{
		WalletId: deposit.WalletId,
		Balance:  float32(1000),
	}

	bp.gkt.SetTableValue(config.BalanceTable, balance.WalletId, &balance)

	bp.gkt.Consume(string(config.TopicDeposit), deposit.WalletId, &deposit)

	value := bp.gkt.TableValue(config.BalanceTable, deposit.WalletId)
	received := value.(*model.Balance)

	// check wallet id
	require.Equal(bp.T(), balance.WalletId, received.WalletId)

	// check if balance is already added by 2000
	balance.Balance += 2000
	require.Equal(bp.T(), balance.Balance, received.Balance)

}
