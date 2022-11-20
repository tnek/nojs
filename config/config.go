package config

import "os"

type AppConfig struct {
	AppName string
	Host    string
	Port    int
	DBPath  string
	Domain  string

	Templates []string
}

var (
	SESSION_KEY = []byte(os.Getenv("SESSION_KEY"))
)
