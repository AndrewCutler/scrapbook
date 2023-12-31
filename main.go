package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func main() {
	fs := http.FileServer(http.Dir("./client"))
	http.HandleFunc("/save", saveFile)
	http.Handle("/", fs)

	if err := http.ListenAndServe(":8000", nil); err != nil {
		panic(err)
	}
}

func saveFile(w http.ResponseWriter, h *http.Request) {
	body, err := io.ReadAll(h.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(body)
}
