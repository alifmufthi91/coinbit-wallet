package logger

import (
	"fmt"

	"github.com/rs/zerolog/log"
)

func Info(format string, v ...interface{}) {
	log.Info().Msg(fmt.Sprintf(format, v...))
}

func Error(format string, v ...interface{}) {
	log.Error().Caller(1).Msg(fmt.Sprintf(format, v...))
}

func Warn(format string, v ...interface{}) {
	log.Warn().Caller(1).Msg(fmt.Sprintf(format, v...))
}
