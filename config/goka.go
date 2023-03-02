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
	logger.Info("brokers : %v", Brokers)

	ctx, cancel := context.WithCancel(context.Background())
	grp, ctx := errgroup.WithContext(ctx)

	defer cancel()

	grp.Go(EnsureStreamExists(string(TopicDeposit), Brokers))
	grp.Go(EnsureTableExists(string(BalanceTable), Brokers))
	grp.Go(EnsureTableExists(string(AboveThresholdTable), Brokers))

	if err := grp.Wait(); err != nil {
		panic(err)
	}

}

func createTopicManager(brokers []string) goka.TopicManager {
	TMC = goka.NewTopicManagerConfig()
	TMC.Table.Replication = 1
	TMC.Stream.Replication = 1

	config := goka.DefaultConfig()
	config.Consumer.IsolationLevel = sarama.ReadCommitted
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	goka.ReplaceGlobalConfig(config)

	tm, err := goka.NewTopicManager(brokers, goka.DefaultConfig(), TMC)
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
			logger.Error("Error creating kafka topic %s: %v", topic, err)
		}
		return err
	}
}

func EnsureTableExists(table string, brokers []string) func() error {
	return func() error {
		logger.Info("Ensuring table %s exists", table)
		tm := createTopicManager(brokers)
		defer tm.Close()
		err := tm.EnsureTableExists(table, 8)
		if err != nil {
			logger.Error("Error creating kafka table %s: %v", table, err)
		}
		return err
	}
}
