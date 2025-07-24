package main

import (
	"fmt"
	"strings"

	gd "github.com/robbiew/godoors"
)

// ScreenLayout manages the persistent UI layout for the Sabacc game
type ScreenLayout struct {
	// Terminal dimensions
	TerminalW int
	TerminalH int

	// Fixed screen positions (for 80x25 terminal)
	HeaderY        int // Line 1 - Game info header
	Player2InfoY   int // Line 2 - Player 2 name/credits
	Player2CardY   int // Line 3 - Player 2 cards start
	Player3InfoY   int // Line 6 - Player 3 area (future expansion)
	Player3CardY   int // Line 7 - Player 3 cards start
	GameLogY       int // Line 10 - Scrolling messages start
	GameLogH       int // Height of game log area
	Player1InfoY   int // Line 14 - Player 1 name/total/credits
	Player1CardY   int // Line 15 - Player 1 cards start
	Player1StaticY int // Line 19 - Static field area
	MenuY          int // Line 21 - Command menu area

	// Card display settings
	CardWidth   int // Width of each card display (7 chars)
	CardHeight  int // Height of each card display (3 lines)
	MaxCardsRow int // Maximum cards per row (10 for 80-char width)

	// Game log system
	GameLog *GameLog
}

// GameLog manages the scrolling message area
type GameLog struct {
	Messages []string // Circular buffer of messages
	MaxLines int      // Maximum lines to display
	StartY   int      // Starting Y position for display
	Width    int      // Width of log area
}

// MenuOption represents a menu choice
type MenuOption struct {
	Key         rune
	Description string
	Enabled     bool
}

// NewScreenLayout creates a new screen layout for the given terminal size
func NewScreenLayout(termW, termH int) *ScreenLayout {
	layout := &ScreenLayout{
		TerminalW: termW,
		TerminalH: termH,

		// Standard 80x25 layout positions
		HeaderY:        1,
		Player2InfoY:   2,
		Player2CardY:   3,
		Player3InfoY:   6, // Reserved for future 3-4 player support
		Player3CardY:   7,
		GameLogY:       10,
		GameLogH:       3, // Compact 3-line log
		Player1InfoY:   14,
		Player1CardY:   15,
		Player1StaticY: 19,
		MenuY:          21,

		// Card display settings
		CardWidth:   7,  // Includes borders: " ┌───┐ "
		CardHeight:  3,  // 3 lines per card
		MaxCardsRow: 10, // Max cards that fit in 80 chars
	}

	// Initialize game log
	layout.GameLog = &GameLog{
		Messages: make([]string, 0),
		MaxLines: layout.GameLogH,
		StartY:   layout.GameLogY,
		Width:    termW - 4, // Leave margins
	}

	return layout
}

// InitializeScreen draws the static UI elements once
func (sl *ScreenLayout) InitializeScreen() {
	gd.ClearScreen()

	// Draw main borders and labels
	sl.drawStaticElements()

	// Initialize all areas with empty content
	sl.ClearPlayerArea(2) // Player 2
	sl.ClearPlayerArea(1) // Player 1 (human)
	sl.GameLog.Clear()
	sl.ClearMenuArea()
}

// drawStaticElements draws borders and static labels
func (sl *ScreenLayout) drawStaticElements() {
	// Draw horizontal separators
	separatorLine := strings.Repeat("─", sl.TerminalW-2)

	// Header separator (line 2)
	gd.MoveCursor(1, sl.Player2InfoY-1)
	fmt.Printf("┌%s┐", separatorLine)

	// Player areas separators
	gd.MoveCursor(1, sl.Player3InfoY-1)
	fmt.Printf("├%s┤", separatorLine)

	gd.MoveCursor(1, sl.GameLogY-1)
	fmt.Printf("├%s┤", separatorLine)

	gd.MoveCursor(1, sl.Player1InfoY-1)
	fmt.Printf("├%s┤", separatorLine)

	gd.MoveCursor(1, sl.MenuY-1)
	fmt.Printf("├%s┤", separatorLine)

	// Bottom border
	gd.MoveCursor(1, sl.TerminalH)
	fmt.Printf("└%s┘", separatorLine)

	// Game log border
	sl.drawGameLogBorder()
}

