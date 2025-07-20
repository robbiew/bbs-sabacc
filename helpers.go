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

// showRules displays the Classic Sabacc rules (updated for West End Games rules)
func showRules() {
	gd.ClearScreen()
	fmt.Print(gd.CyanHi + "═══════════════════════════════════════════\n" + gd.Reset)
	fmt.Print(gd.CyanHi + "           CLASSIC SABACC RULES\n" + gd.Reset)
	fmt.Print(gd.CyanHi + "            West End Games (1989)\n" + gd.Reset)
	fmt.Print(gd.CyanHi + "═══════════════════════════════════════════\n\n" + gd.Reset)

	fmt.Print(gd.Yellow + "OBJECTIVE:\n" + gd.Reset)
	fmt.Print(gd.White + "Get the highest card total that is ≤ 23.\n\n" + gd.Reset)

	fmt.Print(gd.Yellow + "WINNING HANDS:\n" + gd.Reset)
	fmt.Print(gd.Green + "• Pure Sabacc: " + gd.White + "Exactly 23 points (wins Sabacc Pot!)\n" + gd.Reset)
	fmt.Print(gd.Green + "• Idiot's Array: " + gd.White + "Idiot + 2 + 3 (literal 23)\n" + gd.Reset)
	fmt.Print(gd.White + "  " + gd.GreenHi + "Idiot's Array beats Pure Sabacc!" + gd.Reset + "\n")
	fmt.Print(gd.White + "• Highest total ≤ 23 wins Hand Pot\n\n" + gd.Reset)

	fmt.Print(gd.Yellow + "BOMB OUT CONDITIONS:\n" + gd.Reset)
	fmt.Print(gd.Red + "• Over 23 points\n" + gd.Reset)
	fmt.Print(gd.Red + "• Under -23 points\n" + gd.Reset)
	fmt.Print(gd.Red + "• Exactly 0 points\n" + gd.Reset)
	fmt.Print(gd.White + "Penalty: Pay Hand Pot amount to Sabacc Pot\n\n" + gd.Reset)

	// Show second page of rules
	fmt.Print(gd.Yellow + "Press any key for turn structure..." + gd.Reset)
	waitForKey()

	gd.ClearScreen()
	fmt.Print(gd.CyanHi + "═══════════════════════════════════════════\n" + gd.Reset)
	fmt.Print(gd.CyanHi + "            TURN STRUCTURE\n" + gd.Reset)
	fmt.Print(gd.CyanHi + "═══════════════════════════════════════════\n\n" + gd.Reset)

	fmt.Print(gd.Yellow + "EACH TURN HAS 4 PHASES:\n" + gd.Reset)
	fmt.Print(gd.White + "1. " + gd.CyanHi + "BET:" + gd.White + " Check/Call, Raise, or Fold\n" + gd.Reset)
	fmt.Print(gd.White + "2. " + gd.CyanHi + "ROLL:" + gd.White + " Roll dice (doubles = Sabacc Shift!)\n" + gd.Reset)
	fmt.Print(gd.White + "3. " + gd.CyanHi + "CALL:" + gd.White + " Others may call the hand\n" + gd.Reset)
	fmt.Print(gd.White + "4. " + gd.CyanHi + "DRAW:" + gd.White + " Gain, Trade, Stand, or Static Field\n\n" + gd.Reset)

	fmt.Print(gd.Yellow + "CALLING THE HAND:\n" + gd.Reset)
	fmt.Print(gd.White + "• Cannot call until minimum rounds completed\n" + gd.Reset)
	fmt.Print(gd.White + "• Can only call during another player's turn\n" + gd.Reset)
	fmt.Print(gd.White + "• Triggers final showdown\n\n" + gd.Reset)

	fmt.Print(gd.Yellow + "SABACC SHIFT:\n" + gd.Reset)
	fmt.Print(gd.White + "• Triggered by rolling doubles\n" + gd.Reset)
	fmt.Print(gd.White + "• All hands reshuffled and redealt\n" + gd.Reset)
	fmt.Print(gd.White + "• " + gd.MagentaHi + "Static Field cards are protected!" + gd.Reset + "\n\n")

	// Show third page of rules
	fmt.Print(gd.Yellow + "Press any key for deck information..." + gd.Reset)
	waitForKey()

	gd.ClearScreen()
	fmt.Print(gd.CyanHi + "═══════════════════════════════════════════\n" + gd.Reset)
	fmt.Print(gd.CyanHi + "              DECK & CARDS\n" + gd.Reset)
	fmt.Print(gd.CyanHi + "═══════════════════════════════════════════\n\n" + gd.Reset)

	fmt.Print(gd.Yellow + "76-CARD DECK:\n" + gd.Reset)
	fmt.Print(gd.White + "60 Numbered Cards (1-15 in each suit)\n" + gd.Reset)
	fmt.Print(gd.White + "16 Arcana Cards (negative values, 2 copies each)\n\n" + gd.Reset)

	fmt.Print(gd.Yellow + "POSITIVE SUITS:\n" + gd.Reset)
	fmt.Print(gd.Blue + "• Sabers (S): " + gd.White + "1-15 points\n" + gd.Reset)
	fmt.Print(gd.Green + "• Flasks (F): " + gd.White + "1-15 points\n" + gd.Reset)
	fmt.Print(gd.Yellow + "• Coins (C): " + gd.White + "1-15 points\n" + gd.Reset)
	fmt.Print(gd.Red + "• Staves (T): " + gd.White + "1-15 points\n\n" + gd.Reset)

	fmt.Print(gd.Yellow + "ARCANA CARDS (Negative):\n" + gd.Reset)
	fmt.Print(gd.Magenta + "Death(-1), Strength(-2), Moderation(-3), Evil One(-4)\n" + gd.Reset)
	fmt.Print(gd.Magenta + "Justice(-5), Queen(-6), Endurance(-7), Balance(-8)\n" + gd.Reset)
	fmt.Print(gd.Magenta + "Demise(-9), Destruction(-10), Despair(-11), Failure(-12)\n" + gd.Reset)
	fmt.Print(gd.Magenta + "Futility(-13), Mistress(-14), " + gd.MagentaHi + "Idiot(-15)" + gd.Reset + gd.Magenta + ", Star(-17)\n\n" + gd.Reset)

	fmt.Print(gd.Yellow + "ANTE & POTS:\n" + gd.Reset)
	fmt.Print(gd.White + "• Both players ante into " + gd.CyanHi + "both" + gd.White + " pots\n" + gd.Reset)
	fmt.Print(gd.White + "• Hand Pot: Won by best hand ≤ 23\n" + gd.Reset)
	fmt.Print(gd.White + "• Sabacc Pot: Won by Pure Sabacc or Idiot's Array\n" + gd.Reset)
	fmt.Print(gd.White + "• Fold penalty: 1 credit to Sabacc Pot\n\n" + gd.Reset)

	fmt.Print(gd.Yellow + "Press any key to return to menu..." + gd.Reset)
	waitForKey()
}

