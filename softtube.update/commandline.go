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
	for i := 0; i < len(args); i++ {
		switch args[i] {
		// We have a request for version
		case "--version":
			fmt.Println("SoftTube Update Tool")
			fmt.Println("--------------------")
			fmt.Println("softtube.update ", applicationVersion)
			os.Exit(0)
		// Invalid command line arg
		default:
			return errors.New("invalid command line arg : " + args[i])
		}
	}

	return nil
}

func invalidCommandLineArg(err error) {
	fmt.Println(err)
	fmt.Println("Usage : softtube.update [--version]")
	os.Exit(0)
}
