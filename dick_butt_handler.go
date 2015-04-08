package main

import (
	"github.com/gorilla/mux"

	"fmt"
	"html/template"
	"math/rand"
	"net/http"
)

var dickTemplate = `
{{define "page"}}
<html>
	<head>
		<style>
			img {
				position: absolute;
			}
		</style>
	</head>
	<body style="background-image: url('{{.ImgurSource}}'); background-size: cover; background-position: center;">
	<a href='{{.Place}}'><img style="top:{{.Top}}%; left:{{.Bottom}}%" src="/assets/dickbutt.png"/></a>
	</body>
</html>
{{end}}
`

type Page struct {
	ImgurSource string
	Top         int
	Bottom      int
	Place       string
}

func DickButtHandler(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	p := Page{
		ImgurSource: ImgurSearcher(vars["place"]),
		Top:         rand.Intn(80),
		Bottom:      rand.Intn(80),
		Place:       vars["place"],
	}
	fmt.Println(p)
	templ, err := template.New("page").Parse(dickTemplate)
	if err != nil {
		panic(err)
	}
	err = templ.ExecuteTemplate(res, "page", p)
	if err != nil {
		panic(err)
	}
}
