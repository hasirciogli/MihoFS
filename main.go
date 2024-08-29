package main

import (
	"fmt"
	"os"
	"sync"

	"github.com/hasirciogli/MihoFS/cli"
)

func main() {
	command := os.Args[1]

	if command == "" {
		fmt.Println("Command is required")
		os.Exit(1)
	}

	var wg sync.WaitGroup

	switch command {
	case "base":
		cli.SendCommand("base")
	case "cli-server":
		wg.Add(1)
		go cli.RunCliCommands(os.Args, &wg)
	case "hello":
		cli.SendCommand("hello")
	default:
		fmt.Println("Unknown command")
		os.Exit(1)
	}

	wg.Wait()
}
