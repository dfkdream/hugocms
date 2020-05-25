package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"golang.org/x/crypto/acme/autocert"

	"github.com/dfkdream/hugocms/pluginapi"

	"github.com/dfkdream/hugocms/adminapi"

	"github.com/dfkdream/hugocms/admin"
	"github.com/dfkdream/hugocms/hugo"

	"github.com/dfkdream/hugocms/config"
	"github.com/dfkdream/hugocms/session"
	"github.com/dfkdream/hugocms/signin"
	"github.com/dfkdream/hugocms/user"

	"github.com/boltdb/bolt"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	cfg := config.GetConfig()

	fmt.Println(cfg)

	db, err := bolt.Open(cfg.BoltPath, os.FileMode(0644), nil)
	if err != nil {
		log.Fatal(err)
	}

	userBD := user.NewDB(db)

	t, err := template.New("html").ParseGlob("./html/*.html")
	if err != nil {
		log.Fatal(err)
	}

	hg := hugo.New(cfg)

	r := mux.NewRouter().StrictSlash(true)

	s := signin.NewSignInHandler(
		"/admin/signin",
		"/admin/assets/",
		"/admin/api/",
		session.NewDB(true, 10*time.Minute),
		userBD,
		t)

	rAdmin := r.PathPrefix("/admin").Subrouter().StrictSlash(true)

	admin.Admin{SignIn: s, T: t, Config: cfg}.SetupHandlers(rAdmin)

	adminapi.AdminAPI{Conf: cfg, Hugo: hg, UserDB: userBD}.SetupHandlers(rAdmin.PathPrefix("/api").Subrouter().StrictSlash(true))

	pluginapi.PluginAPI{Config: cfg, SignIn: s}.SetupHandlers(r.PathPrefix("/api").Subrouter().StrictSlash(true))

	r.PathPrefix("/").Handler(http.FileServer(http.Dir(cfg.PublicPath)))

	logged := handlers.LoggingHandler(os.Stdout, r)

	log.Println("HTTP Server started at", cfg.Bind)

	if cfg.TLS && cfg.AutoCert {
		if cfg.Domain == "" {
			log.Fatal("autocert: domain must not be empty string")
		}
		log.Fatal(http.Serve(autocert.NewListener(cfg.Domain), logged))
	} else if cfg.TLS {
		if err := http.ListenAndServeTLS(cfg.Bind, cfg.CertPath, cfg.KeyPath, logged); err != nil {
			log.Fatal(err)
		}
	} else {
		if err := http.ListenAndServe(cfg.Bind, logged); err != nil {
			log.Fatal(err)
		}
	}
}
