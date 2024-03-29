package console

import (
	"flag"
	"fmt"
	"time"

	// "github.com/charmbracelet/glamour"
	"log"
	"sort"
	"strings"
)

type Bits uint8

func (bits Bits) Set(flag Bits) Bits    { return bits | flag }
func (bits Bits) Clear(flag Bits) Bits  { return bits &^ flag }
func (bits Bits) Toggle(flag Bits) Bits { return bits ^ flag }
func (bits Bits) Has(flag Bits) bool    { return bits&flag != 0 }

type ExitCode uint8

type Command struct {
	Name          string
	Description   string
	Documentation string
	Hidden        bool
	Run           func(command *Command) ExitCode
	Arguments     []*Argument
	Options       []*Option
	Commands      []*Command

	flags    *flag.FlagSet
	app      *App
	main     bool
	disabled bool
	exitCode ExitCode
}

// Register a new command to the CLI app.
func (command *Command) setup(cli *App) *Command {
	command.app = cli
	command.flags = flag.NewFlagSet(command.Name, flag.ExitOnError)

	for _, option := range command.Options {
		switch option.Type {
		case "array":
			if option.Default == nil {
				option.Default = []string{}
			}

			var value array

			command.flags.Var(&value, option.Name, option.Description)

			option.value = &value

		case "int":
			if option.Default == nil {
				option.Default = 0
			}

			option.value = command.flags.Int(option.Name, option.Default.(int), option.Description)
		case "bool", "boolean":
			if option.Default == nil {
				option.Default = false
			}

			option.value = command.flags.Bool(option.Name, option.Default.(bool), option.Description)
		case "string":
			if option.Default == nil {
				option.Default = ""
			}

			option.value = command.flags.String(option.Name, option.Default.(string), option.Description)
		case "duration":
			if option.Default == nil {
				option.Default = time.Duration(0)
			}

			option.value = command.flags.Duration(option.Name, option.Default.(time.Duration), option.Description)
		default:
			if option.Default == nil {
				option.Default = ""
			}

			option.value = command.flags.String(option.Name, option.Default.(string), option.Description)
		}
	}

	return command
}

// Determine if the command argument & option inputs have been parsed
func (command *Command) parsed() bool {
	return command.flags.Parsed()
}

// parse the given command input
func (command *Command) parse(arguments []string) error {
	err := command.flags.Parse(arguments)

	if err != nil {
		for _, parsing := range command.app.listeners["parse.flags"] {
			parsing(command)
		}

		for _, parsing := range command.app.listeners["parse.args"] {
			parsing(command)
		}
	}

	return err
}

func (command *Command) option(name string) any {
	for _, option := range command.Options {
		if name == option.Name {
			return option.value
		}
	}

	return nil
}

// Read an argument value by it's name.
func (command *Command) Argument(name string) string {
	var position int = -1

	for index, argument := range command.Arguments {
		if name == argument.Name {
			position = index

			break
		}
	}

	if position == -1 {
		log.Fatalf("Unknown argument: %s\n", name)
	}

	return command.flags.Arg(position)
}

func (command *Command) ArgumentRest(name string) []string {
	var position int = -1

	for index, argument := range command.Arguments {
		if name == argument.Name {
			position = index

			break
		}
	}

	if position == -1 {
		log.Fatalf("Unknown argument: %s\n", name)
	}

	args := command.flags.Args()

	return args[position:]

}

func (command *Command) ArgumentArray() []string {
	return command.flags.Args()
}

func (command *Command) OptionStringArray(name string) []string {
	value, ok := command.option(name).(*array)

	if !ok {
		return []string{}
	}

	values := make([]string, len(*value))

	for index, value := range *value {
		values[index] = value
	}

	return values
}

func (command *Command) OptionBool(name string) bool {
	value, ok := command.option(name).(*bool)

	if !ok {
		return false
	}

	return *value
}

func (command *Command) OptionInt(name string) int {
	value, ok := command.option(name).(*int)

	if !ok {
		return 0
	}

	return *value
}

func (command *Command) OptionString(name string) string {
	value, ok := command.option(name).(*string)

	if !ok {
		return ""
	}

	return *value
}

func (command *Command) OptionDuration(name string) time.Duration {
	value, ok := command.option(name).(*time.Duration)

	if !ok {
		return time.Duration(0)
	}

	return *value
}

func (command *Command) ExitCode() ExitCode {
	return command.exitCode
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

	hasOptions := len(help.Options) > 0

	if hasOptions {
		if hasArguments {
			fmt.Println()
		}

		fmt.Println("\033[33mOptions:\033[0m")

		for _, option := range help.Options {
			fmt.Printf("  \033[32m%-18s\033[0m%s\n", option.Name, option.Description)
		}
	}

	if len(help.Commands) > 0 {
		if hasArguments || hasOptions {
			fmt.Println()
		}

		fmt.Println("\033[33mCommands:\033[0m")

		for _, command := range commands {
			if command.Hidden {
				continue
			}

			fmt.Printf("  \033[32m%-18s\033[0m%s\n", command.Name, command.Description)
		}
	}

	return 0
}

var DemandCommand = func(command *Command) ExitCode {
	exit := runSubCommand(command)

	if 0 == exit {
		return 0
	}

	fmt.Println()
	log.Printf("Please execute a subcommand!\n")

	return exit
}

var HelpAndDemandCommand = func(command *Command) ExitCode {
	helpCommand(command)

	return DemandCommand(command)
}

var RunOrHelpCommand = func(command *Command) ExitCode {
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

		subcommand.setup(command.app)

		newArgs := command.flags.Args()[1:]

		if len(newArgs) > 0 {
			subcommand.parse(newArgs)
		}

		subcommand.Name = fmt.Sprintf("%s %s", command.Name, subcommand.Name)
		requiredArguments := []*Argument{}

		for _, a := range subcommand.Arguments {
			if !a.Optional() {
				requiredArguments = append(requiredArguments, a)
			}
		}

		hasSubcommands := len(subcommand.Commands) > 0

		if !hasSubcommands {
			for _, executing := range subcommand.app.listeners["execution"] {
				executing(subcommand)
			}
		}

		var code ExitCode

		if len(subcommand.flags.Args()) < len(requiredArguments) && !hasSubcommands {
			// TODO: error handling/validation message
			code = helpCommand(subcommand)
		} else {
			code = subcommand.Run(subcommand)
		}

		if !hasSubcommands {
			subcommand.exitCode = code

			for _, exiting := range subcommand.app.listeners["exit"] {
				exiting(subcommand)
			}
		}

		return code
	}

	return 1
}
