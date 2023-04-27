package console

import (
	"encoding/base64"
	"fmt"
	"log"
)

var Help = &Command{
	Name:        "help",
	Description: "Display this message or learn about other commands",
	Run: func(command *Command) ExitCode {
		argv := command.flags.Args()
		argc := len(argv)

		if argc == 0 {
			logo, _ := base64.StdEncoding.DecodeString("ICAgICAgICAgICAgICAgCiBfIF8gIG8gIF8gXylfIAopICkgKSAoICggIChfICAKICAgICAgICBfKSAgICAg")

			fmt.Print("\033[3mWelcome to,\033[0m\n")
			fmt.Println(string(logo))
			fmt.Printf("\nVersion \033[34m%s\033[0m, \033[3mmade with \033[31m<3\033[0m \033[3mby jorb\033[0m", command.app.Version)

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
