package main

import (
	"fmt"
	"net/http"
)

func HomeHandler(res http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(res, "Please goto http://www.dickbutt.in/[whatever you want goes here] for some dick butt fun\nI use Imgur api to grab images, credits to come.")
}
