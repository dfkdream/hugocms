package admin

import (
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"path"
	"time"

	"github.com/dfkdream/hugocms/internal"

	"github.com/dfkdream/hugocms/config"

	"github.com/dfkdream/hugocms/signin"

	"github.com/dfkdream/hugocms/user"

	"github.com/gorilla/mux"
)

type Admin struct {
	SignIn *signin.SignInHandler
	T      *template.Template
	Config *config.Config
}

type templateVars struct {
	Title   string
	Plugins []config.PluginData
	Body    template.HTML
	User    *user.User
}

func (a Admin) SetupHandlers(router *mux.Router) {
	router.Use(a.SignIn.Middleware(true))

	router.PathPrefix("/assets").Handler(
		http.StripPrefix("/admin/assets/", http.FileServer(http.Dir("./assets"))))

	router.Handle("/signin", a.SignIn)

	router.HandleFunc("/signout", a.SignIn.SignOut)

	router.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		http.Redirect(res, req, "/admin/list/", http.StatusFound)
	})

	router.PathPrefix("/list/").HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		err := a.T.ExecuteTemplate(res, "list.html", templateVars{Plugins: a.Config.Plugins, User: signin.GetUser(req)})
		if err != nil {
			log.Println(err)
			http.Error(res, "Internal Server Error", http.StatusInternalServerError)
		}
	})

	router.PathPrefix("/edit").HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		err := a.T.ExecuteTemplate(res, "edit.html", templateVars{Plugins: a.Config.Plugins, User: signin.GetUser(req)})
		if err != nil {
			log.Println(err)
			http.Error(res, "Internal Server Error", http.StatusInternalServerError)
		}
	})

	router.HandleFunc("/config", func(res http.ResponseWriter, req *http.Request) {
		err := a.T.ExecuteTemplate(res, "config.html", templateVars{Plugins: a.Config.Plugins, User: signin.GetUser(req)})
		if err != nil {
			log.Println(err)
			http.Error(res, "Internal Server Error", http.StatusInternalServerError)
		}
	})

	router.HandleFunc("/plugins", func(res http.ResponseWriter, req *http.Request) {
		err := a.T.ExecuteTemplate(res, "plugins.html", templateVars{Plugins: a.Config.Plugins, User: signin.GetUser(req)})
		if err != nil {
			log.Println(err)
			http.Error(res, "Internal Server Error", http.StatusInternalServerError)
		}
	})

	for _, v := range a.Config.Plugins {
		router.PathPrefix("/" + v.Metadata.Identifier).Handler(
			http.StripPrefix("/admin/"+v.Metadata.Identifier, http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
				r, err := http.NewRequest("GET", internal.SingleJoiningSlash(v.Addr, path.Join("/admin", req.URL.Path)), nil)
				if err != nil {
					log.Println(err)
					http.Error(res, "Internal Server Error", http.StatusInternalServerError)
					return
				}

				if u, ok := req.Context().Value(signin.ContextKeyUser).(*user.User); ok {
					r.Header.Set("X-HugoCMS-User", u.String())
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
				err = a.T.ExecuteTemplate(res, "plugin.html", templateVars{Plugins: a.Config.Plugins, Title: v.Metadata.Info.Name, Body: template.HTML(body), User: signin.GetUser(req)})
				if err != nil {
					log.Println(err)
					http.Error(res, "Internal Server Error", http.StatusInternalServerError)
					return
				}
			})))
	}

}
