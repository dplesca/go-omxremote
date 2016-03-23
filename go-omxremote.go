package main

import (
	"bytes"
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

	"github.com/julienschmidt/httprouter"
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
	var buf bytes.Buffer

	filename, _ := base64.StdEncoding.DecodeString(ps.ByName("name"))
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

	if _, err := os.Stat(fifo); err == nil {
		os.Remove(fifo)
	}

	fifo_cmd := exec.Command("mkfifo", fifo)
	fifo_cmd.Run()

	cmd := exec.Command("bash", "-c", "omxplayer -o hdmi "+escapedPath+" < "+fifo)
	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}

	buf.Write([]byte{'\033', '[', '3', '4', ';', '1', 'm'})
	fmt.Fprintf(&buf, "%s", string_filename)
	log.Print(buf.String())

	startErr := exec.Command("bash", "-c", "echo . > "+fifo).Run()
	if startErr != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	err = cmd.Wait()

	w.WriteHeader(http.StatusOK)
}

func TogglePlay(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	sendCommand("play", w)
	w.WriteHeader(http.StatusOK)
}

func Stop(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	sendCommand("quit", w)
	os.Remove(fifo)
	w.WriteHeader(http.StatusOK)
}

func ToggleSubs(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	sendCommand("subs", w)
	w.WriteHeader(http.StatusOK)
}

func Forward(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	sendCommand("forward", w)
	w.WriteHeader(http.StatusOK)
}

func Backward(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	sendCommand("backward", w)
	w.WriteHeader(http.StatusOK)
}

func sendCommand(command string, w http.ResponseWriter) {
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
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func main() {
	flag.StringVar(&videosPath, "media", ".", "path to look for videos in")

	router := httprouter.New()
	router.GET("/", Index)
	router.GET("/files", List)

	router.POST("/file/:name/start", Start)
	router.POST("/file/:name/play", TogglePlay)
	router.POST("/file/:name/pause", TogglePlay)
	router.POST("/file/:name/stop", Stop)
	router.POST("/file/:name/subs", ToggleSubs)
	router.POST("/file/:name/forward", Forward)
	router.POST("/file/:name/backward", Backward)

	router.ServeFiles("/assets/*filepath", http.Dir("./assets"))
	log.Fatal(http.ListenAndServe(":8080", router))
}
