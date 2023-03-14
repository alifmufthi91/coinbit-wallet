package processor_test

import (
	"coinbit-wallet/config"
	"coinbit-wallet/generated/model"
	"coinbit-wallet/processor"
	"context"
	"testing"

	"github.com/go-test/deep"
	"github.com/lovoo/goka"
	"github.com/lovoo/goka/tester"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type BalanceProcessorSuite struct {
	suite.Suite
	balanceProcessor *processor.BalanceProcessor
	gkt              *tester.Tester
	proc             *goka.Processor
}

func TestBalanceProcessorSuite(t *testing.T) {
	suite.Run(t, new(BalanceProcessorSuite))
}

func (bp *BalanceProcessorSuite) SetupSuite() {
	bp.gkt = tester.New(bp.T())

	var err error
	bp.balanceProcessor, err = processor.NewBalanceProcessor([]string{}, goka.WithTester(bp.gkt))

	require.Nil(bp.T(), err)

	go bp.balanceProcessor.Run(context.Background())
}

func (bp *BalanceProcessorSuite) BeforeTest(_, _ string) {
	bp.gkt.ClearValues()
}

func (bp *BalanceProcessorSuite) TestBalanceProcessor_Process() {

	testCases := []struct {
		name     string
		deposit  *model.Deposit
		expected *model.Balance
	}{
		{
			name: "first deposit",
			deposit: &model.Deposit{
				WalletId:    "111-222",
				Amount:      2000,
				DepositedAt: timestamppb.Now(),
			},
			expected: &model.Balance{
				WalletId: "111-222",
				Balance:  2000,
			},
		},
		{
			name: "second deposit",
			deposit: &model.Deposit{
				WalletId:    "111-222",
				Amount:      2000,
				DepositedAt: timestamppb.Now(),
			},
			expected: &model.Balance{
				WalletId: "111-222",
				Balance:  4000,
			},
		},
	}

	for i, tc := range testCases {
		bp.Run(tc.name, func() {
			// first deposit
			if i == 0 {
				balance := tc.expected

				bp.gkt.Consume(string(config.TopicDeposit), tc.deposit.WalletId, tc.deposit)
				bp.gkt.SetTableValue(config.BalanceTable, balance.WalletId, balance)

				value := bp.gkt.TableValue(config.BalanceTable, tc.deposit.WalletId)
				received := value.(*model.Balance)

				require.Nil(bp.T(), deep.Equal(tc.expected, received))
			} else {
				bp.gkt.Consume(string(config.TopicDeposit), tc.deposit.WalletId, tc.deposit)

				value := bp.gkt.TableValue(config.BalanceTable, tc.deposit.WalletId)
				received := value.(*model.Balance)

				require.Nil(bp.T(), deep.Equal(tc.expected, received))
			}
		})
	}
}
