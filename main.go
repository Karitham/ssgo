package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/yuin/goldmark"
)

var (
	// ProdDir is where the result will end up
	ProdDir = "public"
	// TemplateDir is where the templates are located
	TemplateDir = "templates"
	// PostDir is where the markdown posts are located
	PostDir = "posts"
)

// Index represent the templated file
type Index struct {
	Body string
}

func main() {
	posts, err := listFiles(PostDir)
	if err != nil {
		log.Println(err)
	}

	t, err := parseTemplates(TemplateDir)
	if err != nil {
		log.Println(err)
	}

	for _, post := range posts {
		var buf bytes.Buffer
		md, err := ioutil.ReadFile(post)
		if err != nil {
			log.Println(err)
		}

		if err := goldmark.Convert(md, &buf); err != nil {
			panic(err)
		}

		f, err := createHTMLPost(ProdDir, post)
		if err != nil {
			log.Println(err)
		}

		t.Execute(f, Index{Body: string(buf.Bytes())})
	}
}

func trimDir(path, dir string) string {
	return strings.TrimPrefix(strings.TrimPrefix(path, dir), string(filepath.Separator))
}

// ConvertExt changes a file's extension to the given ext, for example "opus".
func ConvertExt(file string, ext string) string {
	return file[:len(file)-len(filepath.Ext(file))] + "." + ext
}

func createHTMLPost(publDir string, fp string) (files *os.File, err error) {
	var publpath = ConvertExt(filepath.Join(publDir, trimDir(fp, PostDir)), "html")

	if _, err := os.Stat(publDir); os.IsNotExist(err) {
		os.Mkdir(publDir, 0755)
	}

	f, err := os.Create(publpath)
	if err != nil {
		return nil, err
	}

	return f, nil
}

func listFiles(dir string) (files []string, err error) {
	return files, filepath.Walk(dir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() || strings.HasPrefix(info.Name(), ".") {
				return nil
			}
			files = append(files, path)
			return nil
		},
	)
}

func parseTemplates(TemplateDir string) (*template.Template, error) {
	templates, err := listFiles(TemplateDir)
	if err != nil {
		log.Println(err)
	}

	return template.ParseFiles(templates...)
}
