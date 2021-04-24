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
	log.Logger = log.Level(conf.Logging.Level)

	md := post.NewMarkdown()
	tmpl, err := post.LoadTemplates(conf.Post.TmplPath)
	if err != nil {
		log.Fatal().Err(err).Msg("could not load templates")
	}

	// Get the fs from the postPath
	folder, err := post.Walker(conf.Post.PostPath)
	if err != nil {
		log.Fatal().Err(err).Msg("could not walk post files")
	}

	buildRecursive(conf, tmpl, md, folder)
}

func buildRecursive(conf *cfg.Global, tmpl *template.Template, md goldmark.Markdown, folder *post.Folder) *post.Index {
	if folder.Files == nil {
		return nil
	}
	log.Debug().Strs("files", folder.Files).Msg("Files found")

	// Default index
	i := post.Index{
		Filepath: folder.Path,
		Style:    conf.Post.Index.Style,
		Meta:     post.NewMetaMap(conf.Post.Index.Meta),
		Tree:     make([]post.Article, 0, len(folder.Files)),
	}

	// Recursive calling for each directory
	wg := new(sync.WaitGroup)
	for _, f := range folder.Folders {
		wg.Add(1)
		go func(folder post.Folder) {
			buildRecursive(conf, tmpl, md, &folder)
			wg.Done()
		}(f)
	}
	wg.Wait()

	// Check if there's an _index.md file
	if index := Filtrer(func(f string) bool { return filepath.Base(f) == "_index.md" }, folder.Files...); len(index) == 1 {
		log.Debug().Str("filepath", index[0]).Msg("Found index")
		folder.Path = index[0]
	}

	// Default article
	art := post.Article{
		Style: conf.Post.Article.Style,
	}

	for _, f := range folder.Folders {
		i.Tree = append(i.Tree, post.Article{Meta: map[string]string{"title": f.Path, "url": f.Path}})
	}

	// Run the poster
	return i.Run(
		// Files but filtered on the specified prefix
		Filtrer(
			func(path string) bool {
				return !strings.HasPrefix(filepath.Base(path), conf.Post.DraftPrefix)
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
