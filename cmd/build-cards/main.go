// cmd/build-cards/main.go - Separate build tool
package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"strings"
)

// ANSI color constants for card generation
const (
	ESC = "\x1b["

	// Standard ANSI colors (compatible with CP437)
	BLACK   = "\x1b[30m"
	RED     = "\x1b[31m"
	GREEN   = "\x1b[32m"
	YELLOW  = "\x1b[33m"
	BLUE    = "\x1b[34m"
	MAGENTA = "\x1b[35m"
	CYAN    = "\x1b[36m"
	WHITE   = "\x1b[37m"
	RESET   = "\x1b[0m"

	// Bright colors
	RED_HI     = "\x1b[31;1m"
	GREEN_HI   = "\x1b[32;1m"
	YELLOW_HI  = "\x1b[33;1m"
	BLUE_HI    = "\x1b[34;1m"
	MAGENTA_HI = "\x1b[35;1m"
	CYAN_HI    = "\x1b[36;1m"
	WHITE_HI   = "\x1b[37;1m"

	// Background colors
	BG_BLACK   = "\x1b[40m"
	BG_RED     = "\x1b[41m"
	BG_GREEN   = "\x1b[42m"
	BG_YELLOW  = "\x1b[43m"
	BG_BLUE    = "\x1b[44m"
	BG_MAGENTA = "\x1b[45m"
	BG_CYAN    = "\x1b[46m"
	BG_WHITE   = "\x1b[47m"
)

// Import UI package to access CP437 character constants
// Note: We're using the same constants defined in ui.go to maintain consistency

// Card suit symbols (CP437) - using proper printable characters
const (
	CP437_SPADE   = "\x06" // ♠ (spade) - keeping these for now but will use Unicode fallbacks
	CP437_DIAMOND = "\x04" // ♦ (diamond)
	CP437_CLUB    = "\x05" // ♣ (club)
	CP437_HEART   = "\x03" // ♥ (heart)
	CP437_STAR    = "\x0f" // ☼ (star/sun)

	// Additional decorative characters (CP437)
	CP437_UP_ARROW    = "\x18" // ↑ (up arrow)
	CP437_SMILEY      = "\x01" // ☺ (smiley face)
	CP437_SMILEY_FILL = "\x02" // ☻ (filled smiley face)

	// Half-blocks for card shapes
	TOP_HALF_BLOCK    = "\xdf" // ▀ (top half block)
	BOTTOM_HALF_BLOCK = "\xdc" // ▄ (bottom half block)
	LEFT_HALF_BLOCK   = "\xdd" // ▌ (left half block)
	RIGHT_HALF_BLOCK  = "\xde" // ▐ (right half block)

	// Shading and blocks
	LIGHT_SHADE  = "\xb0" // ░ (light shade)
	MEDIUM_SHADE = "\xb1" // ▒ (medium shade)
	DARK_SHADE   = "\xb2" // ▓ (dark shade)
	SOLID_BLOCK  = "\xdb" // █ (solid block)

	// ANSI cursor control sequences
	CURSOR_DOWN_LEFT_6 = "\x1b[1B\x1b[6D" // Move down 1 row, left 6 columns
	CURSOR_DOWN_LEFT_7 = "\x1b[1B\x1b[7D" // Move down 1 row, left 7 columns
	CURSOR_DOWN_LEFT_5 = "\x1b[1B\x1b[5D" // Move down 1 row, left 5 columns
	CURSOR_SAVE        = "\x1b[s"         // Save cursor position
	CURSOR_RESTORE     = "\x1b[u"         // Restore cursor position
)

