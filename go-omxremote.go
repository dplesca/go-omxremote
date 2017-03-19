package main

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/fatih/color"
	"github.com/julienschmidt/httprouter"
)

var videosPath string
var bindAddr string
var p Player

type Page struct {
	Title string
}

type Video struct {
	File string `json:"file"`
	Hash string `json:"hash"`
}

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	p := &Page{Title: "go-omxremote"}
	tmpl, err := FSString(false, "/views/index.html")
	if err != nil {
		fmt.Println(err.Error())
	}

	t, _ := template.New("index").Parse(tmpl)
	t.Execute(w, p)
}

func List(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var files []*Video
	var root = videosPath
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

func Start(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	filename, _ := base64.StdEncoding.DecodeString(ps.ByName("name"))
	string_filename := string(filename[:])

	err := p.Start(string_filename)
	if err != nil {
		p.Playing = false
		http.Error(w, err.Error(), 500)
		return
	}

	color.Green("Playing media file: %s\n", string_filename)
	startTime := time.Now()
	err = p.Command.Wait()

	color.Red("Stopped media file: %s after %s\n", string_filename, time.Since(startTime))
	p.Playing = false
	w.WriteHeader(http.StatusOK)
}

func SendCommand(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	err := p.SendCommand(ps.ByName("command"))
	if err != nil {
		color.Red(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func main() {
	flag.StringVar(&videosPath, "media", ".", "Path to look for videos in")
	flag.StringVar(&bindAddr, "bind", ":31415", "Address to bind on.")
	flag.Parse()

	router := httprouter.New()
	router.GET("/", Index)
	router.GET("/files", List)

	router.POST("/start/:name", Start)
	router.POST("/file/:name/:command", SendCommand)

	router.ServeFiles("/assets/*filepath", FS(false))
	log.Fatal(http.ListenAndServe(bindAddr, router))
}
