```                                                                              
BBS SABACC - Classic 77-Card Sabacc for Bulletin Board Systems  
                                                                              
[█] Authentic 1989 West End Games Rules [█]
[█] ANSI Color Terminal Support         [█] 
[█] Linux-based BBS Compatibility       [█]
[█] 4 AI Players and 1 Human Player     [█]
```

## ⚠️ DEVELOPMENT STATUS ⚠️

**🚧 This game is currently under active development! 🚧**

While the core Sabacc gameplay mechanics are implemented and functional, there is still significant work to be done on:

- **User Interface & Experience**: Screen layouts, ANSI art integration, menu systems
- **Gameplay Features**: Turns, rules, enhanced AI strategies  
- **Player Statistics**: Persistent stat tracking, leaderboards, achievement system
- **Configuration**: Advanced game settings, difficulty levels, economy balancing
- **Polish & Testing**: Bug fixes, performance optimization, extensive BBS testing

---

# Overview

**"Never tell me the odds!"** - Step into the smoky cantinas of the Outer Rim and experience the most notorious card game in the galaxy! 

Sabacc is the legendary gambling game that Han Solo played to win the *Millennium Falcon* from Lando Calrissian in the Cloud City. This authentic implementation brings the 1989 West End Games Classic Sabacc rules to your BBS, complete with all the excitement, strategy, and unpredictable Sabacc Shifts that make this game a galactic favorite.

## ⭐ What is Sabacc? ⭐

In the Star Wars universe, Sabacc is a high-stakes card game played in cantinas, gambling halls, and starship lounges across the galaxy. Players risk their credits (and sometimes their ships!) trying to achieve the perfect hand of exactly **23 points** - known as "Pure Sabacc" - while avoiding the dreaded "bomb out" that could cost them everything.

But beware the **Sabacc Shift**! At any moment, the mystical interferometer field can scramble all cards not protected in the Static Field, turning victory into defeat in the roll of the dice. Only the most cunning players survive long enough to master this game of skill, luck, and nerves.

---


# BBS Sabacc Installation Guide

This guide will help you install and configure BBS Sabacc on your BBS system.

## Prerequisites

- Linux-based BBS system (Talisman, Mystic, Synchronet, ENiGMA½, WWIV)
- Go 1.19 or later (for building from source)
- ANSI terminal support
- door32.sys drop file support

## Installation Methods

### Method 1: Build from Source (Recommended)

1. **Download or clone the source code**
2. **Navigate to the source directory:**
   ```bash
   cd bbs-sabacc
   ```
3. **Build the card database:**
   ```bash
   # Generate the card database (required first step)
   go run cmd/build-cards/main.go
   
   # This creates sabacc_cards.bin containing the 77-card Sabacc deck
   
   ```

4. **Build the game:**
   ```bash
   # Quick build
   go build -o sabacc .
   
   # Optimized build (smaller executable)
   go build -ldflags="-s -w" -o sabacc .
   
   # Or use the build script
   chmod +x build.sh
   ./build.sh
   ```


### Method 2: Pre-built Binary (if available)

1. Download the appropriate binary for your architecture
2. Extract and copy to your doors directory:
   ```bash
   cp sabacc /path/to/your/bbs/doors/
   chmod +x /path/to/your/bbs/doors/sabacc
   ```

## BBS Configuration

"-path" is the only required parameter. Only Door32.sys is supported.

```bash
./sabacc -path /path/to/drop_file_dir/
```

## Directory Structure

Organize your installation as follows:

```
/path/to/doors/sabacc/
├── sabacc              # Main executable (required)
├── sabacc_cards.bin    # Card database (required)
├── sabacc.conf         # Configuration file (auto-generated, editable)
└── ansi/
    └── portraits.ans   # AI player portraits (make your own!)

```

## Configuration

### BBS Sabacc Portrait System

The BBS Sabacc game features a simple ANSI-based portrait system for AI players.

#### Portrait Specifications

- **Dimensions**: Exactly 9 columns × 6 rows
- **Format**: Single stacked ANSI (.ans) file - `ansi/portraits.ans`
- **Layout**: Portraits stacked vertically in a single file (6 rows each)
- **Default**: Included `ansi/portraits.ans` with character portraits
- **Randomization**: Different portraits selected each game session

---

### Game Settings (sabacc.conf)

The game automatically creates a `sabacc.conf` file on first run with these default settings:

```json
{
  "min_ante": 10,
  "max_ante": 100,
  "max_bet": 500,
  "starting_credits": 1000,
  "shift_probability": 6,
  "min_rounds_to_call": 2,
  "idle_timeout_seconds": 300,
  "enable_statistics": true,
  "ai_personality": "balanced"
}
```

You can edit this file to customize the game behavior:

- **min_ante**: Minimum ante required (both pots)
- **max_ante**: Maximum ante allowed  
- **max_bet**: Maximum betting limit
- **starting_credits**: Credits each player starts with
- **shift_probability**: Dice roll probability (1 in X chance)
- **min_rounds_to_call**: Minimum rounds before calling allowed
- **idle_timeout_seconds**: Player idle timeout in seconds
- **enable_statistics**: Track player statistics
- **ai_personality**: AI behavior ("conservative", "balanced", "aggressive")


## TODO
- Game is over after 1 round. Fix this.
- Fix duplication of AI Player Portraits
- Create Satr Wars Name Generator for AI Players
- Various UI/layout issues
- Better ANSI art for key moments (title, portraits, etc.)
- Congifurable bet amounts
- "All in" moments, drop keys to the Millenium Falcon, etc.
- Leaderboard
- Persistant winnings/losses (reset every X days)
- Get a loan from Jabba