// CP437 character constants (IBM PC character set)
// These match the constants defined in ui.go for consistency
const (
	// Box drawing characters (single line)
	HORIZONTAL_LINE     = "\xc4" // ─ (horizontal line)
	VERTICAL_LINE       = "\xb3" // │ (vertical line)
	TOP_LEFT_CORNER     = "\xda" // ┌ (top-left corner)
	TOP_RIGHT_CORNER    = "\xbf" // ┐ (top-right corner)
	BOTTOM_LEFT_CORNER  = "\xc0" // └ (bottom-left corner)
	BOTTOM_RIGHT_CORNER = "\xd9" // ┘ (bottom-right corner)
	LEFT_T_JUNCTION     = "\xc3" // ├ (left T-junction)
	RIGHT_T_JUNCTION    = "\xb4" // ┤ (right T-junction)

	// Box drawing characters (double line)
	DOUBLE_HORIZONTAL_LINE     = "\xcd" // ═ (double horizontal line)
	DOUBLE_VERTICAL_LINE       = "\xba" // ║ (double vertical line)
	DOUBLE_TOP_LEFT_CORNER     = "\xc9" // ╔ (double top-left corner)
	DOUBLE_TOP_RIGHT_CORNER    = "\xbb" // ╗ (double top-right corner)
	DOUBLE_BOTTOM_LEFT_CORNER  = "\xc8" // ╚ (double bottom-left corner)
	DOUBLE_BOTTOM_RIGHT_CORNER = "\xbc" // ╝ (double bottom-right corner)
)

type CardDefinition struct {
	ID      string
	Value   int
	Suit    string
	Name    string
	Symbol  string
	Color   string
	ColorHi string
}

func main() {
	fmt.Println("Sabacc Card Database Builder")
	fmt.Println("============================")

	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "test":
			if err := testCardDatabase(); err != nil {
				fmt.Printf("Test failed: %v\n", err)
				os.Exit(1)
			}
			fmt.Println("Database test passed!")
			return
		case "preview":
			if err := createANSIPreview(); err != nil {
				fmt.Printf("Preview creation failed: %v\n", err)
				os.Exit(1)
			}
			fmt.Println("Preview created successfully!")
			return
		case "help":
			showHelp()
			return
		}
	}

	// Default action: build database
	fmt.Println("Building card database...")

	if err := buildCardDatabase(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Card database built successfully!")
	fmt.Println("Files created:")
	fmt.Println("  - sabacc_cards.bin (card database)")
	fmt.Println("  - card_index.txt (card reference)")
	fmt.Println()
	fmt.Println("Copy sabacc_cards.bin to your game directory.")
}

func showHelp() {
	fmt.Println("Usage:")
	fmt.Println("  go run cmd/build-cards/main.go        - Build database")
	fmt.Println("  go run cmd/build-cards/main.go test   - Test database")
	fmt.Println("  go run cmd/build-cards/main.go preview - Create ANSI preview")
	fmt.Println("  go run cmd/build-cards/main.go help   - Show this help")
	fmt.Println()
	fmt.Println("Output files:")
	fmt.Println("  sabacc_cards.bin - Binary card database")
	fmt.Println("  card_index.txt   - Human-readable index")
	fmt.Println("  card_preview.ans - ANSI preview file")
}

func buildCardDatabase() error {
	// Define all cards
	cards := generateAllCards()

	// Build the database file
	var buffer bytes.Buffer

	// Write header
	header := make([]byte, 32)
	copy(header[0:4], "SABC")                                      // Magic number
	binary.LittleEndian.PutUint16(header[4:6], 1)                  // Version
	binary.LittleEndian.PutUint16(header[6:8], uint16(len(cards))) // Card count
	binary.LittleEndian.PutUint16(header[8:10], 5)                 // Standard width (5 chars wide)
	binary.LittleEndian.PutUint16(header[10:12], 5)                // Standard height (5 lines tall)
	binary.LittleEndian.PutUint32(header[12:16], 32)               // Index offset

	buffer.Write(header)

	// Calculate offsets
	indexSize := len(cards) * 20
	dataOffset := 32 + indexSize

	// Generate card data and build index
	var dataSection bytes.Buffer
	currentOffset := dataOffset

	for _, card := range cards {
		ansiData := generateCardANSI(card)

		// Write index entry
		indexEntry := make([]byte, 20)

		// Card ID (8 bytes, null-padded)
		idBytes := []byte(card.ID)
		if len(idBytes) > 8 {
			idBytes = idBytes[:8]
		}
		copy(indexEntry[0:8], idBytes)

		// Offset and length
		binary.LittleEndian.PutUint32(indexEntry[8:12], uint32(currentOffset))
		binary.LittleEndian.PutUint32(indexEntry[12:16], uint32(len(ansiData)))
		binary.LittleEndian.PutUint16(indexEntry[16:18], 0) // Use standard width
		binary.LittleEndian.PutUint16(indexEntry[18:20], 0) // Use standard height

		buffer.Write(indexEntry)

		// Add to data section
		dataSection.Write(ansiData)
		currentOffset += len(ansiData)
	}

	// Write data section
	buffer.Write(dataSection.Bytes())

	// Save file
	err := os.WriteFile("sabacc_cards.bin", buffer.Bytes(), 0644)
	if err != nil {
		return err
	}

	fmt.Printf("Created database with %d cards (%d bytes)\n",
		len(cards), buffer.Len())

	// Debug: count card types
	regularCount := 0
	arcanaCount := 0
	backCount := 0
	for _, card := range cards {
		switch card.Suit {
		case "Sabers", "Flasks", "Coins", "Staves":
			regularCount++
		case "Arcana":
			arcanaCount++
		case "Back":
			backCount++
		}
	}
	fmt.Printf("Debug: Regular=%d, Arcana=%d, Back=%d, Total=%d\n", regularCount, arcanaCount, backCount, len(cards))

	// Create a text index for reference
	createCardIndex(cards)

	return nil
}

