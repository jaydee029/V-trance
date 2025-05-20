package service

import (
	"V-trance/worker/Internal/database"
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jaydee029/V-trance/pubsub"
	"go.uber.org/zap"
	"golang.org/x/sync/semaphore"
)

func (h *Handler) EventListner(ctx context.Context, wg *sync.WaitGroup) {

	sem := semaphore.NewWeighted(10)
	//var wg sync.WaitGroup

	err := h.Pb.SubscribeGob(h.Exchange, "jobs_queue", h.key, pubsub.DurableQueue, func(val pubsub.Task) pubsub.Acktype {

		if err := sem.Acquire(ctx, 1); err != nil {
			h.logger.Info("Failed to acquire semaphore: ", zap.Error(err))
			return pubsub.NackRequeue
		}

		wg.Add(1)
		go func(task pubsub.Task) {
			defer sem.Release(1)
			defer wg.Done()

			if err := h.Jobprocesser(task); err != nil {
				h.logger.Info("Failed to process task: ", zap.Error(err))
				// You could implement retry logic or dead-lettering here
			}
		}(val)

		return pubsub.Ack
	})

	if err != nil {
		h.logger.Info("Couldnt subscribe to the queue:", zap.Error(err))
	}
	//wg.Wait()
}

func (h *Handler) Jobprocesser(val pubsub.Task) error {
	// m, ok := val.(map[string]interface{})
	// if !ok {
	// 	h.logger.Info("failed to convert val to map")
	// 	return errors.New("failed to convert val to map")
	// }

	//Manually map to your struct
	// task := &Task{
	// 	VideoID: val.Videoid, //m["Videoid"].(string),
	// 	JobID:   m["Jobid"].(string),
	// }
	// task, ok := val.(*Task)
	// if !ok {
	// 	h.logger.Info("failed to convert val to *Task")
	// 	return errors.New("failed to convert val to *Task")
	// }

	// task := Task{
	// 	VideoID: "string",
	// 	JobID:   "Jobid",
	// }
	var jobid pgtype.UUID

	err := jobid.Scan(val.Jobid)

	fmt.Println(jobid)

	if err != nil {
		h.logger.Info("error converting jobid to pgtype", zap.Error(err))
		return err
	}

	jobDetails, err := h.DB.FetchJob(context.Background(), database.FetchJobParams{
		JobID:  jobid,
		Status: JobKeyInitiated,
	})
	if err != nil {
		h.logger.Info("error fetching the job from the db:", zap.Error(err))
	}
	var options Options
	err = json.Unmarshal(jobDetails.Options, &options)
	if err != nil {
		h.logger.Info("error unmarshalling options:", zap.Error(err))
	}
	videoDetails, err := h.DB.FetchVideo(context.Background(), jobDetails.VideoID)
	if err != nil {
		h.logger.Info("error fetching the video from the db:", zap.Error(err))
	}

	job := &Job{
		Type:              jobDetails.Type,
		VideoId:           jobDetails.VideoID,
		UserId:            videoDetails.UserID,
		VideoUrl:          videoDetails.VideoUrl.String,
		InitialResolution: int(videoDetails.Resolution),
		Options:           options,
	}

	_, err = h.DB.SetStatusJob(context.Background(), database.SetStatusJobParams{
		Status: JobKeyProcessing,
		JobID:  jobid,
	})
	if err != nil {
		h.logger.Info("error setting status to processing:", zap.Error(err))
	}

	err = h.ProcessVideo(job)

	if err != nil {
		_, _ = h.DB.SetStatusJob(context.Background(), database.SetStatusJobParams{
			Status: JobKeyRejected,
			JobID:  jobid,
		})
		return err
	}
	_, _ = h.DB.SetStatusJob(context.Background(), database.SetStatusJobParams{
		Status: JobKeyCompleted,
		JobID:  jobid,
	})

	return nil
}
