package main

import (
	"log"

	"github.com/Karitham/ssgo/pkg/post"
	"github.com/Karitham/ssgo/pkg/server"
)

const (
	// PublDir is where the result will end up
	PublDir = "public"
	// TemplateDir is where the templates are located
	TemplateDir = "assets/templates"
	// PostDir is where the markdown posts are located
	PostDir = "posts"
)

// TODO : Add cobra for CLI usage
func main() {
	post.Run(PostDir, TemplateDir, PublDir)
	go func() {
		err := server.Serve("ssgo", 5050)
		if err != nil {
			log.Println(err)
		}
	}()
	select {}
}
