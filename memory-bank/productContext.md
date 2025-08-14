# BBS Sabacc - Product Context & Technical Specifications

## Product Identity
**Name**: BBS Sabacc  
**Tagline**: "Classic 77-Card Sabacc for Bulletin Board Systems"  
**Genre**: Terminal-based Card Game / BBS Door Game  
**License**: Not specified (appears to be open source)  
**Target Platform**: Linux-based BBS systems (Mystic, Synchronet, Talisman, WWIV, Enigma½)

## Authentic Sabacc Implementation

### Rule Set: 1989 West End Games Classic
- **76-card deck**: 60 numbered cards (1-15 per suit) + 16 Arcana cards
- **Four positive suits**: Sabers, Flasks, Coins, Staves (1-15 points each)
- **Arcana cards**: Death(-1) through Star(-17), single copies
- **Special hands**: Pure Sabacc (exactly 23), Idiot's Array (Idiot+2+3)
- **Bomb conditions**: >23, <-23, or exactly 0 points

### Classic 4-Phase Turn Structure
1. **Betting Phase**: Check/Call, Raise, or Fold
2. **Roll Phase**: Dice roll for Sabacc Shift (doubles trigger)
3. **Call Phase**: Players may call the hand (end game)
4. **Draw Phase**: Draw, Trade, Stand, or use Static Field

## Technical Architecture

### Core Modules
```
main.go      - Game engine, turn management, player logic
cards.go     - Deck system, card graphics, rendering engine
ui.go        - Screen layout, menus, game log system  
config.go    - Settings management, statistics tracking
helpers.go   - BBS integration, terminal control, utilities
portraits.go - AI player portrait system
```

### Key Features

#### BBS Integration
- **door32.sys support**: Complete BBS drop file parsing
- **Terminal compatibility**: ANSI/ASCII mode detection
- **Idle timeout**: Configurable with warnings
- **80x25 display**: Fixed layout optimized for classic terminals

#### Graphics System
- **Binary card database**: `sabacc_cards.bin` with indexed ANSI art
- **CP437 character set**: Full IBM PC character support
- **Diamond-shaped cards**: Authentic Sabacc card design
- **AI player portraits**: External ANSI art system (`ansi/portraits.ans`)

#### Game Features
- **4 AI opponents**: PHOOJA, ASH-TAAC, OOLANGA, KY'ALA
- **Smart AI behavior**: Configurable personalities (conservative, balanced, aggressive)
- **Static Field system**: Card protection from Sabacc Shifts
- **Dual pot system**: Hand Pot + Sabacc Pot with proper distribution
- **Persistent UI**: Real-time updates without screen flicker

### Configuration System
**File**: `sabacc.conf` (JSON format)
```json
{
  "min_ante": 10,
  "max_ante": 100, 
  "max_bet": 500,
  "starting_credits": 1000,
  "shift_probability": 6,
  "min_rounds_to_call": 1,
  "max_rounds": 4,
  "idle_timeout_seconds": 300,
  "enable_statistics": true,
  "ai_personality": "balanced"
}
```

## User Experience Design

### Screen Layout (80x25)
- **AI Players**: 4 corner positions with portraits and stats
- **Central Game Log**: Blue-bordered scrolling message area
- **Human Player**: Bottom center with card display
- **Status Areas**: Pots, round info, deck size
- **Menu System**: Context-sensitive action menus

### Player Interaction
- **Keyboard-driven**: Single keypress commands
- **Visual feedback**: Color-coded messages and card states
- **Turn indicators**: Highlighted AI player borders
- **Card visualization**: ANSI art with suit colors

## AI System Architecture

### Decision Making
- **Hand evaluation**: Statistical analysis of card totals
- **Risk assessment**: Dynamic betting based on hand strength
- **Shift strategy**: Smart Static Field management
- **Calling logic**: Opponent strength analysis from visible cards

### Personality Types
- **Conservative**: Lower risk tolerance, defensive play
- **Balanced**: Moderate risk/reward decision making  
- **Aggressive**: High stakes, aggressive betting patterns

## Build & Deployment

### Dependencies
- **Go 1.19+**: Core language requirement
- **github.com/eiannone/keyboard**: Terminal input handling
- **Linux environment**: BBS system compatibility

### Build Process
```bash
# Generate card database (required first step)
go run cmd/build-cards/main.go

# Build optimized executable  
go build -ldflags="-s -w" -o sabacc .
```

### Installation Structure
```
/path/to/doors/sabacc/
├── sabacc              # Main executable
├── sabacc_cards.bin    # Card graphics database  
├── sabacc.conf         # Configuration file
└── ansi/
    └── portraits.ans   # AI player portraits
```

## Performance Characteristics
- **Memory footprint**: Minimal, suitable for shared BBS systems
- **Startup time**: Fast initialization with binary card database
- **Network efficiency**: Optimized for telnet/SSH connections
- **Error recovery**: Graceful degradation when assets missing

## Compliance & Authenticity
- **Rules accuracy**: Faithful to 1989 West End Games specification
- **BBS standards**: Compatible with major BBS software packages
- **Terminal standards**: Full ANSI/CP437 compliance
- **Drop file format**: Standard door32.sys implementation

---
*Product Context established: 2025-08-14*  
*Architecture documented for Memory Bank reference*