// showStats displays player statistics (updated to mention Classic Sabacc)
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
	fmt.Print(gd.White + "Sabacc Shifts Survived: " + gd.YellowHi + "0\n" + gd.Reset)
	fmt.Print(gd.White + "Credits Won: " + gd.YellowHi + "0\n" + gd.Reset)
	fmt.Print(gd.White + "Credits Lost: " + gd.YellowHi + "0\n\n" + gd.Reset)

	fmt.Print(gd.Red + "Statistics tracking not yet implemented.\n" + gd.Reset)
	fmt.Print(gd.Red + "This feature will be added in a future version.\n\n" + gd.Reset)

	fmt.Print(gd.Cyan + "Classic Sabacc follows the original West End Games\n" + gd.Reset)
	fmt.Print(gd.Cyan + "rules from the 1989 Crisis on Cloud City supplement.\n\n" + gd.Reset)

	fmt.Print(gd.Yellow + "Press any key to return to menu..." + gd.Reset)
	waitForKey()
}

// exitGame handles clean exit (updated message)
func exitGame() {
	gd.ClearScreen()
	gd.MoveCursor(1, game.User.H/2)
	fmt.Print(gd.CyanHi + "Thanks for playing Classic Sabacc!" + gd.Reset)
	gd.MoveCursor(1, game.User.H/2+2)
	fmt.Print(gd.Yellow + "May the Force be with you, " + game.User.Alias + "!" + gd.Reset)
	gd.MoveCursor(1, game.User.H/2+4)
	fmt.Print(gd.White + "West End Games Rules (1989)" + gd.Reset)
	gd.MoveCursor(1, game.User.H-1)
	time.Sleep(2 * time.Second)
	os.Exit(0)
}

// displayWelcomeMessage shows a welcome message with game info
func displayWelcomeMessage() {
	fmt.Printf("%sWelcome to the Classic Sabacc tables, %s%s%s!\n\n",
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
	case "shift":
		fmt.Print("\a") // Sound for Sabacc Shift
	case "sabacc":
		fmt.Print("\a\a") // Double bell for Pure Sabacc
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
	case "shift":
		fmt.Print(gd.YellowHi)
		fmt.Println("  ███████ ██   ██ ██ ███████ ████████ ██")
		fmt.Println("  ██      ██   ██ ██ ██         ██    ██")
		fmt.Println("  ███████ ███████ ██ █████      ██    ██")
		fmt.Println("       ██ ██   ██ ██ ██         ██      ")
		fmt.Println("  ███████ ██   ██ ██ ██         ██    ██")
		fmt.Print(gd.Reset)
	}
}
