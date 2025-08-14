package main

import (
	"bufio"
	_ "embed"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/eiannone/keyboard"
)

// Add a helper to strip ANSI escape codes
var ansiRegexp = regexp.MustCompile(`\x1b\[[0-9;]*m`)

func stripANSI(s string) string {
	return ansiRegexp.ReplaceAllString(s, "")
}

// ANSI Color Constants (replacing godoors colors)
const (
	Reset     = "\x1b[0m"
	Black     = "\x1b[30m"
	Red       = "\x1b[31m"
	Green     = "\x1b[32m"
	Yellow    = "\x1b[33m"
	Blue      = "\x1b[34m"
	Magenta   = "\x1b[35m"
	Cyan      = "\x1b[36m"
	White     = "\x1b[37m"
	RedHi     = "\x1b[91m"
	GreenHi   = "\x1b[92m"
	YellowHi  = "\x1b[93m"
	BlueHi    = "\x1b[94m"
	MagentaHi = "\x1b[95m"
	CyanHi    = "\x1b[96m"
	WhiteHi   = "\x1b[97m"
	EraseLine = "\x1b[2K"

	// Background colors
	BgBlack   = "\x1b[40m"
	BgRed     = "\x1b[41m"
	BgGreen   = "\x1b[42m"
	BgYellow  = "\x1b[43m"
	BgBlue    = "\x1b[44m"
	BgMagenta = "\x1b[45m"
	BgCyan    = "\x1b[46m"
	BgWhite   = "\x1b[47m"

	// Text formatting
	Bold      = "\x1b[1m"
	Dim       = "\x1b[2m"
	Underline = "\x1b[4m"

	// Combined colors for common UI elements
	StatusBarBg          = "\x1b[36;46m"      // Cyan foreground on cyan background
	StatusBarText        = "\x1b[0;30;46m"    // Black text on cyan background
	StatusBarBoldWhite   = "\x1b[1;37m"       // Bold white text
	StatusBarHighlight   = "\x1b[33m"         // Yellow highlight
	StatusBarNormal      = "\x1b[37m"         // Normal white text
	StatusBarTransition  = "\x1b[36;40m"      // Cyan on black (transition)
)

// User represents a BBS user (replacing godoors User)
type User struct {
	Alias       string
	RealName    string
	Handle      string
	TimeLeft    int
	NodeNum     int
	W           int // Terminal width
	H           int // Terminal height
	Emulation   int // 0=ASCII, 1=ANSI
	SecurityLvl int
}

// Timer represents an idle timer (replacing godoors Timer)
type Timer struct {
	duration time.Duration
	callback func()
	timer    *time.Timer
	stopped  bool
}

// Global variables (replacing godoors globals)
var (
	Idle       = 60 // Default idle timeout in seconds
	CurrentUser User
	gameLayout *ScreenLayout // Reference to game layout for timeout warnings

	// Embedded result art files
	//go:embed ansi/result-sabacc.ans
	SabaccResultArt string
	//go:embed ansi/result-bomb.ans
	BombResultArt string
	//go:embed ansi/result-shift.ans
	ShiftResultArt string
)

// Terminal Control Functions (replacing godoors functions)

// ClearScreen clears the terminal screen
func ClearScreen() {
	fmt.Print("\x1b[2J\x1b[H")
}

// MoveCursor moves cursor to x, y position (1-based)
func MoveCursor(x, y int) {
	fmt.Printf("\x1b[%d;%dH", y, x)
}

// NewTimer creates a new timer with callback (replacing godoors NewTimer)
func NewTimer(seconds int, callback func()) *Timer {
	timer := &Timer{
		duration: time.Duration(seconds) * time.Second,
		callback: callback,
		stopped:  false,
	}

	timer.timer = time.AfterFunc(timer.duration, func() {
		if !timer.stopped {
			callback()
		}
	})

	return timer
}

// Stop stops the timer
func (t *Timer) Stop() {
	t.stopped = true
	if t.timer != nil {
		t.timer.Stop()
	}
}

// Initialize initializes the BBS environment (replacing godoors Initialize)
func Initialize(dropPath string) User {
	// Read door32.sys file
	door32Path := dropPath + "/door32.sys"
	data, err := os.ReadFile(door32Path)
	if err != nil {
		fmt.Printf("Error reading door32.sys: %v\n", err)
		os.Exit(1)
	}

	lines := strings.Split(string(data), "\n")
	if len(lines) < 11 {
		fmt.Printf("Invalid door32.sys format\n")
		os.Exit(1)
	}

	// Parse door32.sys fields
	emulation, _ := strconv.Atoi(strings.TrimSpace(lines[0]))
	nodeNum, _ := strconv.Atoi(strings.TrimSpace(lines[1]))
	timeLeft, _ := strconv.Atoi(strings.TrimSpace(lines[9]))
	securityLvl, _ := strconv.Atoi(strings.TrimSpace(lines[10]))

	user := User{
		Alias:       strings.TrimSpace(lines[5]),
		RealName:    strings.TrimSpace(lines[4]),
		Handle:      strings.TrimSpace(lines[6]),
		TimeLeft:    timeLeft,
		NodeNum:     nodeNum,
		W:           80, // Default terminal width
		H:           25, // Default terminal height
		Emulation:   emulation,
		SecurityLvl: securityLvl,
	}

	CurrentUser = user

	return user
}

