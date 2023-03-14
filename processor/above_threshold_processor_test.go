package processor_test

import (
	"coinbit-wallet/config"
	"coinbit-wallet/generated/model"
	"coinbit-wallet/processor"
	"context"
	"testing"
	"time"

	"github.com/go-test/deep"
	"github.com/lovoo/goka"
	"github.com/lovoo/goka/tester"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type AboveThresholdProcessorSuite struct {
	suite.Suite
	aboveThresholdProcessor *processor.AboveThresholdProcessor
	gkt                     *tester.Tester
}

func TestAboveThresholdProcessorSuite(t *testing.T) {
	suite.Run(t, new(AboveThresholdProcessorSuite))
}

func (at *AboveThresholdProcessorSuite) SetupSuite() {
	at.gkt = tester.New(at.T())

	var err error
	at.aboveThresholdProcessor, err = processor.NewAboveThresholdProcessor([]string{}, goka.WithTester(at.gkt))

	require.Nil(at.T(), err)

	go at.aboveThresholdProcessor.Run(context.Background())
}

func (at *AboveThresholdProcessorSuite) BeforeTest(_, _ string) {
	at.gkt.ClearValues()
}

func (at *AboveThresholdProcessorSuite) TestAboveThresholdProcessor_Process() {

	testCases := []struct {
		name     string
		deposit  *model.Deposit
		expected *model.AboveThreshold
	}{
		{
			name: "First deposit within threshold",
			deposit: &model.Deposit{
				WalletId:    "111-222",
				Amount:      7000,
				DepositedAt: timestamppb.Now(),
			},
			expected: &model.AboveThreshold{
				WalletId:            "111-222",
				AmountWithinTwoMins: 7000,
				Status:              false,
				StartPeriod:         timestamppb.Now(),
			},
		},
		{
			name: "Deposit within two mins above threshold",
			deposit: &model.Deposit{
				WalletId:    "111-222",
				Amount:      7000,
				DepositedAt: timestamppb.New(time.Now().Add(5 * time.Second)),
			},
			expected: &model.AboveThreshold{
				WalletId:            "111-222",
				AmountWithinTwoMins: 14000,
				Status:              true,
				StartPeriod:         timestamppb.Now(),
			},
		},
		{
			name: "Deposit after two mins above threshold",
			deposit: &model.Deposit{
				WalletId:    "111-222",
				Amount:      7000,
				DepositedAt: timestamppb.New(time.Now().Add(2 * time.Minute)),
			},
			expected: &model.AboveThreshold{
				WalletId:            "111-222",
				AmountWithinTwoMins: 7000,
				Status:              false,
				StartPeriod:         timestamppb.New(time.Now().Add(2 * time.Minute)),
			},
		},
	}

	for i, tc := range testCases {
		at.Run(tc.name, func() {
			// first deposit
			if i == 0 {
				aboveThreshold := tc.expected

				at.gkt.Consume(string(config.TopicDeposit), tc.deposit.WalletId, tc.deposit)
				at.gkt.SetTableValue(config.AboveThresholdTable, aboveThreshold.WalletId, aboveThreshold)

				value := at.gkt.TableValue(config.AboveThresholdTable, tc.deposit.WalletId)
				received := value.(*model.AboveThreshold)

				require.Nil(at.T(), deep.Equal(tc.expected, received))
			} else {
				at.gkt.Consume(string(config.TopicDeposit), tc.deposit.WalletId, tc.deposit)

				value := at.gkt.TableValue(config.AboveThresholdTable, tc.deposit.WalletId)
				received := value.(*model.AboveThreshold)

				require.Nil(at.T(), deep.Equal(tc.expected, received))
			}
		})
	}

}
