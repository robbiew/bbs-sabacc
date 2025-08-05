package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/eiannone/keyboard"
)

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
	Idle              = 300 // Default idle timeout in seconds
	CurrentUser       User
	idleTimer         *Timer
	terminalInitiated = false
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
	terminalInitiated = true

	return user
}

// getKeyWithTimeout waits for a key press with idle timeout
func getKeyWithTimeout() (rune, keyboard.Key, error) {
	// Start the idle timer
	shortTimer := NewTimer(Idle, func() {
		fmt.Println("\r\nYou've been idle for too long... exiting!")
		time.Sleep(1 * time.Second)
		os.Exit(0)
	})
	defer shortTimer.Stop()

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

	fmt.Print(Yellow + "WINNING HANDS:\n" + Reset)
	fmt.Print(Green + "• Pure Sabacc: " + White + "Exactly 23 points (wins Sabacc Pot!)\n" + Reset)
	fmt.Print(Green + "• Idiot's Array: " + White + "Idiot + 2 + 3 (literal 23)\n" + Reset)
	fmt.Print(White + "  " + GreenHi + "Idiot's Array beats Pure Sabacc!" + Reset + "\n")
	fmt.Print(White + "• Highest total <= 23 wins Hand Pot\n\n" + Reset)

	fmt.Print(Yellow + "BOMB OUT CONDITIONS:\n" + Reset)
	fmt.Print(Red + "• Over 23 points\n" + Reset)
	fmt.Print(Red + "• Under -23 points\n" + Reset)
	fmt.Print(Red + "• Exactly 0 points\n" + Reset)
	fmt.Print(White + "Penalty: Pay Hand Pot amount to Sabacc Pot\n\n" + Reset)

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

	fmt.Print(Yellow + "CALLING THE HAND:\n" + Reset)
	fmt.Print(White + "• Cannot call until minimum rounds completed\n" + Reset)
	fmt.Print(White + "• Can only call during another player's turn\n" + Reset)
	fmt.Print(White + "• Triggers final showdown\n\n" + Reset)

	fmt.Print(Yellow + "SABACC SHIFT:\n" + Reset)
	fmt.Print(White + "• Triggered by rolling doubles\n" + Reset)
	fmt.Print(White + "• All hands reshuffled and redealt\n" + Reset)
	fmt.Print(White + "• " + MagentaHi + "Static Field cards are protected!" + Reset + "\n\n")

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

	fmt.Print(Yellow + "POSITIVE SUITS:\n" + Reset)
	fmt.Print(Blue + "• Sabers (S): " + White + "1-15 points\n" + Reset)
	fmt.Print(Green + "• Flasks (F): " + White + "1-15 points\n" + Reset)
	fmt.Print(Yellow + "• Coins (C): " + White + "1-15 points\n" + Reset)
	fmt.Print(Red + "• Staves (T): " + White + "1-15 points\n\n" + Reset)

	fmt.Print(Yellow + "ARCANA CARDS (Negative):\n" + Reset)
	fmt.Print(Magenta + "Death(-1), Strength(-2), Moderation(-3), Evil One(-4)\n" + Reset)
	fmt.Print(Magenta + "Justice(-5), Queen(-6), Endurance(-7), Balance(-8)\n" + Reset)
	fmt.Print(Magenta + "Demise(-9), Destruction(-10), Despair(-11), Failure(-12)\n" + Reset)
	fmt.Print(Magenta + "Futility(-13), Mistress(-14), " + MagentaHi + "Idiot(-15)" + Reset + Magenta + ", Star(-17)\n\n" + Reset)

	fmt.Print(Yellow + "ANTE & POTS:\n" + Reset)
	fmt.Print(White + "• Both players ante into " + CyanHi + "both" + White + " pots\n" + Reset)
	fmt.Print(White + "• Hand Pot: Won by best hand <= 23\n" + Reset)
	fmt.Print(White + "• Sabacc Pot: Won by Pure Sabacc or Idiot's Array\n" + Reset)
	fmt.Print(White + "• Fold penalty: 1 credit to Sabacc Pot\n\n" + Reset)

	fmt.Print(Yellow + "Press any key to return to menu..." + Reset)
	waitForKey()
}

