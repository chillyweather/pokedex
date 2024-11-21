package main

import (
	"bufio"
	"fmt"
	"os"
)

func runRepl() {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Pokedex > ")

		scanner.Scan()
		text := scanner.Text()

		if len(text) == 0 {
			continue
		}

		if text == "exit" {
			os.Exit(0)
		}
	}

}
