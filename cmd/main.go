package main

import (
	"html/template"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/Karitham/ssgo/cfg"
	"github.com/Karitham/ssgo/post"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/yuin/goldmark"
)

func main() {
	conf := cfg.New()

	// Configure logger
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	log.Logger = log.Level(conf.Level)

	md := post.NewMarkdown()
	tmpl, err := post.LoadTemplates(conf.TmplPath)
	if err != nil {
		log.Fatal().Err(err).Msg("could not load templates")
	}

	// Get the fs from the postPath
	folder, err := post.Walker(conf.PostPath)
	if err != nil {
		log.Fatal().Err(err).Msg("could not walk post files")
	}

	buildRecursive(conf, tmpl, md, folder)
}

func buildRecursive(conf *cfg.Global, tmpl *template.Template, md goldmark.Markdown, folder *post.Folder) {
	if folder.Files == nil {
		return
	}
	log.Debug().Strs("files", folder.Files).Msg("Files found")

	// Recursive calling for each directory
	wg := new(sync.WaitGroup)
	for i := range folder.Folders {
		wg.Add(1)
		go func(folder post.Folder) {
			buildRecursive(conf, tmpl, md, &folder)
			wg.Done()
		}(folder.Folders[i])
	}
	wg.Wait()

	// Check if there's an _index.md file
	if index := Filtrer(func(f string) bool { return filepath.Base(f) == "_index.md" }, folder.Files...); len(index) == 1 {
		log.Debug().Str("filepath", index[0]).Msg("Found index")
		folder.Path = index[0]
	}

	// Default index
	i := post.Index{
		Filepath: folder.Path,
		Style:    conf.Index.Style,
		Tree:     make([]post.Article, 0, len(folder.Files)),
		Meta: map[string]string{
			"title":       "",
			"description": "",
			"date":        "",
			"background":  "",
			"icon":        "",
		},
	}

	// Default article
	art := post.Article{
		Style: conf.Article.Style,
	}

	// Run the poster
	i.Run(
		// Files but filtered on _
		Filtrer(
			func(path string) bool {
				return !strings.HasPrefix(filepath.Base(path), conf.DraftPrefix)
			},
			folder.Files...,
		),
		art,
		tmpl,
		md,
		conf.Post,
	)
}

// Filtrer applies a filterFunc on an input
func Filtrer(filter func(string) bool, input ...string) (filtered []string) {
	for _, inp := range input {
		if filter(inp) {
			filtered = append(filtered, inp)
		}
	}
	return filtered
}
