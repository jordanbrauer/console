# console

A minimalist framework for beautiful CLI applications written in Go.

<img width="818" alt="image" src="https://user-images.githubusercontent.com/18744334/235200395-361b95cc-6f17-44a4-a6d0-eac002eb1efe.png">

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

#### Adding sub commands

Commands are recursive. That is to say that they can have commands of their own.
To do this, simply populate the `Commands` property with one or more `Command`
pointers, instead of defining a `Run` function.

```go
var parent = &console.Command{
    Name: "foo",
    Description: "i am the parent",
    Commands: []*Command{child},
}

var child = &console.Command{
    Name: "bar",
    Description: "i am the child",
    Run: func() ExitCode {
        return 0
    }}
```

which would be used like so

```bash
go run main.go foo bar
```

### Configuring the CLI

Once you have some commands, you can create a new CLI app, give it a version,
and register one or more commands to be executed by the user.

```go
cli := console.New("your version")

cli.Register(noop, greet)

cli.Splash(func () string {
    // logo, author(s), version, licensing, whatever you want...
})
```

### Running the CLI

The easiest part is running the app. It's a good idea to pass the execution
result to `os.Exit` so that the user's operating system can see the exit code
returned by the executed command.

```go
os.Exit(cli.Run())
```
