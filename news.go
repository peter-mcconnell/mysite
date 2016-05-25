package main

import (
	"encoding/xml"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"golang.org/x/net/context"

	"google.golang.org/appengine"
	"google.golang.org/appengine/memcache"
	"google.golang.org/appengine/urlfetch"
)

type rssItem struct {
	Title string `xml:"title"`
	Desc  string `xml:"description"`
	Link  string `xml:"link"`
	Date  string `xml:"pubDate"`
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

func fetchNewsFromRSS(n newsTemplateVars, c context.Context) newsTemplateVars {
	log.Println("fetchNewsFromRSS")
	req, err := http.NewRequest("GET", "http://digg.com/user/2bbe142e874a424cb6c56c3752d62892/diggs.rss", nil)
	if err != nil {
		log.Fatal(err)
		return n
	}
	client := urlfetch.Client(c)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
		return n
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
		return n
	}
	err = xml.Unmarshal(body, &n.Rss)
	if err != nil {
		log.Fatal(err)
		return n
	}
	return n
}

func getNewsFromMemcache(c context.Context) rssStruct {
	log.Println("getNewsFromMemcache")
	var o rssStruct
	memcache.Gob.Get(c, "NewsRss", &o)
	return o
}

func setNewsMemcache(o rssStruct, c context.Context) {
	log.Println("setNewsMemcache")
	item := &memcache.Item{
		Key:        "NewsRss",
		Object:     o,
		Expiration: time.Hour * 24,
	}
	if err := memcache.Gob.Set(c, item); err != nil {
		log.Fatal(err)
	}
}

func fetchNews(n newsTemplateVars, r *http.Request) newsTemplateVars {
	log.Println("fetchNews")
	c := appengine.NewContext(r)
	n.Rss = getNewsFromMemcache(c)
	if len(n.Rss.Channel.Items) == 0 {
		n = fetchNewsFromRSS(n, c)
		setNewsMemcache(n.Rss, c)
	}
	return n
}

func NewsHandler(w http.ResponseWriter, r *http.Request) {
	var n newsTemplateVars
	n = fetchNews(n, r)
	if err := templates.ExecuteTemplate(w, "news", n); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
