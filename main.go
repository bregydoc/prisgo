package main

import (
	"log"
	"os/exec"
	"strings"
)

func main() {

	commandPath := exec.Command("pwd")

	currentPath, err := commandPath.Output()
	if err != nil {
		log.Fatalln(err)
	}

	projectPath := string(currentPath)
	projectPath = strings.Trim(projectPath, " \n")

	log.Println("Workspace: ", projectPath)

	err = InitPrismaProject(projectPath)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Generation step 1 successful")

	// err = FixAndMergerSchema(projectPath)
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	log.Println("Generation step 2 successful")

	ComposeGQLGen()

	log.Println("Generation step 3 successful")

	err = MoveAndFixedPrismaGenerated()
	if err != nil {
		panic(err)
	}

	log.Println("Generation step 4 successful")

	err = DeleteRedeclaredTypes(primaGenerated, modelsFromGQLGen)
	if err != nil {
		panic(err)
	}
	log.Println("Generation step 5 successful")

	err = FixModelsGen()
	if err != nil {
		panic(err)
	}
	log.Println("Generation step 6 successful")

	err = CreatePrismaInitializer()
	if err != nil {
		panic(err)
	}
	log.Println("Generation step 7 successful")
}
