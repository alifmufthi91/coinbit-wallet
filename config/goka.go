package config

import (
	"coinbit-wallet/util/logger"

	"github.com/Shopify/sarama"
	"github.com/lovoo/goka"
)

var (
	Brokers      []string
	TopicDeposit goka.Stream = "deposits"
	TMC          *goka.TopicManagerConfig
)

func InitGoka() {
	logger.Info("Init Goka configuration")

	Brokers = GetEnv().KafkaBrokers

	EnsureStreamExists(string(TopicDeposit), Brokers)
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

func EnsureStreamExists(topic string, brokers []string) {
	tm := createTopicManager(brokers)
	defer tm.Close()
	err := tm.EnsureStreamExists(topic, 8)
	if err != nil {
		logger.Error("Error creating kafka topic %s: %v", TopicDeposit, err)
		panic(err)
	}
}
