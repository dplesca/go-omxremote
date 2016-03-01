package main

import (
	"html/template"
	"net/http"

	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web"
)

type Page struct {
	Title string
}

func home(c web.C, w http.ResponseWriter, r *http.Request) {
	p := &Page{Title: "gomxremote"}
	t, _ := template.ParseFiles("views/index.html")
	t.Execute(w, p)
}

func main() {

	goji.Get("/", home)
	goji.Handle("/assets/*", http.FileServer(http.Dir(".")))

	goji.Serve()
}
