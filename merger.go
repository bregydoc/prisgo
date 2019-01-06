package main

import (
	"fmt"
	"io/ioutil"
	"os"
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
func DeleteRedeclaredTypes(prismaFilename, GQLFilename string) error {
	prismaTypes, err := ExtractTypesFromFile(prismaFilename)
	if err != nil {
		return err
	}

	GQLTypes, err := ExtractTypesFromFile(GQLFilename)
	if err != nil {
		return err
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

	// create regex to delete determitaed struct:

	prismaFileData, err := ioutil.ReadFile(prismaFilename)
	if err != nil {
		return err
	}

	// gqlgenFileData, err := ioutil.ReadFile(GQLFilename)
	// if err != nil {
	// 	panic(err)
	// }

	replacementText := string(prismaFileData)

	totalDeleted := 0
	fmt.Print("\nDeleting structs\n|")
	for _, typeName := range structMatchedTypes {
		structSearcher := `type ` + typeName + ` \{([^\r]*?)\}`
		r := regexp.MustCompile(structSearcher)
		matchs := r.FindAllString(replacementText, -1)
		if len(matchs) == 1 {
			fmt.Print(".")
			new := r.ReplaceAllString(replacementText, "")
			replacementText = new
			totalDeleted++
		} else {
			fmt.Print(" ")
		}

	}
	fmt.Print("| ")
	fmt.Println(totalDeleted)

	totalDeleted = 0
	fmt.Print("\nDeleting generic types\n|")
	for _, typeName := range generalMatchedTypes {
		structSearcher := "type " + typeName
		r := regexp.MustCompile(structSearcher)
		matchs := r.FindAllString(replacementText, -1)
		if len(matchs) == 1 {
			fmt.Print(".")
			new := r.ReplaceAllString(replacementText, "")
			replacementText = new
			totalDeleted++
		} else {
			fmt.Print(" ")
		}

	}
	fmt.Print("| ")
	fmt.Println(totalDeleted)

	totalDeleted = 0
	fmt.Print("\nDeleting consts\n|")
	for _, c := range constsMatchedTypes {
		structSearcher := strings.Replace(c, "const ", "", -1)
		r := regexp.MustCompile(structSearcher)
		matchs := r.FindAllString(replacementText, -1)
		if len(matchs) == 1 {
			fmt.Print(".")
			new := r.ReplaceAllString(replacementText, "")
			replacementText = new
			totalDeleted++
		} else {
			structSearcher := c
			r := regexp.MustCompile(structSearcher)
			matchs := r.FindAllString(replacementText, -1)
			if len(matchs) == 1 {
				fmt.Print(".")
				new := r.ReplaceAllString(replacementText, "")
				replacementText = new
				totalDeleted++
			} else {
				fmt.Print(" ")
			}

		}

	}
	fmt.Print("| ")
	fmt.Println(totalDeleted)

	err = ioutil.WriteFile(prismaFilename, []byte(replacementText), 0446)
	if err != nil {
		return err
	}

	return nil

}

func MoveAndFixedPrismaGenerated() error {
	oldLocation := "generated/prisma/prisma.go"
	newLocation := "prisma.go"
	err := os.Rename(oldLocation, newLocation)
	if err != nil {
		return err
	}
	err = os.Remove("generated/prisma")
	if err != nil {
		return err
	}
	err = os.Remove("generated")
	if err != nil {
		return err
	}

	data, err := ioutil.ReadFile(newLocation)
	if err != nil {
		return err
	}

	prismaText := string(data)
	chuncks := strings.Split(prismaText, "\n")
	prismaText = strings.Join(chuncks[1:], "\n")

	projectName := "prisgo_test"
	prismaText = strings.Replace(prismaText, "package prisma", "package "+projectName, -1)

	err = ioutil.WriteFile(newLocation, []byte(prismaText), 0446)
	if err != nil {
		return err
	}
	return nil
}

func FixModelsGen() error {
	modelsGenFilename := "models_gen.go"

	data, err := ioutil.ReadFile(modelsGenFilename)
	if err != nil {
		return err
	}

	modelsGenText := string(data)

	modelsGenText = strings.Replace(modelsGenText, "Count string `json:\"count\"`", "Count int64 `json:\"count\"`", -1)
	err = ioutil.WriteFile(modelsGenFilename, []byte(modelsGenText), 0446)

	if err != nil {
		return err
	}
	return nil
}