// getKeyWithTimeout waits for a key press with idle timeout and warning
func getKeyWithTimeout() (rune, keyboard.Key, error) {
	// Calculate warning time (10 seconds before timeout)
	warningTime := Idle - 10
	if warningTime < 1 {
		warningTime = Idle / 2 // If timeout is very short, warn at halfway point
	}

	// Start warning timer
	warningTimer := NewTimer(warningTime, func() {
		if gameLayout != nil && game != nil {
			// Send warning to game log
			game.Layout.LogMessage("⚠ Are you still here? Auto-logout in 10 seconds!", "warning")
		} else {
			// Fallback warning for menu screens
			fmt.Print(YellowHi + "\r\n⚠ Are you still here? Auto-logout in 10 seconds!" + Reset + "\r\n")
		}
	})

	// Start the idle timer
	idleTimer := NewTimer(Idle, func() {
		if gameLayout != nil && game != nil {
			game.Layout.LogMessage("Timed out!", "error")
			time.Sleep(1 * time.Second)
		}
		fmt.Println(RedHi + "\r\nYou've been idle for too long... exiting!" + Reset)
		time.Sleep(2 * time.Second)
		os.Exit(0)
	})

	// Clean up timers when we get input
	defer func() {
		warningTimer.Stop()
		idleTimer.Stop()
	}()

	char, key, err := keyboard.GetKey()
	return char, key, err
}

// waitForKey waits for any key press
func waitForKey() {
	getKeyWithTimeout()
}

