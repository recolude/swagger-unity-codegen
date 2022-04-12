package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/ghodss/yaml"
	"github.com/recolude/swagger-unity-codegen/unitygen"
	"github.com/recolude/swagger-unity-codegen/unitygen/convention"
	"github.com/recolude/swagger-unity-codegen/unitygen/model"
	"github.com/recolude/swagger-unity-codegen/unitygen/model/property"
	"github.com/recolude/swagger-unity-codegen/unitygen/security"
	"github.com/spf13/afero"
	"github.com/urfave/cli/v2"
)

func filterSpecForTags(spec unitygen.Spec, tags []string) unitygen.Spec {
	if tags == nil {
		return spec
	}
	filteredServices := make([]unitygen.Service, 0)
	for _, service := range spec.Services {
		for _, tag := range tags {
			if service.Name() == tag {
				filteredServices = append(filteredServices, service)
			}
		}
	}
	spec.Services = filteredServices
	return spec
}

func alreadyRecursed(inquestion string, alreadyAdded []string) bool {
	for _, q := range alreadyAdded {
		if q == inquestion {
			return true
		}
	}
	return false
}

func findReferencePropRecurse(inQuestion model.Property, defs []model.Definition, alreadyAdded []string) []string {
	finalReferences := alreadyAdded

	objectReferenceDefinition, ok := inQuestion.(property.DefinitionReference)
	if ok {
		if alreadyRecursed(objectReferenceDefinition.ToVariableType(), finalReferences) {
			return finalReferences
		}
		for _, def := range defs {
			if def.ToVariableType() == objectReferenceDefinition.ToVariableType() {
				finalReferences = findReferenceRecurse(def, defs, finalReferences)
			}
		}
	}

	arrayDefinition, ok := inQuestion.(property.Array)
	if ok {
		return findReferencePropRecurse(arrayDefinition.Property(), defs, finalReferences)
	}

	objectDefinition, ok := inQuestion.(property.Object)
	if ok {
		finalReferences = findReferenceRecurse(objectDefinition.Object(), defs, finalReferences)
	}

	return finalReferences
}

func findReferenceRecurse(inQuestion model.Definition, defs []model.Definition, alreadyFound []string) []string {
	finalReferences := alreadyFound

	objectDefinition, ok := inQuestion.(model.Object)
	if ok {
		if alreadyRecursed(objectDefinition.ToVariableType(), finalReferences) {
			return finalReferences
		}
		finalReferences = append(finalReferences, objectDefinition.ToVariableType())
		for _, prop := range objectDefinition.Properties() {
			finalReferences = findReferencePropRecurse(prop, defs, finalReferences)
		}
	}

	stringEnumDefinition, ok := inQuestion.(model.StringEnum)
	if ok {
		finalReferences = append(finalReferences, stringEnumDefinition.ToVariableType())
	}

	numEnumDefinition, ok := inQuestion.(model.NumberEnum)
	if ok {
		finalReferences = append(finalReferences, numEnumDefinition.ToVariableType())
	}

	return finalReferences
}

func buildReferenceMapping(defs []model.Definition) map[string][]string {
	refMapping := make(map[string][]string)
	for _, def := range defs {
		refMapping[def.ToVariableType()] = findReferenceRecurse(def, defs, nil)
	}
	return refMapping
}

func filterSpecForUnusedDefinitions(spec unitygen.Spec) unitygen.Spec {
	thingsToKeep := make(map[string]bool)
	for _, def := range spec.Definitions {
		thingsToKeep[def.ToVariableType()] = false
	}

	referenceMapping := buildReferenceMapping(spec.Definitions)

	for _, service := range spec.Services {
		for _, path := range service.Paths() {
			for _, param := range path.Parameters() {
				if param.Schema() != nil {
					if _, ok := thingsToKeep[param.Schema().ToVariableType()]; ok {
						thingsToKeep[param.Schema().ToVariableType()] = true

						for _, reference := range referenceMapping[param.Schema().ToVariableType()] {
							thingsToKeep[reference] = true
						}
					}
				}
			}

			for _, resp := range path.Responses() {
				if resp != nil {
					if _, ok := thingsToKeep[resp.VariableType()]; ok {
						thingsToKeep[resp.VariableType()] = true
						for _, reference := range referenceMapping[resp.VariableType()] {
							thingsToKeep[reference] = true
						}
					}
				}
			}
		}
	}

	filteredDefinitions := make([]model.Definition, 0)
	for _, def := range spec.Definitions {
		if thingsToKeep[def.ToVariableType()] {
			filteredDefinitions = append(filteredDefinitions, def)
		}
	}

	spec.Definitions = filteredDefinitions
	return spec
}

