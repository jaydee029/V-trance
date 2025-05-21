package api

import (
	"V-trance/trance-api/internal/database"
	"V-trance/trance-api/internal/utils"
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/jaydee029/V-trance/pubsub"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/minio/minio-go/v7"
	"go.uber.org/zap"
)

// func (h *Handler) GetUploadUrl(w http.ResponseWriter, r *http.Request) {

// 	useridstr := r.Header.Get("X-User-ID")
// 	decoder := json.NewDecoder(r.Body)
// 	params := UploadUrlInput{}
// 	err := decoder.Decode(&params)
// 	if err != nil {
// 		respondWithError(w, http.StatusInternalServerError, "couldn't decode parameters")
// 		return
// 	}
// 	videoidstr := uuid.New().String()

// 	filename := h.Pathprefix + useridstr + "/" + videoidstr
// 	expiry := time.Second * 60 * 15 // 15 minutes.
// 	presignedURL, err := h.B2Client.PresignedPutObject(context.Background(), h.Bucket, filename, expiry)

// 	var videoid pgtype.UUID

// 	err = videoid.Scan(videoidstr)
// 	if err != nil {
// 		h.logger.Info("Error setting video UUID:", zap.Error(err))
// 		return
// 	}

// 	var timestamp pgtype.Timestamp
// 	err = timestamp.Scan(time.Now().UTC())
// 	if err != nil {
// 		h.logger.Info("Error setting video creation timestamp:", zap.Error(err))
// 		return
// 	}

// 	var userid pgtype.UUID
// 	err = userid.Scan(useridstr)
// 	if err != nil {
// 		h.logger.Info("Error setting user UUID:", zap.Error(err))
// 		return
// 	}

// 	Response, err := h.DB.InsertInitialDetails(r.Context(), database.InsertInitialDetailsParams{
// 		UserID:     userid,
// 		Name:       params.Name,
// 		Type:       params.Type,
// 		Resolution: int32(params.Resolution),
// 		VideoID:    videoid,
// 		CreatedAt:  timestamp,
// 	})
// 	if err != nil {
// 		respondWithError(w, http.StatusInternalServerError, "error inserting details in the db: "+err.Error())
// 		return
// 	}
// 	url := fmt.Sprintf("%s://%s%s?%s", presignedURL.Scheme, presignedURL.Host, presignedURL.Path, presignedURL.RawQuery)

// 	respondWithJson(w, http.StatusAccepted, UploadUrlResponse{
// 		Name:      Response.Name,
// 		Videoid:   Response.VideoID,
// 		UploadUrl: url,
// 	})

// }

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
	// var videoid pgtype.UUID

	// err = videoid.Scan(params.Videoid)
	// if err != nil {
	// 	h.logger.Info("Error setting video UUID:", zap.Error(err))
	// 	return
	// }

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
		VideoID:  params.Videoid,
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

	task := pubsub.Task{
		Videoid: job.VideoID.String(),
		Jobid:   job.JobID.String(),
	}

	err = h.Pbclient.PublishTask(h.Exchange, h.Key, task, h.logger)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "coudn't publish the event:"+err.Error())
		return
	}

	// go func(task pubsub.Task) {
	// 	jobDetails, err := h.DB.FetchJob(context.Background(), database.FetchJobParams{
	// 		JobID:  jobid,
	// 		Status: JobKeyInitiated,
	// 	})
	// 	if err != nil {
	// 		h.logger.Info("error fetching the job from the db:", zap.Error(err))
	// 	}
	// 	var opts Options
	// 	err = json.Unmarshal(jobDetails.Options, &opts)
	// 	if err != nil {
	// 		h.logger.Info("error unmarshalling options:", zap.Error(err))
	// 	}
	// 	videoDetails, err := h.DB.FetchVideo(context.Background(), jobDetails.VideoID)
	// 	if err != nil {
	// 		h.logger.Info("error fetching the video from the db:", zap.Error(err))
	// 	}

	// 	jobd := &Job{
	// 		Type:              jobDetails.Type,
	// 		VideoId:           jobDetails.VideoID,
	// 		UserId:            videoDetails.UserID,
	// 		VideoUrl:          videoDetails.VideoUrl.String,
	// 		InitialResolution: int(videoDetails.Resolution),
	// 		Options:           opts,
	// 	}

	// 	_, err = h.DB.SetStatusJob(context.Background(), database.SetStatusJobParams{
	// 		Status: JobKeyProcessing,
	// 		JobID:  jobid,
	// 	})
	// 	if err != nil {
	// 		h.logger.Info("error setting status to processing:", zap.Error(err))
	// 	}

	// 	err = h.ProcessVideo(jobd)

	// 	if err != nil {
	// 		h.DB.SetStatusJob(context.Background(), database.SetStatusJobParams{
	// 			Status: JobKeyRejected,
	// 			JobID:  jobid,
	// 		})
	// 		log.Printf("error processing video:%v", err)
	// 		return
	// 		//return err
	// 	}
	// 	h.DB.SetStatusJob(context.Background(), database.SetStatusJobParams{
	// 		Status: JobKeyCompleted,
	// 		JobID:  jobid,
	// 	})
	// 	return
	// }(task)

	//return nil

	respondWithJson(w, http.StatusAccepted, NotifyUploadResponse{
		Name:    job.Name,
		Videoid: job.VideoID,
		Jobid:   job.JobID,
	})
}

func (h *Handler) UploadVideo(w http.ResponseWriter, r *http.Request) {

	useridstr := r.Header.Get("X-User-ID")

	err := r.ParseMultipartForm(100 << 20)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error setting memory limits: "+err.Error())
		return
	}

	Name := r.FormValue("name")
	Type := r.FormValue("type")
	Resolutionstr := r.FormValue("resolution")
	Resolution, err := strconv.Atoi(Resolutionstr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "error setting the resolution: "+err.Error())
		return
	}

	videoidstr := uuid.New().String()

	filename := h.Pathprefix + useridstr + "/" + videoidstr

	file, handler, err := r.FormFile("video")
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Failed to read video: "+err.Error())
		return
	}
	defer file.Close()

	_, err = h.B2Client.PutObject(
		context.Background(),
		h.Bucket,
		filename,
		file,
		handler.Size,
		minio.PutObjectOptions{
			ContentType: Type,
			UserMetadata: map[string]string{
				"title":       Name,
				"description": "Unprocessed Video",
			},
		},
	)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to upload video: "+err.Error())
		return
	}

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
		Name:       Name,
		Type:       Type,
		Resolution: int32(Resolution),
		VideoID:    videoid,
		CreatedAt:  timestamp,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error inserting details in the db: "+err.Error())
		return
	}

	respondWithJson(w, http.StatusAccepted, UploadVideoResponse{
		Name:      Response.Name,
		Videoid:   Response.VideoID,
		CreatedAt: timestamp,
	})
}
