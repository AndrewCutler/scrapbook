package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

var fileDir = "./files/"

type FileMeta struct {
	Name      string
	Size      int64
	Thumbnail string
	Height    int
	Width     int
}

func main() {
	r := mux.NewRouter()
	fs := http.FileServer(http.Dir("./client"))
	r.PathPrefix("").Handler(fs)
	r.HandleFunc("/save", saveFileHandler)
	// http.HandleFunc("/files/", getFileDataHandler)
	r.HandleFunc("/files/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("files")
	})

	handler := cors.Default().Handler(r)
	srv := &http.Server{
		Handler: handler,
		// Addr:    "10.0.0.73:8000",
		Addr:         "localhost:8000",
		WriteTimeout: 5 * time.Second,
		ReadTimeout:  5 * time.Second,
	}

	http.Handle("/", r)
	if err := srv.ListenAndServe(); err != nil {
		panic(err)
	}
}

func saveFileHandler(w http.ResponseWriter, r *http.Request) {
	switch origin := r.Header.Get("Origin"); origin {
	case "http://localhost":
		(w).Header().Set("Access-Control-Allow-Origin", "http://localhost")
		fmt.Println("http://localhost")
	}
	// enableCors(&w, r)
	r.ParseMultipartForm(200)

	form := r.MultipartForm
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

func getFileDataHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("get")

	switch origin := r.Header.Get("Origin"); origin {
	case "http://localhost":
		(w).Header().Set("Access-Control-Allow-Origin", "http://localhost")
		fmt.Println("http://localhost")

	case "http://10.0.0.73":
		(w).Header().Set("Access-Control-Allow-Origin", "http://10.0.0.73")
		fmt.Println("http://10.0.0.73")
	}
	(w).Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")

	// enableCors(&w, r)

	name := strings.TrimPrefix(r.URL.Path, "/files/")
	if name != "" {
		getFile(name, w)
		return
	}

	// buildFileList(w)
	dir, err := os.ReadDir(fileDir)
	if err != nil {
		log.Fatal("Cannot read file directory")
	}

	var filemeta []FileMeta
	for _, f := range dir {
		filename := fileDir + f.Name()
		b, err := os.ReadFile("." + getThumbnailPathFromFilename(filename))
		if err != nil {
			fmt.Println(err)
		}

		if !strings.HasSuffix(f.Name(), ".mp4") && !strings.HasSuffix(f.Name(), ".webm") /* || etc. */ {
			continue
		}

		fi, err := os.Stat(filename)
		if err != nil {
			fmt.Println(err)
		}
		var thumbnail string
		if err != nil {
			wd, _ := os.Getwd()
			fmt.Println(wd, "\n", err)
			thumbnail = ""
		} else {
			thumbnail = "data:image/jpeg;base64," + base64.StdEncoding.EncodeToString(b)
		}

		// get image dimensions
		height := 200
		width := 200
		reader, err := os.Open("." + getThumbnailPathFromFilename(filename))
		defer reader.Close()
		if err == nil {
			im, _, err := image.DecodeConfig(reader)
			height = im.Height
			width = im.Width
			if err != nil {
				fmt.Println(err)
			}
		} else {
			fmt.Println(err)
		}

		filemeta = append(filemeta, FileMeta{
			Name:      fi.Name(),
			Size:      fi.Size(),
			Thumbnail: thumbnail,
			Height:    height,
			Width:     width,
		})
	}

	w.Header().Set("Content-Type", "application-json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(filemeta)
}

func enableCors(w *http.ResponseWriter, r *http.Request) {
	switch origin := r.Header.Get("Origin"); origin {
	case "http://localhost":
		(*w).Header().Set("Access-Control-Allow-Origin", "http://localhost")
		fmt.Println("http://localhost")
	}
	// (*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func buildFileList(w http.ResponseWriter) {
	dir, err := os.ReadDir(fileDir)
	if err != nil {
		log.Fatal("Cannot read file directory")
	}

	var filemeta []FileMeta
	for _, f := range dir {
		filename := fileDir + f.Name()
		b, err := os.ReadFile("." + getThumbnailPathFromFilename(filename))
		if err != nil {
			fmt.Println(err)
		}

		if !strings.HasSuffix(f.Name(), ".mp4") && !strings.HasSuffix(f.Name(), ".webm") /* || etc. */ {
			continue
		}

		fi, err := os.Stat(filename)
		if err != nil {
			fmt.Println(err)
		}
		var thumbnail string
		if err != nil {
			wd, _ := os.Getwd()
			fmt.Println(wd, "\n", err)
			thumbnail = ""
		} else {
			thumbnail = "data:image/jpeg;base64," + base64.StdEncoding.EncodeToString(b)
		}

		// get image dimensions
		height := 200
		width := 200
		reader, err := os.Open("." + getThumbnailPathFromFilename(filename))
		defer reader.Close()
		if err == nil {
			im, _, err := image.DecodeConfig(reader)
			height = im.Height
			width = im.Width
			if err != nil {
				fmt.Println(err)
			}
		} else {
			fmt.Println(err)
		}

		filemeta = append(filemeta, FileMeta{
			Name:      fi.Name(),
			Size:      fi.Size(),
			Thumbnail: thumbnail,
			Height:    height,
			Width:     width,
		})
	}

	w.Header().Set("Content-Type", "application-json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(filemeta)

}

func getFile(name string, w http.ResponseWriter) {
	b, err := os.ReadFile(fileDir + name)
	if err != nil {
		fmt.Println(err)
	}
	w.Write(b)
}

func getVideo(w http.ResponseWriter, r *http.Request) {
	enableCors(&w, r)
	fmt.Println("get video")
}

func createThumbnail(filename string) {
	f := strings.Trim(filename, filepath.Ext(filename))
	var errbuff strings.Builder
	// ffmpeg -ss 1 -i .\input.mp4 -qscale:v 4 -frames:v 1 output.jpeg
	cmd := exec.Command("ffmpeg", "-ss", "1", "-i", fileDir+filename, "-qscale:v", "4", "-frames:v", "1", fileDir+getThumbnailPathFromFilename(f))
	cmd.Stderr = &errbuff
	if err := cmd.Run(); err != nil {
		fmt.Println(errbuff.String())
	}
}

func getThumbnailPathFromFilename(filename string) string {
	return strings.Trim(filename, filepath.Ext(filename)) + ".jpeg"
}