// drawGameLogBorder draws the border around the game log area
func (sl *ScreenLayout) drawGameLogBorder() {
	logWidth := sl.TerminalW - 8 // Leave margins
	borderLine := strings.Repeat("═", logWidth-2)

	// Top border
	gd.MoveCursor(3, sl.GameLogY)
	fmt.Printf("╔%s╗", borderLine)

	// Side borders for each line
	for i := 1; i <= sl.GameLogH; i++ {
		gd.MoveCursor(3, sl.GameLogY+i)
		fmt.Print("║")
		gd.MoveCursor(3+logWidth-1, sl.GameLogY+i)
		fmt.Print("║")
	}

	// Bottom border
	gd.MoveCursor(3, sl.GameLogY+sl.GameLogH+1)
	fmt.Printf("╚%s╝", borderLine)
}

// UpdateHeader updates the game information header
func (sl *ScreenLayout) UpdateHeader(round, handPot, sabaccPot, deckSize int, currentPlayer string) {
	gd.MoveCursor(1, sl.HeaderY)

	// Clear the line first
	fmt.Print(gd.EraseLine)

	// Format header with game info
	headerText := fmt.Sprintf("│ SABACC - Round: %d │ Hand: %d │ Sabacc: %d │ Deck: %d │ Turn: %s",
		round, handPot, sabaccPot, deckSize, currentPlayer)

	// Pad to terminal width
	for len(headerText) < sl.TerminalW-1 {
		headerText += " "
	}
	headerText += "│"

	fmt.Print(gd.CyanHi + headerText + gd.Reset)
}

// UpdatePlayerInfo updates a player's name, credits, and total
func (sl *ScreenLayout) UpdatePlayerInfo(playerIndex int, name string, credits, total int, showTotal bool) {
	var infoY int

	switch playerIndex {
	case 0: // Human player (bottom)
		infoY = sl.Player1InfoY
	case 1: // Player 2 (top)
		infoY = sl.Player2InfoY
	case 2: // Player 3 (future)
		infoY = sl.Player3InfoY
	default:
		return
	}

	gd.MoveCursor(1, infoY)
	fmt.Print(gd.EraseLine)

	// Format player info line
	var infoText string
	if showTotal {
		infoText = fmt.Sprintf("│ %s%s%s%s Total: %s Credits: %s%d%s",
			gd.CyanHi, name, gd.Reset,
			strings.Repeat(" ", 40-len(name)),
			displayHandValue(total),
			gd.YellowHi, credits, gd.Reset)
	} else {
		infoText = fmt.Sprintf("│ %s%s%s%s Credits: %s%d%s",
			gd.CyanHi, name, gd.Reset,
			strings.Repeat(" ", 55-len(name)),
			gd.YellowHi, credits, gd.Reset)
	}

	// Pad to terminal width
	for len(infoText) < sl.TerminalW-1 {
		infoText += " "
	}
	infoText += "│"

	fmt.Print(infoText)
}

// ClearPlayerArea clears a player's card display area
func (sl *ScreenLayout) ClearPlayerArea(playerIndex int) {
	var startY int

	switch playerIndex {
	case 0: // Human player
		startY = sl.Player1CardY
	case 1: // Player 2
		startY = sl.Player2CardY
	case 2: // Player 3 (future)
		startY = sl.Player3CardY
	default:
		return
	}

	// Clear card display lines
	for i := 0; i < sl.CardHeight+1; i++ { // +1 for spacing
		gd.MoveCursor(1, startY+i)
		fmt.Print(gd.EraseLine)
		fmt.Print("│" + strings.Repeat(" ", sl.TerminalW-2) + "│")
	}
}