func fileCommentHeader(out io.Writer) {
	fmt.Fprintln(out, "// This code was generated by: ")
	fmt.Fprintln(out, "// https://github.com/recolude/swagger-unity-codegen")
	fmt.Fprintln(out, "// Issues and PRs welcome :)")
	fmt.Fprintln(out, "")
}

func fileImports(out io.Writer) {
	fmt.Fprintln(out, "using UnityEngine;")
	fmt.Fprintln(out, "using UnityEngine.Networking;")
	fmt.Fprintln(out, "using System.Collections;")
	fmt.Fprintln(out, "using System.Text;")
	fmt.Fprintln(out, "using Newtonsoft.Json;")
	fmt.Fprintln(out, "using Newtonsoft.Json.Converters;")
	fmt.Fprintln(out, "using JsonSubTypes;")
	fmt.Fprintln(out, "")
}

func writeDefinitions(out io.Writer, defs []model.Definition) {
	for _, def := range defs {
		fmt.Fprintf(out, "%s\n\n", def.ToCSharp())
	}
}

func writeServices(out io.Writer, services []unitygen.Service, authDefinitions []security.Auth) {
	for _, def := range services {
		fmt.Fprintf(out, "%s\n\n", def.ToCSharp(authDefinitions, "Config"))
	}
}

func openNamespace(out io.Writer, namespace string) {
	if namespace == "" {
		return
	}
	fmt.Fprintf(out, "namespace %s {\n\n", convention.TitleCase(namespace))
}

func closeNamespace(out io.Writer, namespace string) {
	if namespace == "" {
		return
	}
	fmt.Fprint(out, "}")
}

// allOutAtOnce is for when you want the entirety of the swagger in a single file.
func allOutAtOnce(c *cli.Context, spec unitygen.Spec) error {
	fileCommentHeader(c.App.Writer)
	fileImports(c.App.Writer)

	namespace := c.String("namespace")

	openNamespace(c.App.Writer, namespace)

	// Define out all classes!
	fmt.Fprint(c.App.Writer, "#region Definitions\n\n")
	writeDefinitions(c.App.Writer, spec.Definitions)
	fmt.Fprintf(c.App.Writer, "%s\n\n", "#endregion")

	fmt.Fprintf(c.App.Writer, "%s\n\n", "#region Services")
	fmt.Fprintf(c.App.Writer, "%s\n\n", spec.ServiceConfig(c.String("config-name"), c.String("config-menu"), c.Bool("scriptable-object-config")))
	writeServices(c.App.Writer, spec.Services, spec.AuthDefinitions)
	fmt.Fprint(c.App.Writer, "#endregion\n\n")

	closeNamespace(c.App.Writer, namespace)

	return nil
}

func toDir(fs afero.Fs, c *cli.Context, location string, spec unitygen.Spec) error {
	namespace := c.String("namespace")

	err := fs.MkdirAll(location, os.ModePerm)
	if err != nil {
		return fmt.Errorf("error creating folder: %w", err)
	}

	definitionsFile, err := fs.Create(path.Join(location, "Definitions.cs"))
	if err != nil {
		return fmt.Errorf("error creating definitions file: %w", err)
	}
	fileCommentHeader(definitionsFile)
	fileImports(definitionsFile)
	openNamespace(definitionsFile, namespace)
	writeDefinitions(definitionsFile, spec.Definitions)
	closeNamespace(definitionsFile, namespace)

	servicesFile, err := fs.Create(path.Join(location, "Services.cs"))
	if err != nil {
		return fmt.Errorf("error creating services file: %w", err)
	}
	fileCommentHeader(servicesFile)
	fileImports(servicesFile)
	openNamespace(servicesFile, namespace)
	writeServices(servicesFile, spec.Services, spec.AuthDefinitions)
	closeNamespace(servicesFile, namespace)

	configFile, err := fs.Create(path.Join(location, fmt.Sprintf("%s.cs", convention.TitleCase(c.String("config-name")))))
	if err != nil {
		return fmt.Errorf("error creating config file: %w", err)
	}
	fileCommentHeader(configFile)
	fileImports(configFile)
	openNamespace(configFile, namespace)
	fmt.Fprintf(configFile, "%s\n\n", spec.ServiceConfig(c.String("config-name"), c.String("config-menu"), c.Bool("scriptable-object-config")))
	closeNamespace(configFile, namespace)

	return nil
}

