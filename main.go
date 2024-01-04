package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var fileDir = "./files/"

type FileMeta struct {
	Name      string
	Size      int64
	Thumbnail string
}

func main() {
	// createThumbnail("sample-30s.mp4")
	fs := http.FileServer(http.Dir("./client"))
	http.HandleFunc("/save", saveFile)
	http.HandleFunc("/files", getFileData)
	http.Handle("/", fs)

	if err := http.ListenAndServe(":8000", nil); err != nil {
		panic(err)
	}
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func saveFile(w http.ResponseWriter, h *http.Request) {
	enableCors(&w)
	h.ParseMultipartForm(200)

	form := h.MultipartForm
	files := form.File["files"]
	if files == nil {
		log.Fatal("No files received")
	}

	for _, f := range files {
		fmt.Print(f.Filename)
		file, err := f.Open()
		if err != nil {
			log.Fatal(err)
		}
		out, err := os.Create(fileDir + f.Filename)
		if err != nil {
			log.Fatal(err)
		}
		defer out.Close()

		if _, err := io.Copy(out, file); err != nil {
			log.Fatal(err)
		}

		createThumbnail(f.Filename)
	}
}

func getFileData(w http.ResponseWriter, h *http.Request) {
	dir, err := os.ReadDir(fileDir)
	if err != nil {
		log.Fatal("Cannot read file directory")
	}

	var filemeta []FileMeta
	for _, f := range dir {
		fi, err := os.Stat(fileDir + f.Name())
		if err != nil {
			fmt.Println(err)
		}

		// get thumbnail, base64-encoded
		b, err := os.ReadFile(fileDir + f.Name())
		if err != nil {
			fmt.Println(err)
		}
		thumbnail := "data:image/jpeg;base64," + base64.StdEncoding.EncodeToString(b)

		filemeta = append(filemeta, FileMeta{
			Name:      fi.Name(),
			Size:      fi.Size(),
			Thumbnail: thumbnail,
		})
	}

	w.Header().Set("Content-Type", "application-json")
	json.NewEncoder(w).Encode(filemeta)
}

func createThumbnail(filename string) {
	f := strings.Trim(filename, filepath.Ext(filename))
	var errbuff strings.Builder
	// ffmpeg -ss 1 -i .\input.mp4 -qscale:v 4 -frames:v 1 output.jpeg
	cmd := exec.Command("ffmpeg", "-ss", "1", "-i", fileDir+filename, "-qscale:v", "4", "-frames:v", "1", fileDir+f+".jpeg")
	cmd.Stderr = &errbuff
	if err := cmd.Run(); err != nil {
		fmt.Println(errbuff.String())
	}
}