// showStats displays player statistics (updated to mention Classic Sabacc)
func showStats() {
	ClearScreen()
	fmt.Print(CyanHi + strings.Repeat("\xcd", 43) + "\n" + Reset)
	fmt.Print(CyanHi + "            PLAYER STATISTICS\n" + Reset)
	fmt.Print(CyanHi + strings.Repeat("\xcd", 43) + "\n\n" + Reset)

	fmt.Print(Yellow + "Player: " + CyanHi + game.User.Alias + Reset + "\n\n")

	// Placeholder stats - in a real implementation, these would be saved/loaded
	fmt.Print(White + "Games Played: " + YellowHi + "0\n" + Reset)
	fmt.Print(White + "Games Won: " + YellowHi + "0\n" + Reset)
	fmt.Print(White + "Pure Sabaccs: " + YellowHi + "0\n" + Reset)
	fmt.Print(White + "Idiot's Arrays: " + YellowHi + "0\n" + Reset)
	fmt.Print(White + "Bomb Outs: " + YellowHi + "0\n" + Reset)
	fmt.Print(White + "Sabacc Shifts Survived: " + YellowHi + "0\n" + Reset)
	fmt.Print(White + "Credits Won: " + YellowHi + "0\n" + Reset)
	fmt.Print(White + "Credits Lost: " + YellowHi + "0\n\n" + Reset)

	fmt.Print(Red + "Statistics tracking not yet implemented.\n" + Reset)
	fmt.Print(Red + "This feature will be added in a future version.\n\n" + Reset)

	fmt.Print(Cyan + "Classic Sabacc follows the original West End Games\n" + Reset)
	fmt.Print(Cyan + "rules from the 1989 Crisis on Cloud City supplement.\n\n" + Reset)

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

// displayWelcomeMessage shows a welcome message with game info
func displayWelcomeMessage() {
	fmt.Printf("%sWelcome to the Classic Sabacc tables, %s%s%s!\n\n",
		Cyan, CyanHi, game.User.Alias, Reset)
	fmt.Printf("%sTerminal: %s%dx%d%s  ",
		Yellow, YellowHi, game.User.W, game.User.H, Reset)
	fmt.Printf("%sNode: %s%d%s  ",
		Yellow, YellowHi, game.User.NodeNum, Reset)
	fmt.Printf("%sTime: %s%dm%s\n\n",
		Yellow, YellowHi, game.User.TimeLeft, Reset)
}

// clearStatusLine clears the bottom status line
func clearStatusLine() {
	MoveCursor(1, game.User.H)
	fmt.Print(EraseLine)
}

// showStatusLine displays a status message at the bottom
func showStatusLine(message string) {
	clearStatusLine()
	MoveCursor(1, game.User.H)
	fmt.Print(Yellow + message + Reset)
}

// animateCardDeal provides a simple animation effect (placeholder)
func animateCardDeal() {
	fmt.Print(Yellow + "Dealing" + Reset)
	for i := 0; i < 3; i++ {
		time.Sleep(300 * time.Millisecond)
		fmt.Print(".")
	}
	fmt.Println()
	time.Sleep(500 * time.Millisecond)
}

// displayDeckInfo shows remaining cards in deck
func displayDeckInfo() {
	fmt.Printf("%sDeck: %s%d cards remaining%s\n",
		Magenta, MagentaHi, len(game.Deck.Cards), Reset)
}

// checkBombOut checks if a hand is bombed out
func checkBombOut(hand []Card) bool {
	total := calculateHandTotal(hand)
	return total > 23 || total < -23 || total == 0
}

// formatCredits formats credits with proper color coding
func formatCredits(amount int) string {
	if amount > 0 {
		return GreenHi + fmt.Sprintf("+%d", amount) + Reset
	} else if amount < 0 {
		return RedHi + fmt.Sprintf("%d", amount) + Reset
	}
	return Yellow + "0" + Reset
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

// saveGameStats saves game statistics (placeholder for future implementation)
func saveGameStats(won bool, credits int) {
	// TODO: Implement statistics saving to a file
	// This could save to a JSON file or simple text file
	// Stats to track: games played, won, lost, credits, special hands, etc.
}

// loadGameStats loads game statistics (placeholder for future implementation)
func loadGameStats() map[string]int {
	// TODO: Implement statistics loading from a file
	// Return default stats for now
	return map[string]int{
		"games_played": 0,
		"games_won":    0,
		"pure_sabacc":  0,
		"idiots_array": 0,
		"bomb_outs":    0,
		"credits_won":  0,
		"credits_lost": 0,
	}
}

// displayAsciiArt shows ASCII art for special events
func displayAsciiArt(artType string) {
	switch artType {
	case "sabacc":
		fmt.Print(GreenHi)
		fmt.Println("  \xdb\xdb\xdb\xdb\xdb\xdb\xdb  \xdb\xdb\xdb\xdb\xdb  \xdb\xdb\xdb\xdb\xdb\xdb   \xdb\xdb\xdb\xdb\xdb   \xdb\xdb\xdb\xdb\xdb\xdb  \xdb\xdb\xdb\xdb\xdb\xdb ")
		fmt.Println("  \xdb\xdb      \xdb\xdb   \xdb\xdb \xdb\xdb   \xdb\xdb \xdb\xdb   \xdb\xdb \xdb\xdb      \xdb\xdb     ")
		fmt.Println("  \xdb\xdb\xdb\xdb\xdb\xdb\xdb \xdb\xdb\xdb\xdb\xdb\xdb\xdb \xdb\xdb\xdb\xdb\xdb\xdb  \xdb\xdb\xdb\xdb\xdb\xdb\xdb \xdb\xdb      \xdb\xdb     ")
		fmt.Println("       \xdb\xdb \xdb\xdb   \xdb\xdb \xdb\xdb   \xdb\xdb \xdb\xdb   \xdb\xdb \xdb\xdb      \xdb\xdb     ")
		fmt.Println("  \xdb\xdb\xdb\xdb\xdb\xdb\xdb \xdb\xdb   \xdb\xdb \xdb\xdb\xdb\xdb\xdb\xdb  \xdb\xdb   \xdb\xdb  \xdb\xdb\xdb\xdb\xdb\xdb  \xdb\xdb\xdb\xdb\xdb\xdb ")
		fmt.Print(Reset)
	case "bomb":
		fmt.Print(RedHi)
		fmt.Println("  \xdb\xdb\xdb\xdb\xdb\xdb   \xdb\xdb\xdb\xdb\xdb\xdb  \xdb\xdb\xdb    \xdb\xdb\xdb \xdb\xdb\xdb\xdb\xdb\xdb  \xdb\xdb")
		fmt.Println("  \xdb\xdb   \xdb\xdb \xdb\xdb    \xdb\xdb \xdb\xdb\xdb\xdb  \xdb\xdb\xdb\xdb \xdb\xdb   \xdb\xdb \xdb\xdb")
		fmt.Println("  \xdb\xdb\xdb\xdb\xdb\xdb  \xdb\xdb    \xdb\xdb \xdb\xdb \xdb\xdb\xdb\xdb \xdb\xdb \xdb\xdb\xdb\xdb\xdb\xdb  \xdb\xdb")
		fmt.Println("  \xdb\xdb   \xdb\xdb \xdb\xdb    \xdb\xdb \xdb\xdb  \xdb\xdb  \xdb\xdb \xdb\xdb   \xdb\xdb   ")
		fmt.Println("  \xdb\xdb\xdb\xdb\xdb\xdb   \xdb\xdb\xdb\xdb\xdb\xdb  \xdb\xdb      \xdb\xdb \xdb\xdb\xdb\xdb\xdb\xdb  \xdb\xdb")
		fmt.Print(Reset)
	case "shift":
		fmt.Print(YellowHi)
		fmt.Println("  \xdb\xdb\xdb\xdb\xdb\xdb\xdb \xdb\xdb   \xdb\xdb \xdb\xdb \xdb\xdb\xdb\xdb\xdb\xdb\xdb \xdb\xdb\xdb\xdb\xdb\xdb\xdb\xdb \xdb\xdb")
		fmt.Println("  \xdb\xdb      \xdb\xdb   \xdb\xdb \xdb\xdb \xdb\xdb         \xdb\xdb    \xdb\xdb")
		fmt.Println("  \xdb\xdb\xdb\xdb\xdb\xdb\xdb \xdb\xdb\xdb\xdb\xdb\xdb\xdb \xdb\xdb \xdb\xdb\xdb\xdb\xdb      \xdb\xdb    \xdb\xdb")
		fmt.Println("       \xdb\xdb \xdb\xdb   \xdb\xdb \xdb\xdb \xdb\xdb         \xdb\xdb      ")
		fmt.Println("  \xdb\xdb\xdb\xdb\xdb\xdb\xdb \xdb\xdb   \xdb\xdb \xdb\xdb \xdb\xdb         \xdb\xdb    \xdb\xdb")
		fmt.Print(Reset)
	}
}

// Debug function to help identify the formatting issues
func debugGameState() {
	fmt.Printf("DEBUG: HandPot type: %T, value: %v\n", game.HandPot, game.HandPot)
	fmt.Printf("DEBUG: SabaccPot type: %T, value: %v\n", game.SabaccPot, game.SabaccPot)
	fmt.Printf("DEBUG: Round type: %T, value: %v\n", game.Round, game.Round)

	if len(game.Players) > 0 {
		fmt.Printf("DEBUG: Player credits type: %T, value: %v\n",
			game.Players[0].Credits, game.Players[0].Credits)
		fmt.Printf("DEBUG: Player hand length: %d\n", len(game.Players[0].Hand))
	}
}