func buildApp(fs afero.Fs, out io.Writer, errOut io.Writer) *cli.App {
	return &cli.App{
		Name:        "swag3d",
		Description: "Generate C# code specifically for Unity3D from a swagger file",
		Version:     "0.1.0",
		Usage:       "swagger and Unity3D meet",
		Writer:      out,
		ErrWriter:   errOut,
		Authors: []*cli.Author{
			{
				Name: "Elijah C Davis",
			},
		},
		Flags: []cli.Flag{},
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
					&cli.StringFlag{
						Name:        "config-name",
						Usage:       "The name of the config class that will contain all server variables",
						Value:       "ServiceConfig",
						DefaultText: "ServiceConfig",
					},
					&cli.StringFlag{
						Name:        "config-menu",
						Usage:       "Name to be listed in the Assets/Create submenu, so that instances of the server config can be easily created and stored in the project as \".asset\"",
						Value:       "Server/Config",
						DefaultText: "Server/Config",
					},
					&cli.StringFlag{
						Name:    "file",
						Aliases: []string{"f"},
						Usage:   "where to load swagger from",
						Value:   "swagger.json",
					},
					&cli.BoolFlag{
						Name:  "include-unused",
						Usage: "Whether or not to include definitions that where never used in the different services",
						Value: false,
					},
					&cli.BoolFlag{
						Name:  "scriptable-object-config",
						Usage: "Whether or not to generate a scriptable object that contains all server values different services will use.",
						Value: true,
					},
					&cli.StringSliceFlag{
						Name:  "tags",
						Usage: "Specify tags that a route must have to be included in the export. Specifying no tags means include all routes",
					},

					&cli.StringFlag{
						Name:        "out",
						Usage:       "The directory where you want the contents to be written to disc",
						DefaultText: "",
					},
				},
				Action: func(c *cli.Context) error {
					fileToLoad := c.String("file")
					extension := strings.ToLower(filepath.Ext(fileToLoad))

					var jsonStream io.Reader

					switch extension {
					case ".json":
						file, err := fs.Open(fileToLoad)
						if err != nil {
							return fmt.Errorf("error opening swagger file: %w", err)
						}
						defer file.Close()
						jsonStream = file

					case ".yaml":
						fallthrough
					case ".yml":
						file, err := fs.Open(fileToLoad)
						if err != nil {
							return fmt.Errorf("error opening swagger file: %w", err)
						}
						defer file.Close()
						y, err := io.ReadAll(file)
						if err != nil {
							return fmt.Errorf("error reading YAML: %w", err)
						}

						j2, err := yaml.YAMLToJSON(y)
						if err != nil {
							return fmt.Errorf("error translating YAML: %w", err)
						}
						jsonStream = bytes.NewBuffer(j2)

					default:
						return fmt.Errorf("unrecognized swagger file format '%s', please provide either json or yml", extension)
					}

					spec, err := unitygen.NewParser().ParseJSON(jsonStream)
					if err != nil {
						return fmt.Errorf("error reading from swagger file: %w", err)
					}
					spec = filterSpecForTags(spec, c.StringSlice("tags"))
					if !c.Bool("include-unused") {
						spec = filterSpecForUnusedDefinitions(spec)
					}

					if c.IsSet("out") {
						return toDir(fs, c, c.String("out"), spec)
					}

					return allOutAtOnce(c, spec)
				},
			},
		},
	}
}

func main() {
	app := buildApp(afero.NewOsFs(), os.Stdout, os.Stderr)
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
