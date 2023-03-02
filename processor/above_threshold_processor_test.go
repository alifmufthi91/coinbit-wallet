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
	"google.golang.org/protobuf/types/known/timestamppb"
)

func Test_AboveThresholdProcessor(t *testing.T) {

	var (
		expectedStatus bool
		expectedAmount float32
	)

	gkt := tester.New(t)

	aboveThresholdProcessor := processor.NewAboveThresholdProcessor()
	// create a new processor, registering the tester
	proc, _ := goka.NewProcessor([]string{}, aboveThresholdProcessor.Group,
		goka.WithTester(gkt),
	)

	go proc.Run(context.Background())

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

	gkt.SetTableValue(config.AboveThresholdTable, aboveThreshold.WalletId, &aboveThreshold)

	gkt.Consume(string(config.TopicDeposit), deposit.WalletId, &deposit)

	value := gkt.TableValue(config.AboveThresholdTable, aboveThreshold.WalletId)
	received := value.(*model.AboveThreshold)

	// check wallet id
	require.Equal(t, aboveThreshold.WalletId, received.WalletId)

	// check first deposit amount and status
	aboveThreshold.AmountWithinTwoMins += 7000
	expectedStatus = false
	expectedAmount = aboveThreshold.AmountWithinTwoMins
	require.Equal(t, expectedAmount, received.AmountWithinTwoMins)
	require.Equal(t, expectedStatus, received.Status)

	// deposit another to make amount within two mins above 10.000
	deposit.DepositedAt = timestamppb.New(time.Now().Add(5 * time.Second))
	gkt.Consume(string(config.TopicDeposit), deposit.WalletId, &deposit)

	value = gkt.TableValue(config.AboveThresholdTable, aboveThreshold.WalletId)
	received = value.(*model.AboveThreshold)

	expectedAmount += 7000
	expectedStatus = true

	require.Equal(t, expectedAmount, received.AmountWithinTwoMins)
	require.Equal(t, expectedStatus, received.Status)

}
