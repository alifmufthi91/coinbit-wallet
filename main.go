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

	var err error
	var balanceView *view.BalanceView
	var aboveThresholdView *view.AboveThresholdView

	if balanceView, err = view.NewBalanceView(config.Brokers); err != nil {
		panic(err)
	}

	if aboveThresholdView, err = view.NewAboveThresholdView(config.Brokers); err != nil {
		panic(err)
	}

	go RunGokaProcessors(ctx, &wg)
	go RunGokaViewers(balanceView, aboveThresholdView)
	RunServer(balanceView, aboveThresholdView, &wg)
}

func Init() {
	config.InitGoka()
	emitter.InitDepositEmitter(config.Brokers, config.TopicDeposit)
}

func RunServer(balanceView *view.BalanceView, aboveThresholdView *view.AboveThresholdView, wg *sync.WaitGroup) {
	env := config.GetEnv()
	router := server.NewRouter(balanceView, aboveThresholdView)
	srv := http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%s", env.Port),
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

	logger.Info("Running Server on Port: %s", env.Port)
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Printf("listen: %s\n", err)
	}
}

func RunGokaProcessors(ctx context.Context, wg *sync.WaitGroup) {
	go func() {
		err := processor.RunBalanceProcessor(ctx, config.Brokers)
		if err != nil {
			logger.Error("Error running balance processor: %s", err.Error())
		}
		wg.Done()
	}()
	go func() {
		err := processor.RunAboveThresholdProcessor(ctx, config.Brokers)
		if err != nil {
			logger.Error("Error running above threshold processor: %s", err.Error())
		}
		wg.Done()
	}()
}

func RunGokaViewers(balanceView *view.BalanceView, aboveThresholdView *view.AboveThresholdView) {
	go func() {
		for {
			err := balanceView.Run(context.Background())
			if err != nil {
				logger.Error("Error running balance view : %s", err.Error())
			}
		}
	}()
	go func() {
		for {
			err := aboveThresholdView.Run(context.Background())
			if err != nil {
				logger.Error("Error running above threshold view : %s", err.Error())
			}
		}
	}()
}
