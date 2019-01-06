package main

import "io/ioutil"

const prismaInitFilename = "prisma_client.go"

const prismaInitContent = `var client = New(&Options{
	Endpoint: "http://localhost:4466",
})`

func CreatePrismaInitializer() error {
	projectName := "prisgo_test"
	packageHeader := "package " + projectName + "\n"

	finalContent := packageHeader + "\n" + prismaInitContent
	err := ioutil.WriteFile(prismaInitFilename, []byte(finalContent), 0446)
	if err != nil {
		return err
	}

	return nil
}
