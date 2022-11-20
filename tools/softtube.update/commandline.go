package main

import (
	"errors"
	"fmt"
	"os"
)

func processCommandLineArgs() error {
	// Get the command line args (without program name)
	args := os.Args[1:]

	// Loop through the arguments
	for _, arg := range args {
		switch arg {
		case "--version":
			// We have a request for version
			_, _ = fmt.Fprintln(os.Stdout, "SoftTube Update Tool")
			_, _ = fmt.Fprintln(os.Stdout, "--------------------")
			_, _ = fmt.Fprintln(os.Stdout, "softtube.update ", applicationVersion)

			os.Exit(0)
		default:
			// Invalid command line arg
			return errors.New("invalid command line arg : " + arg)
		}
	}

	return nil
}

func invalidCommandLineArg(err error) {
	_, _ = fmt.Fprintln(os.Stderr, err)
	_, _ = fmt.Fprintln(os.Stdout, "Usage : softtube.update [--version]")
	os.Exit(0)
}
