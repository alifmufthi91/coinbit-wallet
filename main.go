package main

import (
	"coinbit-wallet/config"
	"coinbit-wallet/emitter"
	"coinbit-wallet/processor"
	"coinbit-wallet/server"
	"coinbit-wallet/util/logger"
	"coinbit-wallet/view"
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/lovoo/goka"
)

func main() {
	Init()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)
	go func() {
		<-sigint
		logger.Info("Received SIGINT. Cancelling context...")
		cancel()
	}()

	wg := sync.WaitGroup{}
	wg.Add(2)
	defer wg.Wait()

	balanceView := view.CreateBalanceView(config.Brokers)
	aboveThresholdView := view.CreateAboveThresholdView(config.Brokers)

	go RunGokaProcessors(ctx, &wg)
	go RunGokaViewers(balanceView, aboveThresholdView)
	RunServer(balanceView, aboveThresholdView, &wg)
}

func Init() {
	logger.Init()
	config.InitGoka()
	emitter.InitDepositEmitter(config.Brokers, config.TopicDeposit)
}

func RunServer(balanceView *goka.View, aboveThresholdView *goka.View, wg *sync.WaitGroup) {
	env := config.GetEnv()
	router := server.NewRouter(balanceView, aboveThresholdView)
	logger.Info("Running Server on Port: %s", env.Port)
	srv := http.Server{
		Addr:    fmt.Sprintf("localhost:%s", env.Port),
		Handler: router,
	}

	go func() {
		wg.Wait()
		logger.Info("Received cancellation signal. Shutting down server...")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			log.Fatal("Server forced to shutdown:", err)
		}
	}()

	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Printf("listen: %s\n", err)
	}
}

func RunGokaProcessors(ctx context.Context, wg *sync.WaitGroup) {
	go func() {
		processor.RunBalanceProcessor(ctx, config.Brokers)
		wg.Done()
	}()
	go func() {
		processor.RunAboveThresholdProcessor(ctx, config.Brokers)
		wg.Done()
	}()
}

func RunGokaViewers(balanceView *goka.View, aboveThresholdView *goka.View) {
	go view.RunBalanceView(balanceView, context.Background())
	go view.RunAboveThresholdView(aboveThresholdView, context.Background())
}
