package main

import (
	"fmt"
)

func commandHelp(cfg *config, params []string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("--------------------")
	for _, cmd := range getCommands() {
		fmt.Printf("%s | %s \n", cmd.name, cmd.description)
	}
	fmt.Println("--------------------")
	return nil
}
