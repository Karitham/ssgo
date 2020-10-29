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

// Writer is what can be written
type Writer interface {
	Write() error
}

// Post represent the templated file
type Post struct {
	Script    string
	PageTitle string
	Body      string
	Config    *config.General
	WaitGroup *sync.WaitGroup
	Path      string
}

// Index represent the templated Index file
type Index struct {
	FileTree  []IndexTree
	Script    string
	Config    *config.General
	WaitGroup *sync.WaitGroup
	Path      string
}

// IndexTree represent each file present in the current directory by it's URL and name
type IndexTree struct {
	FileURL   string
	FileTitle string
}

// Execute is used to run the whole post making process
func Execute(conf *config.General) error {
	// written count the number of file written
	var written uint
	start := time.Now()
	posts, err := ListFiles(conf.Directories.Post, true)
	if err != nil {
		return err
	}

	conf.Templates, err = ParseTemplates(conf.Directories.Tmpl)
	if err != nil {
		return err
	}

	// make everything
	var wg = new(sync.WaitGroup)
	var toW Writer
	var Indexes []Writer
	for _, p := range posts {
		written++
		if file, err := os.Lstat(p); err == nil && file.IsDir() {
			Indexes = append(
				Indexes,
				&Index{
					Path:      p,
					Config:    conf,
					WaitGroup: wg,
				},
			)
			continue
		} else if !strings.EqualFold(filepath.Ext(p), ".md") {
			continue
		} else {
			toW = &Post{
				Path:      p,
				Config:    conf,
				WaitGroup: wg,
			}
		}

		wg.Add(1)
		go func() {
			if err = toW.Write(); err != nil {
				conf.Log.Println(err)
			}
		}()

	}
	wg.Wait()
	for _, i := range Indexes {
		wg.Add(1)
		err := i.Write()
		if err != nil {
			conf.Log.Println(err)
		}
	}
	wg.Wait()

	conf.Log.Printf("Wrote %d files in %s\n", written, time.Since(start))
	return nil
}

func (i *Index) Write() error {
	defer i.WaitGroup.Done()
	filename := (i.Path + string(filepath.Separator) + "index")
	f, err := createHTMLFile(&i.Config.Directories.Publ, &i.Config.Directories.Post, filename)
	if err != nil {
		return err
	}

	files, err := ioutil.ReadDir(TrimFilename(f.Name()))
	if err != nil {
		return err
	}
	return i.Config.Templates.ExecuteTemplate(
		f, "index.tmpl", Index{
			FileTree: FileTree(files...),
			Script:   i.Script,
		})
}

func (p *Post) Write() error {
	defer p.WaitGroup.Done()

	filecontent, err := ioutil.ReadFile(p.Path)
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	if err := p.Config.Markdown.Convert(filecontent, &buf); err != nil {
		return err
	}

	f, err := createHTMLFile(&p.Config.Directories.Publ, &p.Config.Directories.Post, p.Path)
	defer func(f *os.File) {
		_ = f.Close()
	}(f)
	if err != nil {
		return err
	}

	postName := GetFilename(trimFileExt(p.Path))
	return p.Config.Templates.ExecuteTemplate(
		f, "post.tmpl", Post{
			Body:      buf.String(),
			Script:    p.Config.Server.Script,
			PageTitle: postName,
		})
}

// trimDir trims the directory of the given path
func trimDir(path, dir string) string {
	return strings.TrimPrefix(strings.TrimPrefix(path, dir), string(filepath.Separator))
}

// trimFileExt returns the path without the fileExt
func trimFileExt(path string) string {
	return path[:len(path)-len(filepath.Ext(path))]
}

// TrimFilename returns the directory path
func TrimFilename(path string) string {
	splittedPath := strings.Split(path, string(filepath.Separator))
	return strings.TrimSuffix(path, splittedPath[len(splittedPath)-1])
}

// GetFilename is used to retrieve the the filename of a file from a path, returns the extension too
func GetFilename(path string) string {
	return string(path[len(TrimFilename(path)):])
}

// ConvertExt changes a file's extension to the given ext, for example "html".
func ConvertExt(file string, ext string) string {
	return file[:len(file)-len(filepath.Ext(file))] + "." + ext
}

// createHTMLFile create an html file in the publDir that has the same path and name as the filepath input
func createHTMLFile(publDir, postDir *string, filePath string) (file *os.File, err error) {
	// Convert the `.md` file to a `html` and change the directory
	var publpath = ConvertExt(filepath.Join(*publDir, trimDir(filePath, *postDir)), "html")

	// Get the final directory path
	dir := filepath.Join(*publDir, trimDir(TrimFilename(filePath), *postDir))

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

// ParseTemplates is used to parse all the templates inside the given directory
func ParseTemplates(TemplateDir string) (tpl *template.Template, err error) {
	templates, err := ListFiles(TemplateDir, false)
	if err != nil {
		return nil, err
	}

	return template.ParseFiles(templates...)
}