func generateAllCards() []CardDefinition {
	var cards []CardDefinition

	// 1989 West End Games Classic Sabacc suits
	suits := []struct {
		name    string
		symbol  string
		color   string
		ColorHi string
		letter  string
	}{
		{"Sabers", CP437_UP_ARROW, BLUE, BLUE_HI, "S"}, // ↑ up arrow (representing sword/saber)
		{"Flasks", CP437_SPADE, GREEN, GREEN_HI, "F"},  // ♠ spade (representing a cup)
		{"Coins", CP437_STAR, YELLOW, YELLOW_HI, "C"},  // ☼ sun symbol (representing a coin)
		{"Staves", "|", RED, RED_HI, "T"},              // | vertical line (representing a stick)
	}

	// Regular numbered cards (1-15) for each suit = 60 cards total
	for _, suit := range suits {
		for value := 1; value <= 15; value++ {
			cards = append(cards, CardDefinition{
				ID:      fmt.Sprintf("+%d%s", value, suit.letter),
				Value:   value,
				Suit:    suit.name,
				Symbol:  suit.symbol,
				Color:   suit.color,
				ColorHi: suit.ColorHi,
			})
		}
	}

	// Arcana cards (one copy each) = 16 cards total - 1989 Classic Rules
	arcanaCards := []struct {
		name   string
		value  int
		abbrev string
	}{
		{"Death", -1, "De"},
		{"Strength", -2, "St"},
		{"Moderation", -3, "Mo"},
		{"Evil One", -4, "Ev"},
		{"Justice", -5, "Ju"},
		{"Queen of Air and Darkness", -6, "Qu"},
		{"Endurance", -7, "En"},
		{"Balance", -8, "Ba"},
		{"Demise", -9, "Dm"},
		{"Destruction", -10, "Ds"},
		{"Despair", -11, "Dp"},
		{"Failure", -12, "Fa"},
		{"Futility", -13, "Fu"},
		{"Mistress", -14, "Mi"},
		{"Idiot", -15, "Id"},
		{"The Wheel", -16, "Wh"},
		{"Star", -17, "Sr"},
	}

	for _, arcana := range arcanaCards {
		cards = append(cards, CardDefinition{
			ID:      arcana.abbrev,
			Value:   arcana.value,
			Suit:    "Arcana",
			Name:    arcana.name,
			Symbol:  CP437_STAR,
			Color:   MAGENTA,
			ColorHi: MAGENTA_HI,
		})
	}

	// Face-down card
	cards = append(cards, CardDefinition{
		ID:      "BACK",
		Value:   0,
		Suit:    "Back",
		Symbol:  "░",
		Color:   RED,
		ColorHi: RED_HI,
	})

	return cards
}

