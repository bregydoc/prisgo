package main

func main() {

	// commandPath := exec.Command("pwd")

	// currentPath, err := commandPath.Output()
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	// projectPath := string(currentPath)
	// projectPath = strings.Trim(projectPath, " \n")

	// log.Println("Workspace: ", projectPath)

	// err = InitPrismaProject(projectPath)
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	// log.Println("Generation step 1 successful")

	// // err = FixAndMergerSchema(projectPath)
	// // if err != nil {
	// // 	log.Fatalln(err)
	// // }

	// log.Println("Generation step 2 successful")

	// ComposeGQLGen()

	// log.Println("Generation step 3 successful")

	DeleteRedeclaredTypes(primaGenerated, modelsFromGQLGen)

}
