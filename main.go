package main

import (
	"log"
	"os"

	"github.com/Karitham/ssgo/pkg/config"
	"github.com/Karitham/ssgo/pkg/post"
	"github.com/Karitham/ssgo/pkg/server"
	"github.com/urfave/cli/v2"
)

// conf holds the configuration
var conf *config.General

func main() {
	conf = config.New()
	app := &cli.App{
		Name:  "SSGO",
		Usage: "Generate HTML based on the markdown you write",
		Action: func(_ *cli.Context) error {
			return post.Execute(conf)
		},
		Commands: []*cli.Command{
			{
				Name:    "server",
				Usage:   "serve your files with a liveserver",
				Aliases: []string{"s", "serve"},
				Action: func(_ *cli.Context) error {
					err := post.Execute(conf)
					if err != nil {
						return err
					}
					return server.Serve(conf)
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatalln(err)
	}
}
