package main

import (
	"log"

	"github.com/mefellows/libcompose/docker"
	"github.com/mefellows/libcompose/project"
)

func main() {
	project, err := docker.NewProject(&docker.Context{
		Context: project.Context{
			ComposeFiles: []string{"docker-compose.yml"},
			ProjectName:  "yeah-compose",
		},
	})

	if err != nil {
		log.Fatal(err)
	}

	project.Up()
}
