package server

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/Karitham/ssgo/pkg/config"
	"github.com/Karitham/ssgo/pkg/post"
	"github.com/jaschaephraim/lrserver"
	"gopkg.in/fsnotify.v1"
)

// Serve ...
func Serve(conf *config.General) error {
	go liveReload(conf, "./assets/", "./public/", "./posts/")

	http.Handle("/assets/", http.FileServer(http.Dir(".")))
	http.Handle("/", http.FileServer(http.Dir("public")))

	conf.Log.Printf("Live server listening at http://localhost:%d\n", conf.Server.Port)

	return http.ListenAndServe(fmt.Sprintf(":%d", conf.Server.Port), nil)
}

func liveReload(conf *config.General, directories ...string) {
	// Create file watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		conf.Log.Println(err)
	}
	defer watcher.Close()

	err = watch(watcher, directories...)
	if err != nil {
		conf.Log.Println(err)
	}

	// Create and start LiveReload server
	lr := lrserver.New("SSGO", lrserver.DefaultPort)

	lr.SetErrorLog(nil)
	lr.SetStatusLog(nil)

	go func() {
		conf.Log.Println(lr.ListenAndServe())
	}()

	// Start goroutine that requests reload upon watcher event
	go func() {
		wg := sync.WaitGroup{}
		for {
			select {
			case event := <-watcher.Events:
				if strings.HasPrefix(event.Name, conf.Directories.Post) &&
					!strings.HasPrefix(post.GetFilename(event.Name), "_") {
					wg.Add(1)
					post.MakePost(event.Name, &wg, conf)
				}
				wg.Wait()
				lr.Reload(event.Name)
			case err := <-watcher.Errors:
				conf.Log.Println(err)
			}
		}
	}()

	select {}
}

func watch(w *fsnotify.Watcher, directories ...string) error {
	var dirs []string

	for _, dir := range directories {
		err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if info.IsDir() {
				dirs = append(dirs, path)
			}
			return nil
		})
		if err != nil {
			return err
		}
	}

	for _, d := range dirs {
		err := w.Add(d)
		if err != nil {
			return err
		}
	}

	return nil
}
