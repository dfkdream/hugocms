package main

import (
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/dfkdream/hugocms/plugin"

	"github.com/gorilla/mux"
)

type admin struct {
	signIn *signInHandler
	t      *template.Template
	config *config
}

type templateVars struct {
	Title   string
	Plugins []pluginData
	Body    template.HTML
}

func (a admin) setupHandlers(router *mux.Router) {
	router.Use(a.signIn.middleware(true))

	router.PathPrefix("/assets").Handler(
		http.StripPrefix("/admin/assets/", http.FileServer(http.Dir("./assets"))))

	router.Handle("/signin", a.signIn)

	router.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		http.Redirect(res, req, "/admin/list/", http.StatusFound)
	})

	router.PathPrefix("/list/").HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		err := a.t.ExecuteTemplate(res, "list.html", templateVars{Plugins: a.config.Plugins})
		if err != nil {
			log.Println(err)
			http.Error(res, "Internal Server Error", http.StatusInternalServerError)
		}
	})

	router.PathPrefix("/edit").HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		err := a.t.ExecuteTemplate(res, "edit.html", templateVars{Plugins: a.config.Plugins})
		if err != nil {
			log.Println(err)
			http.Error(res, "Internal Server Error", http.StatusInternalServerError)
		}
	})

	router.HandleFunc("/config", func(res http.ResponseWriter, req *http.Request) {
		err := a.t.ExecuteTemplate(res, "config.html", templateVars{Plugins: a.config.Plugins})
		if err != nil {
			log.Println(err)
			http.Error(res, "Internal Server Error", http.StatusInternalServerError)
		}
	})

	router.HandleFunc("/plugins", func(res http.ResponseWriter, req *http.Request) {
		err := a.t.ExecuteTemplate(res, "plugins.html", templateVars{Plugins: a.config.Plugins})
		if err != nil {
			log.Println(err)
			http.Error(res, "Internal Server Error", http.StatusInternalServerError)
		}
	})

	for _, v := range a.config.Plugins {
		for _, e := range v.Metadata.AdminEndpoints {
			router.HandleFunc(e.Endpoint, func(res http.ResponseWriter, req *http.Request) {
				r, err := http.NewRequest("GET", singleJoiningSlash(v.Addr, e.Endpoint), nil)
				if err != nil {
					log.Println(err)
					http.Error(res, "Internal Server Error", http.StatusInternalServerError)
					return
				}

				if u, ok := req.Context().Value(contextKeyUser).(*user); ok {
					r.Header.Set("X-HugoCMS-User", plugin.User{ID: u.ID, Username: u.Username}.String())
				}

				resp, err := (&http.Client{Timeout: 10 * time.Second}).Do(r)
				if err != nil {
					log.Println(err)
					http.Error(res, "Bad Gateway", http.StatusBadGateway)
					return
				}
				defer func() { _ = resp.Body.Close() }()
				body, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					log.Println(err)
					http.Error(res, "Internal Server Error", http.StatusInternalServerError)
					return
				}
				res.WriteHeader(resp.StatusCode)
				err = a.t.ExecuteTemplate(res, "plugin.html", templateVars{Plugins: a.config.Plugins, Title: v.Metadata.Info.Name, Body: template.HTML(body)})
				if err != nil {
					log.Println(err)
					http.Error(res, "Internal Server Error", http.StatusInternalServerError)
					return
				}
			})
		}
	}

}
