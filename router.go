package main

import (
	"github.com/gorilla/mux"
	"net/http"
)

func setupRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", HomeHandler)
	r.HandleFunc("/{place}", DickButtHandler)
	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets/"))))
	return r
}
