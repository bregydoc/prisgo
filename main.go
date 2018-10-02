package main

import "log"

func main() {
	err := InitPrismaProject("/Users/bregy/go-work/src/prisgo-demo")
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Generation step 1 successful")

	err = FixAndMergerSchema("/Users/bregy/go-work/src/prisgo-demo")
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Generation step 2 successful")

	ComposeGQLGen()

	log.Println("Generation step 3 successful")
}
