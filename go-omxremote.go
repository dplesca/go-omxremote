package main

import (
	"encoding/base32"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/julienschmidt/httprouter"
)

var videosPath string
var bindAddr string
var p Player

// Page is the HTML page struct
type Page struct {
	Title     string
	Timestamp int64
}

// Video struct contains has two fields:
// filename and base32 hash of the filepath
type Video struct {
	File string `json:"file"`
	Hash string `json:"hash"`
}

// Index func that serves the HTML for the home page
func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	p := &Page{Title: "go-omxremote", Timestamp: time.Now().Unix()}
	tmpl, err := FSString(false, "/views/index.html")
	if err != nil {
		fmt.Println(err.Error())
	}

	t, _ := template.New("index").Parse(tmpl)
	t.Execute(w, p)
}

// List function - outputs json with all video files in the videoPath
func List(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var files []*Video
	var root = videosPath
	_ = fastWalk(root, func(path string, f os.FileMode) error {
		if f.IsDir() == false {
			if filepath.Ext(path) == ".mkv" || filepath.Ext(path) == ".mp4" || filepath.Ext(path) == ".avi" {
				files = append(files, &Video{File: filepath.Base(path), Hash: base32.StdEncoding.EncodeToString([]byte(path))})
			}
		}
		return nil
	})
	encoder := json.NewEncoder(w)
	encoder.Encode(files)
}

// Start playback http handler
func Start(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	log.Printf("Start %s", ps.ByName("name"))
	filename, _ := base32.StdEncoding.DecodeString(ps.ByName("name"))
	stringFilename := string(filename[:])

	err := p.Start(stringFilename)
	if err != nil {
		p.Playing = false
		http.Error(w, err.Error(), 500)
		return
	}

	log.Printf("Playing media file: %s\n", stringFilename)
	startTime := time.Now()
	err = p.Command.Wait()

	log.Printf("Stopped media file: %s after %s\n", stringFilename, time.Since(startTime))
	p.Playing = false
	w.WriteHeader(http.StatusOK)
}

// SendCommand is the HTTP handler for sending a command to the player
func SendCommand(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	err := p.SendCommand(ps.ByName("command"))
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
	flag.Parse()

	router := httprouter.New()
	router.GET("/", Index)
	router.GET("/files.json", List)

	router.POST("/start/:name", Start)
	router.POST("/player/:command", SendCommand)

	router.ServeFiles("/dist/*filepath", FS(false))
	log.Fatal(http.ListenAndServe(bindAddr, router))
}
