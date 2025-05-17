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
	Type              string
	VideoId           pgtype.UUID
	UserId            pgtype.UUID
	InitialResolution int
	VideoUrl          string
	Options           Options
}

// type Transcoding struct {
// }

type Dimension struct {
	Height int
	Width  int
}

type Rendition struct {
	Name    string
	Width   int
	Height  int
	Bitrate int // in kbps
}