func generateCardANSI(card CardDefinition) []byte {
	var buffer bytes.Buffer

	// Determine background and foreground colors based on suit
	var bgColor string
	var fgColor string
	var hiColor string
	switch card.Suit {
	case "Sabers":
		bgColor = BG_BLUE
		hiColor = BLUE_HI
		fgColor = WHITE_HI // Bright white for good contrast on blue
	case "Flasks":
		bgColor = BG_GREEN
		fgColor = BLACK // Black for good contrast on green
		hiColor = GREEN_HI
	case "Coins":
		bgColor = BG_YELLOW
		fgColor = BLACK // Black for good contrast on yellow
		hiColor = YELLOW_HI
	case "Staves":
		bgColor = BG_RED
		fgColor = WHITE_HI // Bright white for good contrast on red
		hiColor = RED_HI
	case "Arcana":
		bgColor = BG_MAGENTA
		fgColor = WHITE_HI // Bright white for good contrast on magenta
		hiColor = MAGENTA_HI
	case "Back":
		bgColor = BG_RED
		fgColor = WHITE_HI // Bright white for good contrast on red
		hiColor = RED_HI
	default:
		bgColor = BG_BLACK
		fgColor = WHITE_HI // Bright white for good contrast on black
		hiColor = BLUE_HI
	}

	// Prepare display value for the card
	var displayValue string

	if card.ID == "BACK" {
		displayValue = "?"
	} else if card.Suit == "Arcana" || card.Suit == "Special" {
		displayValue = card.ID
		if len(displayValue) > 3 {
			displayValue = displayValue[:3]
		}
	} else {
		displayValue = fmt.Sprintf("%d", card.Value)
	}

	// Prepare suit symbol for the card
	var suitSymbol string
	if card.ID == "BACK" {
		suitSymbol = "?"
	} else {
		suitSymbol = card.Symbol
	}

	// Create card shape using background colors and spaces - ALL LINES EXACTLY 5 CHARACTERS WIDE
	lines := []string{
		" " + bgColor + fgColor + "  " + bgColor + hiColor + RIGHT_HALF_BLOCK + RESET + " ",        // Line 1: "     " (5 chars)
		bgColor + fgColor + "  " + suitSymbol + " " + bgColor + hiColor + RIGHT_HALF_BLOCK + RESET, // Line 2: "  ♠  " (5 chars)
	}

	// Line 3: " 1▐" (5 chars) - value centered in 3 chars + border
	middlePadding := (3 - len(displayValue)) / 2
	remainingPadding := 3 - len(displayValue) - middlePadding
	middleLine := bgColor + fgColor + " " + strings.Repeat(" ", middlePadding) + displayValue + strings.Repeat(" ", remainingPadding) + bgColor + hiColor + RIGHT_HALF_BLOCK + RESET
	lines = append(lines, middleLine)

	// Bottom lines - mirror the top, all exactly 5 chars
	lines = append(lines,
		bgColor+fgColor+"  "+suitSymbol+" "+bgColor+hiColor+RIGHT_HALF_BLOCK+RESET, // Line 4: "  ♠  " (5 chars)
		" "+bgColor+fgColor+"  "+bgColor+hiColor+RIGHT_HALF_BLOCK+RESET+" ")        // Line 5: "     " (5 chars)

	// Write the lines to buffer
	for i, line := range lines {
		buffer.WriteString(line)
		if i < len(lines)-1 {
			buffer.WriteString("\r\n")
		}
	}

	return buffer.Bytes()
}

func createCardIndex(cards []CardDefinition) {
	var buffer strings.Builder

	buffer.WriteString("SABACC CARD DATABASE INDEX\n")
	buffer.WriteString("==========================\n\n")
	buffer.WriteString("Card ID  | Value | Suit      | Name\n")
	buffer.WriteString("---------|-------|-----------|------------------\n")

	for _, card := range cards {
		name := card.Name
		if name == "" {
			name = fmt.Sprintf("%d of %s", card.Value, card.Suit)
		}

		buffer.WriteString(fmt.Sprintf("%-8s | %5d | %-9s | %s\n",
			card.ID, card.Value, card.Suit, name))
	}

	buffer.WriteString(fmt.Sprintf("\nTotal cards: %d\n", len(cards)))
	buffer.WriteString("Regular cards: 60 (15 per suit × 4 suits)\n")
	buffer.WriteString("Arcana cards: 16 (16 types × 1 copy each)\n")
	buffer.WriteString("Back card: 1\n")
	buffer.WriteString(fmt.Sprintf("Total: %d cards\n", len(cards)))

	err := os.WriteFile("card_index.txt", []byte(buffer.String()), 0644)
	if err != nil {
		fmt.Printf("Warning: Could not create card index: %v\n", err)
	} else {
		fmt.Println("Created card_index.txt")
	}
}

