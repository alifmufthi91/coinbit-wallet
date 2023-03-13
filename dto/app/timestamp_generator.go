package app

import (
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ITimestampGenerator interface {
	Generate() *timestamppb.Timestamp
}

type TimestampGenerator struct {
}

func NewTimeStampGenerator() TimestampGenerator {
	return TimestampGenerator{}
}

func (tg TimestampGenerator) Generate() *timestamppb.Timestamp {
	return timestamppb.Now()
}
