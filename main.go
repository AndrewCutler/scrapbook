package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

var fileDir = "./files"

func main() {
	fs := http.FileServer(http.Dir("./client"))
	http.HandleFunc("/save", saveFile)
	http.Handle("/", fs)

	if err := http.ListenAndServe(":8000", nil); err != nil {
		panic(err)
	}
}

func saveFile(w http.ResponseWriter, h *http.Request) {
	h.ParseMultipartForm(200)

	form := h.MultipartForm

	for key := range form.File {
		file, header, err := h.FormFile(key)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()
		fmt.Println(header.Filename)

		out, err := os.Create(fileDir + header.Filename)
		if err != nil {
			log.Fatal(err)
		}
		defer out.Close()

		if _, err := io.Copy(out, file); err != nil {
			log.Fatal(err)
		}
	}
	// form, err := h.MultipartReader()
	// for {
	// 	part, err_part := form.NextPart()
	// 	if err_part == io.EOF {
	// 		break
	// 	}
	// 	if part.FormName() == "file" {
	// 		x, err := part.Read()
	// 		i
	// 		// go fmt.Println()
	// 	}
	// }
	// body, err := io.ReadAll(h.Body)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// go fmt.Println(body)
}
