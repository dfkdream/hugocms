package signin

import (
	"context"
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/dfkdream/hugocms/internal"

	"github.com/dfkdream/hugocms/session"
	"github.com/dfkdream/hugocms/user"

	"github.com/gorilla/mux"
)

type contextKey string

var (
	ContextKeyUser = contextKey("user")
)

func mustReadCookie(key string, req *http.Request) string {
	if ck, err := req.Cookie(key); err == nil {
		return ck.Value
	}
	return ""
}

func readIP(req *http.Request) string {
	return strings.Split(req.RemoteAddr, ":")[0]
}

type SignInHandler struct {
	signInURL string
	assetsURL string
	apiURL    string
	sessionDB *session.DB
	userDB    *user.DB
	template  *template.Template
}

func NewSignInHandler(signInURL, assetsURL, apiURL string, sessionDB *session.DB, userDB *user.DB, template *template.Template) *SignInHandler {
	return &SignInHandler{
		signInURL: signInURL,
		assetsURL: assetsURL,
		apiURL:    apiURL,
		sessionDB: sessionDB,
		userDB:    userDB,
		template:  template,
	}
}

func (s SignInHandler) Middleware(blocking bool) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {

			if req.URL.Path == s.signInURL || strings.HasPrefix(req.URL.Path, s.assetsURL) {
				next.ServeHTTP(res, req)
				return
			}

			if ok, u := s.sessionDB.Validate(mustReadCookie("sess", req), readIP(req)); ok {
				req = req.WithContext(context.WithValue(req.Context(), ContextKeyUser, u))
				next.ServeHTTP(res, req)
				return
			}

			if strings.HasPrefix(req.URL.Path, s.apiURL) {
				http.Error(res, internal.JsonStatusForbidden, http.StatusForbidden)
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

func (s SignInHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
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

		if s.userDB.Size() == 0 { // Create admin user if DB is empty
			u, err := user.New(id, "admin", password)
			if err != nil {
				log.Println(err)
				http.Redirect(res, req, s.signInURL+"?redirect="+redirect, http.StatusFound)
				return
			}
			err = s.userDB.AddUser(u)
			if err != nil {
				log.Println(err)
				http.Redirect(res, req, s.signInURL+"?redirect="+redirect, http.StatusFound)
				return
			}
		}

		if u := s.userDB.GetUser(id); u != nil {
			if u.Validate(id, password) {
				token := s.sessionDB.Register(u, readIP(req))
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

func GetUser(req *http.Request) *user.User {
	if u, ok := req.Context().Value(ContextKeyUser).(*user.User); ok {
		return u
	}
	return nil
}
