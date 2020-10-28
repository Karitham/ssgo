package post

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"text/template"
	"time"

	"github.com/Karitham/ssgo/pkg/config"
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

// written count the number of file written
var written uint

// Execute is used to run the whole post making process
// TODO : Add more options to run and make it extensible
func Execute(conf *config.General) error {
	start := time.Now()
	posts, err := ListFiles(conf.Directories.Post, true)
	if err != nil {
		return err
	}

	conf.Templates, err = ParseTemplates(conf.Directories.Tmpl)
	if err != nil {
		return err
	}

	// make each post
	directories := makePosts(posts, conf)

	// make index files
	for _, d := range directories {
		err := createIndex(conf, d)
		if err != nil {
			return err
		}
	}

	conf.Log.Printf("Wrote %d files in %s\n", written, time.Since(start))
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
		// Verify if we're reading a markdown file
		if !strings.EqualFold(filepath.Ext(p), ".md") {
			continue
		}

		wg.Add(1)
		go MakePost(p, &wg, conf)
		written++
	}
	wg.Wait()

	return directories
}

// MakePost makes a post and inserts the content
func MakePost(post string, wg *sync.WaitGroup, conf *config.General) {
	defer wg.Done()

	filecontent, err := ioutil.ReadFile(post)
	if err != nil {
		conf.Log.Println(err)
	}

	var buf bytes.Buffer
	if err := conf.Markdown.Convert(filecontent, &buf); err != nil {
		conf.Log.Println(err)
	}

	f, err := createHTMLFile(conf.Directories.Publ, conf.Directories.Post, &post)
	if err != nil {
		conf.Log.Println(err)
	}

	postName := GetFilename(trimFileExt(post))

	err = conf.Templates.ExecuteTemplate(f,
		"post.tmpl",
		Post{
			Script:    conf.Server.Script,
			PageTitle: postName,
			Body:      buf.String(),
		},
	)
	if err != nil {
		conf.Log.Println(err)
	}

	err = f.Close()
	if err != nil {
		conf.Log.Println(err)
	}
}

// trimDir trims the directory of the given path
func trimDir(path, dir string) string {
	return strings.TrimPrefix(strings.TrimPrefix(path, dir), string(filepath.Separator))
}

// trimFileExt returns the path without the fileExt
func trimFileExt(path string) string {
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

// createHTMLFile create an html file in the publDir that has the same path and name as the filepath input
func createHTMLFile(publDir, postDir string, filePath *string) (file *os.File, err error) {
	// Convert the `.md` file to a `html` and change the directory
	var publpath = ConvertExt(filepath.Join(publDir, trimDir(*filePath, postDir)), "html")

	// Get the final directory path
	dir := filepath.Join(publDir, trimDir(trimFilename(*filePath), postDir))

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
				FileTitle: strings.ToUpper(trimFileExt(fn)),
				FileURL:   fn,
			},
		)
	}
	return
}

// createIndex creates an index file in every directory, made for navigation purposes
func createIndex(conf *config.General, directory string) error {
	filename := (directory + string(filepath.Separator) + "index")
	f, err := createHTMLFile(conf.Directories.Publ, conf.Directories.Post, &filename)
	if err != nil {
		return err
	}

	files, err := ioutil.ReadDir(trimFilename(f.Name()))
	if err != nil {
		return err
	}
	written++
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
