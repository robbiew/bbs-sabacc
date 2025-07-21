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
	User          gd.User
	Deck          Deck
	Players       []Player
	HandPot       int
	SabaccPot     int
	CurrentBet    int
	Round         int
	Turn          int
	Dealer        int
	Called        bool
	GameOver      bool
	MinRounds     int  // Minimum rounds before calling allowed
	BettingPhase  bool // True during betting, false during play
	ShiftOccurred bool // Track if shift happened this turn
}

// Player represents a player in the game
type Player struct {
	Name        string
	Hand        []Card
	StaticField []Card
	Credits     int
	Folded      bool
	BombedOut   bool
	HasActed    bool // Track if player has acted this round
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
		MinRounds: 1, // Classic rules: minimum 1-4 rounds
	}

	// Show title screen
	showTitleScreen()

	// Main menu loop
	mainMenu()
}

func showTitleScreen() {
	gd.ClearScreen()
	gd.MoveCursor(0, 0)

	// Show the sabacc ASCII art instead of the simple text title
	centerY := game.User.H / 2
	gd.MoveCursor(1, centerY-6)
	displayAsciiArt("sabacc")

	gd.MoveCursor(1, centerY+2)
	fmt.Print(gd.Cyan + "Classic 76-Card Sabacc for BBS" + gd.Reset)
	gd.MoveCursor(1, centerY+4)
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

	// Initialize deck (76 cards)
	game.Deck = NewDeck()
	game.Deck.Shuffle()

	// Reset game state
	game.Round = 0
	game.Turn = 1 // Start with player to dealer's left
	game.Dealer = 0
	game.Called = false
	game.GameOver = false
	game.BettingPhase = false
	game.ShiftOccurred = false

	// Setup players
	game.Players = []Player{
		{Name: game.User.Alias, Credits: 1000, Hand: []Card{}, StaticField: []Card{}, Folded: false, BombedOut: false},
		{Name: "Droid Dealer", Credits: 1000, Hand: []Card{}, StaticField: []Card{}, Folded: false, BombedOut: false},
	}

	// ANTE PHASE - Both players ante into both pots
	anteAmount := 10
	game.HandPot = anteAmount * 2
	game.SabaccPot = anteAmount * 2

	// Deduct ante from players
	for i := range game.Players {
		game.Players[i].Credits -= anteAmount * 2 // Ante goes to both pots
	}

	fmt.Printf("%sAnting %d credits to each pot...%s\n", gd.Green, anteAmount, gd.Reset)
	time.Sleep(2 * time.Second)

	// DEALING ROUND - Deal 2 cards to each player
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

		// Check if only one player remains (others folded)
		activePlayers := 0
		lastActivePlayer := -1
		for i, player := range game.Players {
			if !player.Folded {
				activePlayers++
				lastActivePlayer = i
			}
		}

		// If only one player left, they win immediately
		if activePlayers <= 1 {
			if lastActivePlayer >= 0 {
				fmt.Printf("\n%s%s wins by default (others folded)!%s\n",
					gd.GreenHi, game.Players[lastActivePlayer].Name, gd.Reset)
				game.Players[lastActivePlayer].Credits += game.HandPot

				// Check if they also get Sabacc pot (if they have special hand)
				if !game.Players[lastActivePlayer].Folded {
					total := calculateHandTotal(game.Players[lastActivePlayer].Hand)
					if total == 23 {
						fmt.Printf("%s%s also wins the Sabacc Pot! (Pure Sabacc)%s\n",
							gd.GreenHi, game.Players[lastActivePlayer].Name, gd.Reset)
						game.Players[lastActivePlayer].Credits += game.SabaccPot
						game.SabaccPot = 0
						displayAsciiArt("sabacc")
					} else if isIdiotsArray(game.Players[lastActivePlayer].Hand) {
						fmt.Printf("%s%s also wins the Sabacc Pot! (Idiot's Array)%s\n",
							gd.GreenHi, game.Players[lastActivePlayer].Name, gd.Reset)
						game.Players[lastActivePlayer].Credits += game.SabaccPot
						game.SabaccPot = 0
						displayAsciiArt("sabacc")
					}
				}

				time.Sleep(3 * time.Second)
			}
			game.GameOver = true
			break // Exit the game loop immediately
		}

		// Check if current player is folded and skip their turn
		if game.Players[game.Turn].Folded {
			nextTurn()
			continue
		}

		if game.Turn == 0 { // Human player
			// Only process turn if player hasn't folded
			if !game.Players[0].Folded {
				handlePlayerTurn()
			} else {
				// Player folded, skip to next turn
				nextTurn()
				continue
			}
		} else { // Computer player
			// Only process turn if computer hasn't folded
			if !game.Players[game.Turn].Folded {
				// Computer can call if it's round 2 or later
				canCall := game.Round >= 2
				handleComputerTurn(canCall)
			} else {
				// Computer folded, skip to next turn
				nextTurn()
				continue
			}
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

// Also fix the handlePlayerTurn function to properly handle folding
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
			return // Don't continue with dice roll when calling
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
		return // Don't continue with dice roll when folding
	}

	// Only roll dice if player didn't fold or call
	if !playerRef.Folded && !game.Called {
		rollForShift()
	}
}