// ClearMenuArea clears the menu/command area
func (sl *ScreenLayout) ClearMenuArea() {
	for i := 0; i < 3; i++ { // Menu area is 3 lines
		gd.MoveCursor(1, sl.MenuY+i)
		fmt.Print(gd.EraseLine)
		fmt.Print("│" + strings.Repeat(" ", sl.TerminalW-2) + "│")
	}
}

// ShowMenu displays the command menu
func (sl *ScreenLayout) ShowMenu(title string, options []MenuOption, prompt string) {
	sl.ClearMenuArea()

	// Show title
	gd.MoveCursor(2, sl.MenuY)
	fmt.Printf("%s%s%s", gd.GreenHi, title, gd.Reset)

	// Show options
	gd.MoveCursor(2, sl.MenuY+1)
	optionTexts := make([]string, 0)
	for _, opt := range options {
		if opt.Enabled {
			optionText := fmt.Sprintf("%s[%s%c%s]%s%s",
				gd.Yellow, gd.YellowHi, opt.Key, gd.Yellow, gd.White, opt.Description)
			optionTexts = append(optionTexts, optionText)
		}
	}
	fmt.Print(strings.Join(optionTexts, " "))

	// Show prompt
	gd.MoveCursor(2, sl.MenuY+2)
	fmt.Printf("%s%s%s ", gd.Green, prompt, gd.Reset)
}

// NewGameLog creates a new game log
func NewGameLog(maxLines, startY, width int) *GameLog {
	return &GameLog{
		Messages: make([]string, 0),
		MaxLines: maxLines,
		StartY:   startY,
		Width:    width,
	}
}

// AddMessage adds a new message to the game log
func (gl *GameLog) AddMessage(message string) {
	// Truncate message if too long
	if len(message) > gl.Width-6 { // Account for borders and margins
		message = message[:gl.Width-9] + "..."
	}

	// Add to circular buffer
	gl.Messages = append(gl.Messages, message)

	// Remove oldest if we exceed max lines
	if len(gl.Messages) > gl.MaxLines {
		gl.Messages = gl.Messages[1:]
	}

	// Render updated log
	gl.Render()
}

// Render displays the current game log messages
func (gl *GameLog) Render() {
	// Clear log area content (not borders)
	for i := 0; i < gl.MaxLines; i++ {
		gd.MoveCursor(5, gl.StartY+1+i)            // Position after left border
		fmt.Print(strings.Repeat(" ", gl.Width-6)) // Clear content area
	}

	// Display messages (most recent at bottom)
	messageCount := len(gl.Messages)
	for i, message := range gl.Messages {
		lineY := gl.StartY + 1 + (gl.MaxLines - messageCount + i)
		if lineY >= gl.StartY+1 && lineY < gl.StartY+1+gl.MaxLines {
			gd.MoveCursor(5, lineY)
			fmt.Print(message)
		}
	}
}

// Clear clears all messages from the game log
func (gl *GameLog) Clear() {
	gl.Messages = make([]string, 0)
	gl.Render()
}

// RenderPlayerCards renders cards for a specific player in their designated area
func (sl *ScreenLayout) RenderPlayerCards(playerIndex int, cards []Card, faceDown bool, cardRenderer *CardRenderer) {
	var startY int

	switch playerIndex {
	case 0: // Human player (bottom)
		startY = sl.Player1CardY
	case 1: // Player 2 (top)
		startY = sl.Player2CardY
	case 2: // Player 3 (future)
		startY = sl.Player3CardY
	default:
		return
	}

	// Clear the card area first
	sl.ClearPlayerArea(playerIndex)

	if len(cards) == 0 {
		return
	}

	// Render cards using the existing CardRenderer if available
	if cardRenderer != nil && cardRenderer.Database != nil {
		if faceDown {
			// Render face-down cards
			for i, _ := range cards {
				cardX := 2 + (i * (sl.CardWidth + 1)) // +1 for spacing
				if cardX+sl.CardWidth > sl.TerminalW-2 {
					break // Don't exceed terminal width
				}
				cardRenderer.RenderFaceDownCard(cardX, startY)
			}
		} else {
			// Render face-up cards using CardRenderer
			cardRenderer.RenderCards(cards, 2, startY)
		}
	} else {
		// ASCII fallback - render simple card representations
		sl.renderCardsASCII(cards, startY, faceDown)
	}
}

