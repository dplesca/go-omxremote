package main

import (
	"encoding/base64"
	"encoding/json"

	"html/template"
	"net/http"
	"os"
	"path/filepath"

	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web"
)

type Page struct {
	Title string
}

type Video struct {
	File string `json:"file"`
	Hash string `json:"hash"`
}

func home(c web.C, w http.ResponseWriter, r *http.Request) {
	p := &Page{Title: "gomxremote"}
	t, _ := template.ParseFiles("views/index.html")
	t.Execute(w, p)
}

func videoFiles(c web.C, w http.ResponseWriter, r *http.Request) {
	var files []*Video
	var root = "."
	_ = filepath.Walk(root, func(path string, f os.FileInfo, _ error) error {
		if f.IsDir() == false {
			if filepath.Ext(path) == ".mkv" || filepath.Ext(path) == ".mp4" || filepath.Ext(path) == ".avi" {
				files = append(files, &Video{File: filepath.Base(path), Hash: base64.StdEncoding.EncodeToString([]byte(path))})
			}
		}
		return nil
	})
	encoder := json.NewEncoder(w)
	encoder.Encode(files)
}

func main() {

	goji.Get("/", home)
	goji.Get("/files", videoFiles)
	goji.Handle("/assets/*", http.FileServer(http.Dir(".")))

	goji.Serve()
}
