package app

import (
	"github.com/codegangsta/cli"
	"github.com/mefellows/libcompose/cli/command"
	"github.com/mefellows/libcompose/cli/logger"
	"github.com/mefellows/libcompose/docker"
	"github.com/mefellows/libcompose/project"
)

// ProjectFactory is a struct that hold the app.ProjectFactory implementation.
type ProjectFactory struct {
}

// Create implements ProjectFactory.Create using docker client.
func (p *ProjectFactory) Create(c *cli.Context) (*project.Project, error) {
	context := &docker.Context{}
	context.LoggerFactory = logger.NewColorLoggerFactory()
	Populate(context, c)
	command.Populate(&context.Context, c)

	return docker.NewProject(context)
}
