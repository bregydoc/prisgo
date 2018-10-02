package main

import (
	"log"
	"os/exec"
	"time"
)

func InitPrismaProject(workspace string) error {
	log.Println("Creating Docker Composer")
	_, err := CreateDockerCompose(workspace, DockerTypePostgres)
	if err != nil {
		return err
	}
	log.Println("Generating Prisma config")
	_, err = GeneratePrismaconfig(workspace)
	if err != nil {
		return err
	}

	log.Println("Verifying if you have the schema seed")

	_, err = VerifyAndCreateTheSeedIfNecessary(workspace)
	if err != nil {
		return err
	}

	log.Println("Inflating all generated code")
	generator := exec.Command(prismaCommand, "generate")
	t1 := time.Now()
	err = generator.Run()
	if err != nil {
		log.Println(generator.Output())
		return err
	}
	t2 := time.Now()
	log.Println("Time elapsed: ", t2.Sub(t1).String())
	return nil
}
