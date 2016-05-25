package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/urlfetch"
)

type gistFile struct {
	Content string `json:"content"`
}

type gist struct {
	Description string              `json:"description"`
	Public      bool                `json:"public"`
	Files       map[string]gistFile `json:"files"`
}

type snippetsTemplateVars struct {
	FormErrors  []string
	FormSuccess []string
	Snippets    []gist
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
			gist := gist{
				Description: r.PostFormValue("desc"),
				Public:      public,
				Files:       map[string]gistFile{},
			}
			gist.Files[r.PostFormValue("filename")] = gistFile{
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
