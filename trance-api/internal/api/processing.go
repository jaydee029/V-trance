package api

import (
	"V-trance/trance-api/internal/database"
	"V-trance/trance-api/internal/utils"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

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
		//_, err = h.B2Client.FPutObject(context.Background(), h.Bucket, h.Pathprefix+job.UserId.String()+"/hls/"+job.VideoId.String(), OutputDir+job.VideoId.String(), minio.PutObjectOptions{})
		err := uploadHLSDirectory(h.B2Client, h.Bucket,
			h.Pathprefix+job.UserId.String()+"/hls/"+job.VideoId.String(),
			OutputDir+job.VideoId.String(),
		)

		if err != nil {
			h.logger.Info("Error uploading HLS playlist", zap.Error(err))
			return err
		}

		videourlstr := "https://snowy-unit-44da.dhruvj797.workers.dev/" + h.Pathprefix + job.UserId.String() + "/hls/" + job.VideoId.String() + "/master.m3u8"
		videourl := utils.ToText(videourlstr)
		_, err = h.DB.InsertVideoUrl(context.Background(), database.InsertVideoUrlParams{
			StreamUrl: videourl,
			VideoID:   job.VideoId,
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
			h.logger.Info("Error uploading transcoded", zap.Error(err))
			return err
		}

		videourlstr := "https://snowy-unit-44da.dhruvj797.workers.dev/" + h.Pathprefix + job.UserId.String() + "/transcode/" + job.VideoId.String()
		videourl := utils.ToText(videourlstr)
		_, err = h.DB.InsertVideoUrl(context.Background(), database.InsertVideoUrlParams{
			StreamUrl: videourl,
			VideoID:   job.VideoId,
		})
		if err != nil {
			h.logger.Info("Error uploading URL to the db", zap.Error(err))
			return err
		}
	}

	//fmt.Println(Avgbitrate)
	return nil
}

// func UploadDirectory(client *minio.Client, bucket, prefix, localDir string) error {
// 	return filepath.WalkDir(localDir, func(path string, d fs.DirEntry, err error) error {
// 		if err != nil {
// 			return err
// 		}
// 		if d.IsDir() {
// 			return nil // skip directories
// 		}

// 		// Create the object name relative to the localDir
// 		relativePath := strings.TrimPrefix(path, localDir)
// 		relativePath = strings.TrimPrefix(relativePath, "/") // prevent double slashes

// 		objectName := prefix + "/" + relativePath

// 		_, err = client.FPutObject(context.Background(), bucket, objectName, path, minio.PutObjectOptions{})
// 		return err
// 	})
// }

func contentTypeFor(fileName string) string {
	if strings.HasSuffix(fileName, ".m3u8") {
		return "application/vnd.apple.mpegurl"

	} else if strings.HasSuffix(fileName, ".ts") {
		return "video/MP2T"
	}
	return "application/octet-stream"
}

func uploadHLSDirectory(minioClient *minio.Client, bucketName, remotePrefix, localDir string) error {
	err := filepath.Walk(localDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		// Relative path to preserve directory structure
		relPath, err := filepath.Rel(localDir, path)
		if err != nil {
			return err
		}

		objectName := filepath.ToSlash(filepath.Join(remotePrefix, relPath))
		contentType := contentTypeFor(path)

		_, err = minioClient.FPutObject(context.Background(), bucketName, objectName, path, minio.PutObjectOptions{
			ContentType: contentType,
		})
		if err != nil {
			return fmt.Errorf("upload failed for %s: %w", path, err)
		}

		fmt.Printf("Uploaded: %s as %s\n", objectName, contentType)
		return nil
	})

	return err
}
