package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/99designs/gqlgen/codegen"
	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v2"
)

const schemaFilename = "schema/prisma.graphql"

// ComposeGQLGen ...
func ComposeGQLGen() {
	var config *codegen.Config
	var err error

	configFilename := ""
	serverFilename := "server/server.go"

	config, err = codegen.LoadConfigFromDefaultLocations()

	if config != nil {
		fmt.Fprintf(os.Stderr, "init failed: a configuration file already exists at %s\n", config.FilePath)
		os.Exit(1)
	}

	if !os.IsNotExist(errors.Cause(err)) {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	configFilename = "gqlgen.yml"

	config = codegen.DefaultConfig()

	config.SchemaFilename = codegen.SchemaFilenames{schemaFilename}

	config.Resolver = codegen.PackageConfig{
		Filename: "resolver.go",
		Type:     "Resolver",
	}

	var buf bytes.Buffer
	{
		var b []byte
		b, err = yaml.Marshal(config)
		if err != nil {
			fmt.Fprintln(os.Stderr, "unable to marshal yaml: "+err.Error())
			os.Exit(1)
		}
		buf.Write(b)
	}

	err = ioutil.WriteFile(configFilename, buf.Bytes(), 0644)
	if err != nil {
		fmt.Fprintln(os.Stderr, "unable to write config file: "+err.Error())
		os.Exit(1)
	}

	schemaRaw, err := ioutil.ReadFile(config.SchemaFilename[0])
	if err != nil {
		fmt.Fprintln(os.Stderr, "unable to open schema: "+err.Error())
		os.Exit(1)
	}
	config.SchemaStr = map[string]string{schemaFilename: string(schemaRaw)}

	if err = config.Check(); err != nil {
		fmt.Fprintln(os.Stderr, "invalid config format: "+err.Error())
		os.Exit(1)
	}

	if err := codegen.Generate(*config); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	if err := codegen.GenerateServer(*config, serverFilename); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	fmt.Printf("Exec \"go run ./%s\" to start GraphQL server\n", serverFilename)
}
