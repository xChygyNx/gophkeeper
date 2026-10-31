package service

import (
	"fmt"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"

	"github.com/golang/protobuf/ptypes/timestamp"
)

func ConvertTimeToTimestamp(t time.Time) *timestamp.Timestamp {
	return timestamppb.New(t)
}

func ConvertTimestampToTime(ts *timestamp.Timestamp) (time.Time, error) {
	err := ts.CheckValid()
	if err != nil {
		return time.Time{}, fmt.Errorf("timestamp is invalid: %w", err)
	}
	return ts.AsTime(), nil
}
