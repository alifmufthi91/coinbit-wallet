package time_util

import (
	"time"
)

type ITimestampGenerator interface {
	Generate() time.Time
}

type TimestampGenerator struct {
}

func NewTimeStampGenerator() TimestampGenerator {
	return TimestampGenerator{}
}

func (tg TimestampGenerator) Generate() time.Time {
	return time.Now()
}
