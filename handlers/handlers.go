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
	"scrapbook/utils"
	"strings"

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
type Config struct {
	Username string
	Password string
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	body, readErr := io.ReadAll(r.Body)
	if readErr != nil {
		fmt.Println(readErr)
	}

	var creds Config
	deserializeErr := json.Unmarshal(body, &creds)
	if deserializeErr != nil {
		fmt.Println(deserializeErr)
	}

	hasher := sha1.New()
	hasher.Write([]byte(creds.Username))
	usernameSha := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
	hasher.Reset()
	hasher.Write([]byte(creds.Password))
	pwSha := base64.URLEncoding.EncodeToString(hasher.Sum(nil))

	config := readConfig()

	if usernameSha == config.Username && pwSha == config.Password {
		fmt.Println("Authenticated.")
		w.WriteHeader(http.StatusOK)
	} else {
		fmt.Println("Not authenticated.")
		w.WriteHeader(http.StatusUnauthorized)
	}
}

func GetTestFileHandler(w http.ResponseWriter, r *http.Request) {
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

func readConfig() Config {
	file, readErr := os.Open("./config.json")
	if readErr != nil {
		fmt.Println(readErr)
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	config := Config{}
	decodeErr := decoder.Decode(&config)
	if decodeErr != nil {
		fmt.Println("decodeErr: ", decodeErr)
	}

	return config
}
