package utils

import (
	"fmt"
	"net/http"
	"os/exec"
	"path/filepath"
	"strings"
)

func GetVideo(w http.ResponseWriter, r *http.Request) {
	fmt.Println("get video")
}

func CreateThumbnail(filename string, fileDir string) {
	f := strings.Trim(filename, filepath.Ext(filename))
	var errbuff strings.Builder
	// ffmpeg -ss 1 -i .\input.mp4 -qscale:v 4 -frames:v 1 output.jpeg
	cmd := exec.Command("ffmpeg", "-ss", "1", "-i", fileDir+filename, "-qscale:v", "4", "-frames:v", "1", fileDir+GetThumbnailPathFromFilename(f))
	cmd.Stderr = &errbuff
	if err := cmd.Run(); err != nil {
		fmt.Println(errbuff.String())
	}
}

func GetThumbnailPathFromFilename(filename string) string {
	return strings.Trim(filename, filepath.Ext(filename)) + ".jpeg"
}
