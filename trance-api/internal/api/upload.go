package api

import (
	"V-trance/trance-api/internal/database"
	"V-trance/trance-api/internal/utils"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"go.uber.org/zap"
)

func (h *Handler) GetUploadUrl(w http.ResponseWriter, r *http.Request) {

	useridstr := r.Header.Get("X-User-ID")
	decoder := json.NewDecoder(r.Body)
	params := UploadUrlInput{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't decode parameters")
		return
	}
	videoidstr := uuid.New().String()

	filename := h.Pathprefix + useridstr + "/" + videoidstr
	expiry := time.Second * 60 * 15 // 15 minutes.
	presignedURL, err := h.B2Client.PresignedPutObject(context.Background(), h.Bucket, filename, expiry)

	var videoid pgtype.UUID

	err = videoid.Scan(videoidstr)
	if err != nil {
		h.logger.Info("Error setting video UUID:", zap.Error(err))
		return
	}

	var timestamp pgtype.Timestamp
	err = timestamp.Scan(time.Now().UTC())
	if err != nil {
		h.logger.Info("Error setting video creation timestamp:", zap.Error(err))
		return
	}

	var userid pgtype.UUID
	err = userid.Scan(useridstr)
	if err != nil {
		h.logger.Info("Error setting user UUID:", zap.Error(err))
		return
	}

	Response, err := h.DB.InsertInitialDetails(r.Context(), database.InsertInitialDetailsParams{
		UserID:     userid,
		Name:       params.Name,
		Type:       params.Type,
		Resolution: int32(params.Resolution),
		VideoID:    videoid,
		CreatedAt:  timestamp,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error inserting details in the db: "+err.Error())
		return
	}
	url := fmt.Sprintf("%s://%s%s?%s", presignedURL.Scheme, presignedURL.Host, presignedURL.Path, presignedURL.RawQuery)

	respondWithJson(w, http.StatusAccepted, UploadUrlResponse{
		Name:      Response.Name,
		Videoid:   Response.VideoID,
		UploadUrl: url,
	})

}

func (h *Handler) NotifyUpload(w http.ResponseWriter, r *http.Request) {

	useridstr := r.Header.Get("X-User-ID")
	decoder := json.NewDecoder(r.Body)
	params := NotifyUploadInput{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't decode parameters")
		return
	}

	if_video_exists, err := h.DB.IfVideoExists(r.Context(), params.Videoid)

	if !if_video_exists {
		respondWithError(w, http.StatusNotFound, "couldnt find video details")
		return
	}

	tx, err := h.DBPool.Begin(r.Context())
	if err != nil {
		h.logger.Info("Error starting the transaction:", zap.Error(err))
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	qtx := h.DB.WithTx(tx)

	defer func() {
		if tx != nil {
			tx.Rollback(r.Context())
		}
	}()

	videoprefix := h.Pathprefix + useridstr + "/" + params.Videoid.String()

	videoprefixpgtype := utils.ToText(videoprefix)

	video_details, err := qtx.InsertFinalVideoDetails(r.Context(), database.InsertFinalVideoDetailsParams{
		VideoUrl: videoprefixpgtype,
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error adding video url: "+err.Error())
		return
	}

	jobidstr := uuid.New().String()
	var jobid pgtype.UUID

	err = jobid.Scan(jobidstr)
	if err != nil {
		h.logger.Info("Error setting job UUID:", zap.Error(err))
		return
	}

	options, err := json.Marshal(params.Options)

	if err != nil {
		h.logger.Info("error marshal options:", zap.Error(err))
		return
	}

	var timestamp pgtype.Timestamp
	err = timestamp.Scan(time.Now().UTC())
	if err != nil {
		h.logger.Info("Error setting video creation timestamp:", zap.Error(err))
		return
	}

	job, err := qtx.CreateJob(r.Context(), database.CreateJobParams{
		JobID:     jobid,
		VideoID:   video_details.VideoID,
		Name:      video_details.Name,
		Type:      params.Type,
		Options:   options,
		Status:    JobKeyInitiated,
		CreatedAt: timestamp,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error creating the job: "+err.Error())
		return
	}

	if err = tx.Commit(r.Context()); err != nil {
		respondWithError(w, http.StatusInternalServerError, "error commmiting the transaction:"+err.Error())
		return
	}
	tx = nil

	task := &Task{
		Videoid: job.VideoID.String(),
		Jobid:   job.JobID.String(),
	}

	err = h.Pbclient.PublishTask(h.Exchange, h.Key, task, h.logger)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "coudn't publish the event:"+err.Error())
		return
	}

	respondWithJson(w, http.StatusAccepted, NotifyUploadResponse{
		Name:    job.Name,
		Videoid: job.VideoID,
		Jobid:   job.JobID,
	})
}
