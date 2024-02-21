package handlers

import (
	"crypto/sha1"
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

	"github.com/google/uuid"
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

type SpaHandler struct {
	StaticPath string
	IndexPath  string
}

func (h SpaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Join internally call path.Clean to prevent directory traversal
	if strings.HasPrefix(r.URL.Path, "/api") {
		return
	}

	path := filepath.Join(h.StaticPath, r.URL.Path)

	// check whether a file exists or is a directory at the given path
	fi, err := os.Stat(path)
	if os.IsNotExist(err) || fi.IsDir() {
		// file does not exist or path is a directory, serve index.html
		http.ServeFile(w, r, filepath.Join(h.StaticPath, h.IndexPath))
		return
	}

	if err != nil {
		// if we got an error (that wasn't that the file doesn't exist) stating the
		// file, return a 500 internal server error and stop
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// otherwise, use http.FileServer to serve the static file
	http.FileServer(http.Dir(h.StaticPath)).ServeHTTP(w, r)
}
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	body, readErr := io.ReadAll(r.Body)
	if readErr != nil {
		http.Error(w, "Error reading body", http.StatusBadRequest)
		return
	}

	var creds utils.Config
	fmt.Println(string(body))
	deserializeErr := json.Unmarshal(body, &creds)
	if deserializeErr != nil {
		fmt.Println(deserializeErr)
		http.Error(w, "Deserialization error", http.StatusBadRequest)
		return
	}

	hasher := sha1.New()
	hasher.Write([]byte(creds.Username))
	usernameSha := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
	hasher.Reset()
	hasher.Write([]byte(creds.Password))
	pwSha := base64.URLEncoding.EncodeToString(hasher.Sum(nil))

	config := utils.ReadConfig()

	if usernameSha == config.Username && pwSha == config.Password {
		fmt.Println("Authenticated.")
		token := uuid.NewString()
		expiration := time.Now().Add(3000 * time.Second)

		http.SetCookie(w, &http.Cookie{
			Name:     "session_token",
			HttpOnly: false,
			Value:    token,
			Expires:  expiration,
		})

		w.WriteHeader(http.StatusOK)
	} else {
		fmt.Println("Not authenticated.")
		w.WriteHeader(http.StatusUnauthorized)
	}
}

func GetTestFileHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("api test")
	b, err := os.ReadFile("./client/test.js")
	if err != nil {
		fmt.Println(err)
		return
	}
	w.Write(b)
	w.Header().Set("Cache-Control", "max-age=10000000")
}

func SaveFileHandler(w http.ResponseWriter, r *http.Request) {
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

func BuildFileListHandler(w http.ResponseWriter, r *http.Request) {
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
	w.Header().Set("Cache-Control", "max-age=600")
	json.NewEncoder(w).Encode(filemeta)

}

func GetFileHandler(w http.ResponseWriter, r *http.Request) {
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
