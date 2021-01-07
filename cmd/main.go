package main

import (
	"html/template"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/Karitham/ssgo/post"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/yuin/goldmark"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	t, err := post.LoadTemplates("../assets/templates")
	if err != nil {
		log.Fatal().Err(err).Msg("could not load templates")
	}

	md := post.NewMarkdown()
	style := "../assets/css/post.css"
	prefix := "_"
	paths := post.Paths{
		Old: "../posts",
		New: "../public",
	}

	folder, err := post.Walker(paths.Old)
	if err != nil {
		log.Fatal().Err(err).Msg("could not walk post files")
	}

	buildRecursive(t, paths, md, prefix, folder, string(style))
}

func buildRecursive(t *template.Template, paths post.Paths, md goldmark.Markdown, prefix string, folder *post.Folder, ArtStyle string) {
	if folder.Files == nil {
		return
	}
	log.Debug().Strs("files", folder.Files).Msg("files to convert")

	// Recursive calling for each directory
	wg := new(sync.WaitGroup)
	for i := range folder.Folders {
		wg.Add(1)
		go func(f post.Folder) {
			buildRecursive(t, paths, md, prefix, &f, ArtStyle)
			wg.Done()
		}(folder.Folders[i])
	}
	wg.Wait()

	// TODO : take care of the _index.md files so metadata is read correctly
	// Default index
	i := post.Index{
		Filepath: folder.Path,
		Style:    ArtStyle,
		Tree:     make([]post.Article, 0, len(folder.Files)),
		Meta: map[string]string{
			"title":       "",
			"description": "",
			"date":        "",
			"image":       "",
			"icon":        "",
		},
	}

	// Default article
	art := post.Article{
		Style: ArtStyle,
	}

	// Run the poster
	i.Run(
		folder.FlattenFilter(func(path string) bool { return !strings.HasPrefix(filepath.Base(path), prefix) }),
		art,
		t,
		paths,
		md,
	)
}
