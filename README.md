```
 ██████╗ ██████╗ ███████╗    ███████╗ █████╗ ██████╗  █████╗  ██████╗ ██████╗
 ██╔══██╗██╔══██╗██╔════╝    ██╔════╝██╔══██╗██╔══██╗██╔══██╗██╔════╝██╔════╝
 ██████╔╝██████╔╝███████╗    ███████╗███████║██████╔╝███████║██║     ██║     
 ██╔══██╗██╔══██╗╚════██║    ╚════██║██╔══██║██╔══██╗██╔══██║██║     ██║     
 ██████╔╝██████╔╝███████║    ███████║██║  ██║██████╔╝██║  ██║╚██████╗╚██████╗
 ╚═════╝ ╚═════╝ ╚══════╝    ╚══════╝╚═╝  ╚═╝╚═════╝ ╚═╝  ╚═╝ ╚═════╝ ╚═════╝
                                                                              
    ♠ ♦ ♣ ♥    Classic 77-Card Sabacc for Bulletin Board Systems    ♠ ♦ ♣ ♥   
                                                                              
                 [█] Authentic 1989 West End Games Rules [█]
                 [█] ANSI Color Terminal Support         [█] 
                 [█] Linux-based BBS Compatibility       [█]
                 [█] 4 AI Players and 1 Human Player     [█]
                 [█] Player Statistics & Rankings (TODO) [█]
```

## ⚠️ DEVELOPMENT STATUS ⚠️

**🚧 This game is currently under active development! 🚧**

While the core Sabacc gameplay mechanics are implemented and functional, there is still significant work to be done on:

- **User Interface & Experience**: Screen layouts, ANSI art integration, menu systems
- **Gameplay Features**: Turns, rules, enhanced AI strategies  
- **Player Statistics**: Persistent stat tracking, leaderboards, achievement system
- **Configuration**: Advanced game settings, difficulty levels, economy balancing
- **Polish & Testing**: Bug fixes, performance optimization, extensive BBS testing

**Current Status**: *Under Development* - Playable but expect ugliness, ongoing changes and improvements!

Feedback, bug reports, and contributions are welcome as we work toward a polished 1.0 release.

---

# Overview

**"Never tell me the odds!"** - Step into the smoky cantinas of the Outer Rim and experience the most notorious card game in the galaxy! 

Sabacc is the legendary gambling game that Han Solo played to win the *Millennium Falcon* from Lando Calrissian in the Cloud City. This authentic implementation brings the 1989 West End Games Classic Sabacc rules to your BBS, complete with all the excitement, strategy, and unpredictable Sabacc Shifts that make this game a galactic favorite.

## ⭐ What is Sabacc? ⭐

In the Star Wars universe, Sabacc is a high-stakes card game played in cantinas, gambling halls, and starship lounges across the galaxy. Players risk their credits (and sometimes their ships!) trying to achieve the perfect hand of exactly **23 points** - known as "Pure Sabacc" - while avoiding the dreaded "bomb out" that could cost them everything.

But beware the **Sabacc Shift**! At any moment, the mystical interferometer field can scramble all cards not protected in the Static Field, turning victory into defeat in the roll of the dice. Only the most cunning players survive long enough to master this game of skill, luck, and nerves.

## 🎮 Key Features 🎮

