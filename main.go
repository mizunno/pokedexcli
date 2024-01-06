package main

import (
	"fmt"
	"strings"
	"time"
	"internal/pokeapihandler"
	"internal/pokecache"
)

type Config struct {
	client pokeapihandler.Client
	cache pokecache.Cache
	nextLocationsURL *string
	previousLocationsURL *string
	pokemonCaught map[string]pokeapihandler.Pokemon
	commandHistory []string
}

func main() {
	commands := getCommands()

	config := Config{
		client: pokeapihandler.NewHTTPClient(
			10 * time.Second,
			pokecache.NewCache(time.Duration(5) * time.Second),
		),
		nextLocationsURL: nil,
		previousLocationsURL: nil,
		pokemonCaught: make(map[string]pokeapihandler.Pokemon),
		commandHistory: make([]string, 0),
	}

	for {
		showPrompt()

		userInput := getUserInput()
		userInput = cleanUserInput(userInput)

		args := strings.Split(userInput, " ")

		command, ok := commands[args[0]]

		if !ok {
			fmt.Printf("Command '%v' does not exist!\nType 'help' for more info.\n", userInput)
		} else {
			err := command.callback(&config, args[1:]...)

			if err != nil {
				fmt.Println(err)
			}

			config.commandHistory = append(config.commandHistory, userInput)
		}
	}
}
