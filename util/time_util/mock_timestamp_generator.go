package time_util

import (
	"time"
)

type MockTimestampGenerator struct {
	Timestamp time.Time
}

func (mtg MockTimestampGenerator) Generate() time.Time {
	return mtg.Timestamp
}
