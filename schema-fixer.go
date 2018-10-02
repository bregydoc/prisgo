package main

import (
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

func FixAndMergerSchema(workspace string) error {
	filePath := join(workspace, "schema/prisma.graphql")
	schema, err := os.Open(filePath)
	defer schema.Close()
	if err != nil {
		return err
	}

	data, err := ioutil.ReadAll(schema)
	if err != nil {
		return err
	}
	dataAsString := string(data)
	reg, err := regexp.Compile("(.*?)_([a-zA-Z])*")
	if err != nil {
		return err
	}

	for len(reg.FindAllStringIndex(dataAsString, -1)) != 0 {
		c := reg.FindStringIndex(dataAsString)

		start := c[0]
		end := c[1]
		value := dataAsString[start:end]
		countSpaces := strings.Count(value, " ")
		valueInCamelCase := ToLowerCamel(value)
		finalValue := strings.Repeat(" ", countSpaces) + valueInCamelCase
		dataAsString = strings.Replace(dataAsString, value, finalValue, -1)
	}

	fixedSchema := join(workspace, "schema/fixed_prisma.graphql")
	err = ioutil.WriteFile(fixedSchema, []byte(dataAsString), 0644)

	if err != nil {
		return err
	}

	return nil

}