func testCardDatabase() error {
	fmt.Println("Testing card database...")

	// Read the file
	data, err := os.ReadFile("sabacc_cards.bin")
	if err != nil {
		return fmt.Errorf("could not read database: %v", err)
	}

	// Check header
	if len(data) < 32 {
		return fmt.Errorf("file too small")
	}

	magic := string(data[0:4])
	if magic != "SABC" {
		return fmt.Errorf("invalid magic number: %s", magic)
	}

	version := binary.LittleEndian.Uint16(data[4:6])
	numCards := binary.LittleEndian.Uint16(data[6:8])
	cardWidth := binary.LittleEndian.Uint16(data[8:10])
	cardHeight := binary.LittleEndian.Uint16(data[10:12])
	indexOffset := binary.LittleEndian.Uint32(data[12:16])

	fmt.Printf("Magic: %s\n", magic)
	fmt.Printf("Version: %d\n", version)
	fmt.Printf("Cards: %d\n", numCards)
	fmt.Printf("Standard size: %dx%d\n", cardWidth, cardHeight)
	fmt.Printf("Index offset: %d\n", indexOffset)
	fmt.Printf("File size: %d bytes\n", len(data))

	// Test a few sample cards
	fmt.Println("\nSample cards:")
	sampleCards := []string{"+1S", "+15S", "Id", "BACK"}

	for _, cardID := range sampleCards {
		found := false
		for i := 0; i < int(numCards); i++ {
			entryOffset := int(indexOffset) + (i * 20)
			if entryOffset+20 > len(data) {
				break
			}

			indexCardID := string(bytes.TrimRight(data[entryOffset:entryOffset+8], "\x00"))
			if indexCardID == cardID {
				dataOffset := binary.LittleEndian.Uint32(data[entryOffset+8 : entryOffset+12])
				dataLength := binary.LittleEndian.Uint32(data[entryOffset+12 : entryOffset+16])

				fmt.Printf("  %s: offset=%d, length=%d bytes\n",
					cardID, dataOffset, dataLength)
				found = true
				break
			}
		}

		if !found {
			fmt.Printf("  %s: NOT FOUND\n", cardID)
		}
	}

	return nil
}

