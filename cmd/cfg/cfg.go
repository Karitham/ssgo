package cfg

import (
	"github.com/rs/zerolog"
)

// Global represents the global struct
type Global struct {
	Server  Server
	Post    Post
	Logging Logging
}

// Logging configures the logger
type Logging struct {
	Level zerolog.Level
}

// Server represents the server configuration
type Server struct {
	Name   string
	Port   uint
	Script string
}

// Post represents the post configuration
type Post struct {
	DraftPrefix string
	PublPath    string
	PostPath    string
	TmplPath    string
	Index       Index
	Article     Article
}

// Index is the index file
type Index struct {
	OptionalFile string
	Meta         map[string]string
	Style        string
}

// Article is the article file
type Article struct {
	Meta  map[string]string
	Style string
}

// New returns a global with the default values
func New() *Global {
	return &Global{
		Post: Post{
			PublPath:    "../public",
			PostPath:    "../posts",
			TmplPath:    "../assets/templates",
			DraftPrefix: "_",
			Index: Index{
				OptionalFile: "_index.md",
				Meta: map[string]string{
					"title":       "",
					"description": "",
					"date":        "",
					"background":  "",
					"icon":        "",
				},
				Style: "/assets/css/post.css",
			},
			Article: Article{
				Style: "/assets/css/post.css",
				Meta: map[string]string{
					"description": "",
					"date":        "",
					"background":  "",
					"icon":        "",
					"title":       "",
					"url":         "",
				},
			},
		},
		Logging: Logging{
			Level: zerolog.DebugLevel,
		},
	}
}
