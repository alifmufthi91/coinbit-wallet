package emitter

import (
	"coinbit-wallet/generated/model"
	"coinbit-wallet/util"
	"coinbit-wallet/util/logger"

	"github.com/lovoo/goka"
)

var (
	depositEmitter *goka.Emitter
)

func InitDepositEmitter(brokers []string, stream goka.Stream) {
	logger.Info("Init deposit emitter..")
	var err error
	depositEmitter, err = goka.NewEmitter(brokers, stream, new(util.DepositCodec))
	if err != nil {
		logger.Error("error creating emitter: %v", err)
		panic(err)
	}
}

func EmitDeposit(deposit *model.Deposit) error {
	logger.Info("emitting deposit request")
	err := depositEmitter.EmitSync(deposit.WalletId, deposit)
	if err != nil {
		logger.Error("error emitting message: %v", err)
		return err
	}
	logger.Info("deposit request emitted")
	return nil
}
