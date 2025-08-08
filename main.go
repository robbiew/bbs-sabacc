package main

import (
	_ "embed"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/eiannone/keyboard"
)

var (
	DropPath string
	game     *SabaccGame
	config   GameConfig

	//go:embed ansi/title.ans
	TitleScreen string
	//go:embed ansi/menu.ans
	MenuScreen string
)

// SabaccGame represents the main game state
type SabaccGame struct {
	User          User
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
	// Load configuration
	config = LoadConfig()
	Idle = config.IdleTimeout // Set timeout from config
}

func main() {
	// Use FLAG to get command line parameters
	pathPtr := flag.String("path", "", "path to door32.sys file")
	flag.Parse()

	if *pathPtr == "" {
		fmt.Fprintf(os.Stderr, "missing path to door32.sys directory: -path\r\n")
		fmt.Fprintf(os.Stderr, "Usage: %s -path /path/to/dropfile/directory/\r\n", os.Args[0])
		os.Exit(2)
	}
	DropPath = *pathPtr

	// Get door32.sys, h, w as user object
	u := Initialize(DropPath)

	// Note: Removing strict ANSI check as some BBS systems may report emulation differently
	// The game will work with any terminal that supports basic ANSI escape sequences
	// Original check: if u.Emulation != 1 { ... }

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
		MinRounds:    config.MinRoundsToCall,    // From configuration
		Layout:       NewScreenLayout(u.W, u.H), // Initialize persistent UI
	}

	// DEBUG: Show portrait loading errors
	if game.Layout.PortraitManager != nil {
		errors := game.Layout.PortraitManager.GetDebugErrors()
		if len(errors) > 0 {
			fmt.Println("=== PORTRAIT DEBUG INFO ===")
			for _, err := range errors {
				fmt.Println(err)
			}
			fmt.Println("============================")
			fmt.Println("Press Enter to continue...")
			fmt.Scanln()
		}
	}

	// Show title screen
	showTitleScreen()

	// Main menu loop
	mainMenu()
}

func showTitleScreen() {
	ClearScreen()
	MoveCursor(0, 0)

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
			MoveCursor(1, centerY-6)
			displayAsciiArt("sabacc")

			MoveCursor(1, centerY+2)
			fmt.Print(Cyan + "Classic 76-Card Sabacc for BBS" + Reset)
		}
	}

	// Add welcome message and continue prompt at bottom
	MoveCursor(1, game.User.H-4)
	fmt.Print(White + "Welcome, " + CyanHi + game.User.Alias + Reset)

	MoveCursor(1, game.User.H-2)
	fmt.Print(Yellow + "Press any key to continue..." + Reset)

	waitForKey()
}

