package adminapi

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/dfkdream/hugocms/pluginapi"

	"github.com/dfkdream/hugocms/article"
	"github.com/dfkdream/hugocms/signin"
	"github.com/dfkdream/hugocms/user"

	"github.com/dfkdream/hugocms/config"
	"github.com/dfkdream/hugocms/hugo"
	"github.com/dfkdream/hugocms/internal"

	"github.com/dfkdream/hugocms/plugin"

	"github.com/gorilla/mux"
)

type AdminAPI struct {
	Conf   *config.Config
	Hugo   *hugo.Hugo
	UserDB *user.DB
}

func (a AdminAPI) postAPI(res http.ResponseWriter, req *http.Request) {
	if !signin.GetUser(req).HasPermission("hugocms:post") {
		http.Error(res, internal.JsonStatusForbidden, http.StatusForbidden)
		return
	}

	switch req.Method {
	case "GET":
		f, err := os.Open(filepath.Join(a.Conf.ContentPath, filepath.Clean("/"+req.URL.Path)))
		if err != nil {
			log.Println(err)
			http.Error(res, internal.JsonStatusNotFound, http.StatusNotFound)
			return
		}
		defer func() { _ = f.Close() }()
		a, err := article.Parse(f)
		if err != nil {
			log.Println(err)
			http.Error(res, internal.JsonStatusNotFound, http.StatusInternalServerError)
			return
		}
		err = json.NewEncoder(res).Encode(a)
		if err != nil {
			log.Println(err)
			http.Error(res, internal.JsonStatusInternalServerError, http.StatusInternalServerError)
			return
		}
	case "POST":
		var articleJSON article.Article
		err := json.NewDecoder(req.Body).Decode(&articleJSON)
		if err != nil {
			log.Println(err)
			http.Error(res, internal.JsonStatusBadRequest, http.StatusBadRequest)
			return
		}
		f, err := os.OpenFile(filepath.Join(a.Conf.ContentPath, filepath.Clean("/"+req.URL.Path)), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.FileMode(0644))
		if err != nil {
			log.Println(err)
			http.Error(res, internal.JsonStatusInternalServerError, http.StatusInternalServerError)
			return
		}
		defer func() { _ = f.Close() }()
		jsonEnc := json.NewEncoder(f)
		jsonEnc.SetIndent("", "    ")
		err = jsonEnc.Encode(articleJSON.FrontMatter)
		if err != nil {
			log.Println(err)
			http.Error(res, internal.JsonStatusInternalServerError, http.StatusInternalServerError)
			return
		}
		_, err = f.Write([]byte(articleJSON.Body))
		if err != nil {
			log.Println(err)
			http.Error(res, internal.JsonStatusInternalServerError, http.StatusInternalServerError)
			return
		}
	default:
		http.Error(res, internal.JsonStatusMethodNotAllowed, http.StatusMethodNotAllowed)
	}
}

type fileInfo struct {
	Name    string    `json:"name"`
	Size    int64     `json:"size"`
	Mode    string    `json:"mode"`
	ModTime time.Time `json:"modTime"`
	IsDir   bool      `json:"isDir"`
}

func (a AdminAPI) listAPI(res http.ResponseWriter, req *http.Request) {
	if !signin.GetUser(req).HasPermission("hugocms:list") {
		http.Error(res, internal.JsonStatusForbidden, http.StatusForbidden)
		return
	}

	switch req.Method {
	case "GET":
		files, err := ioutil.ReadDir(filepath.Join(a.Conf.ContentPath, filepath.Clean("/"+req.URL.Path)))
		if err != nil {
			log.Println(err)
			http.Error(res, internal.JsonStatusNotFound, http.StatusNotFound)
			return
		}
		fJSON := make([]fileInfo, len(files))

		for idx, f := range files {
			fJSON[idx].Name = f.Name()
			fJSON[idx].Size = f.Size()
			fJSON[idx].Mode = f.Mode().String()
			fJSON[idx].ModTime = f.ModTime()
			fJSON[idx].IsDir = f.IsDir()
		}

		err = json.NewEncoder(res).Encode(fJSON)
		if err != nil {
			log.Println(err)
			http.Error(res, internal.JsonStatusInternalServerError, http.StatusInternalServerError)
			return
		}
	case "POST":
		err := os.MkdirAll(filepath.Join(a.Conf.ContentPath, filepath.Clean("/"+req.URL.Path)), os.FileMode(0755))
		if err != nil {
			log.Println(err)
			http.Error(res, internal.JsonStatusInternalServerError, http.StatusInternalServerError)
			return
		}
	case "PUT":
		var path string
		err := json.NewDecoder(req.Body).Decode(&path)
		if err != nil {
			log.Println(err)
			http.Error(res, internal.JsonStatusBadRequest, http.StatusBadRequest)
			return
		}
		err = os.Rename(filepath.Join(a.Conf.ContentPath, filepath.Clean("/"+req.URL.Path)), filepath.Join(a.Conf.ContentPath, filepath.Clean(path)))
		if err != nil {
			log.Println(err)
			http.Error(res, internal.JsonStatusInternalServerError, http.StatusInternalServerError)
		}
	case "DELETE":
		err := os.RemoveAll(filepath.Join(a.Conf.ContentPath, filepath.Clean("/"+req.URL.Path)))
		if err != nil {
			log.Println(err)
			http.Error(res, internal.JsonStatusInternalServerError, http.StatusInternalServerError)
		}
	default:
		http.Error(res, internal.JsonStatusMethodNotAllowed, http.StatusMethodNotAllowed)
	}
}

