package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/dfkdream/hugocms/plugin"
)

type pluginData struct {
	Addr     string          `json:"addr"`
	Metadata plugin.Metadata `json:"metadata"`
}

type config struct {
	Dir         string
	ConfigPath  string
	ContentPath string
	PublicPath  string
	BoltPath    string
	Bind        string
	TLS         bool
	CertPath    string
	KeyPath     string
	Plugins     []pluginData
}

func (c config) String() string {
	if res, err := json.MarshalIndent(c, "", "    "); err == nil {
		return string(res)
	} else {
		panic(err)
	}
}

func getConfig() *config {
	cfg := config{
		Dir:      getEnvStringOr("DIR", "."),
		Bind:     getEnvStringOr("BIND", "0.0.0.0:80"),
		BoltPath: getEnvStringOr("BOLT", "./bolt.db"),
		TLS:      getEnvBoolOr("TLS", false),
		CertPath: getEnvStringOr("CERT", "./cert.pem"),
		KeyPath:  getEnvStringOr("KEY", "./key.pem"),
		Plugins:  make([]pluginData, 0),
	}
	cfg.ConfigPath = getEnvStringOr("CONFIG", filepath.Join(cfg.Dir, "config.yaml"))
	cfg.ContentPath = getEnvStringOr("CONTENT", filepath.Join(cfg.Dir, "/content"))
	cfg.PublicPath = getEnvStringOr("PUBLIC", filepath.Join(cfg.Dir, "/public"))

	ps := getEnvStringOr("PLUGINS", "")
	if ps != "" {
		for _, v := range strings.Split(ps, ",") {
			pluginAddr := strings.TrimSpace(v)
			if !strings.HasPrefix(pluginAddr, "http") {
				pluginAddr = singleJoiningSlash("http://", pluginAddr)
			}

			if !checkPluginLive(pluginAddr) {
				log.Fatalf("plugin %s: plugin does not response", pluginAddr)
			}

			m, err := getPluginMetadata(pluginAddr)
			if err != nil {
				log.Fatal(err)
			}

			cfg.Plugins = append(cfg.Plugins, pluginData{
				Addr:     pluginAddr,
				Metadata: *m,
			})
		}
	}

	return &cfg
}

func getEnvStringOr(key string, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvBoolOr(key string, defaultValue bool) bool {
	if value, err := strconv.ParseBool(os.Getenv(key)); err == nil {
		return value
	}
	return defaultValue
}

func singleJoiningSlash(a, b string) string {
	aSlash := strings.HasSuffix(a, "/")
	bSlash := strings.HasPrefix(b, "/")
	switch {
	case aSlash && bSlash:
		return a + b[1:]
	case !aSlash && !bSlash:
		return a + "/" + b
	}
	return a + b
}

func getPluginMetadata(pluginAddr string) (*plugin.Metadata, error) {
	addr := singleJoiningSlash(pluginAddr, "metadata")
	res, err := http.Get(addr)
	if err != nil {
		return nil, err
	}
	defer func() { _ = res.Body.Close() }()
	m := new(plugin.Metadata)
	err = json.NewDecoder(res.Body).Decode(m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func checkPluginLive(pluginAddr string) bool {
	addr := singleJoiningSlash(pluginAddr, "live")
	res, err := http.Get(addr)
	if err != nil {
		return false
	}
	return res.StatusCode == 200
}
