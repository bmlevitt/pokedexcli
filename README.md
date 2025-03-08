# PokédexCLI

A feature-rich, interactive Pokémon encyclopedia right in your terminal! PokédexCLI is a command-line application that lets you explore the Pokémon world, catch your favorite Pokémon, and build your collection.

<p align="center">
  <img src="https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/other/official-artwork/25.png" alt="Pikachu" width="200">
</p>

## Features

- **Explore Pokémon Locations**: Navigate through different regions and discover where Pokémon can be found
- **Catch Pokémon**: Try your luck at catching wild Pokémon (based on their actual capture rates)
- **Manage Your Pokédex**: Keep track of all your caught Pokémon
- **Pokémon Details**: View detailed information about Pokémon in your collection
- **Evolution**: Evolve your Pokémon to their next forms
- **Pokémon Information**: Get interesting Pokédex entries and descriptions
- **Show Off**: Display your Pokémon's moves in action
- **Persistence**: Your Pokédex is automatically saved between sessions, so you can continue where you left off

## Installation

### Prerequisites

- Go 1.16 or higher

### Building from Source

1. Clone the repository:
   ```
   git clone https://github.com/bmlevitt/pokedexcli.git
   cd pokedexcli
   ```

2. Build the application:
   ```
   go build
   ```

3. Run the application:
   ```
   ./pokedexcli
   ```

## Usage

After starting the application, you'll be presented with a command prompt. Here's a list of available commands:

- `help`: Display a list of all available commands
- `map`: Explore the first page of map locations
- `next`: Navigate to the next page of map locations
- `prev`: Navigate to the previous page of map locations
- `explore [location number]`: List Pokémon that can be found at a specific location
- `catch [pokemon]`: Try to catch a specific Pokémon
- `inspect [pokemon]`: View details about a Pokémon in your collection
- `pokedex`: List all Pokémon in your collection
- `release [pokemon]`: Remove a Pokémon from your collection
- `showoff [pokemon]`: Display one of your Pokémon's moves
- `describe [pokemon]`: Display information and Pokédex entries for a Pokémon
- `evolve [pokemon]`: Evolve a Pokémon from your collection to its next form
- `save`: Manually save your current Pokédex to a file
- `reset`: Clear your Pokédex and start fresh
- `autosave [on/off]`: Enable or disable automatic saving
- `saveinterval [number]`: Set how many changes before auto-saving
- `exit`: Exit the application (automatically saves your Pokédex)

### Example Usage

```
$ ./pokedexcli
Pokedex > map
1. canalave-city-area
2. eterna-city-area
3. pastoria-city-area
4. sunyshore-city-area
5. sinnoh-pokemon-league-area
...

Pokedex > explore 1
Exploring eterna-city-area...
Found Pokemon:
1. drifblim
2. drifloon
3. duskull
4. gastly
5. gengar

Pokedex > catch gastly
Throwing a Pokeball at gastly...
gastly was caught!

Pokedex > inspect gastly
Name: gastly
Height: 13
Weight: 1
Types: ghost, poison
Stats:
- HP: 30
- Attack: 35
- Defense: 30
- Special-Attack: 100
- Special-Defense: 35
- Speed: 80

Pokedex > describe gastly
Gastly, the Gas Pokémon
- With its gas-like body, it can sneak into any place it desires. However, it can be blown away by wind. (From Pokémon Ultra Sun)

Pokedex > evolve gastly
Evolving gastly into haunter...
Congratulations! Your Gastly evolved into Haunter!
```

## Data Persistence

PokédexCLI automatically saves your Pokédex and location data between sessions. This means you can close the application and return later to continue where you left off.

### Automatic Saving

By default, your Pokédex is automatically saved after every change (catching, releasing, or evolving a Pokémon). You can:

- Toggle auto-save on/off with the `autosave` command
- Change how frequently auto-saves occur with the `saveinterval` command
- Manually save at any time with the `save` command
- Reset your Pokédex to start fresh with the `reset` command

Your data is saved to a hidden file in your home directory, so it persists even if you update the application.

## Caching System

PokédexCLI includes a built-in caching system to minimize API calls to the PokeAPI server. Each API response is cached for one hour by default, improving performance and reducing load on the API.

## Credits

- Pokémon data provided by [PokeAPI](https://pokeapi.co/)
- Initial project provided by [Boot.dev](https://boot.dev/)
- Application developed by Ben Levitt

## License

This project is open source and available under the [MIT License](LICENSE). 