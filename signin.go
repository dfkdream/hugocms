package main

import (
	"context"
	"crypto/rand"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

type contextKey string

var (
	contextKeyUser = contextKey("user")
)

func generateRandomKey(bytes int) string {
	buff := make([]byte, bytes)
	_, _ = rand.Read(buff)
	return fmt.Sprintf("%x", buff)
}

func mustReadCookie(key string, req *http.Request) string {
	if ck, err := req.Cookie(key); err == nil {
		return ck.Value
	}
	return ""
}

func readIP(req *http.Request) string {
	return strings.Split(req.RemoteAddr, ":")[0]
}

type signInHandler struct {
	signInURL string
	assetsURL string
	apiURL    string
	sessionDB *sessionDB
	userDB    *userDB
	template  *template.Template
}

func newSignInHandler(signInURL, assetsURL, apiURL string, sessionDB *sessionDB, userDB *userDB, template *template.Template) *signInHandler {
	return &signInHandler{
		signInURL: signInURL,
		assetsURL: assetsURL,
		apiURL:    apiURL,
		sessionDB: sessionDB,
		userDB:    userDB,
		template:  template,
	}
}

func (s signInHandler) middleware(blocking bool) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {

			if req.URL.Path == s.signInURL || strings.HasPrefix(req.URL.Path, s.assetsURL) {
				next.ServeHTTP(res, req)
				return
			}

			if ok, user := s.sessionDB.validate(mustReadCookie("sess", req), readIP(req)); ok {
				req = req.WithContext(context.WithValue(req.Context(), contextKeyUser, user))
				next.ServeHTTP(res, req)
				return
			}

			if strings.HasPrefix(req.URL.Path, s.apiURL) {
				http.Error(res, jsonStatusForbidden, http.StatusForbidden)
				return
			}

			if blocking {
				http.Redirect(res, req, s.signInURL+"?redirect="+req.URL.EscapedPath(), http.StatusFound)
			} else {
				next.ServeHTTP(res, req)
			}
		})
	}
}

func (s signInHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET": // Sign in page
		redirect := req.URL.Query().Get("redirect")
		if redirect == "" || redirect == s.signInURL {
			redirect = "/admin/"
		}
		err := s.template.ExecuteTemplate(res, "signin.html", redirect)
		if err != nil {
			http.Error(res, "Internal Server Error", http.StatusInternalServerError)
			log.Println(err)
			return
		}
	case "POST": // Sign in POST
		if req.ParseForm() != nil {
			http.Redirect(res, req, s.signInURL+"?redirect=/admin/", http.StatusFound)
			return
		}

		id := req.Form.Get("id")
		password := req.Form.Get("password")
		redirect := req.Form.Get("redirect")

		if id == "" || password == "" {
			http.Redirect(res, req, s.signInURL+"?redirect="+redirect, http.StatusFound)
			return
		}

		if s.userDB.size() == 0 { // Create admin user if DB is empty
			u, err := newUser(id, "admin", password)
			if err != nil {
				log.Println(err)
				http.Redirect(res, req, s.signInURL+"?redirect="+redirect, http.StatusFound)
				return
			}
			err = s.userDB.addUser(u)
			if err != nil {
				log.Println(err)
				http.Redirect(res, req, s.signInURL+"?redirect="+redirect, http.StatusFound)
				return
			}
		}

		if user := s.userDB.getUser(id); user != nil {
			if user.validate(id, password) {
				token := s.sessionDB.register(user, readIP(req))
				http.SetCookie(res, &http.Cookie{
					Name:     "sess",
					Value:    token,
					Path:     "/",
					HttpOnly: true,
				})
				http.Redirect(res, req, redirect, http.StatusFound)
				return
			}
		}

		// Delete cookie
		http.SetCookie(res, &http.Cookie{
			Name:     "sess",
			Value:    "",
			Expires:  time.Unix(0, 0),
			Path:     "/",
			HttpOnly: true,
		})
		http.Redirect(res, req, s.signInURL+"?redirect="+redirect, http.StatusFound)
	}
}
