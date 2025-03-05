package main

import "fmt"

func commandPokedex(cfg *config, params []string) error {
	if len(cfg.pokedex) == 0 {
		fmt.Println("you have not caught any pokemon")
	} else {
		fmt.Println("Your Pokedex:")
		for key := range cfg.pokedex {
			fmt.Printf(" - %s\n", key)
		}
	}
	return nil
}
