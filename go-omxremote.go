package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web"
)

const fifo string = "omxcontrol"

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

func startVideo(c web.C, w http.ResponseWriter, r *http.Request) {
	filename, _ := base64.StdEncoding.DecodeString(c.URLParams["name"])
	string_filename := string(filename[:])

	fifo_cmd := exec.Command("mkfifo", fifo)
	fifo_cmd.Run()
	cmd := exec.Command("bash", "-c", "omxplayer -o hdmi "+string_filename+" < "+fifo)
	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}

	exec.Command("bash", "-c", "echo . > "+fifo).Run()
	err = cmd.Wait()

	fmt.Fprintf(w, "%s", string_filename)
}

func togglePlayVideo(c web.C, w http.ResponseWriter, r *http.Request) {

	cmd := exec.Command("bash", "-c", "echo -n p > "+fifo)
	cmd.Run()

	fmt.Fprintf(w, "1")
}

func stopVideo(c web.C, w http.ResponseWriter, r *http.Request) {

	cmd := exec.Command("bash", "-c", "echo -n q > "+fifo)
	cmd.Run()
	os.Remove(fifo)

	fmt.Fprintf(w, "1")
}

func main() {

	goji.Get("/", home)
	goji.Get("/files", videoFiles)

	goji.Post("/file/:name/start", startVideo)
	goji.Post("/file/:name/play", togglePlayVideo)
	goji.Post("/file/:name/pause", togglePlayVideo)
	goji.Post("/file/:name/stop", stopVideo)

	goji.Handle("/assets/*", http.FileServer(FS(false)))

	goji.Serve()
}
