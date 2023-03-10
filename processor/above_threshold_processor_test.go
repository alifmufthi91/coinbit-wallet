package processor_test

import (
	"coinbit-wallet/config"
	"coinbit-wallet/generated/model"
	"coinbit-wallet/processor"
	"context"
	"testing"
	"time"

	"github.com/lovoo/goka"
	"github.com/lovoo/goka/tester"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type AboveThresholdProcessorSuite struct {
	suite.Suite
	aboveThresholdProcessor processor.AboveThresholdProcessor
	gkt                     *tester.Tester
}

func TestAboveThresholdProcessorSuite(t *testing.T) {
	suite.Run(t, new(AboveThresholdProcessorSuite))
}

func (at *AboveThresholdProcessorSuite) SetupSuite() {
	at.aboveThresholdProcessor = processor.NewAboveThresholdProcessor()
	at.gkt = tester.New(at.T())
	proc, err := goka.NewProcessor([]string{}, at.aboveThresholdProcessor.Group,
		goka.WithTester(at.gkt),
	)

	require.Nil(at.T(), err)

	go proc.Run(context.Background())
}

func (at *AboveThresholdProcessorSuite) BeforeTest(_, _ string) {
	at.gkt.ClearValues()
}

func (at *AboveThresholdProcessorSuite) TestAboveThresholdProcessor_Process() {

	var (
		expectedStatus bool
		expectedAmount float32
	)

	aboveThreshold := model.AboveThreshold{
		WalletId:            "111-222",
		AmountWithinTwoMins: float32(0),
		Status:              false,
		StartPeriod:         timestamppb.Now(),
	}

	deposit := model.Deposit{
		WalletId:    aboveThreshold.WalletId,
		Amount:      7000,
		DepositedAt: timestamppb.Now(),
	}

	at.gkt.SetTableValue(config.AboveThresholdTable, aboveThreshold.WalletId, &aboveThreshold)

	at.gkt.Consume(string(config.TopicDeposit), deposit.WalletId, &deposit)

	value := at.gkt.TableValue(config.AboveThresholdTable, aboveThreshold.WalletId)
	received := value.(*model.AboveThreshold)

	// check wallet id
	require.Equal(at.T(), aboveThreshold.WalletId, received.WalletId)

	// check first deposit amount and status
	aboveThreshold.AmountWithinTwoMins += 7000
	expectedStatus = false
	expectedAmount = aboveThreshold.AmountWithinTwoMins
	require.Equal(at.T(), expectedAmount, received.AmountWithinTwoMins)
	require.Equal(at.T(), expectedStatus, received.Status)

	// deposit another to make amount within two mins above 10.000
	deposit.DepositedAt = timestamppb.New(time.Now().Add(5 * time.Second))
	at.gkt.Consume(string(config.TopicDeposit), deposit.WalletId, &deposit)

	value = at.gkt.TableValue(config.AboveThresholdTable, aboveThreshold.WalletId)
	received = value.(*model.AboveThreshold)

	expectedAmount += 7000
	expectedStatus = true

	require.Equal(at.T(), expectedAmount, received.AmountWithinTwoMins)
	require.Equal(at.T(), expectedStatus, received.Status)

}
