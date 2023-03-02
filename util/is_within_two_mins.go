package util

import "google.golang.org/protobuf/types/known/timestamppb"

func IsWithinTwoMins(start *timestamppb.Timestamp, toBeCheck *timestamppb.Timestamp) bool {
	if toBeCheck.Seconds < start.Seconds+120 {
		return true
	}
	return false
}
