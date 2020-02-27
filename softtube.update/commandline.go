package main

import (
	"errors"
	"fmt"
	"os"
)

func processCommandLineArgs(config *Config) error {
	// Get the command line args (without program name)
	args := os.Args[1:]

	// Loop through the arguments
	for i := 0; i < len(args); i++ {
		switch args[i] {
		// We have a log path
		case "-log":
			config.Paths.Log = args[i+1]
			i++

		// We have a database path
		case "-data":
			config.Paths.Database = args[i+1]
			i++

		// Invalid command line arg
		default:
			return errors.New("invalid command line arg : " + args[i])
		}
	}

	return nil
}

func invalidCommandLineArg(err error) {
	fmt.Println(err)
	fmt.Println("Usage : softtube.update [-log path] [-data path]")
	os.Exit(0)
}
