package main

import (
	_ "embed"
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/eiannone/keyboard"
	gd "github.com/robbiew/godoors"
)

var (
	DropPath string
	game     *SabaccGame

	//go:embed ansi/title.ans
	TitleScreen string
	//go:embed ansi/menu.ans
	MenuScreen string
	//go:embed ansi/game.ans
	GameScreen string
)

// SabaccGame represents the main game state
type SabaccGame struct {
	User       gd.User
	Deck       Deck
	Players    []Player
	HandPot    int
	SabaccPot  int
	CurrentBet int
	Round      int
	Turn       int
	Dealer     int
	Called     bool
	GameOver   bool
}

// Player represents a player in the game
type Player struct {
	Name        string
	Hand        []Card
	StaticField []Card
	Credits     int
	Folded      bool
	BombedOut   bool
}

func init() {
	gd.Idle = 300 // 5 minute idle timeout
}

func main() {
	// Use FLAG to get command line parameters
	pathPtr := flag.String("path", "", "path to door32.sys file")
	flag.Parse()

	if *pathPtr == "" {
		fmt.Fprintf(os.Stderr, "missing path to door32.sys directory: -path\n")
		fmt.Fprintf(os.Stderr, "Usage: %s -path /path/to/dropfile/directory/\n", os.Args[0])
		os.Exit(2)
	}
	DropPath = *pathPtr

	// Get door32.sys, h, w as user object
	u := gd.Initialize(DropPath)

	// Exit if no ANSI capabilities
	if u.Emulation != 1 {
		fmt.Println("Sorry, ANSI is required to play Sabacc...")
		time.Sleep(time.Duration(2) * time.Second)
		os.Exit(0)
	}

	// Initialize keyboard
	if err := keyboard.Open(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer func() {
		_ = keyboard.Close()
	}()

	// Initialize game
	game = &SabaccGame{
		User:      u,
		HandPot:   0,
		SabaccPot: 0,
		Round:     0,
		Turn:      0,
		Dealer:    0,
	}

	// Show title screen
	showTitleScreen()

	// Main menu loop
	mainMenu()
}

func showTitleScreen() {
	gd.ClearScreen()
	gd.MoveCursor(0, 0)

	// Center title text since we don't have ANSI art yet
	centerY := game.User.H / 2
	gd.MoveCursor(1, centerY-3)
	fmt.Print(gd.YellowHi + "███████  █████  ██████   █████   ██████  ██████" + gd.Reset)
	gd.MoveCursor(1, centerY-2)
	fmt.Print(gd.YellowHi + "██      ██   ██ ██   ██ ██   ██ ██      ██" + gd.Reset)
	gd.MoveCursor(1, centerY-1)
	fmt.Print(gd.YellowHi + "███████ ███████ ██████  ███████ ██      ██" + gd.Reset)
	gd.MoveCursor(1, centerY)
	fmt.Print(gd.YellowHi + "     ██ ██   ██ ██   ██ ██   ██ ██      ██" + gd.Reset)
	gd.MoveCursor(1, centerY+1)
	fmt.Print(gd.YellowHi + "███████ ██   ██ ██████  ██   ██  ██████  ██████" + gd.Reset)

	gd.MoveCursor(1, centerY+4)
	fmt.Print(gd.Cyan + "Classic 76-Card Sabacc for BBS" + gd.Reset)
	gd.MoveCursor(1, centerY+6)
	fmt.Print(gd.White + "Welcome, " + gd.CyanHi + game.User.Alias + gd.Reset)

	gd.MoveCursor(1, game.User.H-2)
	fmt.Print(gd.Yellow + "Press any key to continue..." + gd.Reset)

	waitForKey()
}

func mainMenu() {
	for {
		gd.ClearScreen()
		gd.MoveCursor(0, 0)

		// Display menu
		fmt.Print(gd.CyanHi + "═══════════════════════════════════════════\n" + gd.Reset)
		fmt.Print(gd.CyanHi + "              SABACC CANTINA\n" + gd.Reset)
		fmt.Print(gd.CyanHi + "═══════════════════════════════════════════\n\n" + gd.Reset)

		fmt.Print(gd.Yellow + "[" + gd.YellowHi + "N" + gd.Yellow + "] " + gd.White + "New Game\n" + gd.Reset)
		fmt.Print(gd.Yellow + "[" + gd.YellowHi + "R" + gd.Yellow + "] " + gd.White + "Rules\n" + gd.Reset)
		fmt.Print(gd.Yellow + "[" + gd.YellowHi + "S" + gd.Yellow + "] " + gd.White + "Statistics\n" + gd.Reset)
		fmt.Print(gd.Yellow + "[" + gd.YellowHi + "Q" + gd.Yellow + "] " + gd.White + "Quit to BBS\n\n" + gd.Reset)

		fmt.Print(gd.Cyan + "Credits: " + gd.CyanHi + "1000" + gd.Reset + "  ")
		fmt.Print(gd.Cyan + "Time Left: " + gd.CyanHi + strconv.Itoa(game.User.TimeLeft) + "m" + gd.Reset + "\n\n")

		fmt.Print(gd.Green + "Choice: " + gd.Reset)

		char, key, err := getKeyWithTimeout()
		if err != nil {
			continue
		}

		switch {
		case char == 'n' || char == 'N':
			startNewGame()
		case char == 'r' || char == 'R':
			showRules()
		case char == 's' || char == 'S':
			showStats()
		case char == 'q' || char == 'Q' || key == keyboard.KeyEsc:
			exitGame()
		}
	}
}

func startNewGame() {
	gd.ClearScreen()

	// Initialize deck
	game.Deck = NewDeck()
	game.Deck.Shuffle()

	// Reset game state
	game.Round = 0
	game.Turn = 1
	game.Dealer = 0
	game.Called = false
	game.GameOver = false

	// Setup players (for now, just player vs computer)
	game.Players = []Player{
		{Name: game.User.Alias, Credits: 1000, Hand: []Card{}, StaticField: []Card{}, Folded: false, BombedOut: false},
		{Name: "Droid Dealer", Credits: 1000, Hand: []Card{}, StaticField: []Card{}, Folded: false, BombedOut: false},
	}

	// Initial ante
	game.HandPot = 20   // 10 credits from each player
	game.SabaccPot = 20 // 10 credits from each player

	// Deduct ante from players
	game.Players[0].Credits -= 10
	game.Players[1].Credits -= 10

	// Deal initial hands
	for i := 0; i < 2; i++ {
		for j := range game.Players {
			if len(game.Deck.Cards) > 0 {
				card := game.Deck.Deal()
				game.Players[j].Hand = append(game.Players[j].Hand, card)
			}
		}
	}

	fmt.Printf("%sDealing cards...%s\n", gd.Green, gd.Reset)
	time.Sleep(2 * time.Second)

	// Start game loop
	gameLoop()
}

func gameLoop() {
	game.Round = 1
	game.Turn = 1 // Start with player after dealer

	for !game.GameOver {
		displayGameScreen()

		// Check if current player is folded and skip their turn
		if game.Players[game.Turn].Folded {
			nextTurn()
			continue
		}

		// Check if only one player remains (others folded)
		activePlayers := 0
		lastActivePlayer := -1
		for i, player := range game.Players {
			if !player.Folded {
				activePlayers++
				lastActivePlayer = i
			}
		}

		// If only one player left, they win
		if activePlayers <= 1 {
			if lastActivePlayer >= 0 {
				fmt.Printf("\n%s%s wins by default (others folded)!%s\n",
					gd.GreenHi, game.Players[lastActivePlayer].Name, gd.Reset)
				game.Players[lastActivePlayer].Credits += game.HandPot
				time.Sleep(2 * time.Second)
			}
			game.GameOver = true
			break
		}

		if game.Turn == 0 { // Human player
			handlePlayerTurn()
		} else { // Computer player
			handleComputerTurn()
		}

		// Check for game end conditions
		if game.Called {
			resolveHand()
			break
		}

		nextTurn()
	}

	// Show results and return to menu
	showGameResults()
	waitForKey()
}

func displayGameScreen() {
	gd.ClearScreen()
	gd.MoveCursor(0, 0)

	// Game header
	fmt.Print(gd.CyanHi + "═══════════════════════════════════════════\n" + gd.Reset)
	fmt.Print(gd.CyanHi + "                SABACC GAME\n" + gd.Reset)
	fmt.Print(gd.CyanHi + "═══════════════════════════════════════════\n\n" + gd.Reset)

	// Pot information
	fmt.Printf("%sHand Pot:%s %s%d%s    %sSabacc Pot:%s %s%d%s    %sRound:%s %s%d%s\n\n",
		gd.Yellow, gd.Reset, gd.YellowHi, game.HandPot, gd.Reset,
		gd.Yellow, gd.Reset, gd.YellowHi, game.SabaccPot, gd.Reset,
		gd.Yellow, gd.Reset, gd.YellowHi, game.Round, gd.Reset)

	// Display opponent's hand (face down)
	opponent := &game.Players[1]
	fmt.Printf("%s%s's Hand:%s ", gd.Cyan, opponent.Name, gd.Reset)
	for i := 0; i < len(opponent.Hand); i++ {
		fmt.Printf("%s[??]%s ", gd.Red, gd.Reset)
	}
	fmt.Printf("  %sCredits:%s %s%d%s\n\n",
		gd.Yellow, gd.Reset, gd.YellowHi, opponent.Credits, gd.Reset)

	// Display player's hand
	playerHand := &game.Players[0]
	fmt.Printf("%sYour Hand:%s ", gd.Cyan, gd.Reset)
	total := 0
	for _, card := range playerHand.Hand {
		fmt.Printf("%s[%s]%s ", getCardColor(card), card.String(), gd.Reset)
		total += card.Value
	}
	fmt.Printf("  %sTotal:%s %s%d%s  %sCredits:%s %s%d%s\n\n",
		gd.Yellow, gd.Reset, gd.YellowHi, total, gd.Reset,
		gd.Yellow, gd.Reset, gd.YellowHi, playerHand.Credits, gd.Reset)

	// Display static field if any
	if len(playerHand.StaticField) > 0 {
		fmt.Printf("%sStatic Field:%s ", gd.Magenta, gd.Reset)
		for _, card := range playerHand.StaticField {
			fmt.Printf("%s[%s]%s ", getCardColor(card), card.String(), gd.Reset)
		}
		fmt.Println()
	}

	fmt.Println()

	// Show whose turn it is
	if game.Turn == 0 {
		fmt.Print(gd.GreenHi + "YOUR TURN\n" + gd.Reset)
	} else {
		fmt.Printf("%s%s's turn...%s\n", gd.YellowHi, game.Players[game.Turn].Name, gd.Reset)
	}
}

func handlePlayerTurn() {
	fmt.Println()
	fmt.Print(gd.Yellow + "[" + gd.YellowHi + "D" + gd.Yellow + "] " + gd.White + "Draw card\n" + gd.Reset)
	fmt.Print(gd.Yellow + "[" + gd.YellowHi + "T" + gd.Yellow + "] " + gd.White + "Trade card\n" + gd.Reset)
	fmt.Print(gd.Yellow + "[" + gd.YellowHi + "S" + gd.Yellow + "] " + gd.White + "Stand (do nothing)\n" + gd.Reset)
	fmt.Print(gd.Yellow + "[" + gd.YellowHi + "F" + gd.Yellow + "] " + gd.White + "Place card in Static Field\n" + gd.Reset)
	fmt.Print(gd.Yellow + "[" + gd.YellowHi + "C" + gd.Yellow + "] " + gd.White + "Call hand\n" + gd.Reset)
	fmt.Print(gd.Yellow + "[" + gd.YellowHi + "Q" + gd.Yellow + "] " + gd.White + "Fold\n\n" + gd.Reset)
	fmt.Print(gd.Green + "Choice: " + gd.Reset)

	char, _, err := getKeyWithTimeout()
	if err != nil {
		return
	}

	playerRef := &game.Players[0]

	switch char {
	case 'd', 'D':
		if len(game.Deck.Cards) > 0 {
			card := game.Deck.Deal()
			playerRef.Hand = append(playerRef.Hand, card)
			fmt.Printf("\n%sYou drew: %s[%s]%s\n", gd.Green, getCardColor(card), card.String(), gd.Reset)
			time.Sleep(1 * time.Second)
		}
	case 't', 'T':
		handleTradeCard()
	case 's', 'S':
		fmt.Printf("\n%sYou stand.%s\n", gd.Green, gd.Reset)
		time.Sleep(1 * time.Second)
	case 'f', 'F':
		handleStaticField()
	case 'c', 'C':
		if game.Round >= 2 { // Can only call after round 2
			game.Called = true
			fmt.Printf("\n%sYou called the hand!%s\n", gd.GreenHi, gd.Reset)
			time.Sleep(1 * time.Second)
		} else {
			fmt.Printf("\n%sCannot call until round 2!%s\n", gd.Red, gd.Reset)
			time.Sleep(1 * time.Second)
		}
	case 'q', 'Q':
		playerRef.Folded = true
		playerRef.Credits -= 1 // Fold penalty
		game.SabaccPot += 1
		fmt.Printf("\n%sYou folded.%s\n", gd.Red, gd.Reset)
		time.Sleep(1 * time.Second)
	}

	// Roll dice for Sabacc Shift
	if !playerRef.Folded {
		rollForShift()
	}
}

func handleComputerTurn() {
	computer := &game.Players[1]

	// Don't process turn if computer has folded
	if computer.Folded {
		return
	}

	time.Sleep(2 * time.Second) // Simulate thinking

	total := calculateHandTotal(computer.Hand)

	// Simple AI logic
	if total > 20 || total < -20 {
		// Risky hand, might fold or try to improve
		if game.Round > 2 && total > 23 {
			computer.Folded = true
			computer.Credits -= 1
			game.SabaccPot += 1
			fmt.Printf("\n%s%s folds.%s\n", gd.Red, computer.Name, gd.Reset)
		} else if len(computer.Hand) > 2 {
			// Trade a card
			computer.Hand = computer.Hand[1:] // Remove first card
			if len(game.Deck.Cards) > 0 {
				card := game.Deck.Deal()
				computer.Hand = append(computer.Hand, card)
				fmt.Printf("\n%s%s trades a card.%s\n", gd.Yellow, computer.Name, gd.Reset)
			}
		} else {
			fmt.Printf("\n%s%s stands.%s\n", gd.Yellow, computer.Name, gd.Reset)
		}
	} else if total >= 20 && total <= 23 {
		// Good hand, call or stand
		if game.Round >= 2 {
			game.Called = true
			fmt.Printf("\n%s%s calls the hand!%s\n", gd.GreenHi, computer.Name, gd.Reset)
		} else {
			fmt.Printf("\n%s%s stands.%s\n", gd.Yellow, computer.Name, gd.Reset)
		}
	} else {
		// Try to improve hand
		if len(game.Deck.Cards) > 0 {
			card := game.Deck.Deal()
			computer.Hand = append(computer.Hand, card)
			fmt.Printf("\n%s%s draws a card.%s\n", gd.Yellow, computer.Name, gd.Reset)
		} else {
			fmt.Printf("\n%s%s stands.%s\n", gd.Yellow, computer.Name, gd.Reset)
		}
	}

	time.Sleep(1 * time.Second)

	// Only roll for shift if computer hasn't folded
	if !computer.Folded {
		rollForShift()
	}
}

func rollForShift() {
	// Simulate dice roll (1/36 chance of shift - rolling double sixes)
	dice1 := (time.Now().UnixNano() % 6) + 1
	dice2 := ((time.Now().UnixNano() / 1000) % 6) + 1

	fmt.Printf("\n%sDice roll:%s %s%d%s, %s%d%s",
		gd.Yellow, gd.Reset, gd.YellowHi, dice1, gd.Reset, gd.YellowHi, dice2, gd.Reset)

	if dice1 == dice2 {
		fmt.Printf("\n%sSABACC SHIFT! All hands shuffled!%s\n", gd.RedHi, gd.Reset)

		// Collect all cards not in static field
		allCards := []Card{}
		for i := range game.Players {
			if !game.Players[i].Folded {
				// Keep static field cards, collect others
				newHand := make([]Card, len(game.Players[i].StaticField))
				copy(newHand, game.Players[i].StaticField)

				// Add non-static cards to shuffle pile
				for _, card := range game.Players[i].Hand {
					isInStatic := false
					for _, staticCard := range game.Players[i].StaticField {
						if card.Value == staticCard.Value && card.Suit == staticCard.Suit {
							isInStatic = true
							break
						}
					}
					if !isInStatic {
						allCards = append(allCards, card)
					}
				}

				game.Players[i].Hand = newHand
			}
		}

		// Shuffle and redistribute
		game.Deck.Cards = append(game.Deck.Cards, allCards...)
		game.Deck.Shuffle()

		// Deal new hands (2 cards base plus any static field cards)
		for i := range game.Players {
			if !game.Players[i].Folded {
				cardsNeeded := 2 // Start with 2 cards minimum
				for j := 0; j < cardsNeeded; j++ {
					if len(game.Deck.Cards) > 0 {
						card := game.Deck.Deal()
						game.Players[i].Hand = append(game.Players[i].Hand, card)
					}
				}
			}
		}

		time.Sleep(3 * time.Second)
	} else {
		fmt.Printf(" %s(No shift)%s\n", gd.Green, gd.Reset)
		time.Sleep(1 * time.Second)
	}
}

func nextTurn() {
	// Keep cycling through players until we find one who hasn't folded
	originalTurn := game.Turn

	for {
		game.Turn = (game.Turn + 1) % len(game.Players)

		// If we've completed a full round (back to dealer)
		if game.Turn == game.Dealer {
			game.Round++
		}

		// If we found an active player or we've checked everyone, break
		if !game.Players[game.Turn].Folded || game.Turn == originalTurn {
			break
		}
	}
}

func resolveHand() {
	gd.ClearScreen()
	fmt.Print(gd.CyanHi + "═══════════════════════════════════════════\n" + gd.Reset)
	fmt.Print(gd.CyanHi + "                HAND RESULTS\n" + gd.Reset)
	fmt.Print(gd.CyanHi + "═══════════════════════════════════════════\n\n" + gd.Reset)

	// Show all hands
	winner := -1
	bestScore := -999

	for i, playerData := range game.Players {
		if playerData.Folded {
			fmt.Printf("%s%s: FOLDED%s\n", gd.Red, playerData.Name, gd.Reset)
			continue
		}

		total := calculateHandTotal(playerData.Hand)
		fmt.Printf("%s%s:%s ", gd.Cyan, playerData.Name, gd.Reset)
		for _, card := range playerData.Hand {
			fmt.Printf("%s[%s]%s ", getCardColor(card), card.String(), gd.Reset)
		}
		fmt.Printf("= %s%d%s", gd.YellowHi, total, gd.Reset)

		// Check for special hands
		if isIdiotsArray(playerData.Hand) {
			fmt.Printf(" %s(IDIOT'S ARRAY!)%s", gd.GreenHi, gd.Reset)
			winner = i
			bestScore = 1000 // Highest priority
		} else if total == 23 {
			fmt.Printf(" %s(PURE SABACC!)%s", gd.GreenHi, gd.Reset)
			if bestScore < 999 {
				winner = i
				bestScore = 999
			}
		} else if total > 23 || total < -23 || total == 0 {
			fmt.Printf(" %s(BOMBED OUT!)%s", gd.Red, gd.Reset)
			game.Players[i].Credits -= game.HandPot
			game.SabaccPot += game.HandPot
		} else if total <= 23 && total > bestScore && bestScore < 100 {
			winner = i
			bestScore = total
		}
		fmt.Println()
	}

	fmt.Println()

	if winner >= 0 {
		fmt.Printf("%s%s wins the hand! (+%d credits)%s\n",
			gd.GreenHi, game.Players[winner].Name, game.HandPot, gd.Reset)
		game.Players[winner].Credits += game.HandPot

		// Check if they also win the Sabacc Pot
		if bestScore >= 999 {
			fmt.Printf("%s%s also wins the Sabacc Pot! (+%d credits)%s\n",
				gd.GreenHi, game.Players[winner].Name, game.SabaccPot, gd.Reset)
			game.Players[winner].Credits += game.SabaccPot
			game.SabaccPot = 0
		}
	} else {
		fmt.Printf("%sNo winner! Hand pot goes to Sabacc pot.%s\n", gd.Yellow, gd.Reset)
		game.SabaccPot += game.HandPot
	}

	game.GameOver = true
}

func showGameResults() {
	fmt.Println()
	fmt.Printf("%sFinal Credits:%s\n", gd.CyanHi, gd.Reset)
	for _, playerData := range game.Players {
		fmt.Printf("%s: %s%d%s credits\n", playerData.Name, gd.YellowHi, playerData.Credits, gd.Reset)
	}
	fmt.Println()
	fmt.Print(gd.Yellow + "Press any key to return to menu..." + gd.Reset)
}

// Helper functions continue in next part...
