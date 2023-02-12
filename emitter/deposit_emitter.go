package emitter

import (
	"coinbit-wallet/config"
	"coinbit-wallet/generated/model"
	"coinbit-wallet/util"
	"coinbit-wallet/util/logger"

	"github.com/lovoo/goka"
)

var (
	depositEmitter *goka.Emitter
)

func RunDepositEmitter() {
	logger.Info("Running deposit emitter..")
	var err error
	depositEmitter, err = goka.NewEmitter(config.Brokers, config.TopicDeposit, new(util.DepositCodec))
	if err != nil {
		logger.Error("error creating emitter: %v", err)
		panic(err)
	}
}

func EmitDeposit(deposit *model.Deposit) {
	logger.Info("emitting deposit request")
	err := depositEmitter.EmitSync(string(config.TopicDeposit), deposit)
	if err != nil {
		logger.Error("error emitting message: %v", err)
	}
	logger.Info("deposit request emitted")
}
