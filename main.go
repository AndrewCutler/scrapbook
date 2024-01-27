package main

import (
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	utils "scrapbook/utils"

	"github.com/gorilla/mux"
)

var fileDir = "./files/"

type FileMeta struct {
	Name      string
	Size      int64
	Thumbnail string
	Height    int
	Width     int
}

type spaHandler struct {
	staticPath string
	indexPath  string
}

func (h spaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Join internally call path.Clean to prevent directory traversal
	if strings.HasPrefix(r.URL.Path, "/api") {
		return
	}

	path := filepath.Join(h.staticPath, r.URL.Path)

	// check whether a file exists or is a directory at the given path
	fi, err := os.Stat(path)
	if os.IsNotExist(err) || fi.IsDir() {
		// file does not exist or path is a directory, serve index.html
		http.ServeFile(w, r, filepath.Join(h.staticPath, h.indexPath))
		return
	}

	if err != nil {
		// if we got an error (that wasn't that the file doesn't exist) stating the
		// file, return a 500 internal server error and stop
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// otherwise, use http.FileServer to serve the static file
	http.FileServer(http.Dir(h.staticPath)).ServeHTTP(w, r)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/api/test", useBasicAuth(getTestFileHandler)).Methods("GET")
	r.HandleFunc("/api/save", saveFileHandler).Methods("POST")
	r.HandleFunc("/api/files", buildFileListHandler).Methods("GET")
	// r.HandleFunc("/api/files", useBasicAuth(buildFileListHandler)).Methods("GET")
	r.HandleFunc("/api/files/{filename}", useBasicAuth(getFileHandler)).Methods("GET")

	spa := spaHandler{staticPath: "client", indexPath: "index.html"}
	r.PathPrefix("/").Handler(spa)
	srv := &http.Server{
		Handler:      r,
		Addr:         "10.0.0.73:8000",
		WriteTimeout: 5 * time.Second,
		ReadTimeout:  5 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}

func getTestFileHandler(w http.ResponseWriter, r *http.Request) {
	b, err := os.ReadFile("./client/test.js")
	if err != nil {
		fmt.Println(err)
		return
	}
	w.Write(b)
}

func saveFileHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(200)

	form := r.MultipartForm
	files := form.File["files"]
	if files == nil {
		log.Fatal("No files received")
	}

	for _, f := range files {
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

		utils.CreateThumbnail(f.Filename, fileDir)
	}
}

func buildFileListHandler(w http.ResponseWriter, r *http.Request) {
	dir, err := os.ReadDir(fileDir)
	if err != nil {
		log.Fatal("Cannot read file directory")
	}

	var filemeta []FileMeta
	for _, f := range dir {
		filename := fileDir + f.Name()
		b, err := os.ReadFile("." + utils.GetThumbnailPathFromFilename(filename))
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
		reader, err := os.Open("." + utils.GetThumbnailPathFromFilename(filename))
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

func getFileHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["filename"]

	if name == "" {
		fmt.Println("empty name received")
		return
	}
	b, err := os.ReadFile(fileDir + name)
	if err != nil {
		fmt.Println(err)
		return
	}
	w.Write(b)
}

func useBasicAuth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, pw, ok := r.BasicAuth()
		if ok {
			usernameHash := sha256.Sum256([]byte(username))
			passwordHash := sha256.Sum256([]byte(pw))
			expectedUsernameHash := sha256.Sum256([]byte("username"))
			expectedPasswordHash := sha256.Sum256([]byte("password"))

			usernameMatch := (subtle.ConstantTimeCompare(usernameHash[:], expectedUsernameHash[:]) == 1)
			passwordMatch := (subtle.ConstantTimeCompare(passwordHash[:], expectedPasswordHash[:]) == 1)

			if usernameMatch && passwordMatch {
				next.ServeHTTP(w, r)
				return
			}
		}

		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	})
}
