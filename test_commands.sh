#!/bin/bash

echo "Building the application..."
go build

echo -e "\nTesting showoff with invalid Pokemon name..."
echo "showoff abcdef" | ./pokedexcli

echo -e "\nTesting evolve with invalid Pokemon name..."
echo "evolve qwerty" | ./pokedexcli

echo -e "\nTesting describe with invalid Pokemon name..."
echo "describe xyz123" | ./pokedexcli

echo -e "\nTesting inspect with invalid Pokemon name..."
echo "inspect nonexistent" | ./pokedexcli

echo -e "\nTesting release with invalid Pokemon name..."
echo "release fakemon" | ./pokedexcli

echo -e "\nDone testing." 