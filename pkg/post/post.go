package post

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"text/template"

	"github.com/Karitham/ssgo/pkg/config"
	"github.com/alecthomas/chroma/formatters/html"
	mathjax "github.com/litao91/goldmark-mathjax"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting"
)

// Post represent the templated file
type Post struct {
	Script    string
	PageTitle string
	Body      string
}

// Index represent the templated Index file
type Index struct {
	FileTree []IndexTree
	Script   string
}

// IndexTree represent each file present in the current directory by it's URL and name
type IndexTree struct {
	FileURL   string
	FileTitle string
}

// Execute is used to run the whole post making process
// TODO : Add more options to run and make it extensible
func Execute(conf *config.General) error {
	posts, err := ListFiles(conf.Directories.PostDir, true)
	if err != nil {
		return err
	}

	conf.Templates, err = ParseTemplates(conf.Directories.TemplateDir)
	if err != nil {
		return err
	}

	conf.Markdown = goldmark.New(
		goldmark.WithExtensions(mathjax.MathJax),
		goldmark.WithExtensions(
			highlighting.NewHighlighting(
				highlighting.WithFormatOptions(html.WithClasses(true)),
			),
		),
	)

	// make each post
	directories := makePosts(posts, conf)

	// make index files
	for _, d := range directories {
		err := CreateIndex(conf, d)
		if err != nil {
			return err
		}
	}

	conf.Log.Printf("Wrote %d files\n", len(posts)+len(directories))
	return nil
}

func makePosts(posts []string, conf *config.General) []string {
	// wg waits for the goroutine to finish making all the files
	// before making the navigation menu
	var directories []string
	var wg sync.WaitGroup
	for _, p := range posts {
		if file, err := os.Lstat(p); err == nil && file.IsDir() {
			directories = append(directories, p)
			continue
		}
		wg.Add(1)
		go MakePost(p, &wg, conf)
	}
	wg.Wait()

	return directories
}

// MakePost makes a post and inserts the content
// TODO : Simplify function signature
// TODO : have the post making extensible via config or such
func MakePost(post string, wg *sync.WaitGroup, conf *config.General) {
	defer wg.Done()

	filecontent, err := ioutil.ReadFile(post)
	if err != nil {
		log.Println(err)
	}

	var buf bytes.Buffer
	if err := conf.Markdown.Convert(filecontent, &buf); err != nil {
		log.Println(err)
	}

	f, err := CreateHTMLFile(conf.Directories.PublDir, conf.Directories.PostDir, &post)
	if err != nil {
		log.Println(err)
	}

	postName := GetFilename(TrimFileExt(post))

	err = conf.Templates.ExecuteTemplate(f,
		"post.tmpl",
		Post{
			Script:    conf.Server.Script,
			PageTitle: postName,
			Body:      buf.String(),
		},
	)
	if err != nil {
		log.Println(err)
	}

	err = f.Close()
	if err != nil {
		log.Println(err)
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
func CreateHTMLFile(publDir, postDir string, filePath *string) (file *os.File, err error) {
	// Convert the `.md` file to a `html` and change the directory
	var publpath = ConvertExt(filepath.Join(publDir, TrimDir(*filePath, postDir)), "html")

	// Get the final directory path
	dir := filepath.Join(publDir, TrimDir(trimFilename(*filePath), postDir))

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

			if (!withDir && info.IsDir()) || strings.HasPrefix(info.Name(), "_") {
				return nil
			}

			files = append(files, path)
			return nil
		},
	)
}

// FileTree returns an IndexTree based on a fileInfo array
func FileTree(f ...os.FileInfo) (tree []IndexTree) {
	for _, file := range f {
		fn := file.Name()
		if fn == "index.html" {
			continue
		}
		tree = append(
			tree,
			IndexTree{
				FileTitle: strings.ToUpper(TrimFileExt(fn)),
				FileURL:   fn,
			},
		)
	}
	return
}

// CreateIndex creates an index file in every directory, made for navigation purposes
func CreateIndex(conf *config.General, directory string) error {
	filename := (directory + string(filepath.Separator) + "index")
	f, err := CreateHTMLFile(conf.Directories.PublDir, conf.Directories.PostDir, &filename)
	if err != nil {
		return err
	}

	files, err := ioutil.ReadDir(trimFilename(f.Name()))
	if err != nil {
		return err
	}

	return conf.Templates.ExecuteTemplate(f, "index.tmpl", Index{FileTree: FileTree(files...), Script: conf.Server.Script})
}

// ParseTemplates is used to parse all the templates inside the given directory
func ParseTemplates(TemplateDir string) (tpl *template.Template, err error) {
	templates, err := ListFiles(TemplateDir, false)
	if err != nil {
		return nil, err
	}

	return template.ParseFiles(templates...)
}
