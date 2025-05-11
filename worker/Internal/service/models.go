package service

import "github.com/jackc/pgx/v5/pgtype"

const (
	JobKeyInitiated  string = "INITIATED"
	JobKeyProcessing string = "PROCESSING"
	JobKeyRejected   string = "REJECTED"
	JobKeyCompleted  string = "COMPLETED"
)

type Task struct {
	VideoID string
	JobID   string
}

type Options struct {
	Output     string
	Codec      string
	Resolution string
}

type Job struct {
	Type     string
	VideoId  pgtype.UUID
	VideoUrl string
	Options  Options
}
