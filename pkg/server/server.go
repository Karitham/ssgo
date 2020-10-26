package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/jaschaephraim/lrserver"
	"gopkg.in/fsnotify.v1"
)

// Serve ...
func Serve(name string, port uint16) error {

	go liveReload("./assets/css/", "./public/")

	http.Handle("/assets/css/", http.FileServer(http.Dir(".")))
	http.Handle("/", http.FileServer(http.Dir("public")))

	return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

func liveReload(directories ...string) {
	// Create file watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Println(err)
	}
	defer watcher.Close()

	err = watch(watcher, directories...)
	if err != nil {
		log.Println(err)
	}

	// Create and start LiveReload server
	lr := lrserver.New("SSGO", lrserver.DefaultPort)

	lr.SetErrorLog(nil)
	lr.SetStatusLog(nil)

	go func() {
		lr.ListenAndServe()
	}()

	// Start goroutine that requests reload upon watcher event
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				lr.Reload(event.Name)
			case err := <-watcher.Errors:
				log.Println(err)
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
