package console

import (
	"fmt"
	"log"
	"strings"
)

var help = &Command{
	Name:        "help",
	Description: "Display this message or learn about other commands",
	Run: func(command *Command) ExitCode {
		argv := command.flags.Args()
		argc := len(argv)

		if argc == 0 {
			if "" != strings.TrimSpace(command.app.header) {
				fmt.Print(command.app.header)
			}

			commands := make([]*Command, len(command.app.Commands))
			var index int

			for _, subcommand := range command.app.Commands {
				commands[index] = subcommand
				index++
			}

			return helpCommand(&Command{
				Commands: commands,
			})
		}

		root := argv[0]
		parent, exists := command.app.Commands[root]

		if !exists {
			log.Panicf("Unknown command given '%s'", root)
		}

		if argc == 1 {
			return helpCommand(parent)
		}

		target := argv[1]

		for _, subcommand := range parent.Commands {
			if target != subcommand.Name {
				continue
			}

			subcommand.Name = fmt.Sprintf("%s %s", parent.Name, target)

			return helpCommand(subcommand)
		}

		return 0
	}}
