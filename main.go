package main

import (
	"bufio"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"os"
	"strings"

	"github.com/mickk/pokedexcli/internal/pokeapi"
)

type config struct {
	Next     string
	Previous string
}

type cliCommand struct {
	name        string
	description string
	config      *config
	callback    func(config *config, args ...string) error
}

var (
	commands      map[string]cliCommand
	caughtPokemon map[string]pokeapi.Pokemon
)

func init() {
	setupCommands()
	caughtPokemon = make(map[string]pokeapi.Pokemon)
}

func main() {
	config := config{
		Next:     "https://pokeapi.co/api/v2/location-area?offset=0&limit=20",
		Previous: "",
	}
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		text := strings.ToLower(scanner.Text())
		textSlice := cleanInput(text)
		args := textSlice[1:]

		command, ok := commands[textSlice[0]]
		if !ok {
			fmt.Println("Unknown Command")
			continue
		}
		err := command.callback(&config, args...)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func cleanInput(text string) []string {
	text = strings.TrimSpace(text)

	result := strings.Fields(text)

	return result
}

func setupCommands() {
	commands = map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "Displays a list of location areas in the Pokemon world",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays previous list of location areas in the Pokemon world",
			callback:    commandMapB,
		},
		"explore": {
			name:        "explore",
			description: "Explore a specific location and list the pokemon that live there",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "Attempt to catch the specified pokemon and store its information in the pokedex",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "Get details on the pokemon you have caught",
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "List all the pokemon currently in your pokedex",
			callback:    commandPokedex,
		},
	}
}

func commandExit(config *config, args ...string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(config *config, args ...string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Print("Usage:\n\n")
	for _, command := range commands {
		fmt.Printf("%v: %v\n", command.name, command.description)
	}
	return nil
}

func commandMap(config *config, args ...string) error {
	if config.Next == "" {
		return errors.New("max page reached")
	}
	res, err := pokeapi.GetLocationAreas(config.Next)
	if err != nil {
		fmt.Printf("Unexpected error when calling %v\n", config.Next)
		return err
	}
	if res.Count == 0 {
		return fmt.Errorf("no results found for url: %v", config.Next)
	}

	for _, resource := range res.Results {
		fmt.Println(resource.Name)
	}
	if res.Next != nil {
		config.Next = *res.Next
	} else {
		config.Next = ""
	}
	if res.Previous != nil {
		config.Previous = *res.Previous
	} else {
		config.Previous = ""
	}

	return nil
}

func commandMapB(config *config, args ...string) error {
	if config.Previous == "" {
		return errors.New("previous page does not exist")
	}
	res, err := pokeapi.GetLocationAreas(config.Previous)
	if err != nil {
		fmt.Printf("Unexpected error when calling %v\n", config.Previous)
		return err
	}
	if res.Count == 0 {
		return fmt.Errorf("no results found for url: %v", config.Previous)
	}

	for _, resource := range res.Results {
		fmt.Println(resource.Name)
	}
	if res.Next != nil {
		config.Next = *res.Next
	} else {
		config.Next = ""
	}
	if res.Previous != nil {
		config.Previous = *res.Previous
	} else {
		config.Previous = ""
	}

	return nil
}

func commandExplore(config *config, args ...string) error {
	if len(args) < 1 {
		return fmt.Errorf("args list contains %v arguments", len(args))
	}
	url := "https://pokeapi.co/api/v2/location-area/" + args[0]
	res, err := pokeapi.GetLocationArea(url)
	if err != nil {
		fmt.Printf("Unexpected error when calling %v\n", url)
		return err
	}

	if len(res.PokemonEncounters) > 0 {
		fmt.Println("Found Pokemon:")
	}

	for _, encounters := range res.PokemonEncounters {
		fmt.Printf("- %v\n", encounters.Pokemon.Name)
	}

	return nil
}

func commandCatch(config *config, args ...string) error {
	if len(args) < 1 {
		return fmt.Errorf("args list contains %v arguments", len(args))
	}
	pokemonStr := args[0]
	url := "https://pokeapi.co/api/v2/pokemon/" + pokemonStr
	res, err := pokeapi.GetPokemonInfo(url)
	if err != nil {
		fmt.Printf("Unexpected error when calling %v\n", url)
		return err
	}
	fmt.Printf("Throwing a Pokeball at %v...\n", pokemonStr)
	catchChance := math.Max(0.1, 1.0-(float64(res.BaseExperience)/500))
	randomValue := rand.Float64()
	if randomValue < catchChance {
		fmt.Printf("%v was caught!\n", pokemonStr)
		caughtPokemon[pokemonStr] = res
	} else {
		fmt.Printf("%v escaped!\n", pokemonStr)
	}

	return nil
}

func commandInspect(config *config, args ...string) error {
	if len(args) < 1 {
		return fmt.Errorf("args list contains %v arguments", len(args))
	}
	pokemonStr := args[0]

	if pokemon, ok := caughtPokemon[pokemonStr]; ok {
		fmt.Printf("Name: %v\n", pokemon.Name)
		fmt.Printf("Height: %v\n", pokemon.Height)
		fmt.Printf("Weight: %v\n", pokemon.Weight)
		fmt.Println("Stats:")
		for _, stat := range pokemon.Stats {
			fmt.Printf("  -%v: %v\n", stat.Stat.Name, stat.BaseStat)
		}
		fmt.Println("Types:")
		for _, item := range pokemon.Types {
			fmt.Printf("  - %v\n", item.Type.Name)
		}
	} else {
		fmt.Println("you have not caught that pokemon")
	}

	return nil
}

func commandPokedex(config *config, args ...string) error {
	if len(caughtPokemon) == 0 {
		fmt.Println("You have not caught any pokemon")
		return nil
	}
	fmt.Println("Your Pokedex:")
	for name := range caughtPokemon {
		fmt.Printf(" - %v\n", name)
	}
	return nil
}
