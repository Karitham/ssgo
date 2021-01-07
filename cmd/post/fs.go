package post

import (
	"io/ioutil"
	"path/filepath"
)

// Folder represent a folder in the structure
type Folder struct {
	Folders []Folder
	Files   []string
	Path    string
}

// Walker walks the file system starting from the root
// and returns all the files
// and directories it can find
func Walker(root string) (*Folder, error) {
	return (&Folder{Path: root}).walker()
}

// See func Walker(root string) (*Folder, error)
func (f *Folder) walker() (*Folder, error) {
	currentDir, err := ioutil.ReadDir(f.Path)
	if err != nil {
		return nil, err
	}

	for _, e := range currentDir {
		fp := filepath.ToSlash(filepath.Join(f.Path, e.Name())) // Clean the filepath

		if e.IsDir() {
			f.Folders = append(f.Folders, Folder{Path: fp})
			if _, err := f.Folders[len(f.Folders)-1].walker(); err != nil {
				return nil, err
			}
			continue
		}

		f.Files = append(f.Files, fp)
	}
	return f, nil
}

// Flatten returns all the files inside a folder structure
// it returns only files and not folders.
// It's returned as a relative path from the root when you created
// The folder `f`
func (f *Folder) Flatten() (files []string) {
	for _, fold := range f.Folders {
		files = append(files, fold.Flatten()...)
	}
	return append(files, f.Files...)
}

// FilterFunc represent a filter
// When true it means it matches the filter
type FilterFunc func(string) bool

// FlattenFilter flattens the directory,
// but only returns the files that match the filter
func (f *Folder) FlattenFilter(filter FilterFunc) (files []string) {
	for _, fold := range f.Folders {
		files = append(files, fold.FlattenFilter(filter)...)
	}
	for _, file := range f.Files {
		if filter(file) {
			files = append(files, file)
		}
	}
	return files
}
