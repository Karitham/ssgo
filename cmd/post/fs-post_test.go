package post

import (
	"reflect"
	"strings"
	"testing"
)

func TestFiltrer(t *testing.T) {
	type args struct {
		filter func(string) bool
		input  []string
	}
	tests := []struct {
		name         string
		args         args
		wantFiltered []string
	}{
		{
			name: "doesn't contain",
			args: args{
				filter: func(s string) bool { return !strings.Contains(s, "whack") },
				input: []string{
					"this doesn't contain the w word",
					"this does contain the word whack",
					"this is an additional one",
				},
			},
			wantFiltered: []string{
				"this doesn't contain the w word",
				"this is an additional one",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotFiltered := Filterer(tt.args.filter, tt.args.input...); !reflect.DeepEqual(gotFiltered, tt.wantFiltered) {
				t.Errorf("Filter() = %v, want %v", gotFiltered, tt.wantFiltered)
			}
		})
	}
}

func TestPathConvert(t *testing.T) {
	type args struct {
		p  string
		fp Paths
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Convert Path",
			args: args{
				p: "~/asset/something.go",
				fp: Paths{
					Old: "~/asset/",
					New: "/dev/null/",
				},
			},
			want: "/dev/null/something.go",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// if got := PathConvert(tt.args.p, tt.args.fp); got != tt.want {
			// 	t.Errorf("PathConvert() = %v, want %v", got, tt.want)
			// }
		})
	}
}

func TestConvertExt(t *testing.T) {
	type args struct {
		file string
		ext  string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Convert Ext",
			args: args{
				file: "post/something.md",
				ext:  ".html",
			},
			want: "post/something.html",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ConvertExt(tt.args.file, tt.args.ext); got != tt.want {
				t.Errorf("ConvertExt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFolder_Flatten(t *testing.T) {
	tests := []struct {
		name      string
		f         *Folder
		wantFiles []string
	}{
		{
			name: "",
			f: &Folder{
				Folders: []Folder{{
					Folders: []Folder{{
						Folders: nil,
						Files:   []string{"deep1\\deep2\\f6", "deep1\\deep2\\f7"},
						Path:    "./deep1\\deep2",
					}},
					Files: []string{"deep1\\f4", "deep1\\f5"},
					Path:  "./deep1",
				}},
				Files: []string{
					"f1",
					"f2",
					"f3",
				},
				Path: "./",
			},
			wantFiles: []string{
				"f1", "f2", "f3", "deep1\\f4", "deep1\\f5", "deep1\\deep2\\f6", "deep1\\deep2\\f7",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotFiles := tt.f.Flatten(); !reflect.DeepEqual(gotFiles, tt.wantFiles) {
				t.Errorf("Folder.Flatten() = %v, want %v", gotFiles, tt.wantFiles)
			}
		})
	}
}
