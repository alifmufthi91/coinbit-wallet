package config

import (
	"coinbit-wallet/util/logger"
	"os"
	"strings"
	"sync"

	"github.com/joho/godotenv"
)

var (
	env  Environment
	once sync.Once
)

type Environment struct {
	Port         string
	KafkaBrokers []string
}

func InitEnv() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}
	env = Environment{
		Port:         os.Getenv("PORT"),
		KafkaBrokers: strings.Split(os.Getenv("KAFKA_BROKERS"), ","),
	}

	logger.Info("Environment config set")
}

func GetEnv() *Environment {
	once.Do(func() {
		InitEnv()
	})
	return &env
}
