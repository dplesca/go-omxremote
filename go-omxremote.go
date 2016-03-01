package main

import (
	"fmt"
	"net/http"

	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web"
)

func home(c web.C, w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s", "hello world")

}

func main() {

	goji.Get("/", home)
	goji.Serve()
}
