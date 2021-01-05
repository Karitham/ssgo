package post

import (
	"html/template"
	"os"
	"path"
	"path/filepath"

	"github.com/alecthomas/chroma/formatters/html"
	mathjax "github.com/litao91/goldmark-mathjax"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"gopkg.in/yaml.v2"
)

// Poster builds a post
type Poster interface {
	Post(t *template.Template, fp Paths, md goldmark.Markdown) Poster
	ParseMetadata(items yaml.MapSlice)
	GetPath() string
}

// Paths represent path that changes between an original post and a converted one
type Paths struct {
	Old string
	New string
}

// NewMarkdown returns a markdown parser with the default configuration
// default config is : autoID, mathjax support, Github flavored markdown,
// metadata as yaml, and syntax highlighting with classes
func NewMarkdown() goldmark.Markdown {
	return goldmark.New(
		goldmark.WithParserOptions(parser.WithAutoHeadingID()),
		goldmark.WithExtensions(mathjax.MathJax),
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithExtensions(meta.Meta),
		goldmark.WithExtensions(
			highlighting.NewHighlighting(
				highlighting.WithFormatOptions(html.WithClasses(true)),
			),
		),
	)
}

// ConvertExt changes a file's extension to the given ext,
// for example ".html"
func ConvertExt(file string, ext string) string {
	return file[:len(file)-len(filepath.Ext(file))] + ext
}

// MakeHTMLFile create an html file on the newPath from the old filename/path
// It also creates directories as needed
func MakeHTMLFile(post Poster, fp Paths) (file *os.File, err error) {
	publpath := filepath.ToSlash(ConvertExt(PathConvert(post.GetPath(), fp), ".html"))

	// Make the final directory if it doesn't exist
	err = os.MkdirAll(path.Dir(publpath), 0755)
	if err != nil {
		return nil, err
	}
	return os.Create(publpath)
}

// LoadTemplates return the templates found in the directory
// You can then select them by filenames.
// The lookup is recursive so if you have a lot you can seperate them
func LoadTemplates(dir string) (*template.Template, error) {
	entries, err := Walker(dir)
	if err != nil {
		return nil, err
	}

	return template.ParseFiles(entries.Flatten()...)
}

// PathConvert converts path of old file to new, for example
// you pass in a post such as ~/asset/something.go
// and it replaces `~/asset/` with `/dev/null/`
// so it returns `/dev/null/something.go`
func PathConvert(p string, fp Paths) string {
	filepath := p[len(fp.Old):]
	return path.Join(fp.New, filepath)
}
