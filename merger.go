package main

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"
)

const modelsFromGQLGen = "models_gen.go"
const primaGenerated = "prisma.go"

var typesStructsModel = regexp.MustCompile(`type ([A-Z])\w+ struct {`)
var typesGeneralModel = regexp.MustCompile(`type ([A-Z])\w+ \w+\n`)
var constsGroupModel = regexp.MustCompile(`const \(([^\r]*?)\)`)

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// ExtractTypesFromFile ...
func ExtractTypesFromFile(filename string) ([]string, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	r := typesStructsModel.FindAllString(string(data), -1)

	finalNames := make([]string, 0)

	for _, v := range r {
		newV := strings.Replace(v, "type", "", -1)
		// newV = strings.Replace(newV, "struct", "", -1)
		newV = strings.Replace(newV, "{", "", -1)
		newV = strings.Trim(newV, "\n ")
		finalNames = append(finalNames, newV)
	}

	r = typesGeneralModel.FindAllString(string(data), -1)

	for _, v := range r {
		newV := strings.Replace(v, "type", "", -1)
		// newV = strings.Replace(newV, "struct", "", -1)
		// newV = strings.Replace(newV, "{", "", -1)
		newV = strings.Trim(newV, "\n ")
		finalNames = append(finalNames, newV)
	}

	r = constsGroupModel.FindAllString(string(data), -1)

	for _, v := range r {
		newV := strings.Replace(v, "const (", "", -1)
		newV = strings.Replace(newV, ")", "", -1)
		chuncks := strings.Split(newV, "\n")
		for _, c := range chuncks {
			newC := strings.Trim(c, "\n\t ")
			if newC != "" {
				finalNames = append(finalNames, "const "+newC)
			}

		}

	}

	return finalNames, nil
}

// DeleteRedeclaredTypes ...
func DeleteRedeclaredTypes(prismaFilename, GQLFilename string) {
	prismaTypes, err := ExtractTypesFromFile(prismaFilename)
	if err != nil {
		panic(err)
	}

	GQLTypes, err := ExtractTypesFromFile(GQLFilename)
	if err != nil {
		panic(err)
	}

	fmt.Println("from GQL: ", len(GQLTypes))
	fmt.Println("from Prisma: ", len(prismaTypes))

	typesMatched := []string{}
	unMatched := []string{}

	for _, s := range GQLTypes {
		if contains(prismaTypes, s) {
			typesMatched = append(typesMatched, s)
		} else {
			unMatched = append(unMatched, s)
		}
	}

	structMatchedTypes := []string{}
	for _, t := range typesMatched {
		if strings.HasSuffix(t, "struct") {
			structMatchedTypes = append(structMatchedTypes, t)
		}

	}

	generalMatchedTypes := []string{}
	for _, t := range typesMatched {
		if !strings.HasSuffix(t, "struct") && !strings.HasPrefix(t, "const") {
			generalMatchedTypes = append(generalMatchedTypes, t)
		}

	}

	constsMatchedTypes := []string{}
	for _, t := range typesMatched {
		if strings.HasPrefix(t, "const") {
			constsMatchedTypes = append(constsMatchedTypes, t)
		}

	}

	fmt.Print("\n")
	fmt.Println("total equal structs: ", len(structMatchedTypes))
	fmt.Println("total equal general: ", len(generalMatchedTypes))
	fmt.Println("total equal consts: ", len(constsMatchedTypes))
	fmt.Println("---------------------------")
	fmt.Println("total matchs: ", len(typesMatched))
}
