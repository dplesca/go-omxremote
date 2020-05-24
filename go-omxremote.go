package main

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"time"
	"strings"
	"os/exec"

	"github.com/julienschmidt/httprouter"
	"github.com/karrick/godirwalk"
)

var videosPath string
var bindAddr string
var omx string
var yt string
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
	_ = godirwalk.Walk(root, &godirwalk.Options{
		Unsorted: true, // set true for faster yet non-deterministic enumeration (see godoc)
		Callback: func(osPathname string, de *godirwalk.Dirent) error {
			if de.IsDir() == false {
				if filepath.Ext(osPathname) == ".mkv" || filepath.Ext(osPathname) == ".mp4" || filepath.Ext(osPathname) == ".avi" {
					files = append(files, &Video{File: filepath.Base(osPathname), Hash: base64.URLEncoding.EncodeToString([]byte(osPathname))})
				}
			}
			return nil
		},
		ErrorCallback: func(osPathname string, err error) godirwalk.ErrorAction {
			return godirwalk.SkipNode
		},
	})

	if len(files) > 0 {
		encoder := json.NewEncoder(w)
		encoder.Encode(files)
	} else {
		w.Write([]byte("[]"))
	}
}

// Start playback http handler
func Start(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	log.Printf("Start %s", ps.ByName("name"))
	filename, _ := base64.URLEncoding.DecodeString(ps.ByName("name"))
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

// Play youtube http handler
func Play(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if err := r.ParseForm(); err != nil {
		p.Playing = false
		http.Error(w, err.Error(), 500)
		return
	}
	log.Printf("Play %s", r.FormValue("link"))
	ytOptions := append(strings.Split(yt, " "), r.FormValue("link"))
	out, err := exec.Command("youtube-dl", ytOptions...).Output()
    if err != nil {
		p.Playing = false
		http.Error(w, err.Error(), 500)
		return
	}
	stringFilename := strings.Trim(string(out), "\r\n")
	omxOptions := append(strings.Split(omx, " "), stringFilename)

	err = p.Start(omxOptions)
	if err != nil {
		p.Playing = false
		http.Error(w, err.Error(), 500)
		return
	}

	log.Printf("Playing media file: %s\n", stringFilename)
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
	flag.StringVar(&omx, "omx", "-o hdmi", "Options to pass to omxplayer")
	flag.StringVar(&yt, "yt", "-g", "Options to pass to youtube-dl")
	flag.Parse()

	router := httprouter.New()
	router.GET("/", Index)
	router.GET("/files.json", List)

	router.POST("/start/:name", Start)
	router.POST("/play/", Play)
	router.POST("/player/:command", SendCommand)

	router.ServeFiles("/dist/*filepath", FS(false))
	log.Fatal(http.ListenAndServe(bindAddr, router))
}
