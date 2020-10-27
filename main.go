package main

import (
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
		Usage: "Generate HTML based on the markdown you provide, easily customizable with your own theme",
		Action: func(_ *cli.Context) error {
			return post.Execute(conf)
		},
		UsageText: "ssgo [global options] command [command options]",
		Commands: []*cli.Command{
			{
				Name:    "server",
				Usage:   "serve your files with a live reloading server",
				Aliases: []string{"s"},
				Action: func(_ *cli.Context) error {
					err := post.Execute(conf)
					if err != nil {
						return err
					}
					return server.Serve(conf)
				},
				Flags: []cli.Flag{
					&cli.UintFlag{
						Name:        "port",
						Aliases:     []string{"p"},
						Value:       conf.Server.Port,
						Usage:       "Change the port of the server",
						Destination: &conf.Server.Port,
					},
				},
			},
		},
		Flags: []cli.Flag{
			&cli.PathFlag{
				Name:        "post",
				Value:       conf.Directories.Post,
				Usage:       "Change the post directory",
				Destination: &conf.Directories.Post,
			},
			&cli.PathFlag{
				Name:        "publ",
				Value:       conf.Directories.Publ,
				Usage:       "Change the publication directory",
				Destination: &conf.Directories.Publ,
			},
			&cli.PathFlag{
				Name:        "tmpl",
				Value:       conf.Directories.Tmpl,
				Usage:       "Change the template directory",
				Destination: &conf.Directories.Tmpl,
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		conf.Log.Fatalln(err)
		os.Exit(1)
	}
}
