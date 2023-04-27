# console

A minimalist framework for CLI applications written in Go.

## Setup

### As a dependency

```bash
go get github.com/jordanbrauer/console
```

### For development

1. Clone (or fork) the repository
2. Checkout a new branch from `trunk`
3. ???
4. PROFIT!!!

Happy hacking!

## Usage

### Creating new commands

Commands are defined as variables using a struct that's properties are act as a
form of configuration. At a minimum, your command should have a `Name` and `Run`
property defined.

```go
var greet = &console.Command{
    Name: "greet",
    Description: "fubar snafu",
    Run: func(fubar *Command) ExitCode {
        println("hello, world!")

        return 0
    }}

var noop = &console.Command{
    Name: "noop",
    Description: "nada, zip, ziltch, nothing",
    Run: func() ExitCode {
        return 0
    }}
```

### Configuring the CLI

Once you have some commands, you can create a new CLI app, give it a version,
and register one or more commands to be executed by the user.

```go
cli := console.New("your version")

cli.Register(noop, greet)
```

### Running the CLI

The easiest part is running the app. It's a good idea to pass the execution
result to `os.Exit` so that the user's operating system can see the exit code
returned by the executed command.

```go
os.Exit(cli.Run())
```