func (a AdminAPI) blobAPI(res http.ResponseWriter, req *http.Request) {
	if !signin.GetUser(req).HasPermission("hugocms:blob") {
		http.Error(res, internal.JsonStatusForbidden, http.StatusForbidden)
		return
	}

	switch req.Method {
	case "GET":
		if strings.HasSuffix("/", req.URL.Path) {
			http.Error(res, "404 page not found", http.StatusNotFound)
			return
		}
		res.Header().Del("Content-Type")
		http.ServeFile(res, req, filepath.Join(a.Conf.ContentPath, filepath.Clean("/"+req.URL.Path)))
	case "POST":
		f, err := os.OpenFile(filepath.Join(a.Conf.ContentPath, filepath.Clean("/"+req.URL.Path)), os.O_WRONLY|os.O_TRUNC|os.O_CREATE, os.FileMode(0644))
		if err != nil {
			log.Println(err)
			http.Error(res, internal.JsonStatusInternalServerError, http.StatusInternalServerError)
			return
		}
		defer func() { _ = f.Close() }()

		_, err = io.Copy(f, req.Body)
		if err != nil {
			log.Println(err)
			http.Error(res, internal.JsonStatusInternalServerError, http.StatusInternalServerError)
		}
	case "PUT":
		var path string
		err := json.NewDecoder(req.Body).Decode(&path)
		if err != nil {
			log.Println(err)
			http.Error(res, internal.JsonStatusBadRequest, http.StatusBadRequest)
			return
		}
		err = os.Rename(filepath.Join(a.Conf.ContentPath, filepath.Clean("/"+req.URL.Path)), filepath.Join(a.Conf.ContentPath, filepath.Clean(path)))
		if err != nil {
			log.Println(err)
			http.Error(res, internal.JsonStatusInternalServerError, http.StatusInternalServerError)
		}
	case "DELETE":
		err := os.Remove(filepath.Join(a.Conf.ContentPath, filepath.Clean("/"+req.URL.Path)))
		if err != nil {
			log.Println(err)
			http.Error(res, internal.JsonStatusInternalServerError, http.StatusInternalServerError)
		}
	default:
		http.Error(res, internal.JsonStatusMethodNotAllowed, http.StatusMethodNotAllowed)
	}
}

func (a AdminAPI) whoamiAPI(res http.ResponseWriter, req *http.Request) {
	if !signin.GetUser(req).HasPermission("hugocms:whoami") {
		http.Error(res, internal.JsonStatusForbidden, http.StatusForbidden)
		return
	}

	switch req.Method {
	case "GET":
		if u, ok := req.Context().Value(signin.ContextKeyUser).(*user.User); ok {
			err := json.NewEncoder(res).Encode(
				struct {
					ID       string `json:"id"`
					Username string `json:"username"`
				}{
					ID:       u.Id,
					Username: u.Username,
				})
			if err != nil {
				log.Println(err)
				http.Error(res, internal.JsonStatusForbidden, http.StatusForbidden)
			}
			return
		}
		http.Error(res, internal.JsonStatusForbidden, http.StatusForbidden)
	case "POST":
		var values struct {
			Username        string `json:"username"`
			CurrentPassword string `json:"currentPassword"`
			NewPassword     string `json:"newPassword"`
		}

		err := json.NewDecoder(req.Body).Decode(&values)
		if err != nil {
			http.Error(res, internal.JsonStatusBadRequest, http.StatusBadRequest)
			return
		}

		u := signin.GetUser(req)
		if u == nil {
			http.Error(res, internal.JsonStatusForbidden, http.StatusForbidden)
			return
		}

		if values.Username != "" {
			u.Username = values.Username
		}

		if values.CurrentPassword != "" && values.NewPassword != "" {
			if !u.Validate(u.Id, values.CurrentPassword) {
				http.Error(res, internal.JsonStatusForbidden, http.StatusForbidden)
				return
			}
			u, err = user.New(u.Id, u.Username, values.NewPassword, u.Permissions)
			if err != nil {
				log.Println(err)
				http.Error(res, internal.JsonStatusInternalServerError, http.StatusInternalServerError)
				return
			}
		}

		a.UserDB.SetUser(u)
	default:
		http.Error(res, internal.JsonStatusMethodNotAllowed, http.StatusMethodNotAllowed)
	}
}