### **Authentic 1989 Classic Rules**
- **76-card deck** with four suits (Sabers, Flasks, Coins, Staves) plus 16 Arcana cards
- **4-phase turn structure**: Betting → Shift Roll → Call → Draw
- **Two pot system**: Hand Pot (best hand ≤23) and Sabacc Pot (Pure Sabacc/Idiot's Array)
- **Static Field mechanics** to protect valuable cards from Sabacc Shifts
- **Authentic hand rankings** where Idiot's Array beats Pure Sabacc

### **Immersive BBS Experience**
- **ANSI color terminal** support with classic BBS aesthetics
- **CP437 character set** for authentic retro terminal compatibility
- **Persistent game state** with proper screen layouts and live updates
- **Custom ANSI art** support for title screens and game backgrounds

### **Intelligent AI Opponents**
Meet your fellow cantina patrons - each with their own personality and playing style:

- **🛸 Phoo_ja** - *Rodian smuggler* with aggressive betting tactics
- **⚡ Rsh-Taac** - *Twi'lek trader* who plays conservatively and uses the Static Field wisely  
- **🎯 Soladi** - *Human bounty hunter* with calculated risk-taking strategies
- **🌟 Ky'Ola** - *Corellian pilot* known for unpredictable bluffs and bold moves

Each AI analyzes visible cards, manages their Static Field strategically, and adapts their betting based on hand strength and opponent behavior - just like seasoned cantina gamblers!

### **Strategic Gameplay**
- **Dynamic betting system** with check, call, raise, and fold options every round
- **Static Field management** - protect your best cards from Sabacc Shifts
- **Hand calling** - end the game early when you're confident in your hand
- **Multiple victory conditions** - Pure Sabacc, Idiot's Array, or simply the best hand
- **Risk vs. reward** - drawing cards can improve your hand or cause you to bomb out

### **Classic BBS Integration**
- **Door32.sys support** for seamless BBS integration
- **Multi-BBS compatibility** (Talisman, Mystic, Synchronet, ENiGMA½, WWIV)
- **Configurable game settings** for different BBS economies
- **Player statistics tracking** (coming soon)
- **No external dependencies** - pure Go implementation

## 🍻 Enter the Cantina 🍻

Whether you're a seasoned smuggler looking to win your next ship or a rookie pilot trying to make some quick credits, BBS Sabacc delivers the authentic Star Wars gambling experience. Feel the tension as the dice roll for a Sabacc Shift, celebrate when you hit Pure Sabacc, and learn to read your opponents' tells as you battle for galactic supremacy.

*"In my experience, there's no such thing as luck."* - Obi-Wan Kenobi

**May the Force be with you... you're going to need it!**

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
   # Verify creation with:
   go run cmd/build-cards/main.go test
   
   # Optional: Generate ANSI preview of card design
   go run cmd/build-cards/main.go preview
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

5. **Copy to your BBS doors directory:**
   ```bash
   # Copy the game executable
   cp sabacc /path/to/your/bbs/doors/
   chmod +x /path/to/your/bbs/doors/sabacc
   
   # Copy the card database (REQUIRED)
   cp sabacc_cards.bin /path/to/your/bbs/doors/
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

### Talisman BBS

Add to your `doors.toml` configuration:

```toml
[[door]]
name = "Sabacc"
command = "/path/to/doors/sabacc -path %3"
directory = "/path/to/doors/"
type = "door"
description = "Classic 76-Card Sabacc Game"
access_level = 10
time_limit = 30
door32 = true
```

### Mystic BBS

Configure in your door setup:

```
Description  : Sabacc Card Game
Access Level : 10
Command Line : /path/to/doors/sabacc -path %3
Optional Data: 
Use Door32   : Yes
Time Limit   : 30
Native Exec  : No
Modify User  : No
```

### Synchronet BBS

Add to your `xtrn.ini` file:

```ini
[Sabacc]
Name=Sabacc Card Game
Command Line=/path/to/doors/sabacc -path %g
Multiple Concurrent Users=No
Intercept Standard I/O=Yes
Native (32-bit) Executable=Yes
Use Shell to Execute=No
Modify User Data=No
BBS Drop File Type=Door32.sys
Place Drop File In=Node Directory
Time Options=Deduct from Time Online
Access Requirements=
```

Alternatively, you can configure it through SCFG (Synchronet Configuration):
1. Run `scfg` as sysop
2. Go to **External Programs** → **Doors (Externals)**
3. Add new door with these settings:
   - **Name**: Sabacc Card Game
   - **Internal Code**: SABACC
   - **Command Line**: `/path/to/doors/sabacc -path %g`
   - **Clean-up Command Line**: (leave blank)
   - **Execution Cost**: 0
   - **Access Requirements**: (set as needed)
   - **Execution Requirements**: (leave blank)
   - **Multiple Concurrent Users**: Yes
   - **Intercept Standard I/O**: Yes
   - **Native (32-bit) Executable**: Yes
   - **Use Shell to Execute**: No
   - **Modify User Data**: No
   - **Execute on Event**: (leave blank)
   - **BBS Drop File Type**: Door32.sys
   - **Place Drop File In**: Node Directory
   - **Time Options**: Deduct from Time Online

### ENiGMA½ BBS

Add to your `menu.hjson` configuration:

```hjson
sabacc: {
    desc: Sabacc Card Game
    module: door
    config: {
        name: Sabacc
        dropFileType: door32.sys
        cmd: /path/to/doors/sabacc
        args: [
            "-path", "{dropFilePath}"
        ]
        cwd: /path/to/doors/
    }
}
```

### WWIV BBS

Add to your `chains.ini` file:

```ini
[Chain]
Description=Sabacc Card Game
Filename=/path/to/doors/sabacc -path %1
Local_only=false
Multi_user=false
Ansir=true
Emulation=door32.sys
```

## Directory Structure

Organize your installation as follows:

```
/path/to/doors/sabacc/
├── sabacc              # Main executable (required)
├── sabacc_cards.bin    # Card database (required)
├── sabacc.conf         # Configuration file (auto-generated)
└── ansi/
    └── portraits.ans   # AI player portraits (optional)

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

#### Portrait Selection:
- Randomly selects different portraits for each AI player
- Ensures each AI player gets a unique portrait (when possible)
- Selection varies each game session

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


## Future Enhancements

The following features are planned for future versions:

### Configuration File Support
- Adjustable ante amounts and credit limits
- Configurable timeout settings  
- Custom AI personality settings

### Enhanced Features
- Player statistics tracking 
- Tournament mode support?
- Multiple AI difficulty levels
- Custom ANSI art integration
- Multiplayer lobby!
