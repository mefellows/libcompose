package main

import (
	"os"

	"github.com/codegangsta/cli"
	cliApp "github.com/mefellows/libcompose/cli/app"
	"github.com/mefellows/libcompose/cli/command"
	dockerApp "github.com/mefellows/libcompose/cli/docker/app"
	"github.com/mefellows/libcompose/version"
)

func main() {
	factory := &dockerApp.ProjectFactory{}

	app := cli.NewApp()
	app.Name = "libcompose-cli"
	app.Usage = "Command line interface for libcompose."
	app.Version = version.VERSION + " (" + version.GITCOMMIT + ")"
	app.Author = "Docker Compose Contributors"
	app.Email = "https://github.com/mefellows/libcompose"
	app.Before = cliApp.BeforeApp
	app.Flags = append(command.CommonFlags(), dockerApp.DockerClientFlags()...)
	app.Commands = []cli.Command{
		command.BuildCommand(factory),
		command.CreateCommand(factory),
		command.UpCommand(factory),
		command.StartCommand(factory),
		command.LogsCommand(factory),
		command.RestartCommand(factory),
		command.StopCommand(factory),
		command.ScaleCommand(factory),
		command.RmCommand(factory),
		command.PullCommand(factory),
		command.KillCommand(factory),
		command.PortCommand(factory),
		command.PsCommand(factory),
		command.PauseCommand(factory),
		command.UnpauseCommand(factory),
	}

	app.Run(os.Args)
}
