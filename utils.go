package main

import (
	"bufio"
	"strings"
	"os"
	"fmt"
)

func getUserInput() string {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	text := scanner.Text()
	return text
}

func cleanUserInput(userInput string) string {
	return strings.TrimSpace(userInput)
}

func showPrompt() {
	fmt.Printf("Pokedex> ")
}
