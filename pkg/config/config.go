package config

import (
	"log"
	"os"
	"text/template"

	"github.com/alecthomas/chroma/formatters/html"
	mathjax "github.com/litao91/goldmark-mathjax"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting"
	"github.com/yuin/goldmark/extension"
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
	Port    uint
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
			Name:   "SSGO",
			Port:   5050,
			Script: "<script src=\"http://localhost:35729/livereload.js\"></script>",
		},
		Directories: Directories{
			PublDir:     "public",
			TemplateDir: "assets/templates",
			PostDir:     "posts",
		},
		Log: log.New(os.Stdout, "[SSGO] ", 0),

		Markdown: goldmark.New(
			goldmark.WithExtensions(mathjax.MathJax),
			goldmark.WithExtensions(extension.GFM),
			goldmark.WithExtensions(
				highlighting.NewHighlighting(
					highlighting.WithFormatOptions(
						html.WithClasses(true),
					),
				),
			),
		),
	}
}
