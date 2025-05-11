package service

import (
	"fmt"
	"os"
	"os/exec"
)

func CreateHls(job *Job, inputpath, outputBase string) error {

	renditions := []Rendition{
		{"1080p", 1920, 1080, 5000},
		{"720p", 1280, 720, 2800},
		{"480p", 854, 480, 1400},
	}
	for _, r := range renditions {
		outDir := fmt.Sprintf("%s/%s", outputBase+job.VideoId.String(), r.Name)
		os.MkdirAll(outDir, 0755)

		cmd := exec.Command("ffmpeg", "-i", inputpath,
			"-codec:", "copy",
			"-start_number", "0",
			"-hls_time", "6",
			"-hls_list_size", "0",
			"-preset", "fast", "-tune", "zerolatency",
			"-c:a", "aac",
			"-f", "hls",
			"-hls_segment_filename", fmt.Sprintf("%s/segment_%%03d.ts", outDir),
			fmt.Sprintf("%s/index.m3u8", outDir),
		)

		err := cmd.Run()
		if err != nil {
			return err
		}
	}

	err := generateMasterPlaylist(outputBase + job.VideoId.String())
	if err != nil {
		return err
	}

	return nil
}

func CreateTranscoding(job *Job, inputpath, outputDir string) error {

	OutputPath := outputDir + job.VideoId.String()
	resolution := map[string]Dimension{
		"360":  {360, 640},
		"480":  {480, 854},
		"720":  {720, 1280},
		"1080": {1080, 1920},
	}

	args := []string{
		"-i", inputpath,
		"-vf", fmt.Sprintf("scale=w=%d:h=%d", resolution[job.Options.Resolution].Width, resolution[job.Options.Resolution].Height),
		"-c:v", job.Options.Codec,
		"-c:a", "aac",
		OutputPath + job.Options.Output,
	}

	cmd := exec.Command("ffmpeg", args...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil

}

// func returnResolutions(res int) ([]int, error) {

// 	if res >= 360 && res < 480 {
// 		return []int{360}, nil
// 	} else if res >= 480 && res < 720 {
// 		return []int{360, 480}, nil
// 	} else if res >= 720 && res < 1080 {
// 		return []int{360, 480, 720}, nil
// 	} else if res >= 1080 {
// 		return []int{480, 720, 1080}, nil
// 	} else {
// 		return []int{}, errors.New("resolution not supported")
// 	}

// }

func generateMasterPlaylist(outputBase string) error {
	renditions := []struct {
		Name       string
		Bandwidth  int // bits per second
		Resolution string
	}{
		{"1080p", 5000000, "1920x1080"},
		{"720p", 2800000, "1280x720"},
		{"480p", 1400000, "854x480"},
		//{"360p", 800000, "512x360"},
	}

	f, err := os.Create(outputBase + "/master.m3u8")
	if err != nil {
		return err
	}
	defer f.Close()

	f.WriteString("#EXTM3U\n")
	f.WriteString("#EXT-X-VERSION:3\n")

	for _, r := range renditions {
		f.WriteString(fmt.Sprintf(
			"#EXT-X-STREAM-INF:BANDWIDTH=%d,RESOLUTION=%s\n%s/index.m3u8\n",
			r.Bandwidth, r.Resolution, r.Name))
	}

	return nil
}
