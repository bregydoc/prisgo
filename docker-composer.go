package main

import (
	"errors"
	"os"
)

const DockerComposePostgresSql = `version: '3'
services:
  prisma:
    image: prismagraphql/prisma:1.17
    restart: always
    ports:
    - "4466:4466"
    environment:
      PRISMA_CONFIG: |
        port: 4466
        # uncomment the next line and provide the env var PRISMA_MANAGEMENT_API_SECRET=my-secret to activate cluster security
        # managementApiSecret: my-secret
        databases:
          default:
            connector: postgres
            host: postgres
            port: 5432
            user: prisma
            password: prisma
            migrations: true
            rawAccess: true
  postgres:
    image: postgres
    restart: always
    environment:
      POSTGRES_USER: prisma
      POSTGRES_PASSWORD: prisma
    volumes:
      - postgres:/var/lib/postgresql/data
volumes:
  postgres:

`
const DockerTypePostgres = "postgres"
const DockerTypeSql = "sql"

func CreateDockerCompose(basePath string, dockerType string) (filePath string, err error) {
	var file *os.File

	switch dockerType {

	case DockerTypePostgres:
		file, err = os.Create(join(basePath, "docker-compose.yml"))
		if err != nil {
			return "", err
		}
		_, err = file.WriteString(DockerComposePostgresSql)
		if err != nil {
			return "", err
		}

		err = file.Close()
		if err != nil {
			return "", err
		}
		break
	case DockerTypeSql:
		return "", errors.New("sql generator not implemented")
		break
	default:
		return "", errors.New("invalid type for code generation")
		break
	}
	return file.Name(), nil
}
