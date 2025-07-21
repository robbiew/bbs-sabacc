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
)

type CardDefinition struct {
	ID     string
	Value  int
	Suit   string
	Name   string
	Symbol string
	Color  string
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
	binary.LittleEndian.PutUint16(header[8:10], 9)                 // Standard width
	binary.LittleEndian.PutUint16(header[10:12], 7)                // Standard height
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

	// Create a text index for reference
	createCardIndex(cards)

	return nil
}

func generateAllCards() []CardDefinition {
	var cards []CardDefinition

	// Regular numbered cards (1-15 for each suit)
	suits := []struct {
		name   string
		symbol string
		color  string
		letter string
	}{
		{"Sabers", "♠", BLUE, "S"},
		{"Flasks", "♦", GREEN, "F"},
		{"Coins", "♣", YELLOW, "C"},
		{"Staves", "♥", RED, "T"},
	}

	for _, suit := range suits {
		for value := 1; value <= 15; value++ {
			cards = append(cards, CardDefinition{
				ID:     fmt.Sprintf("+%d%s", value, suit.letter),
				Value:  value,
				Suit:   suit.name,
				Symbol: suit.symbol,
				Color:  suit.color,
			})
		}
	}

	// Arcana cards (two copies each)
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
		{"Star", -17, "Sr"},
	}

	for _, arcana := range arcanaCards {
		// First copy
		cards = append(cards, CardDefinition{
			ID:     arcana.abbrev,
			Value:  arcana.value,
			Suit:   "Arcana",
			Name:   arcana.name,
			Symbol: "★",
			Color:  MAGENTA,
		})

		// Second copy (with "2" suffix)
		cards = append(cards, CardDefinition{
			ID:     arcana.abbrev + "2",
			Value:  arcana.value,
			Suit:   "Arcana",
			Name:   arcana.name,
			Symbol: "★",
			Color:  MAGENTA,
		})
	}

	// Face-down card
	cards = append(cards, CardDefinition{
		ID:     "BACK",
		Value:  0,
		Suit:   "Back",
		Symbol: "░",
		Color:  RED,
	})

	return cards
}

