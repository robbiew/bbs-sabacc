# BBS Sabacc

A classic 76-card Sabacc game implementation for BBS systems, written in Go using the godoors library.

## Features

- Classic Sabacc rules based on West End Games version
- 76-card deck with four suits (Sabers, Flasks, Coins, Staves) and Arcana cards
- Static Field mechanics to protect cards from Sabacc Shifts
- Terminal-based interface optimized for BBS systems
- ANSI color support
- Idle timeout handling
- Drop file support (door32.sys)

## Requirements

- Go 1.19 or later
- Linux-based BBS system (Talisman, Mystic, Synchronet, etc.)
- ANSI terminal capability
- door32.sys drop file

## Installation

1. Clone or download the source code
2. Run `go mod tidy` to download dependencies
3. Build with `go build -o sabacc`
4. Copy the executable to your BBS doors directory
5. Configure your BBS to launch the game with the drop file path

## Usage

```bash
./sabacc -path /path/to/dropfile/directory/
```

The game expects the drop file directory path with a trailing slash.

## Game Rules

### Objective
Get as close to 23 points as possible without going over.

### Winning Hands
- **Pure Sabacc**: Exactly 23 points
- **Idiot's Array**: Idiot card + 2 + 3 (literal 23) - beats Pure Sabacc!

### Bomb Out Conditions
- Over 23 points
- Under -23 points  
- Exactly 0 points

### Card Values
- **Sabers, Flasks, Coins, Staves**: 1-15 points each
- **Arcana Cards**: Negative values (-1 to -17)

### Special Mechanics
- **Sabacc Shift**: Random event that shuffles all hands
- **Static Field**: Protect cards from being shuffled during shifts
- **Hand Pot**: Won by the best hand each round
- **Sabacc Pot**: Only won with Pure Sabacc or Idiot's Array

### Player Actions
- **Draw**: Take a card from the deck
- **Trade**: Discard a card and draw a new one
- **Stand**: Take no action
- **Static Field**: Place/remove cards for protection
- **Call**: End the hand (available after round 2)
- **Fold**: Give up and pay penalty

## BBS Configuration

### Door Configuration Example (Talisman)
```
[Door]
Name=Sabacc
Command=/path/to/sabacc -path %3
Directory=/path/to/sabacc/
Type=Door
```

### Mystic Configuration Example
```
Description  : Sabacc Card Game
Access Level : 10
Command Line : /path/to/sabacc -path %3
Optional Data: 
```

## File Structure

```
bbs-sabacc/
├── main.go           # Main game logic and menu system
├── cards.go          # Card structures and deck management
├── helpers.go        # Utility functions and UI helpers
├── go.mod           # Go module dependencies
├── README.md        # This file
└── ansi/            # ANSI art files (create this directory)
    ├── title.ans    # Title screen art
    ├── menu.ans     # Menu background art
    └── game.ans     # Game screen art
```

## ANSI Art

The game is designed to use ANSI art files for enhanced visual appeal. Create an `ansi/` directory and add:

- `title.ans` - Title screen artwork
- `menu.ans` - Menu background
- `game.ans` - Game screen background

These files are embedded in the executable but currently show placeholder ASCII art.

## Development Notes

This implementation is based on:
- The Classic Sabacc rules from West End Games (Crisis on Cloud City, 1989)
- The web version at github.com/compycore/sabacc for game logic reference
- The godoors library for BBS-specific terminal handling

### Future Enhancements
- Save/load game statistics
- Multiple computer opponents with different AI personalities
- Sound effects (terminal bell integration)
- More sophisticated ANSI animations
- Network multiplayer support
- Different Sabacc variants (Force Sabacc, Corellian Spike mode)

## Contributing

This is a community project. Feel free to submit issues, suggestions, or pull requests to improve the game.

## License

Released under MIT License. Feel free to use and modify for your BBS.

## Credits

- Based on Sabacc as created by West End Games
- Uses the godoors library by robbiew
- Inspired by the web version by compycore
- Star Wars universe by Lucasfilm Ltd.

*"Never tell me the odds!"*