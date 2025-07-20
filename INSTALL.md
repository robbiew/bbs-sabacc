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
3. **Build the game:**
   ```bash
   # Quick build
   go build -o sabacc .
   
   # Optimized build (smaller executable)
   go build -ldflags="-s -w" -o sabacc .
   
   # Or use the build script
   chmod +x build.sh
   ./build.sh
   ```
4. **Copy to your BBS doors directory:**
   ```bash
   cp sabacc /path/to/your/bbs/doors/
   chmod +x /path/to/your/bbs/doors/sabacc
   ```

### Method 2: Pre-built Binary (if available)

1. Download the appropriate binary for your architecture
2. Extract and copy to your doors directory:
   ```bash
   cp sabacc /path/to/your/bbs/doors/
   chmod +x /path/to/your/bbs/doors/sabacc
   ```

## Testing Installation

Before configuring your BBS, test the game locally:

```bash
# Create a test drop file
echo -e "2\n8\n38400\nTest BBS\n1\nTest Player\nTestUser\n100\n90\n0\n1" > door32.sys

# Test the game
./sabacc -path ./

# You should see the title screen and main menu
```

## BBS Configuration

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
Intercept Standard I/O=No
Native (32-bit) Executable=No
Use Shell to Execute=No
Modify User Data=No
Execute on Event=Logon
BBS Drop File Type=Door32.sys
Place Drop File In=Node Directory
Time Options=Deduct from Time Online
Access Requirements=
```

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
├── sabacc              # Main executable
├── ansi/               # ANSI art files (optional)
│   ├── title.ans
│   ├── menu.ans
│   └── game.ans
├── stats/              # Player statistics (auto-created)
└── sabacc.conf         # Configuration file (auto-created)
```

## Configuration Files

### Game Configuration (sabacc.conf)

The game creates a configuration file automatically with these default settings:

```json
{
  "min_ante": 10,
  "max_ante": 100,
  "max_bet": 500,
  "starting_credits": 1000,
  "shift_probability": 36,
  "min_rounds_to_call": 2,
  "idle_timeout_seconds": 300,
  "enable_sound": true,
  "enable_statistics": true,
  "ai_personality": "balanced"
}
```

### ANSI Art (Optional)

Place custom ANSI art files in the `ansi/` directory:

- **title.ans** - Title screen (displayed on startup)
- **menu.ans** - Menu background  
- **game.ans** - Game screen background

Files should use standard ANSI escape sequences. The game includes built-in ASCII art if these files are missing.

## Verification

After installation, verify everything works:

1. **Check file permissions:**
   ```bash
   ls -la /path/to/doors/sabacc
   # Should show: -rwxr-xr-x (executable)
   ```

2. **Test from command line:**
   ```bash
   /path/to/doors/sabacc -path /tmp/
   # Should show error about missing door32.sys (this is normal)
   ```

3. **Test through your BBS:**
   - Log into your BBS
   - Navigate to the doors/games menu
   - Run Sabacc
   - Verify it displays properly with ANSI colors

## Troubleshooting

### Common Issues

**"No such file or directory" when running sabacc:**
- Check that the executable exists and has proper permissions
- Verify the path in your BBS configuration is correct
- Ensure the executable was built for the correct architecture

**Game starts but shows format errors:**
- Usually indicates a drop file format problem
- Verify your BBS is creating door32.sys correctly
- Check that the path parameter includes trailing slash

**No ANSI colors:**
- Verify the user's terminal supports ANSI
- Check BBS emulation settings
- Ensure door32.sys has emulation=1 for ANSI mode

**Game crashes immediately:**
- Check system logs for error messages
- Verify Go runtime compatibility
- Test with a manual drop file first

### Debug Steps

1. **Create test environment:**
   ```bash
   cd /path/to/doors/
   echo -e "2\n8\n38400\nTest BBS\n1\nTest Player\nTestUser\n100\n90\n0\n1" > door32.sys
   ./sabacc -path ./
   ```

2. **Check dependencies:**
   ```bash
   ldd sabacc  # Shows required libraries
   ```

3. **Verify drop file format:**
   ```bash
   cat door32.sys  # Should show 11 lines of data
   ```

### Performance Notes

- **Memory usage**: ~5-10MB per instance
- **CPU usage**: Minimal (single-threaded, event-driven)
- **Disk usage**: <1MB plus statistics files
- **Concurrent users**: Supported (separate processes)

### Security Considerations

- Runs with BBS user permissions
- No network connections made
- Statistics stored in plain JSON (not sensitive)
- No user input validation beyond game rules

## Maintenance

### Updating

To update to a new version:

1. **Backup your data:**
   ```bash
   cp -r /path/to/doors/sabacc/stats/ ~/sabacc-backup/
   cp /path/to/doors/sabacc/sabacc.conf ~/sabacc-backup/
   ```

2. **Replace executable:**
   ```bash
   # Build new version
   go build -ldflags="-s -w" -o sabacc .
   
   # Stop BBS or disable door temporarily
   # Replace executable
   cp sabacc /path/to/doors/sabacc/
   
   # Restart BBS or re-enable door
   ```

3. **Restore configuration:**
   ```bash
   # Configuration files are forward-compatible
   # Statistics files are preserved automatically
   ```

### Log Monitoring

Monitor your BBS logs for any door-related errors:

```bash
# Common log locations
tail -f /var/log/mystic/mis*.log     # Mystic BBS
tail -f /enigma-bbs/logs/system.log  # ENiGMA½
# Check your specific BBS documentation for log locations
```

### Statistics Management

Player statistics are stored in `stats/` directory as JSON files:

```bash
# View player stats
ls -la stats/
cat stats/player_name.json

# Backup all statistics
tar -czf sabacc-stats-backup.tar.gz stats/

# Reset a player's stats (delete their file)
rm stats/player_name.json
```

## Advanced Configuration

### Custom AI Personalities

Edit `sabacc.conf` to adjust AI behavior:

- **"conservative"** - Cautious, folds early, uses static field often
- **"balanced"** - Default reasonable behavior
- **"aggressive"** - Risk-taking, draws more cards, folds less

### Credit Economy

Adjust starting credits and limits to match your BBS economy:

```json
{
  "starting_credits": 500,    # Lower for tighter economy
  "min_ante": 5,             # Minimum bet
  "max_ante": 50,            # Maximum bet
  "max_bet": 200             # Betting limit
}
```

### Time Limits

Set appropriate idle timeouts:

```json
{
  "idle_timeout_seconds": 300  # 5 minutes default
}
```

