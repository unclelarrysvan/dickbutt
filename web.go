package main

import (
	"net/http"
	"os"
)

func main() {
	r := setupRouter()
	http.Handle("/", r)

	err := http.ListenAndServe(":"+os.Getenv("PORT"), nil)

	if err != nil {
		panic(err)
	}
}
