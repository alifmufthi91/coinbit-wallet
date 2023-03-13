package app

import (
	"google.golang.org/protobuf/types/known/timestamppb"
)

type MockTimestampGenerator struct {
	Timestamp *timestamppb.Timestamp
}

func (mtg MockTimestampGenerator) Generate() *timestamppb.Timestamp {
	return mtg.Timestamp
}