func handleBettingPhase() {
	fmt.Printf("\n%sBETTING PHASE%s\n", gd.Yellow, gd.Reset)
	fmt.Print(gd.Yellow + "[" + gd.YellowHi + "C" + gd.Yellow + "] " + gd.White + "Check/Call\n" + gd.Reset)
	fmt.Print(gd.Yellow + "[" + gd.YellowHi + "R" + gd.Yellow + "] " + gd.White + "Raise\n" + gd.Reset)
	fmt.Print(gd.Yellow + "[" + gd.YellowHi + "F" + gd.Yellow + "] " + gd.White + "Fold\n\n" + gd.Reset)
	fmt.Print(gd.Green + "Choice: " + gd.Reset)

	char, _, err := getKeyWithTimeout()
	if err != nil {
		return
	}

	switch char {
	case 'c', 'C':
		// Call current bet (if any)
		fmt.Printf("\n%sYou check/call.%s\n", gd.Green, gd.Reset)
		time.Sleep(1 * time.Second)
	case 'r', 'R':
		// Raise bet
		raiseAmount := 10 // Simple raise amount
		game.Players[0].Credits -= raiseAmount
		game.HandPot += raiseAmount
		game.CurrentBet += raiseAmount
		fmt.Printf("\n%sYou raise by %d credits.%s\n", gd.Green, raiseAmount, gd.Reset)
		time.Sleep(1 * time.Second)
	case 'f', 'F':
		// Fold
		game.Players[0].Folded = true
		game.Players[0].Credits -= 1 // Fold penalty to Sabacc pot
		game.SabaccPot += 1
		fmt.Printf("\n%sYou folded.%s\n", gd.Red, gd.Reset)
		time.Sleep(1 * time.Second)

		// Check if only one player remains after folding
		checkForGameEnd()
	}
}

// Add this helper function to check if game should end
func checkForGameEnd() {
	activePlayers := 0
	lastActivePlayer := -1

	for i, player := range game.Players {
		if !player.Folded {
			activePlayers++
			lastActivePlayer = i
		}
	}

	// If only one player left, they win immediately
	if activePlayers <= 1 {
		if lastActivePlayer >= 0 {
			fmt.Printf("\n%s%s wins by default (others folded)!%s\n",
				gd.GreenHi, game.Players[lastActivePlayer].Name, gd.Reset)
			game.Players[lastActivePlayer].Credits += game.HandPot

			// Also award Sabacc pot if it exists
			if game.SabaccPot > 0 {
				fmt.Printf("%s%s also wins the Sabacc Pot! (+%d credits)%s\n",
					gd.GreenHi, game.Players[lastActivePlayer].Name, game.SabaccPot, gd.Reset)
				game.Players[lastActivePlayer].Credits += game.SabaccPot
			}

			time.Sleep(2 * time.Second)
		}
		game.GameOver = true
	}
}

