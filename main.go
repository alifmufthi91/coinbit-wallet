package main

import (
	"coinbit-wallet/config"
	"coinbit-wallet/emitter"
	"coinbit-wallet/processor"
	"coinbit-wallet/server"
	"coinbit-wallet/util/logger"
	"coinbit-wallet/view"
	"context"
	"fmt"

	"golang.org/x/sync/errgroup"
)

func main() {
	logger.Init()
	config.InitEnv()
	config.InitGoka()

	env := config.GetEnv()

	balanceView := view.CreateBalanceView(config.Brokers)
	aboveThresholdView := view.CreateAboveThresholdView(config.Brokers)

	ctx, cancel := context.WithCancel(context.Background())
	grp, ctx := errgroup.WithContext(ctx)

	defer cancel()

	emitter.RunDepositEmitter(config.Brokers, config.TopicDeposit)

	grp.Go(processor.RunBalanceProcessor(ctx, config.Brokers))
	grp.Go(processor.RunAboveThresholdProcessor(ctx, config.Brokers))
	grp.Go(view.RunBalanceView(balanceView, ctx))
	grp.Go(view.RunAboveThresholdView(aboveThresholdView, ctx))

	router := server.NewRouter(balanceView, aboveThresholdView)

	logger.Info(fmt.Sprintf("Running Server on Port: %s", env.Port))
	err := router.Run(fmt.Sprintf("localhost:%s", env.Port))
	if err != nil {
		panic(err)
	}
}
