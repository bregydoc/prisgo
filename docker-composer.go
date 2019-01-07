package main

import (
	"errors"
	"os"
)

const DockerComposePostgresSql = `version: '3'
services:
  prisma:
    image: prismagraphql/prisma:1.23
    restart: always
    ports:
    - "4466:4466"
    environment:
      PRISMA_CONFIG: |
        port: 4466
        databases:
          default:
            connector: postgres
            host: postgres
            port: 5432
            user: prisma
            password: prisma
            migrations: true
  postgres:
    image: postgres:10.5
    restart: always
    environment:
      POSTGRES_USER: prisma
      POSTGRES_PASSWORD: prisma
    volumes:
      - postgres:/var/lib/postgresql/data
volumes:
  postgres:

`
const DockerComposeSql = `version: '3'
services:
  prisma:
    image: prismagraphql/prisma:1.23
    restart: always
    ports:
    - "4466:4466"
    environment:
      PRISMA_CONFIG: |
        port: 4466
        databases:
          default:
            connector: mysql
            host: mysql
            port: 3306
            user: root
            password: prisma
            migrations: true
  mysql:
    image: mysql:5.7
    restart: always
    environment:
      MYSQL_ROOT_PASSWORD: prisma
    volumes:
      - mysql:/var/lib/mysql
volumes:
  mysql:
`

const DockerComposeMongoDB = `version: '3'
services:
  prisma:
    image: prismagraphql/prisma:1.23
    restart: always
    ports:
    - "4466:4466"
    environment:
      PRISMA_CONFIG: |
        port: 4466
        databases:
          default:
            connector: mongo
            uri: mongodb://prisma:prisma@mongo
  mongo:
    image: mongo:3.6
    restart: always
    environment:
      MONGO_INITDB_ROOT_USERNAME: prisma
      MONGO_INITDB_ROOT_PASSWORD: prisma
    ports:
      - "27017:27017"
    volumes:
      - mongo:/var/lib/mongo
volumes:
  mongo:
`

const DockerTypePostgres = "postgres"
const DockerTypeSql = "sql"
const DockerMongoType = "mongo"

func CreateDockerCompose(basePath string, dockerType string) (filePath string, err error) {
	var file *os.File
	file, err = os.Create(join(basePath, "docker-compose.yml"))
	if err != nil {
		return "", err
	}
	switch dockerType {
	case DockerTypePostgres:
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
		_, err = file.WriteString(DockerComposeSql)
		if err != nil {
			return "", err
		}

		err = file.Close()
		if err != nil {
			return "", err
		}
		break
	case DockerMongoType:
		_, err = file.WriteString(DockerComposeMongoDB)
		if err != nil {
			return "", err
		}

		err = file.Close()
		if err != nil {
			return "", err
		}
		break
	default:
		return "", errors.New("invalid type for code generation")
		break
	}
	return file.Name(), nil
}
