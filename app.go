package console

import (
	"log"
	"os"
)

type App struct {
	Version  string
	Commands map[string]*Command

	header    string
	listeners map[string][]func(command *Command)
}

// Create a new CLI application with the given release version.
func New(release string) *App {
	return &App{
		Commands: map[string]*Command{
			help.Name:    help,
			version.Name: version,
		},
		Version: release,

		listeners: map[string][]func(command *Command){
			"exit":        make([]func(command *Command), 0),
			"register":    make([]func(command *Command), 0),
			"execution":   make([]func(command *Command), 0),
			"parse.flags": make([]func(command *Command), 0),
			"parse.args":  make([]func(command *Command), 0),
		},
	}
}

// Register a new command and it's subcommands to the CLI app.
func (cli *App) Register(commands ...*Command) {
	for _, command := range commands {
		cli.Commands[command.Name] = command

		for _, registered := range cli.listeners["register"] {
			registered(command)
		}
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

		for _, executing := range cli.listeners["execution"] {
			executing(help)
		}

		code := help.Run(help)
		help.exitCode = code

		for _, exiting := range cli.listeners["exit"] {
			exiting(help)
		}

		return int(code)
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

	hasSubcommands := len(command.Commands) > 0

	if !hasSubcommands {
		for _, executing := range cli.listeners["execution"] {
			executing(command)
		}
	}

	code := command.Run(command)

	if !hasSubcommands {
		command.exitCode = code

		for _, exiting := range cli.listeners["exit"] {
			exiting(command)
		}
	}

	return int(code)
}

// Register an event listener to the CLI
//
// ### Known Events
//
// * `exit`
// * `execution`
// * `register`
// * `parse.flags`
// * `parse.args`
func (cli *App) On(event string, callback func(command *Command)) {
	if _, ok := cli.listeners[event]; !ok {
		keys := make([]string, 0, len(cli.listeners))

		for k := range cli.listeners {
			keys = append(keys, k)
		}

		log.Panicf("Unknown CLI event used: '%s'.\nPlease use one of %s\n", event, keys)
	}

	cli.listeners[event] = append(cli.listeners[event], callback)
}
