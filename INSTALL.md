# BBS Sabacc Installation Guide

This guide will help you install and configure BBS Sabacc on your BBS system.

## Prerequisites

- Linux-based BBS system (Talisman, Mystic, Synchronet, ENiGMA½, WWIV)
- Go 1.19 or later (for building from source)
- ANSI terminal support
- door32.sys drop file support

## Quick Installation

### Option 1: Download Pre-built Binary

1. Download the latest release from the releases page
2. Extract the archive:
   ```bash
   tar -xzf sabacc-1.0.0.tar.gz
   cd sabacc-1.0.0
   ```
3. Copy the binary to your BBS doors directory:
   ```bash
   cp sabacc_linux_amd64 /path/to/your/bbs/doors/sabacc
   chmod +x /path/to/your/bbs/doors/sabacc
   ```

### Option 2: Build from Source

1. Clone or download the source code
2. Build the application:
   ```bash
   make build
   # or manually:
   # go build -o sabacc .
   ```
3. Copy to your doors directory:
   ```bash
   cp bin/sabacc /path/to/your/bbs/doors/
   chmod +x /path/to/your/bbs/doors/sabacc
   ```

## BBS Configuration

### Talisman BBS

Add to your `doors.toml` file:

```toml
[[door]]
name = "Sabacc"
command = "/path/to/doors/sabacc -path %3"
directory = "/path/to/doors/"
type = "door"
description = "Classic 76-Card Sabacc Game"
access_level = 10
time_limit = 30
```

### Mystic BBS

Add to your door configuration:

```
Description  : Sabacc Card Game
Access Level : 10
Command Line : /path/to/doors/sabacc -path %3
Optional Data: 
Use Door32   : Yes
Time Limit   : 30
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
```

### ENiGMA½ BBS

Add to your `menu.hjson` file:

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

Create the following directory structure in your doors directory:

```
/path/to/doors/sabacc/
├── sabacc              # Main executable
├── ansi/               # ANSI art files (optional)
│   ├── title.ans
│   ├── menu.ans
│   └── game.ans
├── stats/              # Player statistics (auto-created)
├── sabacc.conf         # Configuration file (auto-created)
└── README.md           # Documentation
```

## Configuration

### Game Configuration

The game will create a `sabacc.conf` file on first run with default settings:

```json
{
  "min_ante": 10,
  "max_ante": 100,
  "max_bet": 500,
  "starting_credits": 1000,
  "shift_probability": 6,
  "min_rounds_to_call": 2,
  "idle_timeout_seconds": 300,
  "enable_sound": true,
  "enable_statistics": true,
  "ai_personality": "balanced"
}
```

### ANSI Art (Optional)

Place custom ANSI art files in the `ansi/` directory:

- `title.ans` - Title screen (displayed on startup)
- `menu.ans` - Menu background
- `game.ans` - Game screen background

Files should be standard ANSI format. If files are not present, the game will display built-in ASCII art.

## Testing Installation

1. Create a test drop file:
   ```bash
   echo -e "2\n8\n38400\nTest BBS\n1\nTest Player\nTestUser\n100\n90\n0\n1" > door32.sys
   ```

2. Test the game:
   ```bash
   ./sabacc -path ./
   ```

3. Verify the game starts and displays properly

## Troubleshooting

### Common Issues

**Game doesn't start:**
- Check that the executable has proper permissions (`chmod +x sabacc`)
- Verify the drop file path is correct
- Ensure door32.sys exists and is readable

**ANSI colors not displaying:**
- Verify your terminal supports ANSI colors
- Check that the user's terminal is set to ANSI mode
- Ensure the BBS is passing the correct emulation type in door32.sys

**Game crashes on startup:**
- Check system architecture compatibility (x86_64 vs ARM)
- Verify Go runtime dependencies
- Check system logs for error messages

**Drop file errors:**
- Ensure door32.sys format is correct
- Verify file permissions
- Check that the BBS is creating the drop file properly

### Debug Mode

Run the game with verbose output:
```bash
./sabacc -path ./ -debug
```

### Log Files

Check BBS logs for door-related error messages. The game itself doesn't create log files by default.

## Security Considerations

- The game runs with the permissions of the BBS user
- Player statistics are stored in JSON files (not encrypted)
- No network connections are made
- No sensitive data is transmitted

## Performance

- Memory usage: ~5-10MB
- CPU usage: Minimal (single-threaded)
- Disk usage: <1MB plus statistics files
- Supports multiple concurrent users (separate processes)

## Updating

To update to a new version:

1. Stop the BBS or disable the door temporarily
2. Backup your `stats/` directory and `sabacc.conf` file
3. Replace the executable with the new version
4. Restart the BBS or re-enable the door
5. Test functionality

Configuration and statistics files are forward-compatible.

## Uninstalling

1. Remove the door configuration from your BBS
2. Delete the game directory
3. Remove any menu references

Player statistics will be preserved if you keep the `stats/` directory.

## Support

For installation help:
- Check the README.md file
- Review BBS-specific documentation
- Test with the provided door32.sys test file
- Verify ANSI terminal compatibility

## Advanced Configuration

### Custom AI Personalities

Edit `sabacc.conf` to change AI behavior:
- `"conservative"` - Cautious play style
- `"balanced"` - Default behavior
- `"aggressive"` - Risk-taking style

### Time Limits

Adjust the idle timeout in the configuration file. Default is 300 seconds (5 minutes).

### Credit Limits

Modify starting credits and ante limits in the configuration file to match your BBS economy.

### Statistics Database

Player statistics are stored as individual JSON files in the `stats/` directory. These can be backed up, migrated, or reset as needed.

This completes the installation guide. The game should now be ready for your BBS users to enjoy!