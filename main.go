package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/alecthomas/chroma/formatters/html"
	mathjax "github.com/litao91/goldmark-mathjax"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting"
)

var (
	// PublDir is where the result will end up
	PublDir = "public"
	// TemplateDir is where the templates are located
	TemplateDir = "templates"
	// PostDir is where the markdown posts are located
	PostDir = "posts"
)

// Post represent the templated file
type Post struct {
	PageTitle string
	Body      string
}

// Index represent the templated file
type Index struct {
	FileTree []string
}

func main() {
	posts, err := ListFiles(PostDir, true)
	if err != nil {
		log.Println(err)
	}

	t, err := ParseTemplates(TemplateDir)
	if err != nil {
		log.Println(err)
	}

	md := goldmark.New(
		goldmark.WithExtensions(mathjax.MathJax),
		goldmark.WithExtensions(
			highlighting.NewHighlighting(
				highlighting.WithStyle("native"),
				highlighting.WithFormatOptions(
					html.WithClasses(true),
				),
			),
		),
	)

	for _, post := range posts {

		// Used to make index files
		if file, err := os.Lstat(post); err == nil && file.IsDir() {
			f, err := CreateHTMLFile(PublDir, PostDir, post+string(filepath.Separator)+"index")
			if err != nil {
				log.Println("createHTMLFile error main range posts: ", err)
			}

			err = t.ExecuteTemplate(
				f,
				"index.tmpl",
				Index{
					FileTree: posts,
				},
			)
			if err != nil {
				log.Println(err)
			}
			continue
		}

		filecontent, err := ioutil.ReadFile(post)
		if err != nil {
			log.Println(err)
		}

		// Create a corresponding HTML file
		f, err := CreateHTMLFile(PublDir, PostDir, post)
		if err != nil {
			log.Println(err)
		}

		var buf bytes.Buffer
		if err := md.Convert(filecontent, &buf); err != nil {
			log.Println(err)
		}

		postName := GetFilename(TrimFileExt(post))

		err = t.ExecuteTemplate(f,
			"post.tmpl",
			Post{
				PageTitle: postName,
				Body:      buf.String(),
			},
		)
		if err != nil {
			log.Println(err)
		}
	}
}

// TrimDir trims the directory of the given path
func TrimDir(path, dir string) string {
	return strings.TrimPrefix(strings.TrimPrefix(path, dir), string(filepath.Separator))
}

// TrimFileExt returns the path without the fileExt
func TrimFileExt(path string) string {
	return path[:len(path)-len(filepath.Ext(path))]
}

// trimFilename returns the directory path
func trimFilename(path string) string {
	splittedPath := strings.Split(path, string(filepath.Separator))
	return strings.TrimSuffix(path, splittedPath[len(splittedPath)-1])
}

// GetFilename is used to retrieve the the filename of a file from a path, returns the extension too
func GetFilename(path string) string {
	return string(path[len(trimFilename(path)):])
}

// ConvertExt changes a file's extension to the given ext, for example "html".
func ConvertExt(file string, ext string) string {
	return file[:len(file)-len(filepath.Ext(file))] + "." + ext
}

// CreateHTMLFile create an html file in the publDir that has the same path and name as the filepath input
func CreateHTMLFile(publDir, postDir, filePath string) (file *os.File, err error) {
	// Convert the `.md` file to a `html` and change the directory
	var publpath = ConvertExt(filepath.Join(publDir, TrimDir(filePath, postDir)), "html")

	// Get the final directory path
	dir := filepath.Join(publDir, TrimDir(trimFilename(filePath), postDir))

	// Make the final directory if it doesn't exist
	err = os.MkdirAll(dir, 0755)
	if err != nil {
		return nil, err
	}

	// Create the file
	return os.Create(publpath)
}

// ListFiles list all files in a directory and subdirectories. Returns relative path of all files except those starting with a `.`
func ListFiles(dir string, withDir bool) (files []string, err error) {
	return files, filepath.Walk(dir,

		// Walk the file tree while getting both files and directories
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if (!withDir && info.IsDir()) || strings.HasPrefix(info.Name(), ".") {
				return nil
			}

			files = append(files, path)
			return nil
		},
	)
}

// ParseTemplates is used to parse all the templates inside the given directory
func ParseTemplates(TemplateDir string) (tpl *template.Template, err error) {
	templates, err := ListFiles(TemplateDir, false)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return template.ParseFiles(templates...)
}
