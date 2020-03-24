package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/boltdb/bolt"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	cfg := getConfig()

	fmt.Println(cfg)

	db, err := bolt.Open(cfg.BoltPath, os.FileMode(0644), nil)
	if err != nil {
		log.Fatal(err)
	}

	t, err := template.New("html").ParseGlob("./html/*.html")
	if err != nil {
		log.Fatal(err)
	}

	hg := newHugo(cfg)

	r := mux.NewRouter().StrictSlash(true)

	signin := newSignInHandler(
		"/admin/signin",
		"/admin/assets/",
		"/admin/api/",
		newSessionDB(true, 10*time.Minute),
		newUserDB(db),
		t)

	rAdmin := r.PathPrefix("/admin").Subrouter().StrictSlash(true)

	admin{signIn: signin, t: t}.setupHandlers(rAdmin)

	adminAPI{conf: cfg, hugo: hg}.setupHandlers(rAdmin.PathPrefix("/api").Subrouter().StrictSlash(true))

	pluginAPI{config: cfg, signIn: signin}.setupHandlers(r.PathPrefix("/api").Subrouter().StrictSlash(true))

	r.PathPrefix("/").Handler(http.FileServer(http.Dir(cfg.PublicPath)))

	logged := handlers.LoggingHandler(os.Stdout, r)

	log.Println("HTTP Server started at", cfg.Bind)

	if cfg.TLS {
		if err := http.ListenAndServeTLS(cfg.Bind, cfg.CertPath, cfg.KeyPath, logged); err != nil {
			log.Fatal(err)
		}
	} else {
		if err := http.ListenAndServe(cfg.Bind, logged); err != nil {
			log.Fatal(err)
		}
	}
}
