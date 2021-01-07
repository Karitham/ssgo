package post

import (
	"bytes"
	"html/template"
	"io/ioutil"

	"github.com/rs/zerolog/log"
	"github.com/yuin/goldmark"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/parser"
)

// Article is a markdown post
type Article struct {
	Filepath string
	Style    string
	Body     template.HTML
	Meta     map[string]string
}

// GetPath returns the Path
func (art *Article) GetPath() string {
	return art.Filepath
}

// GetMeta returns the Metadata map
func (art *Article) GetMeta() map[string]string {
	return art.Meta
}

// Post builds an HTML post from the article content
// It parses the markdown and gets the metadata tags that it finds and accept
func (art *Article) Post(t *template.Template, fp Paths, md goldmark.Markdown) Poster {
	content, err := ioutil.ReadFile(art.Filepath)
	if err != nil {
		log.Error().Stack().Err(err).Str("filepath", art.Filepath).Msg("error reading file")
		return nil
	}

	// Create the HTML file
	file, err := MakeHTMLFile(art, fp)
	if err != nil {
		log.Err(err).Str("filepath", art.Filepath).Msg("error creating file")
		return nil
	}
	defer file.Close()

	var buf bytes.Buffer
	context := parser.NewContext()

	// Convert markdown contained in the article to the buffer, and use context to get metadata from the article
	if err := md.Convert(content, &buf, parser.WithContext(context)); err != nil {
		log.Err(err).Str("filepath", art.Filepath).Str("content", string(content)).Msg("error converting file")
		return nil
	}

	// Build the article
	art.Meta = ParseMetadata(art, meta.GetItems(context))
	art.Body = template.HTML(buf.String())

	// Run the template building in the html file
	if err := t.ExecuteTemplate(file, "post.tmpl", art); err != nil {
		log.Err(err).Str("filepath", file.Name()).Msgf("article: %v", art)
	}

	log.Info().Str("file", art.Filepath).Msg("Converted article")
	return art
}
