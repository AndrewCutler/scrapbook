package main

import (
	"net/http"
)

func main() {
	fs := http.FileServer(http.Dir("./client"))
	http.Handle("/", fs)

	if err := http.ListenAndServe(":8000", nil); err != nil {
		panic(err)
	}
}
