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
// and directories it can find.
// The files have the full path to access
// them but are respectively in subdirectories
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
		fp := filepath.Join(f.Path, e.Name())

		if e.IsDir() {
			f.Folders = append(f.Folders, Folder{Path: filepath.Join(f.Path, fp)})
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
	return f.FlattenFilter(func(s string) bool { return true })
}

// FlattenFilter flattens the directory,
// but only returns the files that matches the filter.
// The order isn't guaranteed but should be in depth order,
// the closer being first and later being last.
func (f *Folder) FlattenFilter(filter func(string) bool) (files []string) {
	for _, file := range f.Files {
		if filter(file) {
			files = append(files, file)
		}
	}
	for _, fold := range f.Folders {
		files = append(files, fold.FlattenFilter(filter)...)
	}
	return files
}
