package main

import (
	"log"
	"os"

	"github.com/Karitham/ssgo/pkg/post"
	"github.com/Karitham/ssgo/pkg/server"
	"github.com/urfave/cli/v2"
)

// Config holds all of the configuration
type Config struct {
	Server      Server
	Directories Directories
	Log         *log.Logger
}

// Server holds the configuration needed for the server part
type Server struct {
	Enabled bool
	Name    string
	Port    uint16
	Script  string
}

// Directories represent the needed configurations for posts
type Directories struct {
	PublDir     string
	TemplateDir string
	PostDir     string
}

// Conf holds the configuration
var Conf Config

func main() {

	app := &cli.App{
		Name:  "SSGO",
		Usage: "Generate HTML based on the markdown you write",
		Action: func(_ *cli.Context) error {
			return post.Execute(
				Conf.Directories.PostDir,
				Conf.Directories.TemplateDir,
				Conf.Directories.PublDir,
				Conf.Server.Script,
				Conf.Log,
			)
		},
		Commands: []*cli.Command{
			{
				Name:    "server",
				Usage:   "serve your files with a liveserver",
				Aliases: []string{"s", "serve"},
				Action: func(_ *cli.Context) error {
					err := post.Execute(
						Conf.Directories.PostDir,
						Conf.Directories.TemplateDir,
						Conf.Directories.PublDir,
						Conf.Server.Script,
						Conf.Log,
					)
					if err != nil {
						return err
					}
					return server.Serve(&Conf.Server.Port, Conf.Log)
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatalln(err)
	}

}

func init() {
	Conf = Config{
		Server: Server{
			Enabled: false,
			Name:    "SSGO",
			Port:    5050,
			Script:  "<script src=\"http://localhost:35729/livereload.js\"></script>",
		},
		Directories: Directories{
			PublDir:     "public",
			TemplateDir: "assets/templates",
			PostDir:     "posts",
		},
		Log: log.New(os.Stdout, "[SSGO] ", 0),
	}
}