// renderCardsASCII provides ASCII fallback for card rendering
func (sl *ScreenLayout) renderCardsASCII(cards []Card, startY int, faceDown bool) {
	gd.MoveCursor(2, startY)
	fmt.Print("│ ")

	for i, card := range cards {
		if i > 0 {
			fmt.Print(" ")
		}

		if faceDown {
			fmt.Printf("%s[??]%s", gd.Red, gd.Reset)
		} else {
			fmt.Printf("%s[%s]%s", getCardColor(card), card.String(), gd.Reset)
		}

		// Don't exceed terminal width
		if (i+1)*5+2 > sl.TerminalW-10 {
			break
		}
	}
}

// RenderStaticField renders the static field cards for the human player
func (sl *ScreenLayout) RenderStaticField(cards []Card, cardRenderer *CardRenderer) {
	if len(cards) == 0 {
		// Clear static field area if no cards
		gd.MoveCursor(2, sl.Player1StaticY)
		fmt.Print("│" + strings.Repeat(" ", sl.TerminalW-3) + "│")
		gd.MoveCursor(2, sl.Player1StaticY+1)
		fmt.Print("│" + strings.Repeat(" ", sl.TerminalW-3) + "│")
		return
	}

	// Clear static field area
	gd.MoveCursor(2, sl.Player1StaticY)
	fmt.Print("│" + strings.Repeat(" ", sl.TerminalW-3) + "│")
	gd.MoveCursor(2, sl.Player1StaticY+1)
	fmt.Print("│" + strings.Repeat(" ", sl.TerminalW-3) + "│")

	// Show static field label and cards
	gd.MoveCursor(2, sl.Player1StaticY)
	fmt.Printf("│ %sStatic Field (Protected):%s", gd.Magenta, gd.Reset)

	if cardRenderer != nil && cardRenderer.Database != nil {
		// Use CardRenderer for better graphics
		cardRenderer.RenderCards(cards, 2, sl.Player1StaticY+1)
	} else {
		// ASCII fallback
		gd.MoveCursor(2, sl.Player1StaticY+1)
		fmt.Print("│ ")
		for i, card := range cards {
			if i > 0 {
				fmt.Print(" ")
			}
			fmt.Printf("%s[%s]%s", getCardColor(card), card.String(), gd.Reset)

			// Don't exceed terminal width
			if (i+1)*5+2 > sl.TerminalW-10 {
				break
			}
		}
	}
}

// UpdateTurnIndicator updates who's turn it is in the header
func (sl *ScreenLayout) UpdateTurnIndicator(playerName string, isHumanTurn bool) {
	// This function is a placeholder for future use
	// Turn indication is currently handled by UpdateHeader
	// Could be expanded to show turn indicator in a separate area if needed
}

// ShowPlayerTurnMenu shows the menu options for the human player's turn
func (sl *ScreenLayout) ShowPlayerTurnMenu(round int) {
	options := []MenuOption{
		{'D', "raw card", true},
		{'T', "rade card", true},
		{'S', "tand", true},
		{'F', "ield (static)", true},
		{'C', "all hand", round >= 2}, // Only allow calling after round 2
		{'Q', "uit/Fold", true},
	}

	sl.ShowMenu("► YOUR TURN ◄", options, "Choice:")
}