func createANSIPreview() error {
	fmt.Println("Creating ANSI preview file...")

	// Generate all cards for preview (don't require database file)
	allCards := generateAllCards()

	var preview strings.Builder

	// Enhanced header with readable double-line border constants - exactly 40 characters
	preview.WriteString(CYAN_HI + DOUBLE_TOP_LEFT_CORNER + strings.Repeat(DOUBLE_HORIZONTAL_LINE, 38) + DOUBLE_TOP_RIGHT_CORNER + "\n" + RESET)
	preview.WriteString(CYAN_HI + DOUBLE_VERTICAL_LINE + "     " + WHITE_HI + "SABACC CARD DATABASE PREVIEW" + CYAN_HI + "     " + DOUBLE_VERTICAL_LINE + "\n" + RESET)
	preview.WriteString(CYAN_HI + DOUBLE_BOTTOM_LEFT_CORNER + strings.Repeat(DOUBLE_HORIZONTAL_LINE, 38) + DOUBLE_BOTTOM_RIGHT_CORNER + "\n\n" + RESET)

	// Show cards in a grid layout with better organization
	previewSets := []struct {
		title       string
		description string
		cards       []string
		color       string
	}{
		{"Sample Sabers (" + CP437_UP_ARROW + ")", "Positive values from the Sabers suit", []string{"+1S", "+5S", "+10S", "+15S"}, BLUE_HI},
		{"Other Suits:", "One card from each traditional suit", []string{"+7F", "+12C", "+3T", "BACK"}, GREEN_HI},
		{"Arcana Cards:", "Mystical cards with special powers", []string{"Id", "De", "Wh", "Ba"}, MAGENTA_HI},
	}

	// Create a map for quick card lookup
	cardMap := make(map[string]CardDefinition)
	for _, card := range allCards {
		cardMap[card.ID] = card
	}

	for setIdx, set := range previewSets {
		// Hardcoded section headers with exact matching widths - all exactly 40 characters total
		switch setIdx {
		case 0: // Sample Sabers - " Sample Sabers (♠) " = 19 chars, so 38-19=19, 9+9=18, use 9+10
			preview.WriteString(set.color + TOP_LEFT_CORNER + strings.Repeat(HORIZONTAL_LINE, 9) + WHITE_HI + " Sample Sabers (" + CP437_UP_ARROW + ") " + set.color + strings.Repeat(HORIZONTAL_LINE, 10) + TOP_RIGHT_CORNER + "\n" + RESET)
			preview.WriteString(set.color + VERTICAL_LINE + WHITE + " Positive values from the Sabers suit " + set.color + VERTICAL_LINE + "\n" + RESET)
			preview.WriteString(set.color + BOTTOM_LEFT_CORNER + strings.Repeat(HORIZONTAL_LINE, 38) + BOTTOM_RIGHT_CORNER + "\n" + RESET)
		case 1: // Other Suits - " Other Suits: " = 14 chars, so 38-14=24, 12+12=24
			preview.WriteString(set.color + TOP_LEFT_CORNER + strings.Repeat(HORIZONTAL_LINE, 12) + WHITE_HI + " Other Suits: " + set.color + strings.Repeat(HORIZONTAL_LINE, 12) + TOP_RIGHT_CORNER + "\n" + RESET)
			preview.WriteString(set.color + VERTICAL_LINE + WHITE + " One card from each traditional suit  " + set.color + VERTICAL_LINE + "\n" + RESET)
			preview.WriteString(set.color + BOTTOM_LEFT_CORNER + strings.Repeat(HORIZONTAL_LINE, 38) + BOTTOM_RIGHT_CORNER + "\n" + RESET)
		case 2: // Arcana Cards - " Arcana Cards: " = 15 chars, so 38-15=23, 11+12=23
			preview.WriteString(set.color + TOP_LEFT_CORNER + strings.Repeat(HORIZONTAL_LINE, 11) + WHITE_HI + " Arcana Cards: " + set.color + strings.Repeat(HORIZONTAL_LINE, 12) + TOP_RIGHT_CORNER + "\n" + RESET)
			preview.WriteString(set.color + VERTICAL_LINE + WHITE + " Mystical cards with special powers   " + set.color + VERTICAL_LINE + "\n" + RESET)
			preview.WriteString(set.color + BOTTOM_LEFT_CORNER + strings.Repeat(HORIZONTAL_LINE, 38) + BOTTOM_RIGHT_CORNER + "\n" + RESET)
		}

		// Get card data for this set
		var cardLines [][]string
		maxLines := 5 // Standard card height

		for _, targetCard := range set.cards {
			card, exists := cardMap[targetCard]
			var lines []string

			if exists {
				// Generate the card ANSI directly
				cardANSI := generateCardANSI(card)
				lines = strings.Split(string(cardANSI), "\r\n")
			} else {
				// Create a better-looking error card (CP437) - 5 characters wide
				lines = []string{
					RED_HI + TOP_LEFT_CORNER + HORIZONTAL_LINE + HORIZONTAL_LINE + HORIZONTAL_LINE + TOP_RIGHT_CORNER + RESET,
					RED_HI + VERTICAL_LINE + WHITE + "ERR" + RED_HI + VERTICAL_LINE + RESET,
					RED_HI + VERTICAL_LINE + YELLOW + " ? " + RED_HI + VERTICAL_LINE + RESET,
					RED_HI + VERTICAL_LINE + WHITE + "MIS" + RED_HI + VERTICAL_LINE + RESET,
					RED_HI + BOTTOM_LEFT_CORNER + HORIZONTAL_LINE + HORIZONTAL_LINE + HORIZONTAL_LINE + BOTTOM_RIGHT_CORNER + RESET,
				}
			}

			// Pad or trim to exactly maxLines
			for len(lines) < maxLines {
				lines = append(lines, strings.Repeat(" ", 5))
			}
			if len(lines) > maxLines {
				lines = lines[:maxLines]
			}

			cardLines = append(cardLines, lines)
		}

		// Add card labels with better formatting
		labelLine := ""
		for i, cardID := range set.cards {
			if i > 0 {
				labelLine += "  "
			}
			// Center the card ID under each card (5 characters wide)
			padding := (5 - len(cardID)) / 2
			labelLine += strings.Repeat(" ", padding) + YELLOW_HI + cardID + RESET
			remainingPadding := 5 - len(cardID) - padding
			labelLine += strings.Repeat(" ", remainingPadding)
		}
		preview.WriteString(labelLine + "\n")

		// Render cards side by side with better padding
		for row := 0; row < maxLines; row++ {
			line := ""
			for cardIdx, cardData := range cardLines {
				if cardIdx > 0 {
					line += "  " // Space between cards (2 spaces)
				}
				if row < len(cardData) {
					line += cardData[row]
				} else {
					line += strings.Repeat(" ", 5) // Empty space if line missing (5 chars)
				}
			}
			preview.WriteString(line + "\n")
		}

		// Add separator between sections
		if setIdx < len(previewSets)-1 {
			preview.WriteString("\n")
		}
	}

	// Enhanced database information section with readable double-line constants - exactly 40 characters
	preview.WriteString("\n" + CYAN_HI + DOUBLE_TOP_LEFT_CORNER + strings.Repeat(DOUBLE_HORIZONTAL_LINE, 38) + DOUBLE_TOP_RIGHT_CORNER + "\n" + RESET)
	preview.WriteString(CYAN_HI + DOUBLE_VERTICAL_LINE + "         " + WHITE_HI + "DATABASE INFORMATION" + CYAN_HI + "         " + DOUBLE_VERTICAL_LINE + "\n" + RESET)
	preview.WriteString(CYAN_HI + DOUBLE_BOTTOM_LEFT_CORNER + strings.Repeat(DOUBLE_HORIZONTAL_LINE, 38) + DOUBLE_BOTTOM_RIGHT_CORNER + "\n" + RESET)

	// Calculate statistics for 1989 Classic Sabacc
	totalCards := len(allCards)
	regularCards := 0
	arcanaCards := 0
	backCards := 0

	for _, card := range allCards {
		switch card.Suit {
		case "Sabers", "Flasks", "Coins", "Staves":
			regularCards++
		case "Arcana":
			arcanaCards++
		case "Back":
			backCards++
		}
	}

	avgCardSize := 145.6 // Approximate based on ANSI card size

	preview.WriteString(fmt.Sprintf("Total cards: %s%d%s\n", YELLOW_HI, totalCards, RESET))
	preview.WriteString(fmt.Sprintf("File size: %s%d%s bytes\n", YELLOW_HI, totalCards*20+int(avgCardSize*float64(totalCards)), RESET))
	preview.WriteString(fmt.Sprintf("Average card size: %s%.1f%s bytes\n", YELLOW_HI, avgCardSize, RESET))

	// Enhanced card breakdown with readable constants (CP437) - 1989 Classic Rules
	preview.WriteString("\n" + WHITE_HI + "Card Breakdown:\n" + RESET)
	preview.WriteString(fmt.Sprintf("%s%s Regular cards: %s%d%s (15 per suit %s 4 suits)\n", TOP_LEFT_CORNER, HORIZONTAL_LINE, GREEN_HI, regularCards, RESET, CP437_STAR))
	preview.WriteString(fmt.Sprintf("%s%s Arcana cards: %s%d%s (16 types %s 1 copy each)\n", LEFT_T_JUNCTION, HORIZONTAL_LINE, CYAN_HI, arcanaCards, RESET, CP437_STAR))
	preview.WriteString(fmt.Sprintf("%s%s Back card: %s%d%s\n", BOTTOM_LEFT_CORNER, HORIZONTAL_LINE, RED_HI, backCards, RESET))

	// Suit legend with CP437 symbols
	preview.WriteString(fmt.Sprintf("\n%sSuits:%s %sSabers%s(%s%s%s) %sFlasks%s(%s%s%s) %sCoins%s(%s%s%s) %sStaves%s(%s%s%s) %sArcana%s(%s%s%s)\n",
		WHITE_HI, RESET,
		BLUE, RESET, BLUE, CP437_UP_ARROW, RESET,
		GREEN, RESET, GREEN, CP437_SPADE, RESET,
		YELLOW, RESET, YELLOW, CP437_STAR, RESET,
		RED, RESET, RED, "|", RESET,
		MAGENTA, RESET, MAGENTA, CP437_STAR, RESET))

	// Footer with generation info using readable constants - exactly 39 characters
	preview.WriteString(fmt.Sprintf("\n%s%s%s\n", CYAN, strings.Repeat(HORIZONTAL_LINE, 39), RESET))
	preview.WriteString(fmt.Sprintf("%sGenerated by Sabacc Card Database Builder%s\n", WHITE, RESET))

	err := os.WriteFile("card_preview.ans", []byte(preview.String()), 0644)
	if err != nil {
		return err
	}

	fmt.Println("Created enhanced card_preview.ans")
	return nil
}
