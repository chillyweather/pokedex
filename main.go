package main

import (
	"bufio"
	"fmt"
	"os"
)

func getCommand() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("pokedex >")
	command, _ := reader.ReadString('\n')
	fmt.Printf("The command %s unavailable", command)
}

func main() {
	getCommand()
}
