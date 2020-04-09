package config

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/dfkdream/hugocms/protowrapper"

	"github.com/dfkdream/hugocms/internal"

	"github.com/dfkdream/hugocms/plugin"
)

type PluginData struct {
	Addr     string          `json:"addr"`
	Metadata plugin.Metadata `json:"metadata"`
}

type Config struct {
	Dir         string
	ConfigPath  string
	ContentPath string
	PublicPath  string
	BoltPath    string
	Bind        string
	TLS         bool
	CertPath    string
	KeyPath     string
	Plugins     []PluginData
}

func (c Config) String() string {
	if res, err := json.MarshalIndent(c, "", "    "); err == nil {
		return string(res)
	} else {
		panic(err)
	}
}

func GetConfig() *Config {
	cfg := Config{
		Dir:      getEnvStringOr("DIR", "."),
		Bind:     getEnvStringOr("BIND", "0.0.0.0:80"),
		BoltPath: getEnvStringOr("BOLT", "./bolt.db"),
		TLS:      getEnvBoolOr("TLS", false),
		CertPath: getEnvStringOr("CERT", "./cert.pem"),
		KeyPath:  getEnvStringOr("KEY", "./key.pem"),
		Plugins:  make([]PluginData, 0),
	}
	cfg.ConfigPath = getEnvStringOr("CONFIG", filepath.Join(cfg.Dir, "config.yaml"))
	cfg.ContentPath = getEnvStringOr("CONTENT", filepath.Join(cfg.Dir, "/content"))
	cfg.PublicPath = getEnvStringOr("PUBLIC", filepath.Join(cfg.Dir, "/public"))

	ps := getEnvStringOr("PLUGINS", "")
	if ps != "" {
		for _, v := range strings.Split(ps, ",") {
			pluginAddr := strings.TrimSpace(v)
			if !strings.HasPrefix(pluginAddr, "http") {
				pluginAddr = internal.SingleJoiningSlash("http://", pluginAddr)
			}

			if !CheckPluginLive(pluginAddr) {
				log.Fatalf("plugin %s: plugin does not response", pluginAddr)
			}

			m, err := getPluginMetadata(pluginAddr)
			if err != nil {
				log.Fatal(err)
			}

			cfg.Plugins = append(cfg.Plugins, PluginData{
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

func getPluginMetadata(pluginAddr string) (*plugin.Metadata, error) {
	addr := internal.SingleJoiningSlash(pluginAddr, "metadata")
	res, err := http.Get(addr)
	if err != nil {
		return nil, err
	}
	defer func() { _ = res.Body.Close() }()

	m := new(plugin.Metadata)

	err = protowrapper.NewDecoder(res.Body).Decode(m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func CheckPluginLive(pluginAddr string) bool {
	addr := internal.SingleJoiningSlash(pluginAddr, "live")
	res, err := http.Get(addr)
	if err != nil {
		return false
	}
	return res.StatusCode == 200
}
