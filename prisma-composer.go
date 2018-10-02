package main

import "os"

const BasePrimaConfig = `endpoint: http://localhost:4466
datamodel: seed.graphql

generate:
  - generator: go-client
    output: ./generated/prisma
  - generator: graphql-schema
    output: ./schema/
hooks:
  post-deploy:
    - prisma generate
`

const DefaultSeedContent = `type User {
  id: ID! @unique
  name: String!
}
`

func VerifyAndCreateTheSeedIfNecessary(basePath string) (filePath string, err error) {
	fullPath := join(basePath, "seed.graphql")
	seed, err := os.Open(fullPath)
	defer seed.Close()
	if !os.IsExist(err) {
		seed, err = os.Create(fullPath)

		_, err = seed.WriteString(DefaultSeedContent)
		if err != nil {
			return "", err
		}

	} else if err != nil {
		return "", err
	}

	return seed.Name(), nil
}

func GeneratePrismaconfig(basePath string) (filePath string, err error) {
	file, err := os.Create(join(basePath, "prisma.yml"))
	defer file.Close()
	if err != nil {
		return "", err
	}

	_, err = file.WriteString(BasePrimaConfig)
	if err != nil {
		return "", err
	}

	return file.Name(), nil
}
