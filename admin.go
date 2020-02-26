package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type admin struct {
	signIn *signInHandler
	t      *template.Template
}

func (a admin) setupHandlers(router *mux.Router) {
	router.Use(a.signIn.middleware(true))

	router.PathPrefix("/assets").Handler(
		http.StripPrefix("/admin/assets/", http.FileServer(http.Dir("./assets"))))

	router.Handle("/signin/", a.signIn)

	router.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		http.Redirect(res, req, "/admin/list/", http.StatusFound)
	})

	router.HandleFunc("/list/", func(res http.ResponseWriter, req *http.Request) {
		err := a.t.ExecuteTemplate(res, "list.html", nil)
		if err != nil {
			log.Println(err)
			http.Redirect(res, req, "/admin/list/", http.StatusFound)
		}
	})
}
