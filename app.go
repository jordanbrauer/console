package console

import (
	"log"
	"os"
)

type Configuration struct{}

func configure() Configuration {
	return Configuration{}
}

type App struct {
	Version  string
	Commands map[string]*Command
	Config   Configuration
}

type Option struct {
	Name        string
	Description string
	Default     any
	Value       any
}

type Argument struct {
	Name        string
	Description string
}

// New will create a CLI app for the given block chain.
func New(version string) *App {
	return &App{
		Version:  version,
		Commands: make(map[string]*Command),
		Config:   configure(),
	}
}

// Register a new command and it's subcommands to the CLI app.
func (cli *App) Register(commands ...*Command) {
	for _, command := range commands {
		cli.Commands[command.Name] = command
	}
}

// Run the CLI app with any given user input.
func (cli *App) Run() int {
	if len(os.Args) < 2 {
		Help.Setup(cli)

		return int(Help.Run(Help))
	}

	command, exists := cli.Commands[os.Args[1]]

	if !exists {
		log.Panicf("Unknown command given '%s'", os.Args[1])
	}

	command.Setup(cli)
	command.Parse(os.Args[2:])

	if !command.Parsed() {
		log.Panic("uh-oh! can not parse command")
	}

	return int(command.Run(command))
}
