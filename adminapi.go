package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gorilla/mux"
)

type adminAPI struct {
	conf *config
}

func (a adminAPI) postAPI(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		f, err := os.Open(filepath.Join(a.conf.ContentPath, req.URL.Path))
		if err != nil {
			log.Println(err)
			http.Error(res, jsonStatusNotFound, http.StatusNotFound)
			return
		}
		a, err := parseArticle(f)
		if err != nil {
			log.Println(err)
			http.Error(res, jsonStatusNotFound, http.StatusInternalServerError)
			return
		}
		err = json.NewEncoder(res).Encode(a)
		if err != nil {
			log.Println(err)
			http.Error(res, jsonStatusInternalServerError, http.StatusInternalServerError)
			return
		}
	case "POST":
		var articleJSON article
		err := json.NewDecoder(req.Body).Decode(&articleJSON)
		if err != nil {
			log.Println(err)
			http.Error(res, jsonStatusBadRequest, http.StatusBadRequest)
			return
		}
		f, err := os.OpenFile(filepath.Join(a.conf.ContentPath, req.URL.Path), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.FileMode(666))
		if err != nil {
			log.Println(err)
			http.Error(res, jsonStatusInternalServerError, http.StatusInternalServerError)
			return
		}
		jsonEnc := json.NewEncoder(f)
		jsonEnc.SetIndent("", "    ")
		err = jsonEnc.Encode(articleJSON.FrontMatter)
		if err != nil {
			log.Println(err)
			http.Error(res, jsonStatusInternalServerError, http.StatusInternalServerError)
			return
		}
		_, err = f.Write([]byte(articleJSON.Body))
		if err != nil {
			log.Println(err)
			http.Error(res, jsonStatusInternalServerError, http.StatusInternalServerError)
			return
		}
	default:
		http.Error(res, jsonStatusMethodNotAllowed, http.StatusMethodNotAllowed)
	}
}

type fileInfo struct {
	Name    string    `json:"name"`
	Size    int64     `json:"size"`
	Mode    string    `json:"mode"`
	ModTime time.Time `json:"modTime"`
	IsDir   bool      `json:"isDir"`
}

func (a adminAPI) listAPI(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		files, err := ioutil.ReadDir(filepath.Join(a.conf.ContentPath, req.URL.Path))
		if err != nil {
			log.Println(err)
			http.Error(res, jsonStatusNotFound, http.StatusNotFound)
			return
		}
		fJSON := make([]fileInfo, len(files))

		for idx, f := range files {
			fJSON[idx].Name = f.Name()
			fJSON[idx].Size = f.Size()
			fJSON[idx].Mode = f.Mode().String()
			fJSON[idx].ModTime = f.ModTime()
			fJSON[idx].IsDir = f.IsDir()
		}

		err = json.NewEncoder(res).Encode(fJSON)
		if err != nil {
			log.Println(err)
			http.Error(res, jsonStatusInternalServerError, http.StatusInternalServerError)
			return
		}
	case "POST":
		err := os.MkdirAll(filepath.Join(a.conf.ContentPath, req.URL.Path), os.FileMode(0755))
		if err != nil {
			log.Println(err)
			http.Error(res, jsonStatusInternalServerError, http.StatusInternalServerError)
			return
		}
	case "PUT":
	default:
		http.Error(res, jsonStatusMethodNotAllowed, http.StatusMethodNotAllowed)
	}
}

func (a adminAPI) setupHandlers(router *mux.Router) {
	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			res.Header().Set("Content-Type", "application/json")
			next.ServeHTTP(res, req)
		})
	})
	router.PathPrefix("/post").Handler(http.StripPrefix("/admin/api/post", http.HandlerFunc(a.postAPI)))
	router.PathPrefix("/list").Handler(http.StripPrefix("/admin/api/list", http.HandlerFunc(a.listAPI)))
}
