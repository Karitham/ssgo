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
	// PublDir is where the result will end up
	PublDir = "public"
	// TemplateDir is where the templates are located
	TemplateDir = "templates"
	// PostDir is where the markdown posts are located
	PostDir = "posts"
)

// Index represent the templated file
type Index struct {
	PageTitle string
	Body      string
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
			log.Println(err)
		}

		f, err := CreateHTMLFile(PublDir, PostDir, post)
		if err != nil {
			log.Println(err)
		}

		postName := getFileName(post)

		t.Execute(f, Index{PageTitle: postName[:len(postName)-len(filepath.Ext(postName))], Body: string(buf.Bytes())})
	}
}

func trimDir(path, dir string) string {
	return strings.TrimPrefix(strings.TrimPrefix(path, dir), string(filepath.Separator))
}

func trimFilename(path string) string {
	splittedPath := strings.Split(path, string(filepath.Separator))
	return strings.TrimSuffix(path, splittedPath[len(splittedPath)-1])
}

func getFileName(path string) string {
	return string(path[len(trimFilename(path)):])
}

// ConvertExt changes a file's extension to the given ext, for example "opus".
func ConvertExt(file string, ext string) string {
	return file[:len(file)-len(filepath.Ext(file))] + "." + ext
}

// CreateHTMLFile create an html file in the publDir that has the same path and name as the filepath input
func CreateHTMLFile(publDir, postDir, filePath string) (file *os.File, err error) {
	// Convert the `.md` file to a `html` and change the directory
	var publpath = ConvertExt(filepath.Join(publDir, trimDir(filePath, postDir)), "html")

	// Get the final directory path
	dir := filepath.Join(publDir, trimDir(trimFilename(filePath), postDir))

	// Make the final directory if it doesn't exist
	err = os.MkdirAll(dir, 0755)
	if err != nil {
		return nil, err
	}

	// Create the file
	return os.Create(publpath)
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
