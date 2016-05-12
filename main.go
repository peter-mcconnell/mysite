package main

import (
	"fmt"
	"net/http"
	"text/template"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		errorHandler(w, r, http.StatusNotFound)
		return
	}
	t, _ := template.ParseFiles("templates/index.html")
	t.Execute(w, nil)
}

func errorHandler(w http.ResponseWriter, r *http.Request, status int) {
	w.WriteHeader(status)
	if status == http.StatusNotFound {
		fmt.Fprint(w, "idk what you're lookin' for ¯\\_(ツ)_/¯")
	}
}

func init() {
	// static assets
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// handlers
	http.HandleFunc("/", homeHandler)
}
