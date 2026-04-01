package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/dusnm/slack-ips/pkg/container"
)

const (
	serve   Command = "serve"
	initdb  Command = "initdb"
	dumpcfg Command = "dump-config"
	help    Command = "help"
)

var (
	commands = []Command{
		serve,
		initdb,
		dumpcfg,
		help,
	}

	ErrNoCommandSpecified = errors.New("you must specify a command to run")
)

type (
	Command string
)

func (c Command) Usage() string {
	desc := ""
	switch c {
	case serve:
		desc = "start an HTTP server with the configured parameters"
	case initdb:
		desc = "initialize the database"
	case dumpcfg:
		desc = "dump the empty configuration file to STDOUT"
	case help:
		desc = "print this help dialog"
	}

	return fmt.Sprintf("  %s\n    %s\n", c, desc)
}

// usage prints out the list of available options
func usage() {
	name := os.Args[0]
	usage := "Copyright (C) 2026 Dušan Mitrović <dusan@dusanmitrovic.rs>\n" +
		"Licensed under the terms of the GNU AGPL v3 only\n\n" +
		"slack-ips, Easily share your bank account details with others.\n\n" +
		"Usage of %s [COMMAND] [FLAGS]\n\n" +
		"Commands:\n"

	for _, command := range commands {
		usage += command.Usage()
	}

	fmt.Printf(usage, name)
}

// get takes the first argument
// and treats it as the command name
// while performing bounds checking
func get() (Command, error) {
	if len(os.Args) == 1 {
		return "", ErrNoCommandSpecified
	}

	return Command(os.Args[1]), nil
}

func Run(c *container.Container) {
	command, err := get()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n\n", err.Error())
		usage()

		os.Exit(-1)
	}

	switch command {
	case serve:
		Serve(c)
	case initdb:
		InitDB(c)
	case dumpcfg:
		DumpCfg()
	case help:
		usage()
	default:
		fmt.Fprintf(os.Stderr, "Error: unknown command \"%s\"\n\n", command)
		usage()
	}
}
