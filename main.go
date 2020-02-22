package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	cfg := getConfig()

	fmt.Println(cfg)

	t, err := template.New("html").ParseGlob("./html/*.html")
	if err != nil {
		log.Fatal(err)
	}

	r := mux.NewRouter().StrictSlash(true)

	signin := newSignInHandler(
		"/admin/signin/",
		"/admin/assets/",
		newSessionDB(true, 1*time.Minute),
		newUserDB(),
		t)

	admin := r.PathPrefix("/admin").Subrouter().StrictSlash(true)

	admin.Use(signin.middleware)

	admin.PathPrefix("/assets").Handler(
		http.StripPrefix("/admin/assets/", http.FileServer(http.Dir("./assets"))))

	admin.Handle("/signin/", signin)

	admin.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		res.Write([]byte(req.UserAgent() + "\n"))
		if u, ok := req.Context().Value(contextKeyUser).(*user); ok {
			res.Write([]byte("Hello, " + u.username))
		}
	})

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
