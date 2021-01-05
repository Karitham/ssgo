package main

import (
	"html/template"
	"os"
	"path"
	"strings"
	"sync"

	"github.com/Karitham/ssgo/post"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/yuin/goldmark"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	// log.Logger = log.Logger.Level(zerolog.Disabled)

	md := post.NewMarkdown()

	t, err := post.LoadTemplates("../assets/templates")
	if err != nil {
		log.Fatal().Err(err).Msg("could not load templates")
	}

	prefix := "_"

	paths := post.Paths{
		Old: "../posts",
		New: "../public",
	}

	style := "../assets/css/post.css"

	folder, err := post.Walker(paths.Old)
	if err != nil {
		log.Fatal().Err(err).Msg("could not walk post files")
	}

	buildRecursive(t, paths, md, prefix, folder, string(style))
}

func buildRecursive(t *template.Template, paths post.Paths, md goldmark.Markdown, prefix string, folder *post.Folder, style string) {
	for _, fold := range folder.Folders {
		go buildRecursive(t, paths, md, prefix, &fold, style)
	}

	wg := new(sync.WaitGroup)
	var articles []post.Article

	for _, file := range folder.Files {
		if strings.HasPrefix(path.Base(file), prefix) {
			continue
		}
		wg.Add(1)
		go func(file string) {
			article := &post.Article{
				Filepath: file,
				Style:    style,
			}
			if article, ok := article.Post(t, paths, md).(*post.Article); ok {
				articles = append(articles, *article)
			}
			wg.Done()
		}(file)
	}
	wg.Wait()

	var i post.Index = post.Index{
		Filepath: post.PathConvert(folder.Path, paths),
		Style:    style,
		Tree:     articles,
	}

	i.Post(t, paths, md)
}