// showRules displays the Classic Sabacc rules (updated for West End Games rules)
func showRules() {
	ClearScreen()
	fmt.Print(CyanHi + strings.Repeat("\xcd", 43) + "\n" + Reset)
	fmt.Print(CyanHi + "           CLASSIC SABACC RULES\n" + Reset)
	fmt.Print(CyanHi + "            West End Games (1989)\n" + Reset)
	fmt.Print(CyanHi + strings.Repeat("\xcd", 43) + "\n\n" + Reset)

	fmt.Print(Yellow + "OBJECTIVE:\n" + Reset)
	fmt.Print(White + "Get the highest card total that is <= 23.\n\n" + Reset)

	fmt.Print(Yellow + "WINNING HANDS:\r\n" + Reset)
	fmt.Print(Green + "\xfa Pure Sabacc: " + White + "Exactly 23 points (wins Sabacc Pot!)\r\n" + Reset)
	fmt.Print(Green + "\xfa Idiot's Array: " + White + "Idiot + 2 + 3 (literal 23)\r\n" + Reset)
	fmt.Print(White + "  " + GreenHi + "Idiot's Array beats Pure Sabacc!" + Reset + "\r\n")
	fmt.Print(White + "\xfa Highest total <= 23 wins Hand Pot\r\n\r\n" + Reset)

	fmt.Print(Yellow + "BOMB OUT CONDITIONS:\r\n" + Reset)
	fmt.Print(Red + "\xfa Over 23 points\r\n" + Reset)
	fmt.Print(Red + "\xfa Under -23 points\r\n" + Reset)
	fmt.Print(Red + "\xfa Exactly 0 points\r\n" + Reset)
	fmt.Print(White + "Penalty: Pay Hand Pot amount to Sabacc Pot\r\n\r\n" + Reset)

	// Show second page of rules
	fmt.Print(Yellow + "Press any key for turn structure..." + Reset)
	waitForKey()

	ClearScreen()
	fmt.Print(CyanHi + strings.Repeat("\xcd", 43) + "\n" + Reset)
	fmt.Print(CyanHi + "            TURN STRUCTURE\n" + Reset)
	fmt.Print(CyanHi + strings.Repeat("\xcd", 43) + "\n\n" + Reset)

	fmt.Print(Yellow + "EACH TURN HAS 4 PHASES:\n" + Reset)
	fmt.Print(White + "1. " + CyanHi + "BET:" + White + " Check/Call, Raise, or Fold\n" + Reset)
	fmt.Print(White + "2. " + CyanHi + "ROLL:" + White + " Roll dice (doubles = Sabacc Shift!)\n" + Reset)
	fmt.Print(White + "3. " + CyanHi + "CALL:" + White + " Others may call the hand\n" + Reset)
	fmt.Print(White + "4. " + CyanHi + "DRAW:" + White + " Gain, Trade, Stand, or Static Field\n\n" + Reset)

	fmt.Print(Yellow + "CALLING THE HAND:\r\n" + Reset)
	fmt.Print(White + "\xfa Cannot call until minimum rounds completed\r\n" + Reset)
	fmt.Print(White + "\xfa Can only call during another player's turn\r\n" + Reset)
	fmt.Print(White + "\xfa Triggers final showdown\r\n\r\n" + Reset)

	fmt.Print(Yellow + "SABACC SHIFT:\r\n" + Reset)
	fmt.Print(White + "\xfa Triggered by rolling doubles\r\n" + Reset)
	fmt.Print(White + "\xfa All hands reshuffled and redealt\r\n" + Reset)
	fmt.Print(White + "\xfa " + MagentaHi + "Static Field cards are protected!" + Reset + "\r\n\r\n")

	// Show third page of rules
	fmt.Print(Yellow + "Press any key for deck information..." + Reset)
	waitForKey()

	ClearScreen()
	fmt.Print(CyanHi + strings.Repeat("\xcd", 43) + "\n" + Reset)
	fmt.Print(CyanHi + "              DECK & CARDS\n" + Reset)
	fmt.Print(CyanHi + strings.Repeat("\xcd", 43) + "\n\n" + Reset)

	fmt.Print(Yellow + "76-CARD DECK:\n" + Reset)
	fmt.Print(White + "60 Numbered Cards (1-15 in each suit)\n" + Reset)
	fmt.Print(White + "16 Arcana Cards (negative values, 2 copies each)\n\n" + Reset)

	fmt.Print(Yellow + "POSITIVE SUITS:\r\n" + Reset)
	fmt.Print(Blue + "\xfa Sabers (S): " + White + "1-15 points\r\n" + Reset)
	fmt.Print(Green + "\xfa Flasks (F): " + White + "1-15 points\r\n" + Reset)
	fmt.Print(Yellow + "\xfa Coins (C): " + White + "1-15 points\r\n" + Reset)
	fmt.Print(Red + "\xfa Staves (T): " + White + "1-15 points\r\n\r\n" + Reset)

	fmt.Print(Yellow + "ARCANA CARDS (Negative):\n" + Reset)
	fmt.Print(Magenta + "Death(-1), Strength(-2), Moderation(-3), Evil One(-4)\n" + Reset)
	fmt.Print(Magenta + "Justice(-5), Queen(-6), Endurance(-7), Balance(-8)\n" + Reset)
	fmt.Print(Magenta + "Demise(-9), Destruction(-10), Despair(-11), Failure(-12)\n" + Reset)
	fmt.Print(Magenta + "Futility(-13), Mistress(-14), " + MagentaHi + "Idiot(-15)" + Reset + Magenta + ", Star(-17)\n\n" + Reset)

	fmt.Print(Yellow + "ANTE & POTS:\r\n" + Reset)
	fmt.Print(White + "\xfa Both players ante into " + CyanHi + "both" + White + " pots\r\n" + Reset)
	fmt.Print(White + "\xfa Hand Pot: Won by best hand <= 23\r\n" + Reset)
	fmt.Print(White + "\xfa Sabacc Pot: Won by Pure Sabacc or Idiot's Array\r\n" + Reset)
	fmt.Print(White + "\xfa Fold penalty: 1 credit to Sabacc Pot\r\n\r\n" + Reset)

	fmt.Print(Yellow + "Press any key to return to menu..." + Reset)
	waitForKey()
}



// exitGame handles clean exit (updated message)
func exitGame() {
	ClearScreen()
	MoveCursor(1, game.User.H/2)
	fmt.Print(CyanHi + "Thanks for playing Classic Sabacc!" + Reset)
	MoveCursor(1, game.User.H/2+2)
	fmt.Print(Yellow + "May the Force be with you, " + game.User.Alias + "!" + Reset)
	MoveCursor(1, game.User.H/2+4)
	fmt.Print(White + "West End Games Rules (1989)" + Reset)
	MoveCursor(1, game.User.H-1)
	time.Sleep(2 * time.Second)
	os.Exit(0)
}


// displayHandValue shows the current hand value with color coding
func displayHandValue(total int) string {
	if total == 23 {
		return GreenHi + "23 (SABACC!)" + Reset
	} else if total > 23 || total < -23 || total == 0 {
		return RedHi + fmt.Sprintf("%d (BOMB!)", total) + Reset
	} else if total >= 20 && total <= 22 {
		return YellowHi + fmt.Sprintf("%d", total) + Reset
	} else if total >= -22 && total <= -20 {
		return YellowHi + fmt.Sprintf("%d", total) + Reset
	} else {
		return White + fmt.Sprintf("%d", total) + Reset
	}
}