func handleCallPhase() {
	// In single player vs computer, computer decides whether to call
	if game.Turn == 0 {
		fmt.Printf("\n%sAsking if anyone wants to call the hand...%s\n", gd.Yellow, gd.Reset)
		time.Sleep(1 * time.Second)
		// Computer AI decides - simplified logic
		computerTotal := calculateHandTotal(game.Players[1].Hand)
		if computerTotal >= 20 && computerTotal <= 23 {
			game.Called = true
			fmt.Printf("\n%s%s calls the hand!%s\n", gd.GreenHi, game.Players[1].Name, gd.Reset)
			time.Sleep(1 * time.Second)
		}
	}
}

func handleDrawPhase() {
	fmt.Printf("\n%sDRAW PHASE%s\n", gd.Yellow, gd.Reset)
	fmt.Print(gd.Yellow + "[" + gd.YellowHi + "G" + gd.Yellow + "] " + gd.White + "Gain (draw card)\n" + gd.Reset)
	fmt.Print(gd.Yellow + "[" + gd.YellowHi + "T" + gd.Yellow + "] " + gd.White + "Trade card\n" + gd.Reset)
	fmt.Print(gd.Yellow + "[" + gd.YellowHi + "S" + gd.Yellow + "] " + gd.White + "Stand (do nothing)\n" + gd.Reset)
	fmt.Print(gd.Yellow + "[" + gd.YellowHi + "F" + gd.Yellow + "] " + gd.White + "Static Field\n\n" + gd.Reset)
	fmt.Print(gd.Green + "Choice: " + gd.Reset)

	char, _, err := getKeyWithTimeout()
	if err != nil {
		return
	}

	playerRef := &game.Players[0]

	switch char {
	case 'g', 'G':
		// Gain: Draw one card
		if len(game.Deck.Cards) > 0 {
			card := game.Deck.Deal()
			playerRef.Hand = append(playerRef.Hand, card)
			fmt.Printf("\n%sYou drew: %s[%s]%s\n", gd.Green, getCardColor(card), card.String(), gd.Reset)
			time.Sleep(1 * time.Second)
		}
	case 't', 'T':
		// Trade: Discard one card, draw one card
		handleTradeCard()
	case 's', 'S':
		// Stand: Do nothing
		fmt.Printf("\n%sYou stand.%s\n", gd.Green, gd.Reset)
		time.Sleep(1 * time.Second)
	case 'f', 'F':
		// Static Field management
		handleStaticField()
	}
}

