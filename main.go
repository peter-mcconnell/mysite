package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"text/template"
	"encoding/xml"
	"google.golang.org/appengine"
	"google.golang.org/appengine/urlfetch"
)

var templates = template.Must(template.ParseGlob("templates/*.html"))

func newsHandler(w http.ResponseWriter, r *http.Request) {
	type rssItem struct {
		Title string `xml:"title"`
		Desc string `xml:"description"`
		Link string `xml:"link"`
		Date string `xml:"pubDate"`
	}
	type rssChannel struct {
		Items []rssItem `xml:"item"`
	}
	type rssStruct struct {
		Channel rssChannel `xml:"channel"`
	}
	type newsTemplateVars struct {
		Rss rssStruct
	}
	req, err := http.NewRequest("GET", "http://digg.com/user/2bbe142e874a424cb6c56c3752d62892/diggs.rss", nil)
	if err != nil {
		log.Fatal(err)
		return
	}
	c := appengine.NewContext(r)
	client := urlfetch.Client(c)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
		return
	}
	var n newsTemplateVars
	err = xml.Unmarshal(body, &n.Rss)
	if err != nil {
		log.Fatal(err)
		return
	}
	log.Print(n)
	if err = templates.ExecuteTemplate(w, "news", n); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("home handler triggered")
	if r.URL.Path != "/" {
		errorHandler(w, r, http.StatusNotFound)
		return
	}
	if err := templates.ExecuteTemplate(w, "home", nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func errorHandler(w http.ResponseWriter, r *http.Request, status int) {
	w.WriteHeader(status)
	if status == http.StatusNotFound {
		fmt.Fprint(w, "idk what you're lookin' for ¯\\_(ツ)_/¯")
	}
}

func init() {
	log.Println("initiated ...")
	// static assets
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// handlers
	http.HandleFunc("/news", newsHandler)
	http.HandleFunc("/", homeHandler)
}
