package main

import (
	"github.com/gorilla/mux"

	"fmt"
	"net/http"
)

func DickButtHandler(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	fmt.Fprintln(res, "<html><head><style>img { position: absolute; }</style></head><body><img src=\""+ImgurSearcher(vars["place"])+"\"></body><img src=\"/assets/dickbutt.png\"/></html>")
}