// ShowBettingMenu shows betting options (for future expansion)
func (sl *ScreenLayout) ShowBettingMenu() {
	options := []MenuOption{
		{'C', "heck/Call", true},
		{'R', "aise", true},
		{'F', "old", true},
	}

	sl.ShowMenu("BETTING PHASE", options, "Choice:")
}

// ShowTradeMenu shows available cards for trading
func (sl *ScreenLayout) ShowTradeMenu(cards []Card) {
	gd.MoveCursor(2, sl.MenuY)
	fmt.Printf("%sTRADE A CARD - Select card to discard:%s", gd.GreenHi, gd.Reset)

	gd.MoveCursor(2, sl.MenuY+1)
	for i, card := range cards {
		fmt.Printf("%s[%d]%s%s[%s]%s ",
			gd.Green, i+1, gd.Reset,
			getCardColor(card), card.String(), gd.Reset)
	}
	fmt.Printf("%s[0]%sCancel", gd.Red, gd.Reset)

	gd.MoveCursor(2, sl.MenuY+2)
	fmt.Printf("%sChoice:%s ", gd.Green, gd.Reset)
}

// ShowStaticFieldMenu shows static field management options
func (sl *ScreenLayout) ShowStaticFieldMenu() {
	options := []MenuOption{
		{'1', " Place card in Static Field", true},
		{'2', " Remove card from Static Field", true},
		{'0', " Cancel", true},
	}

	sl.ShowMenu("STATIC FIELD MANAGEMENT", options, "Choice:")
}

// DisplayMessage shows a temporary message in the menu area
func (sl *ScreenLayout) DisplayMessage(message string, messageType string, pauseSeconds int) {
	sl.ClearMenuArea()

	var coloredMessage string
	switch messageType {
	case "success":
		coloredMessage = gd.GreenHi + message + gd.Reset
	case "error":
		coloredMessage = gd.RedHi + message + gd.Reset
	case "warning":
		coloredMessage = gd.YellowHi + message + gd.Reset
	case "info":
		coloredMessage = gd.CyanHi + message + gd.Reset
	default:
		coloredMessage = gd.White + message + gd.Reset
	}

	gd.MoveCursor(2, sl.MenuY+1)
	fmt.Print(coloredMessage)

	if pauseSeconds > 0 {
		gd.MoveCursor(2, sl.MenuY+2)
		fmt.Printf("%sPress any key to continue...%s", gd.Yellow, gd.Reset)
	}
}

// ShowGameResults displays the final game results
func (sl *ScreenLayout) ShowGameResults(players []Player) {
	sl.ClearMenuArea()

	gd.MoveCursor(2, sl.MenuY)
	fmt.Printf("%sFINAL RESULTS:%s", gd.CyanHi, gd.Reset)

	gd.MoveCursor(2, sl.MenuY+1)
	for i, player := range players {
		if i > 0 {
			fmt.Print("  ")
		}
		fmt.Printf("%s: %s%d%s credits",
			player.Name, gd.YellowHi, player.Credits, gd.Reset)
	}

	gd.MoveCursor(2, sl.MenuY+2)
	fmt.Printf("%sPress any key to return to menu...%s", gd.Yellow, gd.Reset)
}

// Utility function to add colored messages to game log
func (sl *ScreenLayout) LogMessage(message string, messageType string) {
	var coloredMessage string

	switch messageType {
	case "info":
		coloredMessage = gd.White + message + gd.Reset
	case "action":
		coloredMessage = gd.Yellow + message + gd.Reset
	case "important":
		coloredMessage = gd.RedHi + message + gd.Reset
	case "success":
		coloredMessage = gd.GreenHi + message + gd.Reset
	case "warning":
		coloredMessage = gd.YellowHi + message + gd.Reset
	default:
		coloredMessage = message
	}

	sl.GameLog.AddMessage(coloredMessage)
}