func mainMenu() {
	for {
		ClearScreen()
		MoveCursor(0, 0)

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
		MoveCursor(1, 4)
		fmt.Print(Yellow + "[" + YellowHi + "N" + Reset + Yellow + "] " + White + "New Game\r\n" + Reset)
		fmt.Print(Yellow + "[" + YellowHi + "R" + Reset + Yellow + "] " + White + "Rules\r\n" + Reset)
		fmt.Print(Yellow + "[" + YellowHi + "S" + Reset + Yellow + "] " + White + "Statistics\r\n" + Reset)
		fmt.Print(Yellow + "[" + YellowHi + "Q" + Reset + Yellow + "] " + White + "Quit to BBS\r\n\r\n" + Reset)

		fmt.Print(Cyan + "Credits: " + CyanHi + "1000" + Reset + "  ")
		fmt.Print(Cyan + "Time Left: " + CyanHi + strconv.Itoa(game.User.TimeLeft) + "m" + Reset + "\r\n\r\n")

		fmt.Print(Green + "Choice: " + Reset)

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
	fmt.Print(CyanHi + "-------------------------------------------\r\n" + Reset)
	fmt.Print(CyanHi + "              SABACC CANTINA\r\n" + Reset)
	fmt.Print(CyanHi + "-------------------------------------------\r\n\r\n" + Reset)
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
		{Name: game.User.Alias, Credits: config.StartingCredits, Hand: []Card{}, StaticField: []Card{}, Folded: false, BombedOut: false}, // Human player
		{Name: "PHOOJA", Credits: config.StartingCredits, Hand: []Card{}, StaticField: []Card{}, Folded: false, BombedOut: false},        // AI 1 (top-left)
		{Name: "ASH-TAAC", Credits: config.StartingCredits, Hand: []Card{}, StaticField: []Card{}, Folded: false, BombedOut: false},      // AI 2 (top-right)
		{Name: "OOLANGA", Credits: config.StartingCredits, Hand: []Card{}, StaticField: []Card{}, Folded: false, BombedOut: false},       // AI 3 (bottom-left)
		{Name: "KY'ALA", Credits: config.StartingCredits, Hand: []Card{}, StaticField: []Card{}, Folded: false, BombedOut: false},        // AI 4 (bottom-right)
	}

	// ANTE PHASE - Both players ante into both pots
	anteAmount := config.MinAnte
	game.HandPot = anteAmount * 2
	game.SabaccPot = anteAmount * 2

	// Deduct ante from players
	for i := range game.Players {
		game.Players[i].Credits -= anteAmount * 2 // Ante goes to both pots
	}

	// Show ante message in game log
	game.Layout.LogMessage(fmt.Sprintf("All players ante %d credits to each pot", anteAmount), "info")
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
				fmt.Printf("\r\n%s%s wins by default (others folded)!%s\r\n",
					GreenHi, game.Players[lastActivePlayer].Name, Reset)
				game.Players[lastActivePlayer].Credits += game.HandPot

				// Check if they also get Sabacc pot (if they have special hand)
				if !game.Players[lastActivePlayer].Folded {
					total := calculateHandTotal(game.Players[lastActivePlayer].Hand)
					if total == 23 {
						fmt.Printf("%s%s also wins the Sabacc Pot! (Pure Sabacc)%s\r\n",
							GreenHi, game.Players[lastActivePlayer].Name, Reset)
						game.Players[lastActivePlayer].Credits += game.SabaccPot
						game.SabaccPot = 0
						displayAsciiArt("sabacc")
					} else if isIdiotsArray(game.Players[lastActivePlayer].Hand) {
						fmt.Printf("%s%s also wins the Sabacc Pot! (Idiot's Array)%s\r\n",
							GreenHi, game.Players[lastActivePlayer].Name, Reset)
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

// Updated handlePlayerTurn function using Classic Sabacc 4-phase turn structure
func handlePlayerTurn() {
	// PHASE 1: BETTING PHASE
	if !handlePlayerBetting() {
		return // Player folded during betting
	}

	// PHASE 2: ROLL FOR SHIFT
	rollForShift()
	if game.ShiftOccurred {
		return // Turn ends after shift
	}

	// PHASE 3: CALL PHASE (only if round 2+)
	if game.Round >= game.MinRounds {
		if handlePlayerCall() {
			return // Hand was called
		}
	}

	// PHASE 4: DRAW PHASE
	handlePlayerDraw()
}

// PHASE 1: Betting Phase (Classic Sabacc)
func handlePlayerBetting() bool {
	game.Layout.DisplayMessage("BETTING PHASE", "info", 0)
	time.Sleep(1 * time.Second)

	// Use the same menu system as player turn menu
	bettingOptions := []MenuOption{
		{'C', "Check/Call current bet", true},
		{'R', "Raise bet", true},
		{'F', "Fold (1 credit penalty)", true},
	}

	game.Layout.ShowMenu("Betting Phase", bettingOptions, "Betting choice: ")

	char, _, err := getKeyWithTimeout()
	if err != nil {
		return true
	}

	playerRef := &game.Players[0]

	switch char {
	case 'c', 'C':
		// Check or call current bet
		if game.CurrentBet > 0 {
			betAmount := game.CurrentBet
			if playerRef.Credits >= betAmount {
				playerRef.Credits -= betAmount
				game.HandPot += betAmount
				game.Layout.LogMessage(fmt.Sprintf("You call %d credits", betAmount), "action")
				game.Layout.DisplayMessage(fmt.Sprintf("You call %d credits", betAmount), "success", 0)
			} else {
				game.Layout.DisplayMessage("Not enough credits to call!", "error", 0)
				return false
			}
		} else {
			game.Layout.LogMessage("You check (no bet)", "action")
			game.Layout.DisplayMessage("You check", "info", 0)
		}
		time.Sleep(1 * time.Second)

	case 'r', 'R':
		// Raise bet
		raiseAmount := 10 // Standard raise amount, could be configurable
		totalBet := game.CurrentBet + raiseAmount

		if playerRef.Credits >= totalBet {
			playerRef.Credits -= totalBet
			game.HandPot += totalBet
			game.CurrentBet = totalBet // Set new bet amount for other players
			game.Layout.LogMessage(fmt.Sprintf("You raise to %d credits", totalBet), "action")
			game.Layout.DisplayMessage(fmt.Sprintf("You raise to %d credits", totalBet), "success", 0)
		} else {
			game.Layout.DisplayMessage("Not enough credits to raise!", "error", 0)
			return false
		}
		time.Sleep(1 * time.Second)

	case 'f', 'F':
		// Fold
		playerRef.Folded = true
		playerRef.Credits -= 1 // Fold penalty goes to Sabacc Pot
		game.SabaccPot += 1
		game.Layout.LogMessage("You folded (-1 credit penalty)", "important")
		game.Layout.DisplayMessage("You folded", "warning", 0)
		time.Sleep(2 * time.Second)
		return false // Player folded, end turn

	default:
		game.Layout.DisplayMessage("Invalid choice!", "error", 0)
		time.Sleep(1 * time.Second)
		return true // Try again
	}

	return true // Continue to next phase
}

// PHASE 3: Call Phase (Classic Sabacc)
func handlePlayerCall() bool {
	game.Layout.DisplayMessage("CALL PHASE", "info", 0)
	time.Sleep(1 * time.Second)

	fmt.Print(Yellow + "[" + YellowHi + "C" + Yellow + "] " + White + "Call the hand (end game)\r\n" + Reset)
	fmt.Print(Yellow + "[" + YellowHi + "N" + Yellow + "] " + White + "No call (continue)\r\n\r\n" + Reset)
	fmt.Print(Green + "Call choice: " + Reset)

	char, _, err := getKeyWithTimeout()
	if err != nil {
		return false
	}

	switch char {
	case 'c', 'C':
		game.Called = true
		game.Layout.LogMessage("You called the hand!", "important")
		game.Layout.DisplayMessage("You called the hand!", "success", 0)
		time.Sleep(2 * time.Second)
		return true // Hand was called

	case 'n', 'N':
		game.Layout.LogMessage("You choose not to call", "action")
		game.Layout.DisplayMessage("No call", "info", 0)
		time.Sleep(1 * time.Second)
		return false // Continue game

	default:
		game.Layout.DisplayMessage("Invalid choice!", "error", 0)
		time.Sleep(1 * time.Second)
		return false
	}
}

// PHASE 4: Draw Phase (Classic Sabacc)
func handlePlayerDraw() {
	game.Layout.DisplayMessage("DRAW PHASE", "info", 0)
	time.Sleep(1 * time.Second)

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
	default:
		game.Layout.DisplayMessage("Invalid choice!", "error", 0)
		time.Sleep(1 * time.Second)
	}
}

func handleComputerTurn(canCall bool) {
	computer := &game.Players[game.Turn]

	if computer.Folded {
		return
	}

	game.Layout.LogMessage(computer.Name+" is thinking...", "info")
	time.Sleep(1 * time.Second)

	// PHASE 1: BETTING PHASE (Classic Sabacc AI)
	if !handleComputerBetting(game.Turn) {
		return // Computer folded during betting
	}

	// PHASE 2: ROLL FOR SHIFT
	rollForShift()
	if game.ShiftOccurred {
		return // Turn ends after shift
	}

	// PHASE 3: CALL PHASE (only if round 2+)
	if game.Round >= game.MinRounds {
		if handleComputerCall(game.Turn) {
			return // Hand was called
		}
	}

	// PHASE 4: DRAW PHASE
	handleComputerDraw(game.Turn)
}

// AI PHASE 1: Computer Betting Phase
func handleComputerBetting(playerIndex int) bool {
	computer := &game.Players[playerIndex]
	total := calculateHandTotal(computer.Hand) + calculateHandTotal(computer.StaticField)

	// AI betting logic based on hand strength
	if total > 23 || total < -23 || total == 0 {
		// Bombed out - must fold
		computer.Folded = true
		computer.Credits -= 1
		game.SabaccPot += 1
		game.Layout.LogMessage(computer.Name+" folds (bombed out)", "important")
		return false
	}

	// Strong hands (20-23): Raise or call aggressively
	if total >= 20 && total <= 23 {
		if game.CurrentBet == 0 {
			// Raise with strong hand
			raiseAmount := 10
			if computer.Credits >= raiseAmount {
				computer.Credits -= raiseAmount
				game.HandPot += raiseAmount
				game.CurrentBet = raiseAmount
				game.Layout.LogMessage(computer.Name+" raises to "+fmt.Sprintf("%d", raiseAmount)+" credits", "action")
			} else {
				game.Layout.LogMessage(computer.Name+" checks", "action")
			}
		} else {
			// Call existing bet
			if computer.Credits >= game.CurrentBet {
				computer.Credits -= game.CurrentBet
				game.HandPot += game.CurrentBet
				game.Layout.LogMessage(computer.Name+" calls "+fmt.Sprintf("%d", game.CurrentBet)+" credits", "action")
			} else {
				// Can't afford to call - fold
				computer.Folded = true
				computer.Credits -= 1
				game.SabaccPot += 1
				game.Layout.LogMessage(computer.Name+" folds (insufficient credits)", "important")
				return false
			}
		}
	} else if total >= 15 && total <= 19 {
		// Medium hands: Call small bets, fold large ones
		if game.CurrentBet <= 5 {
			if computer.Credits >= game.CurrentBet {
				computer.Credits -= game.CurrentBet
				game.HandPot += game.CurrentBet
				game.Layout.LogMessage(computer.Name+" calls "+fmt.Sprintf("%d", game.CurrentBet)+" credits", "action")
			} else {
				game.Layout.LogMessage(computer.Name+" checks", "action")
			}
		} else {
			// Bet too high for medium hand - fold
			computer.Folded = true
			computer.Credits -= 1
			game.SabaccPot += 1
			game.Layout.LogMessage(computer.Name+" folds (bet too high)", "important")
			return false
		}
	} else {
		// Weak hands: Check/fold
		if game.CurrentBet == 0 {
			game.Layout.LogMessage(computer.Name+" checks", "action")
		} else {
			// Any bet with weak hand - fold
			computer.Folded = true
			computer.Credits -= 1
			game.SabaccPot += 1
			game.Layout.LogMessage(computer.Name+" folds (weak hand)", "important")
			return false
		}
	}

	time.Sleep(1 * time.Second)
	return true
}

// AI PHASE 3: Computer Call Phase
func handleComputerCall(playerIndex int) bool {
	computer := &game.Players[playerIndex]
	total := calculateHandTotal(computer.Hand) + calculateHandTotal(computer.StaticField)

	shouldCall := evaluateAICallDecision(playerIndex, total)

	if shouldCall {
		game.Called = true
		game.Layout.LogMessage(computer.Name+" calls the hand!", "important")
		time.Sleep(2 * time.Second)
		return true
	}

	game.Layout.LogMessage(computer.Name+" chooses not to call", "action")
	time.Sleep(500 * time.Millisecond)
	return false
}

// AI PHASE 4: Computer Draw Phase
func handleComputerDraw(playerIndex int) {
	computer := &game.Players[playerIndex]
	total := calculateHandTotal(computer.Hand) + calculateHandTotal(computer.StaticField)

	// First: Static field management
	handleAIStaticField(playerIndex)

	// Then: Draw decision
	if total >= 20 && total <= 23 {
		// Excellent hand - stand
		game.Layout.LogMessage(computer.Name+" stands", "action")
	} else if total >= 15 && total <= 19 {
		// Good hand - might draw if calculated risk is worth it
		drawDecision := evaluateAIDrawDecision(playerIndex, total)
		if drawDecision && len(game.Deck.Cards) > 0 {
			card := game.Deck.Deal()
			computer.Hand = append(computer.Hand, card)
			game.Layout.LogMessage(computer.Name+" draws a card", "action")
		} else {
			game.Layout.LogMessage(computer.Name+" stands", "action")
		}
	} else if total >= 5 && total <= 14 {
		// Medium hand - usually draw
		if len(game.Deck.Cards) > 0 {
			card := game.Deck.Deal()
			computer.Hand = append(computer.Hand, card)
			game.Layout.LogMessage(computer.Name+" draws a card", "action")
		} else {
			game.Layout.LogMessage(computer.Name+" stands (no cards left)", "action")
		}
	} else if total < 5 && total > -15 {
		// Poor hand - aggressive improvement needed
		if len(computer.Hand) > 2 && len(game.Deck.Cards) > 0 {
			// Trade worst card
			computer.Hand = computer.Hand[1:] // Remove card (simplified)
			card := game.Deck.Deal()
			computer.Hand = append(computer.Hand, card)
			game.Layout.LogMessage(computer.Name+" trades a card", "action")
		} else if len(game.Deck.Cards) > 0 {
			card := game.Deck.Deal()
			computer.Hand = append(computer.Hand, card)
			game.Layout.LogMessage(computer.Name+" draws a card", "action")
		} else {
			game.Layout.LogMessage(computer.Name+" stands", "action")
		}
	} else {
		// Very poor hand or close to bombing - stand
		game.Layout.LogMessage(computer.Name+" stands", "action")
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
	ClearScreen()
	fmt.Print(CyanHi + strings.Repeat(DOUBLE_HORIZONTAL_LINE, 39) + "\r\n" + Reset)
	fmt.Print(CyanHi + "                HAND RESULTS\r\n" + Reset)
	fmt.Print(CyanHi + strings.Repeat(DOUBLE_HORIZONTAL_LINE, 39) + "\r\n\r\n" + Reset)

	// Show all hands and determine winner
	winner := -1
	bestScore := -999
	sabaccWinner := -1
	bombedOutPlayers := []int{}

	for i, playerData := range game.Players {
		if playerData.Folded {
			fmt.Printf("%s%s: FOLDED%s\r\n", Red, playerData.Name, Reset)
			continue
		}

		total := calculateHandTotal(playerData.Hand)
		fmt.Printf("%s%s:%s ", Cyan, playerData.Name, Reset)
		for _, card := range playerData.Hand {
			fmt.Printf("%s[%s]%s ", getCardColor(card), card.String(), Reset)
		}
		fmt.Printf("= %s%d%s", YellowHi, total, Reset)

		// Check for special hands (Sabacc Pot winners)
		if isIdiotsArray(playerData.Hand) {
			fmt.Printf(" %s(IDIOT'S ARRAY!)%s\r\n", GreenHi, Reset)
			displayAsciiArt("sabacc") // Show sabacc art for special hands
			bestScore = 1000
			sabaccWinner = i            // Idiot's Array beats Pure Sabacc
			time.Sleep(3 * time.Second) // Let player see the art
		} else if total == 23 && sabaccWinner == -1 {
			fmt.Printf(" %s(PURE SABACC!)%s\r\n", GreenHi, Reset)
			displayAsciiArt("sabacc") // Show sabacc art for Pure Sabacc
			sabaccWinner = i
			time.Sleep(3 * time.Second) // Let player see the art
		} else if total > 23 || total < -23 || total == 0 {
			fmt.Printf(" %s(BOMBED OUT!)%s\r\n", Red, Reset)
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
		fmt.Printf("%s%s wins both pots with a special hand!%s\r\n",
			GreenHi, game.Players[sabaccWinner].Name, Reset)
		fmt.Printf("%s+%d credits (Hand Pot) +%d credits (Sabacc Pot)%s\r\n",
			GreenHi, game.HandPot, game.SabaccPot, Reset)

		game.Players[sabaccWinner].Credits += game.HandPot + game.SabaccPot
		game.SabaccPot = 0 // Reset Sabacc pot

		// Display celebration art one more time
		displayAsciiArt("sabacc")
		time.Sleep(2 * time.Second)

	} else if winner >= 0 {
		// Regular hand winner
		fmt.Printf("%s%s wins the hand! (+%d credits)%s\r\n",
			GreenHi, game.Players[winner].Name, game.HandPot, Reset)
		game.Players[winner].Credits += game.HandPot

	} else {
		// Everyone bombed out or folded
		fmt.Printf("%sNo winner! Hand pot goes to Sabacc pot.%s\r\n", Yellow, Reset)
		game.SabaccPot += game.HandPot

		// Show bomb art for the chaos
		if len(bombedOutPlayers) > 1 {
			fmt.Printf("\r\n%sEveryone bombed out!%s\r\n", RedHi, Reset)
			displayAsciiArt("bomb")
			time.Sleep(2 * time.Second)
		}
	}

	// Show penalty summary if anyone bombed out
	if len(bombedOutPlayers) > 0 {
		fmt.Printf("\r\n%sBomb Out Penalties:%s\r\n", Red, Reset)
		for _, playerIndex := range bombedOutPlayers {
			fmt.Printf("%s%s paid %d credits to Sabacc Pot%s\r\n",
				Red, game.Players[playerIndex].Name, game.HandPot, Reset)
		}
	}

	game.GameOver = true
}

// Fixed display functions for Sabacc game

func displayGameScreen() {
	// Update header with current game info
	currentPlayerName := ""
	if game.Turn == 0 {
		currentPlayerName = "\x10 YOUR TURN \x11" // CP437 arrows
	} else if game.Turn < len(game.Players) {
		currentPlayerName = game.Players[game.Turn].Name + " thinking..."
	}

	game.Layout.UpdateHeader(game.Round, game.HandPot, game.SabaccPot, len(game.Deck.Cards), currentPlayerName)

	// Update all AI players (indexes 1-4)
	for i := 1; i < len(game.Players) && i <= 4; i++ {
		aiPlayer := &game.Players[i]

		// Update AI player info (don't show total)
		game.Layout.UpdatePlayerInfo(i, aiPlayer.Name, aiPlayer.Credits, 0, false)

		// Render AI player's cards (as CP437 blocks)
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

// Enhanced AI functions for 1989 West End Games rules compliance

// evaluateAICallDecision determines if AI should call based on visible cards and hand strength
func evaluateAICallDecision(playerIndex int, combinedTotal int) bool {
	computer := &game.Players[playerIndex]

	// Don't call if bombed out
	if combinedTotal > 23 || combinedTotal < -23 || combinedTotal == 0 {
		return false
	}

	// Strong call for Pure Sabacc or near-perfect hands
	if combinedTotal == 23 {
		return true // Pure Sabacc - definitely call
	}

	// Check for Idiot's Array potential
	if isIdiotsArray(computer.Hand) {
		return true // Idiot's Array beats Pure Sabacc
	}

	// Analyze visible cards from all players' static fields
	visibleCards := getVisibleCardsFromStaticFields()

	// Strong hands worth calling (20-22)
	if combinedTotal >= 20 && combinedTotal <= 22 {
		// More aggressive calling if opponents have visible weak cards
		opponentStrength := assessOpponentStrengthFromVisibleCards(playerIndex, visibleCards)
		if opponentStrength == "weak" {
			return combinedTotal >= 20
		} else if opponentStrength == "strong" {
			return combinedTotal >= 22 // Only call with very strong hands
		}
		return combinedTotal >= 21 // Moderate calling threshold
	}

	// Conservative approach for medium hands (15-19)
	if combinedTotal >= 15 && combinedTotal <= 19 {
		// Only call if we can see opponents have weak visible cards
		opponentStrength := assessOpponentStrengthFromVisibleCards(playerIndex, visibleCards)
		return opponentStrength == "weak" && combinedTotal >= 18
	}

	return false // Don't call with weak hands
}

// handleAIStaticField manages AI static field placement and removal
func handleAIStaticField(playerIndex int) {
	computer := &game.Players[playerIndex]

	if len(computer.Hand) == 0 {
		return
	}

	// Strategy 1: Protect valuable cards before shifts
	valuableCards := findValuableCardsToProtect(computer.Hand, computer.StaticField)

	for _, cardIndex := range valuableCards {
		if cardIndex < len(computer.Hand) {
			// Move valuable card to static field
			card := computer.Hand[cardIndex]
			computer.Hand = append(computer.Hand[:cardIndex], computer.Hand[cardIndex+1:]...)
			computer.StaticField = append(computer.StaticField, card)

			game.Layout.LogMessage(fmt.Sprintf("%s places a card in Static Field", computer.Name), "action")
			break // Only move one card per turn
		}
	}

	// Strategy 2: Remove cards from static field if hand composition changed
	if len(computer.StaticField) > 0 {
		shouldRemove := evaluateStaticFieldRemoval(computer.Hand, computer.StaticField)
		if shouldRemove >= 0 {
			// Move card back to hand
			card := computer.StaticField[shouldRemove]
			computer.StaticField = append(computer.StaticField[:shouldRemove], computer.StaticField[shouldRemove+1:]...)
			computer.Hand = append(computer.Hand, card)

			game.Layout.LogMessage(fmt.Sprintf("%s removes a card from Static Field", computer.Name), "action")
		}
	}
}

// evaluateAIDrawDecision determines if AI should draw based on strategy and visible cards
func evaluateAIDrawDecision(playerIndex int, combinedTotal int) bool {
	computer := &game.Players[playerIndex]

	// Never draw if already at optimal scores
	if combinedTotal == 23 || isIdiotsArray(computer.Hand) {
		return false
	}

	// Never draw if close to bombing out
	if combinedTotal > 20 || combinedTotal < -20 {
		return false
	}

	// Analyze visible cards to make informed decisions
	visibleCards := getVisibleCardsFromStaticFields()
	availableCardTypes := analyzeAvailableCards(visibleCards)

	// Draw if we need specific cards and they're likely available
	neededValue := 23 - combinedTotal

	// Conservative drawing for good hands (15-20)
	if combinedTotal >= 15 && combinedTotal <= 20 {
		// Only draw if we can see beneficial cards haven't been taken
		if neededValue > 0 && neededValue <= 8 {
			return availableCardTypes[neededValue] > 0 // Draw if needed cards are available
		}
		return false
	}

	// More aggressive drawing for medium hands (5-14)
	if combinedTotal >= 5 && combinedTotal <= 14 {
		// Draw if it could help reach good range
		return neededValue >= 3 && neededValue <= 18
	}

	// Desperate drawing for poor hands (below 5)
	if combinedTotal < 5 {
		return true // Need to improve urgently
	}

	return false
}

// Helper functions for AI intelligence

// getVisibleCardsFromStaticFields returns all cards visible in static fields
func getVisibleCardsFromStaticFields() []Card {
	var visibleCards []Card

	for _, player := range game.Players {
		if !player.Folded {
			// All static field cards are visible to opponents
			visibleCards = append(visibleCards, player.StaticField...)
		}
	}

	return visibleCards
}

// assessOpponentStrengthFromVisibleCards analyzes opponent strength based on visible cards
func assessOpponentStrengthFromVisibleCards(playerIndex int, visibleCards []Card) string {
	strongCardCount := 0
	weakCardCount := 0

	for i, player := range game.Players {
		if i != playerIndex && !player.Folded && len(player.StaticField) > 0 {
			// Analyze visible cards in opponent's static field
			for _, card := range player.StaticField {
				if card.Value >= 10 || card.Value <= -10 {
					strongCardCount++
				} else if card.Value >= -5 && card.Value <= 5 {
					weakCardCount++
				}
			}
		}
	}

	if strongCardCount > weakCardCount {
		return "strong"
	} else if weakCardCount > strongCardCount {
		return "weak"
	}
	return "moderate"
}

// findValuableCardsToProtect identifies cards worth protecting in static field
func findValuableCardsToProtect(hand []Card, staticField []Card) []int {
	var valuableIndices []int

	// Don't protect if already have too many cards in static field
	if len(staticField) >= 2 {
		return valuableIndices
	}

	for i, card := range hand {
		// Protect high-value positive cards
		if card.Value >= 10 && card.Value <= 15 {
			valuableIndices = append(valuableIndices, i)
		}
		// Protect Idiot card (very valuable)
		if card.Name == "Idiot" {
			valuableIndices = append(valuableIndices, i)
		}
		// Protect cards that complete good combinations
		if isPartOfGoodCombination(card, hand) {
			valuableIndices = append(valuableIndices, i)
		}
	}

	return valuableIndices
}

// evaluateStaticFieldRemoval determines if cards should be removed from static field
func evaluateStaticFieldRemoval(hand []Card, staticField []Card) int {
	// Remove cards if they're now less valuable
	for i, staticCard := range staticField {
		// Remove low-value cards if hand composition changed
		if staticCard.Value > 0 && staticCard.Value < 5 {
			// Check if we now have better cards in hand
			handTotal := calculateHandTotal(hand)
			if handTotal+staticCard.Value > 23 {
				return i // Remove this card to avoid bombing out
			}
		}

		// Remove cards that no longer serve strategic purpose
		if !isPartOfGoodCombination(staticCard, append(hand, staticField...)) {
			if len(hand) <= 2 { // Only if we have room in hand
				return i
			}
		}
	}

	return -1 // Don't remove any cards
}

// analyzeAvailableCards estimates what cards might still be available
func analyzeAvailableCards(visibleCards []Card) map[int]int {
	cardCounts := make(map[int]int)

	// Initialize with 1989 Classic Sabacc deck composition knowledge
	for value := 1; value <= 15; value++ {
		cardCounts[value] = 4 // 4 suits (60 cards total)
	}

	// Arcana cards: one copy each (16 cards total) - 1989 Classic Rules
	arcanaValues := []int{-1, -2, -3, -4, -5, -6, -7, -8, -9, -10, -11, -12, -13, -14, -15, -17}
	for _, value := range arcanaValues {
		cardCounts[value] = 1 // Only 1 copy of each Arcana card
	}

	// Subtract visible cards from static fields
	for _, card := range visibleCards {
		if cardCounts[card.Value] > 0 {
			cardCounts[card.Value]--
		}
	}

	return cardCounts
}

// isPartOfGoodCombination checks if card contributes to valuable hand combinations
func isPartOfGoodCombination(card Card, allCards []Card) bool {
	// Check for Idiot's Array components
	if card.Name == "Idiot" {
		return true
	}
	if card.Value == 2 || card.Value == 3 {
		// Check if we have other Idiot's Array components
		hasIdiot := false
		hasTwoOrThree := false
		for _, c := range allCards {
			if c.Name == "Idiot" {
				hasIdiot = true
			}
			if (c.Value == 2 || c.Value == 3) && !(c.Value == card.Value && c.Suit == card.Suit) {
				hasTwoOrThree = true
			}
		}
		return hasIdiot || hasTwoOrThree
	}

	// High-value cards that help reach Pure Sabacc
	total := calculateHandTotal(allCards)
	if total+card.Value == 23 {
		return true
	}

	return false
}

func showGameResults() {
	fmt.Println()
	fmt.Printf("%sFinal Credits:%s\r\n", CyanHi, Reset)
	for _, playerData := range game.Players {
		fmt.Printf("%s: %s%d%s credits\r\n", playerData.Name, YellowHi, playerData.Credits, Reset)
	}
	fmt.Println()
	fmt.Print(Yellow + "Press any key to return to menu..." + Reset)
}
