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
)

// SabaccGame represents the main game state
type SabaccGame struct {
	User          gd.User
	CardRenderer  *CardRenderer
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
	MinRounds     int           // Minimum rounds before calling allowed
	BettingPhase  bool          // True during betting, false during play
	ShiftOccurred bool          // Track if shift happened this turn
	Layout        *ScreenLayout // New persistent UI system
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

	cardRenderer := NewCardRenderer()
	if cardRenderer == nil {
		fmt.Println("Warning: Card graphics not available, using ASCII fallback")
		// Game still works without graphics
	}

	// Initialize game with new UI system
	game = &SabaccGame{
		User:         u,
		CardRenderer: cardRenderer,
		HandPot:      0,
		SabaccPot:    0,
		Round:        0,
		Turn:         0,
		Dealer:       0,
		MinRounds:    1,                         // Classic rules: minimum 1-4 rounds
		Layout:       NewScreenLayout(u.W, u.H), // Initialize persistent UI
	}

	// Show title screen
	showTitleScreen()

	// Main menu loop
	mainMenu()
}

func showTitleScreen() {
	gd.ClearScreen()
	gd.MoveCursor(0, 0)

	// Try to display the ANSI title screen first
	if _, err := os.Stat("ansi/title.ans"); err == nil {
		// ANSI file exists, use it
		file, err := os.ReadFile("ansi/title.ans")
		if err == nil {
			fmt.Print(string(file))
		} else {
			// Fall back to embedded title if file read fails
			fmt.Print(TitleScreen)
		}
	} else {
		// ANSI file doesn't exist, use embedded title or simple fallback
		if TitleScreen != "" {
			fmt.Print(TitleScreen)
		} else {
			// Simple fallback if no embedded title
			centerY := game.User.H / 2
			gd.MoveCursor(1, centerY-6)
			displayAsciiArt("sabacc")

			gd.MoveCursor(1, centerY+2)
			fmt.Print(gd.Cyan + "Classic 76-Card Sabacc for BBS" + gd.Reset)
		}
	}

	// Add welcome message and continue prompt at bottom
	gd.MoveCursor(1, game.User.H-4)
	fmt.Print(gd.White + "Welcome, " + gd.CyanHi + game.User.Alias + gd.Reset)

	gd.MoveCursor(1, game.User.H-2)
	fmt.Print(gd.Yellow + "Press any key to continue..." + gd.Reset)

	waitForKey()
}

