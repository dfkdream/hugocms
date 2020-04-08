package pluginapi

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/dfkdream/hugocms/config"

	"github.com/dfkdream/hugocms/internal"
	"github.com/dfkdream/hugocms/signin"
	"github.com/dfkdream/hugocms/user"

	"github.com/gorilla/mux"

	"github.com/dfkdream/hugocms/plugin"
)

func NewAuthenticatedReverseProxy(path string) *httputil.ReverseProxy {
	return &httputil.ReverseProxy{Director: func(req *http.Request) {
		target, err := url.Parse(path)
		if err != nil {
			log.Fatal(err)
		}
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.URL.Path = internal.SingleJoiningSlash(target.Path, req.URL.Path)

		if target.RawQuery == "" || req.URL.RawPath == "" {
			req.URL.RawQuery = target.RawQuery + req.URL.RawQuery
		} else {
			req.URL.RawQuery = target.RawQuery + "&" + req.URL.RawQuery
		}

		if _, ok := req.Header["User-Agent"]; !ok {
			req.Header.Set("User-Agent", "")
		}

		req.Header.Del("X-HugoCMS-User")

		if u, ok := req.Context().Value(signin.ContextKeyUser).(*user.User); ok {
			req.Header.Set("X-HugoCMS-User", plugin.User{ID: u.Id, Username: u.Username}.String())
		}
	}}
}

type PluginAPI struct {
	Config *config.Config
	SignIn *signin.SignInHandler
}

func (p PluginAPI) SetupHandlers(router *mux.Router) {
	router.Use(p.SignIn.Middleware(false))
	for _, v := range p.Config.Plugins {
		router.PathPrefix("/" + v.Metadata.Identifier).Handler(
			http.StripPrefix("/api/"+v.Metadata.Identifier,
				NewAuthenticatedReverseProxy(internal.SingleJoiningSlash(v.Addr, "/api"))))
	}
}
