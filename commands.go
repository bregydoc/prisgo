package main

import "strings"

const prismaCommand = "prisma"

func join(base string, leaf string) string {
	if strings.HasSuffix(base, "/") {
		return base + leaf
	}

	return base + "/" + leaf
}
