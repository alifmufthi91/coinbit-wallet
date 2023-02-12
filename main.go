package main

import (
	"coinbit-wallet/config"
	"coinbit-wallet/emitter"
	"coinbit-wallet/processor"
	"coinbit-wallet/server"
	"coinbit-wallet/util/logger"
	"coinbit-wallet/view"
	"fmt"
)

func main() {
	logger.Init()
	config.EnvInit()
	config.InitGoka()

	env := config.GetEnv()

	go emitter.RunDepositEmitter()
	go processor.RunBalanceProcessor()
	go processor.RunAboveThresholdProcessor()
	go view.RunAboveThresholdView()
	go view.RunBalanceView()

	router := server.NewRouter()

	logger.Info(fmt.Sprintf("Running Server on Port: %s", env.Port))
	err := router.Run(fmt.Sprintf("localhost:%s", env.Port))
	if err != nil {
		panic(err)
	}
}
