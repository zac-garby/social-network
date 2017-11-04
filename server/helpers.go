package server

import (
	"html/template"
	"io/ioutil"
	"log"
	"mime"
	"net/http"
	"path/filepath"
)

func handleRequest(w http.ResponseWriter, name, base, req string, data interface{}) {
	if len(req) <= len(base)+2 {
		t, err := template.New(name).ParseFiles(filepath.Join("public", base, "index.html"))
		if err != nil {
			log.Println(err)
			return
		}

		if err := t.ExecuteTemplate(w, name, data); err != nil {
			log.Println(err)
		}
	} else {
		t := mime.TypeByExtension(filepath.Ext(req))
		text, err := ioutil.ReadFile(filepath.Join("public", req))
		if err != nil {
			log.Println(err)
		}

		w.Header().Set("Content-Type", t)
		w.Write(text)
	}
}
