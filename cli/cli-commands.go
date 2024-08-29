package cli

import (
	"fmt"
	"sync"
)

func RunCliCommands(commands []string, wg *sync.WaitGroup) {
	// ./main [1] [2]
	switch commands[2] {
	case "start":
		StartCliServer(wg)
	default:
		fmt.Println("Invalid Command")
		wg.Done()
	}

}
