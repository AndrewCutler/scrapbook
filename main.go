package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

var fileDir = "./files/"

type FileMeta struct {
	Name string
	Size int64
}

func main() {
	createThumbnail("sample-30s.mp4")
	// fs := http.FileServer(http.Dir("./client"))
	// http.HandleFunc("/save", saveFile)
	// http.HandleFunc("/files", getFileData)
	// http.Handle("/", fs)

	// if err := http.ListenAndServe(":8000", nil); err != nil {
	// 	panic(err)
	// }
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

		filemeta = append(filemeta, FileMeta{
			Name: fi.Name(),
			Size: fi.Size(),
		})
	}

	w.Header().Set("Content-Type", "application-json")
	json.NewEncoder(w).Encode(filemeta)
}

func createThumbnail(filename string) {
	var buffer bytes.Buffer
	fmt.Println(fileDir + strings.Replace(filename, ".mp4", ".jpg", 1))
	// ffmpeg -ss 1 -i .\input.mp4 -qscale:v 4 -frames:v 1 output.jpg
	cmd := exec.Command("ffmpeg", "-ss", "1", "-i", fileDir+filename, "-qscale:v", "4", "-frames:v", "1", fileDir+strings.Replace(filename, ".mp4", ".jpg", 1))
	cmd.Stdout = &buffer
	if err := cmd.Run(); err != nil {
		fmt.Println(err)
	}
}
