package main

import (
	"fmt"
	"os"
)

type Config struct {
	ApiHost string
	Port    string
	Sqlite  struct {
		DbDir   string
		MigrDir string
	}
	TLS struct {
		CertDir string
		KeyDir  string
	}
	UI struct {
		TmplDir string
		CSSDir  string
	}
	Github struct {
		ClientSecret string
		ClientID     string
	}
	Google struct {
		ClientSecret string
		ClientID     string
	}
	PostImagesDir string
}

func newConfig() *Config {
	apiPort := os.Getenv("PORT")
	if len(apiPort) == 0 {
		apiPort = "8080"
	}

	apiHost := os.Getenv("API_HOST")
	if len(apiHost) == 0 {
		apiHost = fmt.Sprintf("http://localhost:%s", apiPort)
	}

	return &Config{
		ApiHost: apiHost,
		Port:    apiPort,
		Sqlite: struct {
			DbDir   string
			MigrDir string
		}{DbDir: "./storage/storage.db?parseTime=true", MigrDir: "./migrations/sqlite/00001_initial.up.sql"},
		TLS: struct {
			CertDir string
			KeyDir  string
		}{CertDir: "./tls/cert.pem", KeyDir: "./tls/key.pem"},
		UI: struct {
			TmplDir string
			CSSDir  string
		}{TmplDir: "./ui/html/", CSSDir: "./ui/static/"},
		Github: struct {
			ClientSecret string
			ClientID     string
		}{ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"), ClientID: os.Getenv("GITHUB_CLIENT_ID")},
		Google: struct {
			ClientSecret string
			ClientID     string
		}{ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"), ClientID: os.Getenv("GOOGLE_CLIENT_ID")},
		PostImagesDir: "./post_images/",
	}
}
