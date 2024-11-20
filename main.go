package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func getCommand() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("pokedex >")
		command, _ := reader.ReadString('\n')
		command = strings.TrimSpace(command)

		if command == "help" {
			fmt.Printf("The %s is on the way\n", command)
		}
		if command == "exit" {
			fmt.Printf("There is no %s\n", command)
		}
	}

}

func main() {
	getCommand()
}