func (a AdminAPI) buildAPI(res http.ResponseWriter, req *http.Request) {
	if !signin.GetUser(req).HasPermission("hugocms:build") {
		http.Error(res, internal.JsonStatusForbidden, http.StatusForbidden)
		return
	}

	switch req.Method {
	case "POST":
		r := a.Hugo.Build()
		if r.Err != nil {
			log.Println(r.Err)
		}
		fmt.Println(r.Result)
		err := json.NewEncoder(res).Encode(r)
		if err != nil {
			http.Error(res, internal.JsonStatusInternalServerError, http.StatusInternalServerError)
		}
	default:
		http.Error(res, internal.JsonStatusMethodNotAllowed, http.StatusMethodNotAllowed)
	}
}

func (a AdminAPI) configAPI(res http.ResponseWriter, req *http.Request) {
	if !signin.GetUser(req).HasPermission("hugocms:config") {
		http.Error(res, internal.JsonStatusForbidden, http.StatusForbidden)
		return
	}

	switch req.Method {
	case "GET":
		res.Header().Del("Content-Type")
		http.ServeFile(res, req, a.Conf.ConfigPath)
	case "POST":
		f, err := os.OpenFile(a.Conf.ConfigPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.FileMode(0644))
		if err != nil {
			log.Println(err)
			http.Error(res, internal.JsonStatusInternalServerError, http.StatusInternalServerError)
			return
		}
		defer func() { _ = f.Close() }()
		_, err = io.Copy(f, req.Body)
		if err != nil {
			log.Println(err)
			http.Error(res, internal.JsonStatusInternalServerError, http.StatusInternalServerError)
			return
		}
	default:
		http.Error(res, internal.JsonStatusMethodNotAllowed, http.StatusMethodNotAllowed)
	}
}

type pluginInfo struct {
	Info   plugin.Info `json:"info"`
	IsLive bool        `json:"isLive"`
}

func (a AdminAPI) pluginsAPI(res http.ResponseWriter, req *http.Request) {
	if !signin.GetUser(req).HasPermission("hugocms:plugins") {
		http.Error(res, internal.JsonStatusForbidden, http.StatusForbidden)
		return
	}

	switch req.Method {
	case "GET":
		i := make([]pluginInfo, len(a.Conf.Plugins))
		for idx, v := range a.Conf.Plugins {
			i[idx] = pluginInfo{*v.Metadata.Info, config.CheckPluginLive(v.Addr)}
		}
		err := json.NewEncoder(res).Encode(i)
		if err != nil {
			log.Println(err)
			http.Error(res, internal.JsonStatusInternalServerError, http.StatusInternalServerError)
		}
	default:
		http.Error(res, internal.JsonStatusMethodNotAllowed, http.StatusMethodNotAllowed)
	}
}

