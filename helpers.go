package main

import (
	"fmt"
	"os"
	"time"

	"github.com/eiannone/keyboard"
	gd "github.com/robbiew/godoors"
)

// getKeyWithTimeout waits for a key press with idle timeout
func getKeyWithTimeout() (rune, keyboard.Key, error) {
	// Start the idle timer
	shortTimer := gd.NewTimer(gd.Idle, func() {
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

// showRules displays the game rules
func showRules() {
	gd.ClearScreen()
	fmt.Print(gd.CyanHi + "═══════════════════════════════════════════\n" + gd.Reset)
	fmt.Print(gd.CyanHi + "              SABACC RULES\n" + gd.Reset)
	fmt.Print(gd.CyanHi + "═══════════════════════════════════════════\n\n" + gd.Reset)

	fmt.Print(gd.Yellow + "OBJECTIVE:\n" + gd.Reset)
	fmt.Print(gd.White + "Get as close to 23 as possible without going over.\n\n" + gd.Reset)

	fmt.Print(gd.Yellow + "WINNING HANDS:\n" + gd.Reset)
	fmt.Print(gd.Green + "• Pure Sabacc: " + gd.White + "Exactly 23 points\n" + gd.Reset)
	fmt.Print(gd.Green + "• Idiot's Array: " + gd.White + "Idiot card + 2 + 3 (literal 23)\n" + gd.Reset)
	fmt.Print(gd.White + "  Idiot's Array beats Pure Sabacc!\n\n" + gd.Reset)

	fmt.Print(gd.Yellow + "BOMB OUT CONDITIONS:\n" + gd.Reset)
	fmt.Print(gd.Red + "• Over 23 points\n" + gd.Reset)
	fmt.Print(gd.Red + "• Under -23 points\n" + gd.Reset)
	fmt.Print(gd.Red + "• Exactly 0 points\n\n" + gd.Reset)

	fmt.Print(gd.Yellow + "SABACC SHIFT:\n" + gd.Reset)
	fmt.Print(gd.White + "Random event that shuffles all hands!\n" + gd.Reset)
	fmt.Print(gd.White + "Use Static Field to protect important cards.\n\n" + gd.Reset)

	fmt.Print(gd.Yellow + "CARD SUITS:\n" + gd.Reset)
	fmt.Print(gd.Blue + "• Sabers (S): " + gd.White + "1-15 points\n" + gd.Reset)
	fmt.Print(gd.Green + "• Flasks (F): " + gd.White + "1-15 points\n" + gd.Reset)
	fmt.Print(gd.Yellow + "• Coins (C): " + gd.White + "1-15 points\n" + gd.Reset)
	fmt.Print(gd.Red + "• Staves (T): " + gd.White + "1-15 points\n" + gd.Reset)
	fmt.Print(gd.Magenta + "• Arcana: " + gd.White + "Negative value cards\n\n" + gd.Reset)

	fmt.Print(gd.Yellow + "Press any key to return to menu..." + gd.Reset)
	waitForKey()
}

// showStats displays player statistics (placeholder)
func showStats() {
	gd.ClearScreen()
	fmt.Print(gd.CyanHi + "═══════════════════════════════════════════\n" + gd.Reset)
	fmt.Print(gd.CyanHi + "            PLAYER STATISTICS\n" + gd.Reset)
	fmt.Print(gd.CyanHi + "═══════════════════════════════════════════\n\n" + gd.Reset)

	fmt.Print(gd.Yellow + "Player: " + gd.CyanHi + game.User.Alias + gd.Reset + "\n\n")

	// Placeholder stats - in a real implementation, these would be saved/loaded
	fmt.Print(gd.White + "Games Played: " + gd.YellowHi + "0\n" + gd.Reset)
	fmt.Print(gd.White + "Games Won: " + gd.YellowHi + "0\n" + gd.Reset)
	fmt.Print(gd.White + "Pure Sabaccs: " + gd.YellowHi + "0\n" + gd.Reset)
	fmt.Print(gd.White + "Idiot's Arrays: " + gd.YellowHi + "0\n" + gd.Reset)
	fmt.Print(gd.White + "Bomb Outs: " + gd.YellowHi + "0\n" + gd.Reset)
	fmt.Print(gd.White + "Credits Won: " + gd.YellowHi + "0\n" + gd.Reset)
	fmt.Print(gd.White + "Credits Lost: " + gd.YellowHi + "0\n\n" + gd.Reset)

	fmt.Print(gd.Red + "Statistics tracking not yet implemented.\n" + gd.Reset)
	fmt.Print(gd.Red + "This feature will be added in a future version.\n\n" + gd.Reset)

	fmt.Print(gd.Yellow + "Press any key to return to menu..." + gd.Reset)
	waitForKey()
}

// exitGame handles clean exit
func exitGame() {
	gd.ClearScreen()
	gd.MoveCursor(1, game.User.H/2)
	fmt.Print(gd.CyanHi + "Thanks for playing Sabacc!" + gd.Reset)
	gd.MoveCursor(1, game.User.H/2+2)
	fmt.Print(gd.Yellow + "May the Force be with you, " + game.User.Alias + "!" + gd.Reset)
	gd.MoveCursor(1, game.User.H-1)
	time.Sleep(2 * time.Second)
	os.Exit(0)
}

// displayWelcomeMessage shows a welcome message with game info
func displayWelcomeMessage() {
	fmt.Printf("%sWelcome to the Sabacc tables, %s%s%s!\n\n",
		gd.Cyan, gd.CyanHi, game.User.Alias, gd.Reset)
	fmt.Printf("%sTerminal: %s%dx%d%s  ",
		gd.Yellow, gd.YellowHi, game.User.W, game.User.H, gd.Reset)
	fmt.Printf("%sNode: %s%d%s  ",
		gd.Yellow, gd.YellowHi, game.User.NodeNum, gd.Reset)
	fmt.Printf("%sTime: %s%dm%s\n\n",
		gd.Yellow, gd.YellowHi, game.User.TimeLeft, gd.Reset)
}

// clearStatusLine clears the bottom status line
func clearStatusLine() {
	gd.MoveCursor(1, game.User.H)
	fmt.Print(gd.EraseLine)
}

// showStatusLine displays a status message at the bottom
func showStatusLine(message string) {
	clearStatusLine()
	gd.MoveCursor(1, game.User.H)
	fmt.Print(gd.Yellow + message + gd.Reset)
}

// animateCardDeal provides a simple animation effect (placeholder)
func animateCardDeal() {
	fmt.Print(gd.Yellow + "Dealing" + gd.Reset)
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
		gd.Magenta, gd.MagentaHi, len(game.Deck.Cards), gd.Reset)
}

// checkBombOut checks if a hand is bombed out
func checkBombOut(hand []Card) bool {
	total := calculateHandTotal(hand)
	return total > 23 || total < -23 || total == 0
}

// formatCredits formats credits with proper color coding
func formatCredits(amount int) string {
	if amount > 0 {
		return gd.GreenHi + fmt.Sprintf("+%d", amount) + gd.Reset
	} else if amount < 0 {
		return gd.RedHi + fmt.Sprintf("%d", amount) + gd.Reset
	}
	return gd.Yellow + "0" + gd.Reset
}

// displayHandValue shows the current hand value with color coding
func displayHandValue(total int) string {
	if total == 23 {
		return gd.GreenHi + "23 (SABACC!)" + gd.Reset
	} else if total > 23 || total < -23 || total == 0 {
		return gd.RedHi + fmt.Sprintf("%d (BOMB!)", total) + gd.Reset
	} else if total >= 20 && total <= 22 {
		return gd.YellowHi + fmt.Sprintf("%d", total) + gd.Reset
	} else if total >= -22 && total <= -20 {
		return gd.YellowHi + fmt.Sprintf("%d", total) + gd.Reset
	} else {
		return gd.White + fmt.Sprintf("%d", total) + gd.Reset
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

// playSound plays a sound effect (placeholder - would need audio library)
func playSound(soundType string) {
	// Placeholder for sound effects
	// Could use terminal bell or external sound commands
	switch soundType {
	case "deal":
		fmt.Print("\a") // Terminal bell
	case "win":
		fmt.Print("\a")
	case "lose":
		fmt.Print("\a")
	}
}

// displayAsciiArt shows ASCII art for special events
func displayAsciiArt(artType string) {
	switch artType {
	case "sabacc":
		fmt.Print(gd.GreenHi)
		fmt.Println("  ███████  █████  ██████   █████   ██████  ██████ ")
		fmt.Println("  ██      ██   ██ ██   ██ ██   ██ ██      ██     ")
		fmt.Println("  ███████ ███████ ██████  ███████ ██      ██     ")
		fmt.Println("       ██ ██   ██ ██   ██ ██   ██ ██      ██     ")
		fmt.Println("  ███████ ██   ██ ██████  ██   ██  ██████  ██████ ")
		fmt.Print(gd.Reset)
	case "bomb":
		fmt.Print(gd.RedHi)
		fmt.Println("  ██████   ██████  ███    ███ ██████  ██")
		fmt.Println("  ██   ██ ██    ██ ████  ████ ██   ██ ██")
		fmt.Println("  ██████  ██    ██ ██ ████ ██ ██████  ██")
		fmt.Println("  ██   ██ ██    ██ ██  ██  ██ ██   ██   ")
		fmt.Println("  ██████   ██████  ██      ██ ██████  ██")
		fmt.Print(gd.Reset)
	}
}
