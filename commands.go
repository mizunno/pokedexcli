package main

import (
	"fmt"
	"os"
	"errors"
	"math/rand"
)

type Command struct {
	name		string
	description string
	callback	func(*Config, ...string) error
}

func commandHelp(c *Config, args ...string) error {
	fmt.Printf("Welcome to the Pokédex!\nUsage:\n\n")

	commands := getCommands()

	for _, command := range commands {
		fmt.Printf("%v: %v\n", command.name, command.description)
	}

	fmt.Printf("\nGotta Catch 'Em All!\n")
	return nil
}

func commandExit(c *Config, args ...string) error {
	os.Exit(0)
	return nil
}

func commandMap(c *Config, args ...string) error {
	locations, err := c.client.ListLocations(nil, nil, c.nextLocationsURL)

	if err != nil {
		return err
	}

	c.nextLocationsURL = locations.Next
	c.previousLocationsURL = locations.Previous

	for _, location := range locations.Results {
		fmt.Printf("- %v\n", location.Name)
	}

	return nil
}

func commandMapb(c *Config, args ...string) error {

	if c.previousLocationsURL == nil {
		return errors.New("You are already on the first page.")
	}

	locations, err := c.client.ListLocations(nil, nil, c.previousLocationsURL)

	if err != nil {
		return err
	}

	c.nextLocationsURL = locations.Next
	c.previousLocationsURL = locations.Previous

	for _, location := range locations.Results {
		fmt.Printf("- %v\n", location.Name)
	}

	return nil
}

func commandExplore(c *Config, args ...string) error {
	if len(args) != 1{
		return errors.New("Provide a valid location name.")
	}

	location, err := c.client.GetLocation(args[0])

	if err != nil {
		return err
	}

	fmt.Printf("Exploring %v...\n", args[0])
	fmt.Println("Found Pokémon:")

	for _, pokemonEncounter := range location.PokemonEncounters {
		fmt.Printf("- %v\n", pokemonEncounter.Pokemon.Name)
	}

	return nil
}

func commandCatch(c *Config, args ...string) error {
	if len(args) != 1 {
		return errors.New("Provide a valid Pokémon name.")
	}

	pokemon, err := c.client.GetPokemon(args[0])

	if err != nil {
		return err
	}
	fmt.Printf("Throwing a Pokeball at %v...\n", pokemon.Name)

	y := float64(50 - (((pokemon.BaseExperience - 1) * 50) / 340))

	// Random float in [0.0, 100.0]
	roll := rand.Float64() * 100

	//fmt.Printf("You need: %v\nYou rolled: %v\n", y, roll)

	if roll > y {
		// You don't catch the Pokemon.
		fmt.Printf("You didn't catch %v. Try again.\n", pokemon.Name)
	} else {
		// You catch the Pokemon.
		fmt.Printf("You caught %v. Gz!\n", pokemon.Name)
		c.pokemonCaught[pokemon.Name] = pokemon
		fmt.Println(len(c.pokemonCaught))
	}

	return nil
}

func commandInspect(c *Config, args ...string) error {
	if len(args) != 1 {
		return errors.New("Provide a valid Pokémon name.")
	}

	pokemon, ok := c.pokemonCaught[args[0]]

	if !ok {
		return errors.New("You didn't catch that Pokémon.")
	}

	fmt.Printf("Name: %v\n", pokemon.Name)
	fmt.Printf("Height: %v\n", pokemon.Height)
	fmt.Printf("Weight: %v\n", pokemon.Weight)
	fmt.Printf("Stats:\n")
	for _, stat := range pokemon.Stats{
		fmt.Printf(" - %v: %v\n", stat.Stat.Name, stat.BaseStat)
	}

	fmt.Printf("Types:\n")
	for _, pokemonType := range pokemon.Types{
		fmt.Printf(" - %v\n", pokemonType.Type.Name)
	}

	return nil
}

func commandPokedex(c *Config, args ...string) error {
	if len(c.pokemonCaught) == 0 {
		return errors.New("You don't have any Pokémon.")
	}

	fmt.Println("Your Pokedex:")

	for key, _ := range c.pokemonCaught {
		fmt.Printf(" - %v\n", key)
	}

	return nil
}

func commandHistory(c *Config, args ...string) error {
	for i, cmd := range c.commandHistory {
		fmt.Printf("%v %v\n", i, cmd)
	}
	return nil
}

func getCommands() map[string]Command {
	return map[string]Command{
		"help": {
			name:		 	"help",
			description: 	"Displays a help message.",
			callback: 		commandHelp,
		},
		"exit": {
			name:			"exit",
			description: 	"Exits the pokedex.",
			callback: 		commandExit,
		},
		"map": {
			name:			"map",
			description:	"Displays the name of the next 20 location areas.",
			callback:		commandMap,
		},
		"mapb": {
			name:			"mapb",
			description:	"Displays the name of the previous 20 location areas.",
			callback:		commandMapb,
		},
		"explore": {
			name:			"explore",
			description:	"Displays all the Pokémon in a given area.",
			callback:		commandExplore,
		},
		"catch": {
			name:			"catch",
			description:	"Tries to catch the given Pokémon.",
			callback:		commandCatch,
		},
		"inspect": {
			name:			"inspect",
			description:	"Inspects the given Pokémon and returns its information.",
			callback:		commandInspect,
		},
		"pokedex": {
			name:			"pokedex",
			description:	"Shows all the Pokémon you have caught.",
			callback:		commandPokedex,
		},
		"history": {
			name:			"history",
			description:	"Prints the last commands you have run.",
			callback:		commandHistory,
		},
	}
}