func generateCardANSI(card CardDefinition) []byte {
	var buffer bytes.Buffer

	if card.ID == "BACK" {
		// Special case for back card using CP437 characters
		lines := []string{
			"\x1b[31m\xda\xc4\xc4\xc4\xc4\xc4\xc4\xc4\xc4\xbf\x1b[0m",
			"\x1b[31m\xb3\x1b[37;1m \xb0\xb0\xb0\xb0\xb0 \x1b[31m\xb3\x1b[0m",
			"\x1b[31m\xb3\x1b[37;1m \xb0\xdb\xdb\xdb\xb0 \x1b[31m\xb3\x1b[0m",
			"\x1b[31m\xb3\x1b[37;1m \xb0\xb0\xb0\xb0\xb0 \x1b[31m\xb3\x1b[0m",
			"\x1b[31m\xb3\x1b[37;1m \xb0\xdb\xdb\xdb\xb0 \x1b[31m\xb3\x1b[0m",
			"\x1b[31m\xb3\x1b[37;1m \xb0\xb0\xb0\xb0\xb0 \x1b[31m\xb3\x1b[0m",
			"\x1b[31m\xc0\xc4\xc4\xc4\xc4\xc4\xc4\xc4\xc4\xd9\x1b[0m",
		}

		for i, line := range lines {
			buffer.WriteString(line)
			if i < len(lines)-1 {
				buffer.WriteString("\r\n")
			}
		}

		return buffer.Bytes()
	}

	// Regular card using proper CP437 box drawing characters
	var displayValue string
	if card.Suit == "Arcana" {
		displayValue = card.ID
		if len(displayValue) > 2 {
			displayValue = displayValue[:2]
		}
	} else {
		if card.Value < 10 {
			displayValue = fmt.Sprintf(" %d", card.Value)
		} else {
			displayValue = fmt.Sprintf("%d", card.Value)
		}
	}

	// Ensure 3-character width for display
	for len(displayValue) < 3 {
		displayValue = " " + displayValue
	}
	if len(displayValue) > 3 {
		displayValue = displayValue[:3]
	}

	// Get proper suit symbol using CP437
	var suitChar string
	switch card.Suit {
	case "Sabers":
		suitChar = "\x06" // CP437 spade ♠
	case "Flasks":
		suitChar = "\x04" // CP437 diamond ♦
	case "Coins":
		suitChar = "\x05" // CP437 club ♣
	case "Staves":
		suitChar = "\x03" // CP437 heart ♥
	case "Arcana":
		suitChar = "\x0f" // CP437 star ☼
	default:
		suitChar = "?"
	}

	// Build card using CP437 box drawing characters
	lines := []string{
		card.Color + "\xda\xc4\xc4\xc4\xc4\xc4\xc4\xc4\xbf" + RESET,
		card.Color + "\xb3" + WHITE_HI + displayValue + "    " + card.Color + "\xb3" + RESET,
		card.Color + "\xb3       \xb3" + RESET,
		card.Color + "\xb3   " + card.Color + suitChar + "   \xb3" + RESET,
		card.Color + "\xb3       \xb3" + RESET,
		card.Color + "\xb3    " + WHITE_HI + displayValue + card.Color + "\xb3" + RESET,
		card.Color + "\xc0\xc4\xc4\xc4\xc4\xc4\xc4\xc4\xd9" + RESET,
	}

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
	buffer.WriteString("Arcana cards: 32 (16 types × 2 copies)\n")
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

	// Read database
	data, err := os.ReadFile("sabacc_cards.bin")
	if err != nil {
		return err
	}

	if len(data) < 32 {
		return fmt.Errorf("invalid database file")
	}

	numCards := binary.LittleEndian.Uint16(data[6:8])
	indexOffset := binary.LittleEndian.Uint32(data[12:16])

	var preview strings.Builder
	preview.WriteString(CYAN_HI + "SABACC CARD DATABASE PREVIEW\n" + RESET)
	preview.WriteString(CYAN + "============================\n\n" + RESET)

	// Show cards in a grid layout (4 cards per row)
	previewSets := []struct {
		title string
		cards []string
	}{
		{"Sample Sabers (♠)", []string{"+1S", "+5S", "+10S", "+15S"}},
		{"Other Suits", []string{"+1F", "+1C", "+1T", "BACK"}},
		{"Arcana Cards", []string{"Id", "De", "St", "Ba"}},
	}

	for _, set := range previewSets {
		preview.WriteString(YELLOW_HI + set.title + ":\n" + RESET)
		preview.WriteString(strings.Repeat("─", len(set.title)+1) + "\n")

		// Get card data for this set
		var cardLines [][]string
		cardFound := make([]bool, len(set.cards))

		for i, targetCard := range set.cards {
			// Find card in database
			for j := 0; j < int(numCards); j++ {
				entryOffset := int(indexOffset) + (j * 20)
				if entryOffset+20 > len(data) {
					break
				}

				cardID := string(bytes.TrimRight(data[entryOffset:entryOffset+8], "\x00"))
				if cardID == targetCard {
					dataOffset := binary.LittleEndian.Uint32(data[entryOffset+8 : entryOffset+12])
					dataLength := binary.LittleEndian.Uint32(data[entryOffset+12 : entryOffset+16])

					if int(dataOffset)+int(dataLength) <= len(data) {
						cardData := data[dataOffset : dataOffset+dataLength]
						lines := strings.Split(string(cardData), "\r\n")

						// Ensure we have exactly 7 lines
						for len(lines) < 7 {
							lines = append(lines, "         ")
						}

						if len(cardLines) <= i {
							for k := len(cardLines); k <= i; k++ {
								cardLines = append(cardLines, make([]string, 7))
							}
						}

						cardLines[i] = lines
						cardFound[i] = true
					}
					break
				}
			}

			// If card not found, create placeholder
			if !cardFound[i] {
				placeholder := []string{
					"┌───────┐",
					"│ ERROR │",
					"│       │",
					"│   ?   │",
					"│       │",
					"│ ERROR │",
					"└───────┘",
				}
				if len(cardLines) <= i {
					for k := len(cardLines); k <= i; k++ {
						cardLines = append(cardLines, make([]string, 7))
					}
				}
				cardLines[i] = placeholder
			}
		}

		// Add card labels
		labelLine := ""
		for i, cardID := range set.cards {
			if i > 0 {
				labelLine += "  "
			}
			// Center the card ID under each card (9 chars wide)
			padding := (9 - len(cardID)) / 2
			labelLine += strings.Repeat(" ", padding) + WHITE_HI + cardID + RESET
			if padding*2+len(cardID) < 9 {
				labelLine += " "
			}
		}
		preview.WriteString(labelLine + "\n")

		// Render cards side by side
		for row := 0; row < 7; row++ {
			line := ""
			for cardIdx := 0; cardIdx < len(cardLines) && cardIdx < len(set.cards); cardIdx++ {
				if cardIdx > 0 {
					line += "  " // Space between cards
				}
				if row < len(cardLines[cardIdx]) {
					line += cardLines[cardIdx][row]
				} else {
					line += "         " // Empty space if line missing
				}
			}
			preview.WriteString(line + "\n")
		}

		preview.WriteString("\n")
	}

	// Add detailed summary
	preview.WriteString(CYAN_HI + "DATABASE INFORMATION:\n" + RESET)
	preview.WriteString(CYAN + "====================\n" + RESET)
	preview.WriteString(fmt.Sprintf("Total cards: %s%d%s\n", YELLOW_HI, numCards, RESET))
	preview.WriteString(fmt.Sprintf("File size: %s%d%s bytes\n", YELLOW_HI, len(data), RESET))

	// Calculate average card size
	dataSize := len(data) - int(indexOffset) - (int(numCards) * 20)
	avgSize := float64(dataSize) / float64(numCards)
	preview.WriteString(fmt.Sprintf("Average card size: %s%.1f%s bytes\n", YELLOW_HI, avgSize, RESET))

	// Card breakdown
	preview.WriteString("\nCard Breakdown:\n")
	preview.WriteString(fmt.Sprintf("• Regular cards: %s60%s (15 per suit × 4 suits)\n", GREEN_HI, RESET))
	preview.WriteString(fmt.Sprintf("• Arcana cards: %s32%s (16 types × 2 copies)\n", MAGENTA_HI, RESET))
	preview.WriteString(fmt.Sprintf("• Back card: %s1%s\n", RED_HI, RESET))

	preview.WriteString(fmt.Sprintf("\nSuits: %sSabers(♠)%s %sFlasks(♦)%s %sCoins(♣)%s %sStaves(♥)%s %sArcana(★)%s\n",
		BLUE, RESET, GREEN, RESET, YELLOW, RESET, RED, RESET, MAGENTA, RESET))

	err = os.WriteFile("card_preview.ans", []byte(preview.String()), 0644)
	if err != nil {
		return err
	}

	fmt.Println("Created card_preview.ans")
	return nil
}
