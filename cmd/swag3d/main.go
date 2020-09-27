package main

import (
	"fmt"
	"log"
	"os"

	"github.com/recolude/swagger-unity-codegen/unitygen"
	"github.com/urfave/cli/v2"
)

func buildApp() *cli.App {
	return &cli.App{
		Name:        "swag3d",
		Description: "Generate C# code specifically for Unity3D from a swagger file",
		Version:     "0.1.0",
		Usage:       "swagger and Unity3D meet",
		Authors: []*cli.Author{
			{
				Name: "Eli C Davis",
			},
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "file, f",
				Usage: "where to load swagger from",
				Value: "swagger.json",
			},
		},
		Commands: []*cli.Command{
			{
				Name:    "generate",
				Aliases: []string{"g"},
				Usage:   "Generate c# code from a swagger file",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "namespace",
						Usage: "The namespace the code will be wrapped in",
					},
				},
				Action: func(c *cli.Context) error {
					file, err := os.Open(c.String("file"))
					if err != nil {
						fmt.Fprintf(cli.ErrWriter, "Error opening swagger file: %s", err.Error())
					}

					spec, err := unitygen.Parser{}.ParseJSON(file)
					if err != nil {
						fmt.Fprintf(cli.ErrWriter, "Error opening swagger file: %s", err.Error())
					}

					namespace := c.String("namespace")

					if namespace != "" {
						fmt.Fprintf(c.App.Writer, "namespace %s {\n\n", namespace)
					}

					for _, def := range spec.Definitions {
						fmt.Fprintf(c.App.Writer, "%s\n\n", def.ToCSharp())
					}

					if namespace != "" {
						fmt.Fprint(c.App.Writer, "}")
					}

					return nil
				},
			},
		},
	}
}

func main() {
	app := buildApp()

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
