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
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web"
)

const fifo string = "omxcontrol"

var videosPath string

type Page struct {
	Title string
}

type Video struct {
	File string `json:"file"`
	Hash string `json:"hash"`
}

func home(c web.C, w http.ResponseWriter, r *http.Request) {
	p := &Page{Title: "gomxremote"}
	tmpl, err := FSString(false, "/views/index.html")
	if err != nil {
		fmt.Println(err.Error())
	}

	t, _ := template.New("index").Parse(tmpl)
	t.Execute(w, p)
}

func videoFiles(c web.C, w http.ResponseWriter, r *http.Request) {
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

func startVideo(c web.C, w http.ResponseWriter, r *http.Request) {
	filename, _ := base64.StdEncoding.DecodeString(c.URLParams["name"])
	string_filename := string(filename[:])
	escapePathReplacer := strings.NewReplacer(
		"[", "\\[",
		"]", "\\]",
		"(", "\\(",
		")", "\\)",
		"'", "\\'",
		" ", "\\ ",
		"*", "\\*",
		"?", "\\?",
	)
	escapedPath := escapePathReplacer.Replace(string_filename)

	fifo_cmd := exec.Command("mkfifo", fifo)
	fifo_cmd.Run()

	cmd := exec.Command("bash", "-c", "omxplayer -o hdmi "+escapedPath+" < "+fifo)
	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}

	startErr := exec.Command("bash", "-c", "echo . > "+fifo).Run()
	if startErr != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	err = cmd.Wait()

	fmt.Fprintf(w, "%s", string_filename)
}

func togglePlayVideo(c web.C, w http.ResponseWriter, r *http.Request) {

	err := sendCommand("play")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	fmt.Fprintf(w, "1")
}

func stopVideo(c web.C, w http.ResponseWriter, r *http.Request) {

	err := sendCommand("quit")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	os.Remove(fifo)

	fmt.Fprintf(w, "1")
}

func toggleSubsVideo(c web.C, w http.ResponseWriter, r *http.Request) {

	err := sendCommand("subs")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	fmt.Fprintf(w, "1")
}

func forwardVideo(c web.C, w http.ResponseWriter, r *http.Request) {

	err := sendCommand("forward")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	fmt.Fprintf(w, "1")
}

func backwardVideo(c web.C, w http.ResponseWriter, r *http.Request) {

	err := sendCommand("backward")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	fmt.Fprintf(w, "1")
}

func sendCommand(command string) error {
	commands := strings.NewReplacer(
		"play", "p",
		"pause", "p",
		"subs", "m",
		"quit", "q",
		"forward", "\x5b\x43",
		"backward", "\x5b\x44",
	)

	commandString := "echo -n " + commands.Replace(command) + " > " + fifo
	cmd := exec.Command("bash", "-c", commandString)
	err := cmd.Run()
	return err
}

func main() {

	flag.StringVar(&videosPath, "media", ".", "path to look for videos in")

	goji.Get("/", home)
	goji.Get("/files", videoFiles)

	goji.Post("/file/:name/start", startVideo)
	goji.Post("/file/:name/play", togglePlayVideo)
	goji.Post("/file/:name/pause", togglePlayVideo)
	goji.Post("/file/:name/stop", stopVideo)
	goji.Post("/file/:name/subs", toggleSubsVideo)
	goji.Post("/file/:name/forward", forwardVideo)
	goji.Post("/file/:name/backward", backwardVideo)

	goji.Handle("/assets/*", http.FileServer(FS(false)))

	goji.Serve()
}
