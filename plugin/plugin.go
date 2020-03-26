package plugin

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type contextKey string

var (
	// ContextKeyUser is context key for user data
	ContextKeyUser = contextKey("user")
)

// User contains user information
type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

// String converts user to json string
func (u User) String() string {
	if res, err := json.Marshal(u); err == nil {
		return string(res)
	} else {
		return ""
	}
}

// Info contains information about plugin which will be displayed on HugoCMS dashboard.
type Info struct {
	Name        string `json:"name"`
	Author      string `json:"author"`
	Description string `json:"description"`
	Version     string `json:"version"`
}

type adminMenuItem struct {
	MenuName string `json:"menuName"`
	Endpoint string `json:"endpoint"`
}

// Metadata contains metadata about plugin
type Metadata struct {
	Identifier     string          `json:"identifier"`
	Info           Info            `json:"info"`
	AdminMenuItems []adminMenuItem `json:"adminMenuItems"`
}

// Plugin is HugoCMS Plugin which implements http.Handler.
type Plugin struct {
	router         *mux.Router
	adminRouter    *mux.Router
	adminAPIRouter *mux.Router
	apiRouter      *mux.Router
	metadata       *Metadata
}

// New creates new plugin.
func New(Info Info, Identifier string) *Plugin {
	p := &Plugin{
		router: mux.NewRouter().StrictSlash(true),
		metadata: &Metadata{
			Identifier:     Identifier,
			Info:           Info,
			AdminMenuItems: make([]adminMenuItem, 0),
		},
	}
	p.adminRouter = p.router.PathPrefix("/admin").Subrouter().StrictSlash(true)
	p.adminAPIRouter = p.router.PathPrefix("/admin_api").Subrouter().StrictSlash(true)
	p.apiRouter = p.router.PathPrefix("/api").Subrouter().StrictSlash(true)

	p.router.HandleFunc("/metadata", func(res http.ResponseWriter, req *http.Request) {
		err := json.NewEncoder(res).Encode(p.metadata)
		if err != nil {
			http.Error(res, "Internal Server Error", http.StatusInternalServerError)
			log.Println(err)
			return
		}
	})

	p.router.HandleFunc("/live", func(res http.ResponseWriter, req *http.Request) {
	})
	return p
}

// HandleAdminPage handles admin page handlers.
// menuName will be displayed on navigation bar.
// Handler should write HTML document.
func (p *Plugin) HandleAdminPage(path, menuName string, handler http.Handler) {
	p.metadata.AdminMenuItems = append(p.metadata.AdminMenuItems, adminMenuItem{Endpoint: path, MenuName: menuName})
	p.adminRouter.Handle(path, handler)
}

// AdminPageRouter returns admin page router.
func (p *Plugin) AdminPageRouter() *mux.Router {
	return p.adminRouter
}

// HandleAdminAPI handles admin only API handlers.
// Non logged in requests will be rejected.
func (p *Plugin) HandleAdminAPI(path string, handler http.Handler) {
	p.adminAPIRouter.Handle(path, handler)
}

// AdminAPIRouter returns admin API router.
func (p *Plugin) AdminAPIRouter() *mux.Router {
	return p.adminAPIRouter
}

// HandleAPI handles API handlers.
// Non logged in users can access these APIs.
func (p *Plugin) HandleAPI(path string, handler http.Handler) {
	p.apiRouter.Handle(path, handler)
}

// APIRouter returns API router.
func (p *Plugin) APIRouter() *mux.Router {
	return p.adminAPIRouter
}

// ServeHTTP dispatches the requests to plugin.
func (p *Plugin) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	if h := req.Header.Get("X-HugoCMS-User"); h != "" {
		u := new(User)
		err := json.Unmarshal([]byte(h), &u)
		if err != nil {
			http.Error(res, "Bad Request", http.StatusBadRequest)
			log.Println(err)
			return
		}
		req = req.WithContext(context.WithValue(req.Context(), ContextKeyUser, u))
	}
	p.router.ServeHTTP(res, req)
}
