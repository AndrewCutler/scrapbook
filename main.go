package main

import (
	"log"
	"net/http"
	"time"

	"scrapbook/auth"
	handlers "scrapbook/handlers"

	"github.com/gorilla/mux"
)

var fileDir = "./files/"

func main() {
	r := mux.NewRouter()
	session := new(auth.Session)

	r.HandleFunc("/api/test", auth.Authenticate(handlers.GetTestFileHandler, session)).Methods("GET")
	r.HandleFunc("/api/save", auth.Authenticate(handlers.SaveFileHandler, session)).Methods("POST")
	r.HandleFunc("/api/files", auth.Authenticate(handlers.BuildFileListHandler, session)).Methods("GET")
	r.HandleFunc("/api/login", auth.WithSession(handlers.LoginHandler, session)).Methods("POST")
	r.HandleFunc("/api/files/{filename}", auth.Authenticate(handlers.GetFileHandler, session)).Methods("GET")

	spa := handlers.SpaHandler{StaticPath: "client", IndexPath: "index.html"}
	r.PathPrefix("/").Handler(spa)
	srv := &http.Server{
		Handler:      r,
		Addr:         "10.0.0.73:8000",
		WriteTimeout: 5 * time.Second,
		ReadTimeout:  5 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
