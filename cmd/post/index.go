package post

import (
	"bytes"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"

	"github.com/Karitham/ssgo/cfg"
	"github.com/rs/zerolog/log"
	"github.com/yuin/goldmark"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/parser"
)

// Index is a directory where an index file is generated
type Index struct {
	Filepath string
	Style    string
	Meta     map[string]string
	Tree     []Article
}

// GetPath returns the Path
func (i *Index) GetPath() string {
	return i.Filepath
}

// GetMeta returns the Metadata map
func (i *Index) GetMeta() map[string]string {
	return i.Meta
}

// Post builds an HTML index
func (i *Index) Post(t *template.Template, fp Paths, md goldmark.Markdown) Poster {
	// Read the md file if it exists
	var content []byte

	if filepath.Ext(i.Filepath) == ".md" {
		var err error
		content, err = ioutil.ReadFile(i.Filepath)
		if err != nil {
			log.Error().Err(err).Str("filepath", i.Filepath).Msg("error reading file")
			return nil
		}
	}

	// Create the HTML file
	file, err := os.Create(filepath.Join(filepath.Dir((PathConvert(i.Filepath, fp))), "index.html"))
	if err != nil {
		log.Error().Err(err).Str("filepath", i.Filepath).Msg("error creating file")
		return nil
	}
	defer file.Close()

	context := parser.NewContext()
	// Get the metadata contained in the file if there is some,
	// the buffer will be emptied
	// so we can build the body
	if err := md.Convert(content, bytes.NewBuffer(content), parser.WithContext(context)); err != nil {
		log.Err(err).Str("filepath", i.Filepath).Str("content", string(content)).Msg("error converting file")
		return nil
	}

	// Get the index metadata
	i.Meta = ParseMetadata(i, meta.GetItems(context))

	// Run the template building in the html file
	if err := t.ExecuteTemplate(file, "index.tmpl", i); err != nil {
		log.Err(err).Str("file", file.Name()).Msgf("article: %v", i)
		return nil
	}

	log.Info().Str("filepath", i.Filepath).Msg("Built index")
	return i
}

// Run builds each file then builds up the index
func (i *Index) Run(files []string, artTemplate Article, t *template.Template, md goldmark.Markdown, conf cfg.Post) {
	wg := new(sync.WaitGroup)
	mu := new(sync.Mutex)
	paths := Paths{
		Old: conf.PostPath,
		New: conf.PublPath,
	}

	for _, file := range files {
		wg.Add(1)

		// Rebuild the article
		go func(file string) {
			art := artTemplate
			art.Filepath = file
			art.Meta = map[string]string{
				"description": "",
				"date":        "",
				"background":  "",
				"icon":        "",
				"title":       ConvertExt(filepath.Base(file), ""),
				"url":         filepath.Base(ConvertExt(filepath.Base(file), ".html")),
			}

			if art, ok := art.Post(t, paths, md).(*Article); ok {
				mu.Lock()
				i.Tree = append(i.Tree, *art)
				mu.Unlock()
			}

			wg.Done()
		}(file + "")
	}
	wg.Wait()

	i.Post(t, paths, md)
}
