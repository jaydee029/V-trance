package service

import (
	"V-trance/worker/Internal/database"
	"V-trance/worker/Internal/utils"
	"context"
	"os"

	"github.com/minio/minio-go/v7"
	"go.uber.org/zap"
)

func (h *Handler) ProcessVideo(job *Job) error {
	// Avgbitrate := map[int]int{
	// 	360:  506,
	// 	480:  960,
	// 	720:  1726,
	// 	1080: 3399,
	// }

	//InputDir := "tmp/"
	OutputDir := "tmp/output/"
	filepath := "tmp/" + job.VideoId.String()

	// err := os.MkdirAll(InputDir, 0755)
	// if err != nil {
	// 	h.logger.Info("Error creating input directory", zap.Error(err))
	// 	return
	// }
	err := os.MkdirAll(OutputDir, 0755)
	if err != nil {
		h.logger.Info("Error creating output directory", zap.Error(err))
		return err
	}
	err = h.B2Client.FGetObject(context.Background(), h.Bucket, h.Pathprefix+job.UserId.String()+"/"+job.VideoId.String(), filepath, minio.GetObjectOptions{})
	if err != nil {
		h.logger.Info("Error downloading video from b2", zap.Error(err))
		return err
	}

	if job.Type == "STREAMING" {
		err = CreateHls(job, filepath, OutputDir)
		if err != nil {
			h.logger.Info("Error Creating HLS playlist", zap.Error(err))
			return err
		}
		_, err = h.B2Client.FPutObject(context.Background(), h.Bucket, h.Pathprefix+job.UserId.String()+"/hls/"+job.VideoId.String(), OutputDir+job.VideoId.String(), minio.PutObjectOptions{})
		if err != nil {
			h.logger.Info("Error uploading HLS playlist", zap.Error(err))
			return err
		}

		videourlstr := "https://snowy-unit-44da.dhruvj797.workers.dev/" + h.Pathprefix + job.UserId.String() + "/hls/" + job.VideoId.String() + "/master.m3u8"
		videourl := utils.ToText(videourlstr)
		_, err = h.DB.InsertVideoUrl(context.Background(), database.InsertVideoUrlParams{
			StreamUrl: videourl,
		})
		if err != nil {
			h.logger.Info("Error uploading HLS URL to the db", zap.Error(err))
			return err
		}
	} else if job.Type == "TRANSCODING" {
		err = CreateTranscoding(job, filepath, OutputDir)
		if err != nil {
			h.logger.Info("Error creating transcoded video", zap.Error(err))
			return err
		}
		_, err = h.B2Client.FPutObject(context.Background(), h.Bucket, h.Pathprefix+job.UserId.String()+"/transcode/"+job.VideoId.String(), OutputDir+job.VideoId.String()+job.Options.Output, minio.PutObjectOptions{})
		if err != nil {
			h.logger.Info("Error uploading HLS playlist", zap.Error(err))
			return err
		}

		videourlstr := "https://snowy-unit-44da.dhruvj797.workers.dev/" + h.Pathprefix + job.UserId.String() + "/transcode" + job.VideoId.String()
		videourl := utils.ToText(videourlstr)
		_, err = h.DB.InsertVideoUrl(context.Background(), database.InsertVideoUrlParams{
			StreamUrl: videourl,
		})
		if err != nil {
			h.logger.Info("Error uploading URL to the db", zap.Error(err))
			return err
		}
	}

	//fmt.Println(Avgbitrate)
	return nil
}
