package config

import (
	"log"
	"os"
	"text/template"

	"github.com/yuin/goldmark"
)

// General holds all of the configuration
type General struct {
	Server      Server
	Directories Directories
	Markdown    goldmark.Markdown
	Templates   *template.Template
	Log         *log.Logger
}

// Server holds the configuration needed for the server part
type Server struct {
	Enabled bool
	Name    string
	Port    uint16
	Script  string
}

// Directories represent the needed configurations for posts
type Directories struct {
	PublDir     string
	TemplateDir string
	PostDir     string
}

// New returns a default config
func New() *General {
	return &General{
		Server: Server{
			Enabled: false,
			Name:    "SSGO",
			Port:    5050,
			Script:  "<script src=\"http://localhost:35729/livereload.js\"></script>",
		},
		Directories: Directories{
			PublDir:     "public",
			TemplateDir: "assets/templates",
			PostDir:     "posts",
		},
		Log: log.New(os.Stdout, "[SSGO] ", 0),
	}
}
