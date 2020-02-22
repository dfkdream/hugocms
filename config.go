package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strconv"
)

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
	}
	cfg.ConfigPath = getEnvStringOr("CONFIG", filepath.Join(cfg.Dir, "config.yaml"))
	cfg.ContentPath = getEnvStringOr("CONTENT", filepath.Join(cfg.Dir, "/content"))
	cfg.PublicPath = getEnvStringOr("PUBLIC", filepath.Join(cfg.Dir, "/public"))
	return &cfg
}

func getEnvStringOr(key string, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvBoolOr(key string, defaultValue bool) bool {
	if value, error := strconv.ParseBool(os.Getenv(key)); error != nil {
		return value
	}
	return defaultValue
}
