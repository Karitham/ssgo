package post

import (
	"bytes"
	"html/template"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/yuin/goldmark"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/parser"
	"gopkg.in/yaml.v2"
)

// Article is a markdown post
type Article struct {
	Filepath    string
	Style       string
	Body        template.HTML
	Description string
	Title       template.HTML
	URL         string
	Date        string
}

// GetPath returns the Path
func (art *Article) GetPath() string {
	return art.Filepath
}

// ParseMetadata modifies the article with the metadata contained in the mapslice
func (art *Article) ParseMetadata(items yaml.MapSlice) {
	for _, m := range items {
		key, okK := m.Key.(string)
		value, okV := m.Value.(string)
		if !okK || !okV {
			continue
		}

		log.Trace().Str("key", key).Str("value", value).Msg("parsing metadata")

		switch {
		case strings.EqualFold(key, "title"):
			art.Title = template.HTML(value)
		case strings.EqualFold(key, "url"):
			art.URL = value
		case strings.EqualFold(key, "date"):
			art.Date = value
		case strings.EqualFold(key, "description"):
			art.Description = value
		}
	}
}

// Post builds an HTML post from the article content
// It parses the markdown and gets the metadata tags that it finds and accept
func (art *Article) Post(t *template.Template, fp Paths, md goldmark.Markdown) Poster {
	// Rebuild the article
	art = &Article{
		Filepath: art.Filepath,
		Style:    art.Style,
		Title:    template.HTML(ConvertExt(filepath.Base(art.Filepath), "")),
		URL:      filepath.Base(ConvertExt(filepath.Base(art.Filepath), ".html")),
		Body:     "",
		Date:     "",
	}

	// Read the md file
	content, err := ioutil.ReadFile(art.Filepath)
	if err != nil {
		log.Err(err).Str("file", art.Filepath).Msg("error reading file")
		return nil
	}

	// Create the HTML file
	file, err := MakeHTMLFile(art, fp)
	if err != nil {
		log.Err(err).Str("file", art.Filepath).Msg("error creating file")
		return nil
	}
	defer file.Close()

	var buf bytes.Buffer
	context := parser.NewContext()

	// Convert markdown contained in the article to the buffer, and use context to get metadata from the article
	if err := md.Convert(content, &buf, parser.WithContext(context)); err != nil {
		log.Err(err).Str("file", art.Filepath).Str("content", string(content)).Msg("error converting file")
		return nil
	}

	// Build the article
	art.ParseMetadata(meta.GetItems(context))
	art.Body = template.HTML(buf.String())

	// Run the template building in the html file
	if err := t.ExecuteTemplate(file, "post.tmpl", art); err != nil {
		log.Err(err).Str("file", file.Name()).Msgf("article: %v", art)
	}

	log.Debug().Str("file", art.Filepath).Msg("Converted article")
	return art
}