func mainMenu() {
	for {
		gd.ClearScreen()
		gd.MoveCursor(0, 0)

		// Try to display the ANSI menu screen first
		if _, err := os.Stat("ansi/menu.ans"); err == nil {
			// ANSI file exists, use it
			file, err := os.ReadFile("ansi/menu.ans")
			if err == nil {
				fmt.Print(string(file))
			} else {
				// Fall back to embedded menu if file read fails
				if MenuScreen != "" {
					fmt.Print(MenuScreen)
				} else {
					// Simple fallback menu
					displaySimpleMenu()
				}
			}
		} else {
			// ANSI file doesn't exist, use embedded menu or simple fallback
			if MenuScreen != "" {
				fmt.Print(MenuScreen)
			} else {
				// Simple fallback menu
				displaySimpleMenu()
			}
		}

		// Position cursor for menu options (after ANSI art)
		gd.MoveCursor(1, 4)
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

// Helper function for simple fallback menu
func displaySimpleMenu() {
	fmt.Print(gd.CyanHi + "-------------------------------------------\n" + gd.Reset)
	fmt.Print(gd.CyanHi + "              SABACC CANTINA\n" + gd.Reset)
	fmt.Print(gd.CyanHi + "-------------------------------------------\n\n" + gd.Reset)
}

func startNewGame() {
	// Initialize the persistent UI layout
	game.Layout.InitializeScreen()

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

	// Setup players: 1 human + 4 AI players to match UI layout
	game.Players = []Player{
		{Name: game.User.Alias, Credits: 1000, Hand: []Card{}, StaticField: []Card{}, Folded: false, BombedOut: false}, // Human player
		{Name: "Phoo_ja", Credits: 1000, Hand: []Card{}, StaticField: []Card{}, Folded: false, BombedOut: false},       // AI 1 (top-left)
		{Name: "Rsh-Taac", Credits: 1000, Hand: []Card{}, StaticField: []Card{}, Folded: false, BombedOut: false},      // AI 2 (top-right)
		{Name: "Soladi", Credits: 1000, Hand: []Card{}, StaticField: []Card{}, Folded: false, BombedOut: false},        // AI 3 (bottom-left)
		{Name: "Ky'Ola", Credits: 1000, Hand: []Card{}, StaticField: []Card{}, Folded: false, BombedOut: false},        // AI 4 (bottom-right)
	}

	// ANTE PHASE - Both players ante into both pots
	anteAmount := 10
	game.HandPot = anteAmount * 2
	game.SabaccPot = anteAmount * 2

	// Deduct ante from players
	for i := range game.Players {
		game.Players[i].Credits -= anteAmount * 2 // Ante goes to both pots
	}

	// Show ante message in game log
	game.Layout.LogMessage("Both players ante 10 credits to each pot", "info")
	game.Layout.DisplayMessage("Anting credits to pots...", "info", 0)
	time.Sleep(2 * time.Second)

	// DEALING ROUND - Deal 2 cards to each player
	game.Layout.LogMessage("Dealing 2 cards to each player", "info")
	for i := 0; i < 2; i++ {
		for j := range game.Players {
			if len(game.Deck.Cards) > 0 {
				card := game.Deck.Deal()
				game.Players[j].Hand = append(game.Players[j].Hand, card)
			}
		}
	}

	// Show dealing completion message in game log
	game.Layout.LogMessage("Cards dealt. Round 1 begins.", "info")
	game.Layout.DisplayMessage("Cards dealt. Game begins...", "info", 0)
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

// Updated handlePlayerTurn function using persistent UI
func handlePlayerTurn() {
	// Show the player turn menu using our new UI system
	game.Layout.ShowPlayerTurnMenu(game.Round)

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
			game.Layout.LogMessage("You drew: ["+card.String()+"]", "action")
			game.Layout.DisplayMessage("You drew: ["+card.String()+"]", "success", 0)
			time.Sleep(2 * time.Second)
		} else {
			game.Layout.DisplayMessage("No more cards in deck!", "error", 0)
			time.Sleep(1 * time.Second)
		}
	case 't', 'T':
		handleTradeCard()
	case 's', 'S':
		game.Layout.LogMessage("You stand (no action)", "action")
		game.Layout.DisplayMessage("You stand.", "info", 0)
		time.Sleep(1 * time.Second)
	case 'f', 'F':
		handleStaticField()
	case 'c', 'C':
		if game.Round >= 2 { // Can only call after round 2
			game.Called = true
			game.Layout.LogMessage("You called the hand!", "important")
			game.Layout.DisplayMessage("You called the hand!", "success", 0)
			time.Sleep(2 * time.Second)
			return // Don't continue with dice roll when calling
		} else {
			game.Layout.DisplayMessage("Cannot call until round 2!", "error", 0)
			time.Sleep(2 * time.Second)
		}
	case 'q', 'Q':
		playerRef.Folded = true
		playerRef.Credits -= 1 // Fold penalty
		game.SabaccPot += 1
		game.Layout.LogMessage("You folded (-1 credit penalty)", "important")
		game.Layout.DisplayMessage("You folded.", "warning", 0)
		time.Sleep(2 * time.Second)
		return // Don't continue with dice roll when folding
	default:
		game.Layout.DisplayMessage("Invalid choice!", "error", 0)
		time.Sleep(1 * time.Second)
		return
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
	computer := &game.Players[game.Turn]

	if computer.Folded {
		return
	}

	time.Sleep(1 * time.Second) // Simulate thinking

	// Log the AI action in the game log
	game.Layout.LogMessage(computer.Name+" is thinking...", "info")
	time.Sleep(1 * time.Second)

	// PHASE 1: BET (simplified AI)
	game.Layout.LogMessage(computer.Name+" checks.", "action")
	time.Sleep(500 * time.Millisecond)

	// PHASE 2: ROLL
	rollForShift()

	// PHASE 3: CALL
	if canCall {
		total := calculateHandTotal(computer.Hand)
		if total >= 20 && total <= 23 && game.Round >= 2 {
			game.Called = true
			game.Layout.LogMessage(computer.Name+" calls the hand!", "important")
			time.Sleep(2 * time.Second)
			return
		}
	}

	// PHASE 4: DRAW (AI decision logic)
	if !game.Called {
		total := calculateHandTotal(computer.Hand)

		if total > 20 || total < -20 {
			// Risky hand, might fold or try to improve
			if total > 23 || total < -23 {
				computer.Folded = true
				computer.Credits -= 1
				game.SabaccPot += 1
				game.Layout.LogMessage(computer.Name+" folds.", "important")
			} else if len(computer.Hand) > 2 {
				// Trade a card
				computer.Hand = computer.Hand[1:] // Remove first card (simplified)
				if len(game.Deck.Cards) > 0 {
					card := game.Deck.Deal()
					computer.Hand = append(computer.Hand, card)
					game.Layout.LogMessage(computer.Name+" trades a card.", "action")
				}
			} else {
				game.Layout.LogMessage(computer.Name+" stands.", "action")
			}
		} else {
			// Try to improve hand
			if len(game.Deck.Cards) > 0 {
				card := game.Deck.Deal()
				computer.Hand = append(computer.Hand, card)
				game.Layout.LogMessage(computer.Name+" draws a card.", "action")
			} else {
				game.Layout.LogMessage(computer.Name+" stands.", "action")
			}
		}
	}

	time.Sleep(1 * time.Second)
}

func rollForShift() {
	// Roll two dice - shift occurs on doubles
	dice1 := (time.Now().UnixNano() % 6) + 1
	dice2 := ((time.Now().UnixNano() / 1000) % 6) + 1

	// Send dice roll message to game log
	game.Layout.LogMessage(fmt.Sprintf("Rolling dice: %d, %d (No shift)", dice1, dice2), "info")

	if dice1 == dice2 {
		// Update the message to show shift occurred
		game.Layout.LogMessage(fmt.Sprintf("Rolling dice: %d, %d", dice1, dice2), "info")
		game.Layout.LogMessage("SABACC SHIFT! All hands shuffled!", "important")

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
	// Update header with current game info
	currentPlayerName := ""
	if game.Turn == 0 {
		currentPlayerName = "► YOUR TURN ◄"
	} else if game.Turn < len(game.Players) {
		currentPlayerName = game.Players[game.Turn].Name + " thinking..."
	}

	game.Layout.UpdateHeader(game.Round, game.HandPot, game.SabaccPot, len(game.Deck.Cards), currentPlayerName)

	// Update all AI players (indexes 1-4)
	for i := 1; i < len(game.Players) && i <= 4; i++ {
		aiPlayer := &game.Players[i]

		// Update AI player info (don't show total)
		game.Layout.UpdatePlayerInfo(i, aiPlayer.Name, aiPlayer.Credits, 0, false)

		// Render AI player's cards (as CP$37 blocks)
		game.Layout.RenderPlayerCards(i, aiPlayer.Hand, true, game.CardRenderer)
	}

	// Update human player (index 0)
	if len(game.Players) > 0 {
		player := &game.Players[0]
		total := calculateHandTotal(player.Hand)
		// Add static field cards to total
		for _, card := range player.StaticField {
			total += card.Value
		}

		game.Layout.UpdatePlayerInfo(0, player.Name, player.Credits, total, true)

		// Render player's cards (face up)
		game.Layout.RenderPlayerCards(0, player.Hand, false, game.CardRenderer)

		// Render static field
		game.Layout.RenderStaticField(player.StaticField, game.CardRenderer)
	}
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
