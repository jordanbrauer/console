package console

import "fmt"

var Version = &Command{
	Name:        "version",
	Description: "Show the currently installed version of Mist",
	Run: func(command *Command) ExitCode {
		fmt.Println(command.app.Version)

		return 0
	},
}
