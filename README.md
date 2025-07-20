# BBS Sabacc

A classic 76-card Sabacc game implementation for BBS systems, written in Go using the godoors library.

## Features

- Classic Sabacc rules based on West End Games version
- 76-card deck with four suits (Sabers, Flasks, Coins, Staves) and Arcana cards
- Static Field mechanics to protect cards from Sabacc Shifts
- Terminal-based interface optimized for BBS systems
- ANSI color support with custom art
- Idle timeout handling
- Drop file support (door32.sys)
- AI opponent with configurable personalities

## Requirements

- Go 1.19 or later (for building)
- Linux-based BBS system (Talisman, Mystic, Synchronet, etc.)
- ANSI terminal capability
- door32.sys drop file

## Quick Start

### Building from Source

1. **Clone or download the source code**
2. **Build the game:**
   ```bash
   # Simple build
   go build -o sabacc .
   
   # Optimized build (recommended)
   go build -ldflags="-s -w" -o sabacc .
   
   # Or use the build script
   ./build.sh
   ```

3. **Test locally:**
   ```bash
   # Create test drop file
   echo -e "2\n8\n38400\nTest BBS\n1\nTest Player\nTestUser\n100\n90\n0\n1" > door32.sys
   
   # Run the game
   ./sabacc -path ./
   ```

### Installation on BBS

1. **Copy to your BBS doors directory:**
   ```bash
   cp sabacc /path/to/your/bbs/doors/
   chmod +x /path/to/your/bbs/doors/sabacc
   ```

2. **Configure your BBS software** (see INSTALL.md for specific examples)

## Game Rules

### Objective
Get as close to 23 points as possible without going over.

### Winning Hands
- **Pure Sabacc**: Exactly 23 points (wins Sabacc Pot)
- **Idiot's Array**: Idiot card + 2 + 3 (literal 23) - beats Pure Sabacc!

### Bomb Out Conditions
- Over 23 points
- Under -23 points  
- Exactly 0 points

### Card Values
- **Sabers, Flasks, Coins, Staves**: 1-15 points each
- **Arcana Cards**: Negative values (-1 to -17)

### Special Mechanics
- **Sabacc Shift**: Random event (rolling doubles) that shuffles all hands
- **Static Field**: Protect important cards from being shuffled during shifts
- **Hand Pot**: Won by the best hand each round
- **Sabacc Pot**: Only won with Pure Sabacc or Idiot's Array

### Player Actions
- **[D] Draw**: Take a card from the deck
- **[T] Trade**: Discard a card and draw a new one
- **[S] Stand**: Take no action
- **[F] Static Field**: Place/remove cards for protection
- **[C] Call**: End the hand (available after round 2)
- **[Q] Fold**: Give up and pay penalty

## File Structure

### Required Files
```
main.go          # Core game logic
cards.go         # Card and deck management
helpers.go       # UI and utility functions
go.mod          # Go module dependencies
```

### Optional Files
```
config.go        # Game configuration system
main_test.go     # Test suite
build.sh         # Build script
README.md        # This file
INSTALL.md       # Detailed installation guide
ansi/            # ANSI art files
├── title.ans    # Title screen
├── menu.ans     # Menu background
└── game.ans     # Game screen
```

## Development

### Running Tests
```bash
go test              # Run all tests
go test -v           # Verbose output
go test -cover       # With coverage report
```

### Building
```bash
# Development build
go build -o sabacc .

# Production build (smaller executable)
go build -ldflags="-s -w" -o sabacc .

# Using build script
./build.sh
```

### Adding Features
The codebase is modular and well-tested, making it easy to add new features:
- AI personalities are configurable in `config.go`
- Game rules are in `cards.go` and `main.go`
- UI functions are in `helpers.go`
- All major functions have tests in `main_test.go`

## BBS Configuration Examples

### Talisman BBS
```toml
[[door]]
name = "Sabacc"
command = "/path/to/doors/sabacc -path %3"
directory = "/path/to/doors/"
description = "Classic 76-Card Sabacc Game"
```

### Mystic BBS
```
Description  : Sabacc Card Game
Command Line : /path/to/doors/sabacc -path %3
Use Door32   : Yes
```

### Synchronet BBS
```ini
[Sabacc]
Name=Sabacc Card Game
Command Line=/path/to/doors/sabacc -path %g
BBS Drop File Type=Door32.sys
```

See INSTALL.md for more detailed configuration examples.

## Troubleshooting

### Common Issues

**Build fails:**
- Ensure Go 1.19+ is installed: `go version`
- Check dependencies: `go mod tidy`

**Game doesn't start:**
- Verify executable permissions: `chmod +x sabacc`
- Check drop file path and format
- Ensure ANSI terminal support

**Display issues:**
- Verify ANSI color support in terminal
- Check terminal size (80x25 recommended)
- Ensure proper drop file emulation setting

### Debug Mode
Create a test environment:
```bash
# Test drop file
echo -e "2\n8\n38400\nTest BBS\n1\nTest Player\nTestUser\n100\n90\n0\n1" > door32.sys

# Run with explicit path
./sabacc -path ./
```

## Credits

- Based on Sabacc as created by West End Games
- Uses the [godoors](https://github.com/robbiew/godoors) library by robbiew
- Inspired by the web version at [compycore/sabacc](https://github.com/compycore/sabacc)
- Star Wars universe by Lucasfilm Ltd.

## License

Released under MIT License. See LICENSE file for details.

## Contributing

Feel free to submit issues, suggestions, or pull requests to improve the game. The test suite helps ensure changes don't break existing functionality.

*"Never tell me the odds!" - Han Solo*