func (a AdminAPI) usersAPI(res http.ResponseWriter, req *http.Request) {
	if !signin.GetUser(req).HasPermission("hugocms:user") {
		http.Error(res, internal.JsonStatusForbidden, http.StatusForbidden)
		return
	}

	switch req.Method {
	case "GET":
		u := a.UserDB.GetAllUsers()
		for idx := range u {
			u[idx].Hash = ""
			u[idx].Salt = ""
		}

		err := json.NewEncoder(res).Encode(u)
		if err != nil {
			log.Println(err)
			http.Error(res, internal.JsonStatusInternalServerError, http.StatusInternalServerError)
		}
	case "POST":
		var u struct {
			Id          string   `json:"id"`
			Username    string   `json:"username"`
			Password    string   `json:"password"`
			Permissions []string `json:"permissions"`
		}

		err := json.NewDecoder(req.Body).Decode(&u)
		if err != nil {
			http.Error(res, internal.JsonStatusBadRequest, http.StatusBadRequest)
			return
		}

		usr, err := user.New(u.Id, u.Username, u.Password, u.Permissions)
		if err != nil {
			log.Println(err)
			http.Error(res, internal.JsonStatusInternalServerError, http.StatusInternalServerError)
			return
		}

		err = a.UserDB.AddUser(usr)
		if err != nil {
			if err == user.ErrDuplicatedUser {
				http.Error(res, internal.JsonStatusConflict, http.StatusConflict)
				return
			}
			log.Println(err)
			http.Error(res, internal.JsonStatusInternalServerError, http.StatusInternalServerError)
			return
		}
	default:
		http.Error(res, internal.JsonStatusMethodNotAllowed, http.StatusMethodNotAllowed)
	}
}

func (a AdminAPI) userAPI(res http.ResponseWriter, req *http.Request) {
	if !signin.GetUser(req).HasPermission("hugocms:user") {
		http.Error(res, internal.JsonStatusForbidden, http.StatusForbidden)
		return
	}

	id := mux.Vars(req)["id"]
	u := a.UserDB.GetUser(id)
	if u == nil {
		http.Error(res, internal.JsonStatusNotFound, http.StatusNotFound)
		return
	}

	switch req.Method {
	case "GET":
		u.Hash = ""
		u.Salt = ""

		err := json.NewEncoder(res).Encode(u)
		if err != nil {
			log.Println(err)
			http.Error(res, internal.JsonStatusInternalServerError, http.StatusInternalServerError)
		}
	case "POST":
		var value struct {
			Username    string   `json:"username"`
			Password    string   `json:"password"`
			Permissions []string `json:"permissions"`
		}

		err := json.NewDecoder(req.Body).Decode(&value)
		if err != nil {
			http.Error(res, internal.JsonStatusBadRequest, http.StatusBadRequest)
			return
		}

		if value.Username != "" {
			u.Username = value.Username
		}

		if len(value.Permissions) > 0 {
			u.Permissions = value.Permissions
		}

		if value.Password != "" {
			u, err = user.New(u.Id, u.Username, value.Password, u.Permissions)
			if err != nil {
				log.Println(err)
				http.Error(res, internal.JsonStatusInternalServerError, http.StatusInternalServerError)
				return
			}
		}

		a.UserDB.SetUser(u)
	case "DELETE":
		err := a.UserDB.DeleteUser(u.Id)
		if err != nil {
			log.Println(err)
			http.Error(res, internal.JsonStatusInternalServerError, http.StatusInternalServerError)
		}
	default:
		http.Error(res, internal.JsonStatusMethodNotAllowed, http.StatusMethodNotAllowed)
	}
}

func (a AdminAPI) SetupHandlers(router *mux.Router) {
	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			res.Header().Set("Content-Type", "application/json")
			res.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
			next.ServeHTTP(res, req)
		})
	})

	router.PathPrefix("/post").Handler(http.StripPrefix("/admin/api/post", http.HandlerFunc(a.postAPI)))
	router.PathPrefix("/list").Handler(http.StripPrefix("/admin/api/list", http.HandlerFunc(a.listAPI)))
	router.PathPrefix("/blob").Handler(http.StripPrefix("/admin/api/blob", http.HandlerFunc(a.blobAPI)))
	router.HandleFunc("/whoami", a.whoamiAPI)
	router.HandleFunc("/build", a.buildAPI)
	router.HandleFunc("/config", a.configAPI)
	router.HandleFunc("/plugins", a.pluginsAPI)
	router.HandleFunc("/users", a.usersAPI)
	router.HandleFunc("/user/{id}", a.userAPI)

	for _, v := range a.Conf.Plugins {
		v := v
		router.PathPrefix("/" + v.Metadata.Identifier).HandlerFunc(
			func(res http.ResponseWriter, req *http.Request) {
				if signin.GetUser(req).HasPermission("plugin:" + v.Metadata.Identifier) {
					http.StripPrefix("/admin/api/"+v.Metadata.Identifier,
						pluginapi.NewAuthenticatedReverseProxy(internal.SingleJoiningSlash(v.Addr, "/admin_api"), true)).ServeHTTP(res, req)
				} else {
					http.Error(res, internal.JsonStatusForbidden, http.StatusForbidden)
				}
			})
	}
}
