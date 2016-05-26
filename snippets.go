package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"golang.org/x/net/context"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/urlfetch"
)

type GistFile struct {
	Content string `json:"content"`
}

type Gist struct {
	Id          int64               `json:"id" datastore:"-"`
	Description string              `json:"description"`
	Public      bool                `json:"public"`
	Files       map[string]GistFile `json:"files"`
}

type snippetsTemplateVars struct {
	FormErrors  []string
	FormSuccess []string
	Snippets    []Gist
}

func (gist *Gist) key(c context.Context) *datastore.Key {
	if gist.Id == 0 {
		return datastore.NewIncompleteKey(c, "Gist", nil)
	}
	return datastore.NewKey(c, "Gist", "", gist.Id, nil)
}

func (gist *Gist) save(c context.Context) error {
	k, err := datastore.Put(c, gist.key(c), gist)
	if err != nil {
		return err
	}
	gist.Id = k.IntID()
	return nil
}

func GetGists(c context.Context) ([]Gist, error) {
	q := datastore.NewQuery("Gist").Order("Description")

	var gists []Gist
	k, err := q.GetAll(c, &gists)
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(gists); i++ {
		gists[i].Id = k[i].IntID()
	}

	return gists, nil
}

func SnippetsHandler(w http.ResponseWriter, r *http.Request) {
	templateVars := snippetsTemplateVars{}
	c := appengine.NewContext(r)
	if r.Method == "POST" {
		r.ParseForm()
		var (
			responseObj map[string]interface{}
			public      bool = true
		)
		log.Println("snippet handler triggered")
		// validation
		if r.PostFormValue("filename") == "" || r.PostFormValue("snippet") == "" {
			templateVars.FormErrors = append(templateVars.FormErrors, "Need snippet & filename")
		} else {
			if r.PostFormValue("where") == "2" {
				public = false
			}
			gist := Gist{
				Description: r.PostFormValue("desc"),
				Public:      public,
				Files:       map[string]GistFile{},
			}
			gist.Files[r.PostFormValue("filename")] = GistFile{
				Content: r.PostFormValue("snippet"),
			}
			b, err := json.Marshal(gist)
			if err != nil {
				log.Fatal("JSON Error: ", err)
			}
			req, err := http.NewRequest("POST", "https://api.github.com/gists", bytes.NewBuffer(b))
			req.Header.Set("Content-Type", "application/json")
			client := urlfetch.Client(c)
			resp, err := client.Do(req)
			if err != nil {
				log.Fatal("HTTP Error: ", err)
			}
			defer resp.Body.Close()
			for k, v := range resp.Header {
				log.Println("key:", k, "value:", v)
			}
			err = json.NewDecoder(resp.Body).Decode(&responseObj)
			if err != nil {
				log.Fatal("Response JSON Error: ", err)
			}
			// @todo: write to datastore
			success := "Successfully created: <a target=\"_blank\" href=\"" + responseObj["html_url"].(string) + "\">" + responseObj["html_url"].(string) + "</a>"
			templateVars.FormSuccess = append(templateVars.FormSuccess, success)
		}
	}
	// @todo grab snippets from datastore
	if err := templates.ExecuteTemplate(w, "snippets", templateVars); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
