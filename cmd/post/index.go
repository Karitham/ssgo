package post

import (
	"bytes"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/yuin/goldmark"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/parser"
	"gopkg.in/yaml.v2"
)

// Index is a directory where an index file is generated
type Index struct {
	Filepath    string
	Style       string
	Title       string
	Description string
	Tree        []Article
}

// GetPath returns the Path
func (i *Index) GetPath() string {
	return i.Filepath
}

// ParseMetadata parse the index metadata
func (i *Index) ParseMetadata(items yaml.MapSlice) {
	for _, m := range items {
		key, okK := m.Key.(string)
		value, okV := m.Value.(string)
		if !okK && !okV {
			log.Warn().Msg("metadata doesn't conform to the standard")
			continue
		}

		log.Trace().Str("key", key).Str("value", value).Msg("parsing metadata")

		switch {
		case strings.EqualFold(key, "title"):
			i.Title = value
			continue
		case strings.EqualFold(key, "description"):
			i.Description = value
			continue
		}
	}
}

// Post builds an HTML index
func (i *Index) Post(t *template.Template, fp Paths, md goldmark.Markdown) Poster {
	// Read the md file if it exists
	var content []byte

	if filepath.Ext(i.Filepath) == ".md" {
		var err error
		content, err = ioutil.ReadFile(i.Filepath)
		if err != nil {
			log.Err(err).Str("file", i.Filepath).Msg("error reading file")
			return nil
		}
	}

	// Create the HTML file
	file, err := os.Create(filepath.Join(i.Filepath, "index.html"))
	if err != nil {
		log.Err(err).Str("file", i.Filepath).Msg("error creating file")
		return nil
	}
	defer file.Close()

	var mdBuf bytes.Buffer
	context := parser.NewContext()

	// Get the metadata contained in the file if there is some,
	// the buffer will be emptied
	// so we can build the body
	if err := md.Convert(content, &mdBuf, parser.WithContext(context)); err != nil {
		log.Err(err).Str("file", i.Filepath).Str("content", string(content)).Msg("error converting file")
		return nil
	}

	// Get the index metadata
	i.ParseMetadata(meta.GetItems(context))

	// Run the template building in the html file
	if err := t.ExecuteTemplate(file, "index.tmpl", i); err != nil {
		log.Err(err).Str("file", file.Name()).Msgf("article: %v", i)
		return nil
	}

	log.Debug().Str("file", i.Filepath).Msg("built index")

	return i
}
