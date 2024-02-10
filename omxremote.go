package main

import (
	"embed"
	"encoding/base64"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/lithammer/fuzzysearch/fuzzy"
)

var videosPath string
var bindAddr string
var omx string

var p Player

//go:embed public
var static embed.FS

//go:embed views
var tplFolder embed.FS

// Page is the HTML page struct
type Page struct {
	Title, Search string
	Timestamp     int64
	Files         []Video
}

// Video struct contains has two fields:
// filename and base32 hash of the filepath
type Video struct {
	File string `json:"file"`
	Hash string `json:"hash"`
}

// Index func that serves the HTML for the home page
func Index(w http.ResponseWriter, r *http.Request) {
	var files []Video
	var root = videosPath
	s := r.URL.Query().Get("search")
	_ = filepath.Walk(root, func(path string, f os.FileInfo, _ error) error {
		if !f.IsDir() {
			if filepath.Ext(path) == ".mkv" || filepath.Ext(path) == ".mp4" || filepath.Ext(path) == ".avi" {
				fn := filepath.Base(path)
				if s != "" {
					if fuzzy.Match(s, fn) {
						files = append(files, Video{File: filepath.Base(path), Hash: base64.URLEncoding.EncodeToString([]byte(path))})
					}
				} else {
					files = append(files, Video{File: filepath.Base(path), Hash: base64.URLEncoding.EncodeToString([]byte(path))})
				}
			}
		}
		return nil
	})

	p := &Page{Title: "go-omxremote", Timestamp: time.Now().Unix(), Files: files, Search: s}
	tmpl, err := template.ParseFS(tplFolder, "views/index.html")
	if err != nil {
		fmt.Println(err.Error())
	}

	tmpl.Execute(w, p)
}

// Start playback http handler
func Start(w http.ResponseWriter, r *http.Request) {
	filename, _ := base64.URLEncoding.DecodeString(r.PathValue("name"))
	stringFilename := string(filename[:])
	omxOptions := append(strings.Split(omx, " "), stringFilename)

	err := p.Start(omxOptions)
	if err != nil {
		p.Playing = false
		http.Error(w, err.Error(), 500)
		return
	}

	log.Printf("Playing media file: %s\n", stringFilename)
	w.WriteHeader(http.StatusOK)
}

// SendCommand is the HTTP handler for sending a command to the player
func SendCommand(w http.ResponseWriter, r *http.Request) {
	err := p.SendCommand(r.PathValue("command"))
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func main() {
	flag.StringVar(&videosPath, "media", ".", "Path to look for videos in")
	flag.StringVar(&bindAddr, "bind", ":31415", "Address to bind on.")
	flag.StringVar(&omx, "omx", "-o hdmi", "Options to pass to omxplayer")
	flag.Parse()

	mux := http.NewServeMux()
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServerFS(static)))
	mux.HandleFunc("GET /", Index)
	mux.HandleFunc("POST /start/{name}", Start)
	mux.HandleFunc("POST /player/{command}", SendCommand)
	http.ListenAndServe(bindAddr, mux)
}
