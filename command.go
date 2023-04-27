package console

import (
	"flag"
	"fmt"

	// "github.com/charmbracelet/glamour"
	"log"
	"sort"
	"strings"
)

type ExitCode uint8

type Command struct {
	Name          string
	Description   string
	Documentation string
	Run           func(command *Command) ExitCode
	Arguments     []*Argument
	Options       []*Option
	Commands      []*Command

	flags    *flag.FlagSet
	app      *App
	main     bool
	hidden   bool
	disabled bool
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

// Register a new command to the CLI app.
func (command *Command) Setup(cli *App) *Command {
	command.app = cli

	// for _, argument := range command.Arguments {
	// }

	command.flags = flag.NewFlagSet(command.Name, flag.ExitOnError)

	for _, option := range command.Options {
		command.flags.String(option.Name, "", option.Description)
	}

	return command
}

func (command *Command) Option(name string) flag.Value {
	return command.flags.Lookup(name).Value
}

func (command *Command) Parsed() bool {
	return command.flags.Parsed()
}

func (command *Command) Parse(arguments []string) error {
	return command.flags.Parse(arguments)
}

// var markdown, _ = glamour.NewTermRenderer(
// 	// detect background color and pick either the default dark or light theme
// 	glamour.WithAutoStyle(),
// 	// wrap output at specific width
// 	glamour.WithWordWrap(80),
// )

var helpCommand = func(help *Command) ExitCode {
	commands := make([]*Command, len(help.Commands))

	var index int

	for _, command := range help.Commands {
		commands[index] = command
		index++
	}

	name := "{command}"
	positionals := []string{"argument[]"}

	if "" != help.Name {
		name = help.Name
	}

	hasArguments := len(help.Arguments) > 0

	if hasArguments {
		positionals = []string{}

		for _, positional := range help.Arguments {
			positionals = append(positionals, fmt.Sprintf("{%s}", positional.Name))
		}
	}

	sort.Slice(commands, func(current, next int) bool {
		return commands[current].Name < commands[next].Name
	})

	if "" != help.Documentation {
		// out, _ := markdown.Render(help.Documentation)
		out := "MISSING"
		fmt.Print(out)
	}

	fmt.Printf("%s\n\n", help.Description)
	fmt.Println("\033[33mUsage:\033[0m")
	fmt.Printf("  %s [--flag(s)] %s\n\n", name, strings.Join(positionals, " "))

	if hasArguments {
		fmt.Println("\033[33mArguments:\033[0m")

		for _, argument := range help.Arguments {
			fmt.Printf("  \033[32m%-18s\033[0m%s\n", argument.Name, argument.Description)
		}
	}

	if len(help.Options) > 0 {
		fmt.Println("\033[33mOptions:\033[0m")
	}

	if len(help.Commands) > 0 {
		fmt.Println("\033[33mCommands:\033[0m")

		for _, command := range commands {
			fmt.Printf("  \033[32m%-18s\033[0m%s\n", command.Name, command.Description)
		}
	}

	return 0
}

var demandCommand = func(command *Command) ExitCode {
	exit := runSubCommand(command)

	if 0 == exit {
		return 0
	}

	fmt.Println()
	log.Printf("Please execute a subcommand!\n")

	return exit
}

var demandHelpCommand = func(command *Command) ExitCode {
	helpCommand(command)

	return demandCommand(command)
}

var runHelpCommand = func(command *Command) ExitCode {
	exit := runSubCommand(command)

	if 0 == exit {
		return 0
	}

	return helpCommand(command)
}

func runSubCommand(command *Command) ExitCode {
	if len(command.flags.Args()) <= 0 {
		return 1
	}

	for _, subcommand := range command.Commands {
		if subcommand.Name != command.flags.Arg(0) {
			continue
		}

		subcommand.Setup(command.app)

		newArgs := command.flags.Args()[1:]

		if len(newArgs) > 0 {
			subcommand.Parse(newArgs)
		}

		if len(newArgs) != len(subcommand.Arguments) {
			subcommand.Name = fmt.Sprintf("%s %s", command.Name, subcommand.Name)

			return helpCommand(subcommand)
		}

		return subcommand.Run(subcommand)
	}

	return 1
}
