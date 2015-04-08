package main

import (
	"fmt"
	"net/http"
)

func HomeHandler(res http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(res, "Dick butt")
}
