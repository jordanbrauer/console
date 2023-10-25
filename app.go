package console

import (
	"log"
	"os"
)

type App struct {
	Version  string
	Commands map[string]*Command
	Config   configuration

	header string
}

// Create a new CLI application with the given release version.
func New(release string) *App {
	return &App{
		Commands: map[string]*Command{
			help.Name:    help,
			version.Name: version,
		},
		Config:  configure(),
		Version: release,
	}
}

// Register a new command and it's subcommands to the CLI app.
func (cli *App) Register(commands ...*Command) {
	for _, command := range commands {
		cli.Commands[command.Name] = command
	}
}

// Build a string to be printed above the main help list of available commands.
func (cli *App) Splash(header func() string) {
	cli.header = header()
}

// Run the CLI app with any given user input.
func (cli *App) Run() int {
	if len(os.Args) < 2 {
		help.setup(cli)

		return int(help.Run(help))
	}

	command, exists := cli.Commands[os.Args[1]]

	if !exists {
		log.Panicf("Unknown command given '%s'", os.Args[1])
	}

	command.setup(cli)
	command.parse(os.Args[2:])

	if !command.parsed() {
		log.Panic("uh-oh! can not parse command")
	}

	return int(command.Run(command))
}
