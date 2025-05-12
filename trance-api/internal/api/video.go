package api

import (
	"V-trance/trance-api/internal/database"
	"context"
	"net/http"
	"net/url"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"go.uber.org/zap"
)

func (h *Handler) GetVideos(w http.ResponseWriter, r *http.Request) {

	useridstr := r.Header.Get("X-User-ID")

	var userid pgtype.UUID
	err := userid.Scan(useridstr)
	if err != nil {
		h.logger.Info("Error setting user UUID:", zap.Error(err))
		return
	}

	videos, err := h.DB.GetVideos(r.Context(), userid)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error fetching videos from the db: "+err.Error())
		return
	}

	var videosRes []GetVideosResponse
	for _, video := range videos {
		videosRes = append(videosRes, GetVideosResponse{
			Name: video.Name,
			Url:  video.StreamUrl.String,
		})
	}

	respondWithJson(w, http.StatusAccepted, videosRes)
}

func (h *Handler) GetStatus(w http.ResponseWriter, r *http.Request) {

	//useridstr := r.Context().Value(middleware.UserIDKey).(string)
	jobidstr := chi.URLParam(r, "jobid")

	var jobid pgtype.UUID
	err := jobid.Scan(jobidstr)
	if err != nil {
		h.logger.Info("Error setting job UUID:", zap.Error(err))
		return
	}

	timeout := time.After(time.Second * 40)
	ticker := time.NewTicker(time.Second * 8)
	for {
		select {
		case <-r.Context().Done():
			return
		case <-timeout:
			respondWithJson(w, http.StatusNoContent, GetStatusResponse{
				Status: JobKeyProcessing,
			})
			return
		case <-ticker.C:
			job, err := h.DB.FetchStatus(r.Context(), jobid)
			if err != nil {
				respondWithError(w, http.StatusInternalServerError, "error fetching the job status: "+err.Error())
				return
			}
			if job.Status == JobKeyRejected {
				respondWithJson(w, http.StatusBadRequest, GetStatusResponse{
					Name:    job.Name,
					Videoid: job.VideoID,
					Status:  JobKeyRejected,
				})
				return
			}
			if job.Status == JobKeyCompleted {
				respondWithJson(w, http.StatusOK, GetStatusResponse{
					Name:    job.Name,
					Videoid: job.VideoID,
					Status:  JobKeyCompleted,
				})
				return
			}

		}
	}
}

func (h *Handler) GetStreamUrl(w http.ResponseWriter, r *http.Request) {

	useridstr := r.Header.Get("X-User-ID")
	videoidstr := chi.URLParam(r, "videoid")

	var videoid pgtype.UUID
	err := videoid.Scan(videoidstr)
	if err != nil {
		h.logger.Info("Error setting Video UUID:", zap.Error(err))
		return
	}

	var userid pgtype.UUID
	err = userid.Scan(useridstr)
	if err != nil {
		h.logger.Info("Error setting user UUID:", zap.Error(err))
		return
	}

	video, err := h.DB.GetStreamurl(r.Context(), database.GetStreamurlParams{
		UserID:  userid,
		VideoID: videoid,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error fetching video url from the db: "+err.Error())
		return
	}
	respondWithJson(w, http.StatusAccepted, GetStreamUrlResponse{
		Name:      video.Name,
		StreamUrl: video.StreamUrl.String,
	})

}

func (h *Handler) GetDownloadUrl(w http.ResponseWriter, r *http.Request) {
	useridstr := r.Header.Get("X-User-ID")
	videoidstr := chi.URLParam(r, "videoid")

	var videoid pgtype.UUID
	err := videoid.Scan(videoidstr)
	if err != nil {
		h.logger.Info("Error setting video UUID:", zap.Error(err))
		return
	}

	var userid pgtype.UUID
	err = userid.Scan(useridstr)
	if err != nil {
		h.logger.Info("Error setting user UUID:", zap.Error(err))
		return
	}

	reqParams := make(url.Values)
	expiry := time.Second * 60 * 15
	objectName := h.Pathprefix + useridstr + "/transcode/" + videoidstr
	presignedURL, err := h.B2Client.PresignedGetObject(context.Background(), h.Bucket, objectName, expiry, reqParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error generating donwload url: "+err.Error())
		return
	}

	respondWithJson(w, http.StatusAccepted, GetDownloadUrlResponse{
		DownloadUrl: presignedURL.String(),
	})
}
