package api

import (
	"net/url"

	"github.com/jackc/pgx/v5/pgtype"
)

const (
	JobKeyInitiated  string = "INITIATED"
	JobKeyProcessing string = "PROCESSING"
	JobKeyRejected   string = "REJECTED"
	JobKeyCompleted  string = "COMPLETED"
)

type options struct {
	Output     *string `json:"output"`
	Codec      *string `json:"codec"`
	Resolution *string `json:"resolution"`
}

type NotifyUploadInput struct {
	Videoid pgtype.UUID `json:"videoid"`
	Type    string      `json:"type"`
	Options *options    `json:"options"`
}

type NotifyUploadResponse struct {
	Name    string      `json:"name"`
	Videoid pgtype.UUID `json:"videoid"`
	Jobid   pgtype.UUID `json:"jobid"`
}

type UploadUrlInput struct {
	Name       string `json:"Name"`
	Type       string `json:"type"`
	Resolution int    `json:"resolution"`
}

type UploadUrlResponse struct {
	Name      string      `json:"name"`
	Videoid   pgtype.UUID `json:"videoid"`
	UploadUrl *url.URL    `json:"uploadurl"`
}

type GetVideosResponse struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type GetStreamUrlResponse struct {
	Name      string `json:"name"`
	StreamUrl string `json:"streamUrl"`
}

type GetDownloadUrlResponse struct {
	DownloadUrl string `json:"downloadUrl"`
}

type GetStatusResponse struct {
	Name    string      `json:"name"`
	Videoid pgtype.UUID `json:"videoid"`
	Status  string      `json:"status"`
}
