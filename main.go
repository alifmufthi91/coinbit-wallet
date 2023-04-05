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

var (
	balanceProcessor        *processor.BalanceProcessor
	aboveThresholdProcessor *processor.AboveThresholdProcessor
	balanceView             *view.BalanceView
	aboveThresholdView      *view.AboveThresholdView
	depositEmitter          *emitter.DepositEmitter
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

	go RunGokaProcessors(ctx, &wg)
	go RunGokaViewers(balanceView, aboveThresholdView)
	RunServer(balanceView, aboveThresholdView, depositEmitter, &wg)
}

func Init() {
	config.InitEnv()
	err := config.InitGoka()
	if err != nil {
		panic(err)
	}

	if balanceView, err = view.NewBalanceView(config.Brokers); err != nil {
		panic(err)
	}

	if aboveThresholdView, err = view.NewAboveThresholdView(config.Brokers); err != nil {
		panic(err)
	}

	if depositEmitter, err = emitter.NewDepositEmitter(config.Brokers); err != nil {
		panic(err)
	}
}

func RunServer(bv *view.BalanceView, atv *view.AboveThresholdView, de *emitter.DepositEmitter, wg *sync.WaitGroup) {
	env := config.GetEnv()
	router := server.NewRouter(bv, atv, de)
	srv := http.Server{
		Addr:    fmt.Sprintf("%s:%s", env.Host, env.Port),
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
		var err error
		balanceProcessor, err = processor.NewBalanceProcessor(config.Brokers,
			goka.WithTopicManagerBuilder(goka.TopicManagerBuilderWithTopicManagerConfig(config.TMC)),
			goka.WithConsumerGroupBuilder(goka.DefaultConsumerGroupBuilder))
		if err != nil {
			panic(err)
		}
		err = balanceProcessor.Run(ctx)
		if err != nil {
			logger.Error("Error running balance processor: %s", err.Error())
		}
		wg.Done()
	}()
	go func() {
		var err error
		aboveThresholdProcessor, err = processor.NewAboveThresholdProcessor(config.Brokers,
			goka.WithTopicManagerBuilder(goka.TopicManagerBuilderWithTopicManagerConfig(config.TMC)),
			goka.WithConsumerGroupBuilder(goka.DefaultConsumerGroupBuilder))
		if err != nil {
			panic(err)
		}
		err = aboveThresholdProcessor.Run(ctx)
		if err != nil {
			logger.Error("Error running above threshold processor: %s", err.Error())
		}
		wg.Done()
	}()
}

func RunGokaViewers(bv *view.BalanceView, atv *view.AboveThresholdView) {
	go func() {
		for {
			err := bv.Run(context.Background())
			if err != nil {
				logger.Error("Error running balance view : %s", err.Error())
			}
		}
	}()
	go func() {
		for {
			err := atv.Run(context.Background())
			if err != nil {
				logger.Error("Error running above threshold view : %s", err.Error())
			}
		}
	}()
}
