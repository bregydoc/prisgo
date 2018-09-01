package main

import (
	"errors"
	"regexp"
	"strings"
)

var testData = `
type Todo {
  id: ID!
  text: String!
  done: Boolean!
  user: User!
}

type User {
  id: ID!
  name: String!
}

type Query {
  todos: [Todo!]!
}

input NewTodo {
  text: String!
  userId: String!
}

type Mutation {
  createTodo(input: NewTodo!): Todo!
}
`

var basicTypes = map[string]string{
	"String":  "string",
	"Boolean": "bool",
	"ID":      "string",
	"Int":     "int",
	"Float":   "float",
}

func parseSDLProp(data string) (*SDLProp, error) {
	rProps, err := regexp.Compile(`[^type] \w+:* *[A-Z]*[a-z]*\(*\w*!*\)*:*.*!*`)
	if err != nil {
		return nil, err
	}

	props := rProps.FindAllString(data, -1)
	for _, prop := range props {
		parts := strings.Split(prop, ":")
		if len(parts) != 2 {
			return nil, errors.New("invalid prop on SDLType")
		}
		nProp := strings.TrimSpace(parts[0])

		vProp := strings.TrimSpace(parts[1])

		prop := SDLProp{
			Name: nProp,
			Type: GQLType{
				Name:   strings.Replace(vProp, "!", "", -1),
				GoType: basicTypes[vProp],
			},
			Required: strings.Contains(vProp, "!"),
		}

	}
}

func parseType(data string) (*SDLType, error) {
	rName, err := regexp.Compile(`type [A-Z]\w+`)
	if err != nil {
		return nil, err
	}

	name := strings.TrimSpace(strings.Replace(rName.FindString(data), "type", "", -1))
	if name == "" {
		return nil, errors.New("invalid name on SDLType")
	}

}

func getTypesFromText(data string) ([]*SDLType, error) {
	r, err := regexp.Compile(`type ([A-Z])\w+[\ ]*{[^}\r]*}`)
	if err != nil {
		return nil, err
	}

	types := r.FindAllString(data, -1)
	typesArray := make([]*SDLType, 0)
	for i := 0; i < len(types); i++ {
		t, err := parseType(string(types[i]))

	}
}
