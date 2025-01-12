package server

import (
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgtype" // Import sqlc-generated code
	_type "github.com/thylong/go-templates/06-grpc-sqlc/pkg/proto/google/type"
)

// Helpers for DateTime Conversion

func convertToProtoDateTime(t pgtype.Timestamp) (*_type.DateTime, error) {
	if !t.Valid {
		return nil, fmt.Errorf("timestamp is not valid")
	}

	timeValue := t.Time
	return &_type.DateTime{
		Year:    int32(timeValue.Year()),
		Month:   int32(timeValue.Month()),
		Day:     int32(timeValue.Day()),
		Hours:   int32(timeValue.Hour()),
		Minutes: int32(timeValue.Minute()),
		Seconds: int32(timeValue.Second()),
	}, nil
}

func convertFromProtoDateTime(dt *_type.DateTime) time.Time {
	return time.Date(
		int(dt.Year), time.Month(dt.Month), int(dt.Day),
		int(dt.Hours), int(dt.Minutes), int(dt.Seconds), 0, time.UTC,
	)
}
