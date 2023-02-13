package config

import (
	"coinbit-wallet/util/logger"
	"context"

	"github.com/Shopify/sarama"
	"github.com/lovoo/goka"
	"golang.org/x/sync/errgroup"
)

var (
	Brokers             []string
	TopicDeposit        goka.Stream = "deposits"
	TMC                 *goka.TopicManagerConfig
	BalanceGroup        goka.Group = "balance"
	BalanceTable        goka.Table = goka.GroupTable(BalanceGroup)
	AboveThresholdGroup goka.Group = "aboveThreshold"
	AboveThresholdTable goka.Table = goka.GroupTable(AboveThresholdGroup)
)

func InitGoka() {
	logger.Info("Init Goka configuration")

	Brokers = GetEnv().KafkaBrokers

	ctx, cancel := context.WithCancel(context.Background())
	grp, ctx := errgroup.WithContext(ctx)

	defer cancel()

	grp.Go(EnsureStreamExists(string(TopicDeposit), Brokers))
	grp.Go(EnsureStreamExists(string(BalanceTable), Brokers))
	grp.Go(EnsureStreamExists(string(AboveThresholdTable), Brokers))

	if err := grp.Wait(); err != nil {
		panic(err)
	}

}

func createTopicManager(brokers []string) goka.TopicManager {
	TMC = goka.NewTopicManagerConfig()
	TMC.Table.Replication = 1
	TMC.Stream.Replication = 1

	config := goka.DefaultConfig()

	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	goka.ReplaceGlobalConfig(config)

	tm, err := goka.NewTopicManager(Brokers, goka.DefaultConfig(), TMC)
	if err != nil {
		logger.Error("Error creating topic manager: %v", err)
		panic(err)
	}
	return tm
}

func EnsureStreamExists(topic string, brokers []string) func() error {
	return func() error {
		logger.Info("Ensuring topic %s exists", topic)
		tm := createTopicManager(brokers)
		defer tm.Close()
		err := tm.EnsureStreamExists(topic, 8)
		if err != nil {
			logger.Error("Error creating kafka topic %s: %v", TopicDeposit, err)
		}
		return err
	}
}
