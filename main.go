package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

var fileDir = "./files/"

func main() {
	fs := http.FileServer(http.Dir("./client"))
	http.HandleFunc("/save", saveFile)
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
	for key := range form.File {
		file, header, err := h.FormFile(key)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		fmt.Printf("Writing file %s to disk\n", header.Filename)

		out, err := os.Create(fileDir + header.Filename)
		if err != nil {
			log.Fatal(err)
		}
		defer out.Close()

		if _, err := io.Copy(out, file); err != nil {
			log.Fatal(err)
		}
	}
}
