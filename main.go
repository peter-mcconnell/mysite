package main

import (
	"fmt"
	"log"
	"net/http"
	"text/template"
)

var templates = template.Must(template.ParseGlob("templates/*.html"))

func errorHandler(w http.ResponseWriter, r *http.Request, status int) {
	w.WriteHeader(status)
	if status == http.StatusNotFound {
		fmt.Fprint(w, "idk what you're lookin' for ¯\\_(ツ)_/¯")
	}
}

func main() {
	log.Println("initiated ...")
	// static assets
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// public handlers
	http.HandleFunc("/", HomeHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