func handleComputerTurn(canCall bool) {
	computer := &game.Players[1]

	if computer.Folded {
		return
	}

	time.Sleep(2 * time.Second) // Simulate thinking

	fmt.Printf("\n%s=== %s's Turn ===%s\n", gd.CyanHi, computer.Name, gd.Reset)

	// PHASE 1: BET (simplified AI)
	fmt.Printf("\n%s%s checks.%s\n", gd.Yellow, computer.Name, gd.Reset)
	time.Sleep(1 * time.Second)

	// PHASE 2: ROLL
	rollForShift()

	// PHASE 3: CALL
	if canCall {
		total := calculateHandTotal(computer.Hand)
		if total >= 20 && total <= 23 && game.Round >= 2 {
			game.Called = true
			fmt.Printf("\n%s%s calls the hand!%s\n", gd.GreenHi, computer.Name, gd.Reset)
			time.Sleep(1 * time.Second)
			return
		}
	}

	// PHASE 4: DRAW
	if !game.Called {
		total := calculateHandTotal(computer.Hand)

		if total > 20 || total < -20 {
			// Risky hand, might fold or try to improve
			if total > 23 || total < -23 {
				computer.Folded = true
				computer.Credits -= 1
				game.SabaccPot += 1
				fmt.Printf("\n%s%s folds.%s\n", gd.Red, computer.Name, gd.Reset)
			} else if len(computer.Hand) > 2 {
				// Trade a card
				computer.Hand = computer.Hand[1:] // Remove first card (simplified)
				if len(game.Deck.Cards) > 0 {
					card := game.Deck.Deal()
					computer.Hand = append(computer.Hand, card)
					fmt.Printf("\n%s%s trades a card.%s\n", gd.Yellow, computer.Name, gd.Reset)
				}
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
	}

	time.Sleep(1 * time.Second)
}

func rollForShift() {
	// Roll two dice - shift occurs on doubles
	dice1 := (time.Now().UnixNano() % 6) + 1
	dice2 := ((time.Now().UnixNano() / 1000) % 6) + 1

	fmt.Printf("\n%sRolling dice:%s %s%d%s, %s%d%s",
		gd.Yellow, gd.Reset, gd.YellowHi, dice1, gd.Reset, gd.YellowHi, dice2, gd.Reset)

	if dice1 == dice2 {
		fmt.Printf("\n%sSABACC SHIFT! All hands shuffled!%s\n", gd.RedHi, gd.Reset)

		// Collect all cards not in static field
		allCards := []Card{}
		for i := range game.Players {
			if !game.Players[i].Folded {
				// Cards in static field are protected
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

		// Deal new hands (maintain original hand size)
		for i := range game.Players {
			if !game.Players[i].Folded {
				originalSize := 2                                       // Original hand size was 2
				cardsNeeded := originalSize - len(game.Players[i].Hand) // Subtract static field cards
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
	// Move to next player
	originalTurn := game.Turn

	for {
		game.Turn = (game.Turn + 1) % len(game.Players)

		// If we've completed a full round (back to player after dealer)
		if game.Turn == (game.Dealer+1)%len(game.Players) {
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

	// Show all hands and determine winner
	winner := -1
	bestScore := -999
	sabaccWinner := -1
	bombedOutPlayers := []int{}

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

		// Check for special hands (Sabacc Pot winners)
		if isIdiotsArray(playerData.Hand) {
			fmt.Printf(" %s(IDIOT'S ARRAY!)%s\n", gd.GreenHi, gd.Reset)
			displayAsciiArt("sabacc") // Show sabacc art for special hands
			bestScore = 1000
			sabaccWinner = i            // Idiot's Array beats Pure Sabacc
			time.Sleep(3 * time.Second) // Let player see the art
		} else if total == 23 && sabaccWinner == -1 {
			fmt.Printf(" %s(PURE SABACC!)%s\n", gd.GreenHi, gd.Reset)
			displayAsciiArt("sabacc") // Show sabacc art for Pure Sabacc
			sabaccWinner = i
			time.Sleep(3 * time.Second) // Let player see the art
		} else if total > 23 || total < -23 || total == 0 {
			fmt.Printf(" %s(BOMBED OUT!)%s\n", gd.Red, gd.Reset)
			displayAsciiArt("bomb") // Show bomb art for bomb outs
			game.Players[i].BombedOut = true
			bombedOutPlayers = append(bombedOutPlayers, i)
			// Bombed out player pays Hand Pot amount to Sabacc Pot
			penalty := game.HandPot
			game.Players[i].Credits -= penalty
			game.SabaccPot += penalty
			time.Sleep(2 * time.Second) // Let player see the bomb art
		} else if total <= 23 && total > bestScore {
			winner = i
			bestScore = total
			fmt.Println()
		} else {
			fmt.Println()
		}
	}

	fmt.Println()

	// Determine winners and distribute pots
	if sabaccWinner >= 0 {
		// Someone won with Pure Sabacc or Idiot's Array
		fmt.Printf("%s%s wins both pots with a special hand!%s\n",
			gd.GreenHi, game.Players[sabaccWinner].Name, gd.Reset)
		fmt.Printf("%s+%d credits (Hand Pot) +%d credits (Sabacc Pot)%s\n",
			gd.GreenHi, game.HandPot, game.SabaccPot, gd.Reset)

		game.Players[sabaccWinner].Credits += game.HandPot + game.SabaccPot
		game.SabaccPot = 0 // Reset Sabacc pot

		// Display celebration art one more time
		displayAsciiArt("sabacc")
		time.Sleep(2 * time.Second)

	} else if winner >= 0 {
		// Regular hand winner
		fmt.Printf("%s%s wins the hand! (+%d credits)%s\n",
			gd.GreenHi, game.Players[winner].Name, game.HandPot, gd.Reset)
		game.Players[winner].Credits += game.HandPot

	} else {
		// Everyone bombed out or folded
		fmt.Printf("%sNo winner! Hand pot goes to Sabacc pot.%s\n", gd.Yellow, gd.Reset)
		game.SabaccPot += game.HandPot

		// Show bomb art for the chaos
		if len(bombedOutPlayers) > 1 {
			fmt.Printf("\n%sEveryone bombed out!%s\n", gd.RedHi, gd.Reset)
			displayAsciiArt("bomb")
			time.Sleep(2 * time.Second)
		}
	}

	// Show penalty summary if anyone bombed out
	if len(bombedOutPlayers) > 0 {
		fmt.Printf("\n%sBomb Out Penalties:%s\n", gd.Red, gd.Reset)
		for _, playerIndex := range bombedOutPlayers {
			fmt.Printf("%s%s paid %d credits to Sabacc Pot%s\n",
				gd.Red, game.Players[playerIndex].Name, game.HandPot, gd.Reset)
		}
	}

	game.GameOver = true
}

// Fixed display functions for Sabacc game

func displayGameScreen() {
	gd.ClearScreen()
	gd.MoveCursor(0, 0)

	// Game header
	fmt.Print(gd.CyanHi + "═══════════════════════════════════════════\n" + gd.Reset)
	fmt.Print(gd.CyanHi + "                SABACC GAME\n" + gd.Reset)
	fmt.Print(gd.CyanHi + "═══════════════════════════════════════════\n\n" + gd.Reset)

	// Pot information - Fix: ensure all values are integers
	fmt.Printf("%sHand Pot:%s %s%d%s    %sSabacc Pot:%s %s%d%s    %sRound:%s %s%d%s\n\n",
		gd.Yellow, gd.Reset, gd.YellowHi, game.HandPot, gd.Reset,
		gd.Yellow, gd.Reset, gd.YellowHi, game.SabaccPot, gd.Reset,
		gd.Yellow, gd.Reset, gd.YellowHi, game.Round, gd.Reset)

	// Display opponent's hand (face down)
	if len(game.Players) > 1 {
		opponent := &game.Players[1]
		fmt.Printf("%s%s's Hand:%s ", gd.Cyan, opponent.Name, gd.Reset)
		for i := 0; i < len(opponent.Hand); i++ {
			fmt.Printf("%s[??]%s ", gd.Red, gd.Reset)
		}
		fmt.Printf("  %sCredits:%s %s%d%s\n\n",
			gd.Yellow, gd.Reset, gd.YellowHi, opponent.Credits, gd.Reset)
	}

	// Display player's hand - Fix: ensure proper card display and total calculation
	if len(game.Players) > 0 {
		playerHand := &game.Players[0]
		fmt.Printf("%sYour Hand:%s ", gd.Cyan, gd.Reset)

		// Calculate total while displaying cards
		total := 0
		for _, card := range playerHand.Hand {
			fmt.Printf("%s[%s]%s ", getCardColor(card), card.String(), gd.Reset)
			total += card.Value
		}

		// Add static field cards to total
		for _, card := range playerHand.StaticField {
			total += card.Value
		}

		fmt.Printf("  %sTotal:%s %s%d%s  %sCredits:%s %s%d%s\n\n",
			gd.Yellow, gd.Reset, displayHandValue(total), total, gd.Reset,
			gd.Yellow, gd.Reset, gd.YellowHi, playerHand.Credits, gd.Reset)

		// Display static field if any
		if len(playerHand.StaticField) > 0 {
			fmt.Printf("%sStatic Field:%s ", gd.Magenta, gd.Reset)
			for _, card := range playerHand.StaticField {
				fmt.Printf("%s[%s]%s ", getCardColor(card), card.String(), gd.Reset)
			}
			fmt.Println()
		}
	}

	fmt.Println()

	// Show whose turn it is
	if game.Turn == 0 {
		fmt.Print(gd.GreenHi + "YOUR TURN\n" + gd.Reset)
	} else if len(game.Players) > game.Turn {
		fmt.Printf("%s%s's turn...%s\n", gd.YellowHi, game.Players[game.Turn].Name, gd.Reset)
	}

	// Show remaining deck info
	fmt.Printf("%sDeck: %s%d cards remaining%s\n",
		gd.Magenta, gd.MagentaHi, len(game.Deck.Cards), gd.Reset)
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