// displayAsciiArt shows ASCII art for special events using embedded ANSI files
func displayAsciiArt(artType string) {
	var artContent string
	
	// Select the appropriate embedded art
	switch artType {
	case "sabacc":
		artContent = SabaccResultArt
	case "bomb":
		artContent = BombResultArt
	case "shift":
		artContent = ShiftResultArt
	default:
		// Unknown art type, just return
		return
	}
	
	// Display the embedded ANSI art (already centered in the .ans files)
	fmt.Print(artContent)
}


func TrimStringFromSauce(s string) string {
	if idx := strings.Index(s, "COMNT"); idx != -1 {
		string := s
		delimiter := "COMNT"
		leftOfDelimiter := strings.Split(string, delimiter)[0]
		trim := TrimLastChar(leftOfDelimiter)
		return trim
	}
	if idx := strings.Index(s, "SAUCE00"); idx != -1 {
		string := s
		delimiter := "SAUCE00"
		leftOfDelimiter := strings.Split(string, delimiter)[0]
		trim := TrimLastChar(leftOfDelimiter)
		return trim
	}
	return s
}

func TrimLastChar(s string) string {
	r, size := utf8.DecodeLastRuneInString(s)
	if r == utf8.RuneError && (size == 0 || size == 1) {
		size = 0
	}
	return s[:len(s)-size]
}

func PrintAnsiLoc(artfile string, x int, y int) {
	yLoc := y

	noSauce := TrimStringFromSauce(artfile) // strip off the SAUCE metadata
	s := bufio.NewScanner(strings.NewReader(string(noSauce)))

	for s.Scan() {
		fmt.Fprintf(os.Stdout, Esc+strconv.Itoa(yLoc)+";"+strconv.Itoa(x)+"f"+s.Text())
		yLoc++
	}
}

// createStatusBar creates a formatted status bar with embedded positioning using layout coordinates
func createStatusBar(sl *ScreenLayout, round, deckSize int) string {
	var sb strings.Builder
	
	// Set color and position for "Round:" label - use layout coordinates
	sb.WriteString(Cyan)
	sb.WriteString(fmt.Sprintf("\x1b[%d;%dH", sl.StatusBarY, sl.StatusBarX))
	sb.WriteString("Round: ")
	
	// Yellow for the round number
	sb.WriteString(Yellow)
	sb.WriteString(strconv.Itoa(round))

	// Set color and position for "Deck:" label
	sb.WriteString(Cyan)
	sb.WriteString(fmt.Sprintf("\x1b[%d;%dH", sl.StatusBarY+1, sl.StatusBarX))
	sb.WriteString("Deck: ")
	
	// Yellow for deck size
	sb.WriteString(Yellow)
	sb.WriteString(strconv.Itoa(deckSize))
	sb.WriteString(" ")
	
	// Reset colors
	sb.WriteString(Reset)
		
	return sb.String()
}

// createPotInfo creates a formatted pot info display with embedded positioning using layout coordinates
func createPotInfo(sl *ScreenLayout, gamePot, sabaccPot, sidePot int) string {
	var sb strings.Builder
	
	// Set color and position for "Game Pot:" label - use layout coordinates
	sb.WriteString(White)
	sb.WriteString(fmt.Sprintf("\x1b[%d;%dH", sl.PotInfoY, sl.PotInfoX))
	sb.WriteString("Game Pot: ")
	
	// Yellow for the game pot value
	sb.WriteString(Yellow)
	sb.WriteString(strconv.Itoa(gamePot))

	// Set color and position for "Sabacc Pot:" label
	sb.WriteString(White)
	sb.WriteString(fmt.Sprintf("\x1b[%d;%dH", sl.PotInfoY+1, sl.PotInfoX))
	sb.WriteString("Sabacc Pot: ")
	
	// Yellow for sabacc pot value
	sb.WriteString(Yellow)
	sb.WriteString(strconv.Itoa(sabaccPot))

	// Set color and position for "Side Pot:" label
	sb.WriteString(White)
	sb.WriteString(fmt.Sprintf("\x1b[%d;%dH", sl.PotInfoY+2, sl.PotInfoX))
	sb.WriteString("Side Pot: ")
	
	// Yellow for side pot value
	sb.WriteString(Yellow)
	sb.WriteString(strconv.Itoa(sidePot))
	sb.WriteString(" ")
	
	// Reset colors
	sb.WriteString(Reset)
		
	return sb.String()
}
