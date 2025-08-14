package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

// Card represents a Sabacc card
type Card struct {
	Value int
	Suit  string
	Name  string
}

// Deck represents a deck of Sabacc cards
type Deck struct {
	Cards []Card
}

// Add these structs at the top of cards.go
type CardDatabase struct {
	Filename   string
	CardWidth  int
	CardHeight int
	CardIndex  map[string]CardEntry
	FileData   []byte
}

type CardEntry struct {
	Offset int
	Length int
	Width  int
	Height int
}

type CardRenderer struct {
	Database    *CardDatabase
	CardSpacing int
}

// Card suits
const (
	SuitSabers = "Sabers"
	SuitFlasks = "Flasks"
	SuitCoins  = "Coins"
	SuitStaves = "Staves"
)

// Arcana cards (negative values)
var ArcanaCards = []Card{
	{Value: -1, Suit: "Arcana", Name: "Death"},
	{Value: -2, Suit: "Arcana", Name: "Strength"},
	{Value: -3, Suit: "Arcana", Name: "Moderation"},
	{Value: -4, Suit: "Arcana", Name: "Evil One"},
	{Value: -5, Suit: "Arcana", Name: "Justice"},
	{Value: -6, Suit: "Arcana", Name: "Queen of Air and Darkness"},
	{Value: -7, Suit: "Arcana", Name: "Endurance"},
	{Value: -8, Suit: "Arcana", Name: "Balance"},
	{Value: -9, Suit: "Arcana", Name: "Demise"},
	{Value: -10, Suit: "Arcana", Name: "Destruction"},
	{Value: -11, Suit: "Arcana", Name: "Despair"},
	{Value: -12, Suit: "Arcana", Name: "Failure"},
	{Value: -13, Suit: "Arcana", Name: "Futility"},
	{Value: -14, Suit: "Arcana", Name: "Mistress"},
	{Value: -15, Suit: "Arcana", Name: "Idiot"},
	{Value: -17, Suit: "Arcana", Name: "Star"},
}

// File format for sabacc_cards.bin:
/*
SABACC CARD DATABASE FORMAT:

HEADER (Fixed size):
- Bytes 0-3:   Magic number "SABC" (0x53414243)
- Bytes 4-5:   Version number (currently 0x0001)
- Bytes 6-7:   Number of cards in database
- Bytes 8-9:   Standard card width
- Bytes 10-11: Standard card height
- Bytes 12-15: Offset to card index table
- Bytes 16-31: Reserved for future use

CARD INDEX TABLE:
For each card:
- 8 bytes: Card ID (null-terminated, e.g. "+1S\0\0\0\0\0")
- 4 bytes: Data offset (from start of file)
- 4 bytes: Data length
- 2 bytes: Card width (if different from standard)
- 2 bytes: Card height (if different from standard)

CARD DATA SECTION:
- Raw ANSI data for each card, referenced by index table
*/

// NewCardDatabase creates or loads the card database
func NewCardDatabase(filename string) (*CardDatabase, error) {
	db := &CardDatabase{
		Filename:  filename,
		CardIndex: make(map[string]CardEntry),
	}

	// Try to load existing file
	if err := db.Load(); err != nil {
		// If file doesn't exist, create a new one with defaults
		return db.CreateDefault()
	}

	return db, nil
}

// Load reads the card database from file
func (db *CardDatabase) Load() error {
	data, err := os.ReadFile(db.Filename)
	if err != nil {
		return err
	}

	db.FileData = data

	if len(data) < 32 {
		return fmt.Errorf("invalid card database file")
	}

	// Read header
	magic := string(data[0:4])
	if magic != "SABC" {
		return fmt.Errorf("invalid magic number")
	}

	version := binary.LittleEndian.Uint16(data[4:6])
	if version != 1 {
		return fmt.Errorf("unsupported version: %d", version)
	}

	numCards := binary.LittleEndian.Uint16(data[6:8])
	db.CardWidth = int(binary.LittleEndian.Uint16(data[8:10]))
	db.CardHeight = int(binary.LittleEndian.Uint16(data[10:12]))
	indexOffset := binary.LittleEndian.Uint32(data[12:16])

	// Read card index
	db.CardIndex = make(map[string]CardEntry)

	for i := 0; i < int(numCards); i++ {
		entryOffset := int(indexOffset) + (i * 20) // 20 bytes per entry

		if entryOffset+20 > len(data) {
			break
		}

		// Read card ID (8 bytes, null-terminated)
		cardID := string(bytes.TrimRight(data[entryOffset:entryOffset+8], "\x00"))

		// Read entry data
		dataOffset := binary.LittleEndian.Uint32(data[entryOffset+8 : entryOffset+12])
		dataLength := binary.LittleEndian.Uint32(data[entryOffset+12 : entryOffset+16])
		cardWidth := binary.LittleEndian.Uint16(data[entryOffset+16 : entryOffset+18])
		cardHeight := binary.LittleEndian.Uint16(data[entryOffset+18 : entryOffset+20])

		// Use standard dimensions if card-specific ones are 0
		if cardWidth == 0 {
			cardWidth = uint16(db.CardWidth)
		}
		if cardHeight == 0 {
			cardHeight = uint16(db.CardHeight)
		}

		db.CardIndex[cardID] = CardEntry{
			Offset: int(dataOffset),
			Length: int(dataLength),
			Width:  int(cardWidth),
			Height: int(cardHeight),
		}
	}

	return nil
}

// GetCardData retrieves ANSI data for a specific card
func (db *CardDatabase) GetCardData(cardID string) ([]byte, int, int, error) {
	entry, exists := db.CardIndex[cardID]
	if !exists {
		return nil, 0, 0, fmt.Errorf("card not found: %s", cardID)
	}

	if entry.Offset+entry.Length > len(db.FileData) {
		return nil, 0, 0, fmt.Errorf("invalid card data offset")
	}

	data := db.FileData[entry.Offset : entry.Offset+entry.Length]
	return data, entry.Width, entry.Height, nil
}

// CreateDefault creates a new card database with all Sabacc cards
func (db *CardDatabase) CreateDefault() (*CardDatabase, error) {
	var buffer bytes.Buffer
	cardData := make(map[string][]byte)

	// Generate all cards
	suits := []string{SuitSabers, SuitFlasks, SuitCoins, SuitStaves}

	// Regular numbered cards (1-15 for each suit)
	for _, suit := range suits {
		for value := 1; value <= 15; value++ {
			card := Card{Value: value, Suit: suit}
			cardID := card.String()
			ansiData := db.generateCardANSI(card)
			cardData[cardID] = ansiData
		}
	}

	// Arcana cards
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
		// One copy of each arcana card (1989 Classic Rules)
		card := Card{Value: arcana.value, Suit: "Arcana", Name: arcana.name}
		cardID := arcana.abbrev
		ansiData := db.generateCardANSI(card)
		cardData[cardID] = ansiData
	}

	// Add face-down card
	cardData["BACK"] = db.generateBackCardANSI()

	// Set standard dimensions - diamond shaped Sabacc cards
	db.CardWidth = 6  
	db.CardHeight = 5 

	// Write header
	header := make([]byte, 32)
	copy(header[0:4], "SABC")                                           // Magic
	binary.LittleEndian.PutUint16(header[4:6], 1)                       // Version
	binary.LittleEndian.PutUint16(header[6:8], uint16(len(cardData)))   // Number of cards
	binary.LittleEndian.PutUint16(header[8:10], uint16(db.CardWidth))   // Standard width
	binary.LittleEndian.PutUint16(header[10:12], uint16(db.CardHeight)) // Standard height
	binary.LittleEndian.PutUint32(header[12:16], 32)                    // Index offset (after header)

	buffer.Write(header)

	// Calculate data section offset
	indexSize := len(cardData) * 20 // 20 bytes per index entry
	dataOffset := 32 + indexSize

	// Write index table and collect data for data section
	var dataSection bytes.Buffer
	currentDataOffset := dataOffset

	db.CardIndex = make(map[string]CardEntry)

	for cardID, ansiData := range cardData {
		// Write index entry
		indexEntry := make([]byte, 20)

		// Card ID (8 bytes, null-padded)
		idBytes := []byte(cardID)
		if len(idBytes) > 8 {
			idBytes = idBytes[:8]
		}
		copy(indexEntry[0:8], idBytes)

		// Data offset and length
		binary.LittleEndian.PutUint32(indexEntry[8:12], uint32(currentDataOffset))
		binary.LittleEndian.PutUint32(indexEntry[12:16], uint32(len(ansiData)))
		binary.LittleEndian.PutUint16(indexEntry[16:18], 0) // Use standard width
		binary.LittleEndian.PutUint16(indexEntry[18:20], 0) // Use standard height

		buffer.Write(indexEntry)

		// Add to data section
		dataSection.Write(ansiData)

		// Store in index
		db.CardIndex[cardID] = CardEntry{
			Offset: currentDataOffset,
			Length: len(ansiData),
			Width:  db.CardWidth,
			Height: db.CardHeight,
		}

		currentDataOffset += len(ansiData)
	}

	// Write data section
	buffer.Write(dataSection.Bytes())

	// Save to file
	db.FileData = buffer.Bytes()
	err := os.WriteFile(db.Filename, db.FileData, 0644)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Created card database with %d cards (%d bytes)\n",
		len(cardData), len(db.FileData))

	return db, nil
}

// generateCardANSI creates ANSI art for a card
// generateCardANSI creates ANSI art for a card using the hexagonal Sabacc shape
func (db *CardDatabase) generateCardANSI(card Card) []byte {
	var buffer bytes.Buffer

	if card.Suit == "BACK" || card.String() == "BACK" {
		// Diamond-shaped back card
		lines := []string{
			"  \x1b[31m\xdb\xdb\x1b[1m\xdb\x1b[0m",
			" \x1b[31m\xdb\xdb\xdb\xdb\x1b[1m\xdb\x1b[0m",
			"\x1b[31m\xdb\x1b[30;41m  \x1b[1;37m?\x1b[0;30;41m  \x1b[1;31;40m\xdb\x1b[0m",
			" \x1b[31m\xdb\xdb\xdb\xdb\x1b[1m\xdb\x1b[0m",
			"  \x1b[31m\xdb\xdb\x1b[1m\xdb\x1b[0m",
		}

		for i, line := range lines {
			buffer.WriteString(line)
			if i < len(lines)-1 {
				buffer.WriteString("\r\n")
			}
		}

		return buffer.Bytes()
	}

	// Regular card using diamond/hexagonal shape like reference
	var displayValue string
	if card.Suit == "Arcana" {
		displayValue = card.String() // This will be "Id", "De", etc.
		if len(displayValue) > 2 {
			displayValue = displayValue[:2]
		}
	} else {
		displayValue = fmt.Sprintf("%d", card.Value)
	}

	// Get proper color codes for suit
	var colorCode string
	var bgColorCode string
	switch card.Suit {
	case SuitSabers:
		colorCode = "34"   // Blue text
		bgColorCode = "44" // Blue background
	case SuitFlasks:
		colorCode = "32"   // Green text
		bgColorCode = "42" // Green background
	case SuitCoins:
		colorCode = "33"   // Yellow text
		bgColorCode = "43" // Yellow background
	case SuitStaves:
		colorCode = "31"   // Red text
		bgColorCode = "41" // Red background
	case "Arcana":
		colorCode = "35"   // Magenta text
		bgColorCode = "45" // Magenta background
	default:
		colorCode = "37"   // White text
		bgColorCode = "47" // White background
	}

	// Create diamond-shaped card like the reference (5 lines)
	// Line format matches the ANSI reference provided
	lines := []string{
		fmt.Sprintf("  \x1b[%sm\xdb\xdb\x1b[1m\xdb\x1b[0m", colorCode),
		fmt.Sprintf(" \x1b[%sm\xdb\xdb\xdb\xdb\x1b[1m\xdb\x1b[0m", colorCode),
		fmt.Sprintf("\x1b[%sm\xdb\x1b[30;%sm  \x1b[1;33m%s\x1b[0;30;%sm  \x1b[1;%s;40m\xdb\x1b[0m", colorCode, bgColorCode, displayValue, bgColorCode, colorCode),
		fmt.Sprintf(" \x1b[%sm\xdb\xdb\xdb\xdb\x1b[1m\xdb\x1b[0m", colorCode),
		fmt.Sprintf("  \x1b[%sm\xdb\xdb\x1b[1m\xdb\x1b[0m", colorCode),
	}

	for i, line := range lines {
		buffer.WriteString(line)
		if i < len(lines)-1 {
			buffer.WriteString("\r\n")
		}
	}

	return buffer.Bytes()
}

// generateBackCardANSI creates ANSI art for face-down card using CP437
func (db *CardDatabase) generateBackCardANSI() []byte {
	var buffer bytes.Buffer

	lines := []string{
		"\x1b[31m\xda\xc4\xc4\xc4\xc4\xc4\xc4\xc4\xbf\x1b[0m",
		"\x1b[31m\xb3\x1b[37;1m \xb0\xb0\xb0\xb0\xb0 \x1b[31m\xb3\x1b[0m",
		"\x1b[31m\xb3\x1b[37;1m \xb0\xdb\xdb\xdb\xb0 \x1b[31m\xb3\x1b[0m",
		"\x1b[31m\xb3\x1b[37;1m \xb0\xb0\xb0\xb0\xb0 \x1b[31m\xb3\x1b[0m",
		"\x1b[31m\xb3\x1b[37;1m \xb0\xdb\xdb\xdb\xb0 \x1b[31m\xb3\x1b[0m",
		"\x1b[31m\xc0\xc4\xc4\xc4\xc4\xc4\xc4\xc4\xd9\x1b[0m",
	}

	for i, line := range lines {
		buffer.WriteString(line)
		if i < len(lines)-1 {
			buffer.WriteString("\r\n")
		}
	}

	return buffer.Bytes()
}

// NewCardRenderer creates renderer with single file database
func NewCardRenderer() *CardRenderer {
	db, err := NewCardDatabase("sabacc_cards.bin")
	if err != nil {
		fmt.Printf("Warning: Could not load card database: %v\n", err)
		// Create default database
		db, _ = (&CardDatabase{}).CreateDefault()
	}

	return &CardRenderer{
		Database:    db,
		CardSpacing: 1,
	}
}

// RenderCard renders a single card at current position
func (cr *CardRenderer) RenderCard(card Card) {
	cardID := card.String()

	data, _, height, err := cr.Database.GetCardData(cardID)
	if err != nil {
		// Fallback to ASCII if card not found
		cr.renderFallbackCard(card)
		return
	}

	// Parse and render ANSI data
	lines := strings.Split(string(data), "\r\n")

	startX, startY := cr.getCurrentPos()

	for i, line := range lines {
		if i >= height {
			break
		}
		MoveCursor(startX, startY+i)
		fmt.Print(line)
	}
}

// RenderCards renders multiple cards horizontally
func (cr *CardRenderer) RenderCards(cards []Card, startX, startY int) {
	for i, card := range cards {
		cardX := startX + (i * (cr.Database.CardWidth + cr.CardSpacing))
		MoveCursor(cardX, startY)

		cardID := card.String()
		data, _, height, err := cr.Database.GetCardData(cardID)
		if err != nil {
			cr.renderFallbackCardAt(card, cardX, startY)
			continue
		}

		// Render card
		lines := strings.Split(string(data), "\r\n")
		for row, line := range lines {
			if row >= height {
				break
			}
			MoveCursor(cardX, startY+row)
			fmt.Print(line)
		}
	}
}

// RenderFaceDownCard renders a face-down card
func (cr *CardRenderer) RenderFaceDownCard(x, y int) {
	data, _, height, err := cr.Database.GetCardData("BACK")
	if err != nil {
		// Fallback ASCII back card
		cr.renderFallbackBack(x, y)
		return
	}

	lines := strings.Split(string(data), "\r\n")
	for row, line := range lines {
		if row >= height {
			break
		}
		MoveCursor(x, y+row)
		fmt.Print(line)
	}
}

// Utility functions
func (cr *CardRenderer) getCurrentPos() (int, int) {
	// Placeholder - would query terminal cursor position
	return 1, 1
}

func (cr *CardRenderer) renderFallbackCard(card Card) {
	// Simple ASCII fallback
	fmt.Printf("[%s]", card.String())
}

func (cr *CardRenderer) renderFallbackCardAt(card Card, x, y int) {
	MoveCursor(x, y)
	fmt.Printf("[%s]", card.String())
}

func (cr *CardRenderer) renderFallbackBack(x, y int) {
	MoveCursor(x, y)
	fmt.Print("[??]")
}


// Also check that your cards.go has proper String() method:
func (c Card) String() string {
	if c.Suit == "Arcana" {
		// Show arcana cards with their names abbreviated
		switch c.Name {
		case "Death":
			return "De"
		case "Strength":
			return "St"
		case "Moderation":
			return "Mo"
		case "Evil One":
			return "Ev"
		case "Justice":
			return "Ju"
		case "Queen of Air and Darkness":
			return "Qu"
		case "Endurance":
			return "En"
		case "Balance":
			return "Ba"
		case "Demise":
			return "Dm"
		case "Destruction":
			return "Ds"
		case "Despair":
			return "Dp"
		case "Failure":
			return "Fa"
		case "Futility":
			return "Fu"
		case "Mistress":
			return "Mi"
		case "Idiot":
			return "Id"
		case "Star":
			return "Sr"
		default:
			return "??"
		}
	}

	// Regular numbered cards
	valueStr := fmt.Sprintf("%d", c.Value)
	if c.Value > 0 {
		valueStr = "+" + valueStr
	}

	switch c.Suit {
	case SuitSabers:
		return valueStr + "S"
	case SuitFlasks:
		return valueStr + "F"
	case SuitCoins:
		return valueStr + "C"
	case SuitStaves:
		return valueStr + "T"
	default:
		return valueStr
	}
}

// NewDeck creates a new 76-card Sabacc deck (1989 West End Games Classic Rules)
func NewDeck() Deck {
	var cards []Card

	// Add numbered cards (1-15) for each suit (positive values) = 60 cards
	suits := []string{SuitSabers, SuitFlasks, SuitCoins, SuitStaves}
	for _, suit := range suits {
		for value := 1; value <= 15; value++ {
			cards = append(cards, Card{Value: value, Suit: suit, Name: fmt.Sprintf("%d", value)})
		}
	}

	// Add one copy of each Arcana card (16 cards) - 1989 Classic Rules
	cards = append(cards, ArcanaCards...)

	// Total: 60 + 16 = 76 cards (authentic 1989 West End Games deck)
	return Deck{Cards: cards}
}

// Shuffle randomizes the deck using modern best practices
func (d *Deck) Shuffle() {
	// Early return for empty or single-card decks
	if len(d.Cards) <= 1 {
		return
	}
	
	// Create a new random generator with time-based seed
	// This avoids the deprecated global rand.Seed() pattern
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	
	// Use the built-in Fisher-Yates shuffle algorithm
	r.Shuffle(len(d.Cards), func(i, j int) {
		d.Cards[i], d.Cards[j] = d.Cards[j], d.Cards[i]
	})
}

// ShuffleWithSource allows using a custom random source for testing
func (d *Deck) ShuffleWithSource(source rand.Source) {
	if len(d.Cards) <= 1 {
		return
	}
	
	r := rand.New(source)
	r.Shuffle(len(d.Cards), func(i, j int) {
		d.Cards[i], d.Cards[j] = d.Cards[j], d.Cards[i]
	})
}

// ShuffleMultiple performs multiple shuffles for extra randomness
func (d *Deck) ShuffleMultiple(times int) {
	if times < 1 {
		times = 1
	}
	
	for i := 0; i < times; i++ {
		d.Shuffle()
		// Small delay to ensure different seeds
		time.Sleep(1 * time.Nanosecond)
	}
}

// Deal removes and returns the top card from the deck
func (d *Deck) Deal() Card {
	if len(d.Cards) == 0 {
		// Return a dummy card if deck is empty
		return Card{Value: 0, Suit: "Empty", Name: "Empty"}
	}

	card := d.Cards[0]
	d.Cards = d.Cards[1:]
	return card
}

// Helper functions for the main game

func calculateHandTotal(hand []Card) int {
	total := 0
	for _, card := range hand {
		total += card.Value
	}
	return total
}

// And ensure getCardColor() returns proper ANSI codes:
func getCardColor(card Card) string {
	switch card.Suit {
	case SuitSabers:
		return Blue
	case SuitFlasks:
		return Green
	case SuitCoins:
		return Yellow
	case SuitStaves:
		return Red
	case "Arcana":
		return Magenta
	default:
		// Debug: print unknown suit
		fmt.Printf("[DEBUG] Unknown suit: '%s'\n", card.Suit)
		return White
	}
}

func isIdiotsArray(hand []Card) bool {
	hasIdiot := false
	hasTwo := false
	hasThree := false

	for _, card := range hand {
		if card.Name == "Idiot" {
			hasIdiot = true
		} else if card.Value == 2 {
			hasTwo = true
		} else if card.Value == 3 {
			hasThree = true
		}
	}

	return hasIdiot && hasTwo && hasThree && len(hand) == 3
}

func handleTradeCard() {
	player := &game.Players[0]

	if len(player.Hand) < 2 {
		game.Layout.DisplayMessage("You need at least 2 cards to trade!", "error", 0)
		time.Sleep(1 * time.Second)
		return
	}

	// Create menu options for each card in hand
	var tradeOptions []MenuOption
	for i, card := range player.Hand {
		tradeOptions = append(tradeOptions, MenuOption{
			Key:         rune('1' + i),
			Description: fmt.Sprintf("Trade [%s]", card.String()),
			Enabled:     true,
		})
	}
	tradeOptions = append(tradeOptions, MenuOption{'0', "Cancel", true})

	game.Layout.ShowCompactMenu(tradeOptions)

	char, _, err := getKeyWithTimeout()
	if err != nil {
		return
	}

	// Clear the compact menu after selection
	game.Layout.ClearCompactMenu()

	choice := int(char - '0')
	if choice == 0 {
		return
	}

	if choice < 1 || choice > len(player.Hand) {
		game.Layout.DisplayMessage("Invalid choice! Enter 1-"+fmt.Sprintf("%d", len(player.Hand))+" or 0", "error", 0)
		time.Sleep(1 * time.Second)
		handleTradeCard() // Loop back for valid input
		return
	}

	// Remove the selected card
	tradedCard := player.Hand[choice-1]
	player.Hand = append(player.Hand[:choice-1], player.Hand[choice:]...)

	// Draw a new card
	if len(game.Deck.Cards) > 0 {
		newCard := game.Deck.Deal()
		player.Hand = append(player.Hand, newCard)

		game.Layout.DisplayMessage(fmt.Sprintf("You traded [%s] for [%s]", tradedCard.String(), newCard.String()), "success", 0)
	} else {
		game.Layout.DisplayMessage("No more cards in deck!", "error", 0)
	}

	time.Sleep(2 * time.Second)
}

func handleStaticField() {
	// UI Tweak 3: Use compact menu for Static Field (single row format)
	staticOptions := []MenuOption{
		{'1', "Place card in Static Field", true},
		{'2', "Remove card from Static Field", true},
		{'0', "Cancel", true},
	}

	game.Layout.ShowCompactMenu(staticOptions)

	char, _, err := getKeyWithTimeout()
	if err != nil {
		return
	}
	
	// Clear the compact menu after selection
	game.Layout.ClearCompactMenu()

	switch char {
	case '1':
		placeInStaticField()
	case '2':
		removeFromStaticField()
	case '0':
		return
	default:
		game.Layout.DisplayMessage("Invalid choice! Press 1/2/0", "error", 0)
		time.Sleep(1 * time.Second)
		handleStaticField() // Loop back for valid input
	}
}

func placeInStaticField() {
	player := &game.Players[0]

	if len(player.Hand) == 0 {
		game.Layout.DisplayMessage("No cards in hand!", "error", 0)
		time.Sleep(1 * time.Second)
		return
	}

	// Clear the previous menu before showing the new one
	game.Layout.ClearMenuArea()

	// Create menu options for each card in hand
	var cardOptions []MenuOption
	for i, card := range player.Hand {
		cardOptions = append(cardOptions, MenuOption{
			Key:         rune('1' + i),
			Description: fmt.Sprintf("Place [%s]", card.String()),
			Enabled:     true,
		})
	}
	cardOptions = append(cardOptions, MenuOption{'0', "Cancel", true})

	game.Layout.ShowMenu("Place in Static Field", cardOptions, "Card choice: ")

	char, _, err := getKeyWithTimeout()
	if err != nil {
		return
	}

	choice := int(char - '0')
	if choice == 0 {
		return
	}

	if choice < 1 || choice > len(player.Hand) {
		game.Layout.DisplayMessage("Invalid choice! Enter 1-"+fmt.Sprintf("%d", len(player.Hand))+" or 0", "error", 0)
		time.Sleep(1 * time.Second)
		placeInStaticField() // Loop back for valid input
		return
	}

	// Mark card as protected in static field (card stays in hand)
	card := player.Hand[choice-1]
	// Check if card is already in static field
	alreadyStatic := false
	for _, staticCard := range player.StaticField {
		if card.Value == staticCard.Value && card.Suit == staticCard.Suit && card.Name == staticCard.Name {
			alreadyStatic = true
			break
		}
	}
	
	if !alreadyStatic {
		player.StaticField = append(player.StaticField, card) // Add to static field (card remains in hand)
	}

	game.Layout.DisplayMessage(fmt.Sprintf("Placed [%s] in Static Field (protected from shifts)", card.String()), "success", 0)
	time.Sleep(2 * time.Second)

	// Refresh the game display to show updated static field
	displayGameScreen()
}

func removeFromStaticField() {
	player := &game.Players[0]

	if len(player.StaticField) == 0 {
		game.Layout.DisplayMessage("No cards in Static Field!", "error", 0)
		time.Sleep(1 * time.Second)
		return
	}

	// Clear the previous menu before showing the new one
	game.Layout.ClearMenuArea()

	// Create menu options for each card in static field
	var cardOptions []MenuOption
	for i, card := range player.StaticField {
		cardOptions = append(cardOptions, MenuOption{
			Key:         rune('1' + i),
			Description: fmt.Sprintf("Remove [%s]", card.String()),
			Enabled:     true,
		})
	}
	cardOptions = append(cardOptions, MenuOption{'0', "Cancel", true})

	game.Layout.ShowMenu("Remove from Static Field", cardOptions, "Card choice: ")

	char, _, err := getKeyWithTimeout()
	if err != nil {
		return
	}

	choice := int(char - '0')
	if choice == 0 {
		return
	}

	if choice < 1 || choice > len(player.StaticField) {
		game.Layout.DisplayMessage("Invalid choice! Enter 1-"+fmt.Sprintf("%d", len(player.StaticField))+" or 0", "error", 0)
		time.Sleep(1 * time.Second)
		removeFromStaticField() // Loop back for valid input
		return
	}

	// Remove card from static field (card stays in hand)
	card := player.StaticField[choice-1]
	player.StaticField = append(player.StaticField[:choice-1], player.StaticField[choice:]...)
	// Card remains in hand - it was never removed

	game.Layout.DisplayMessage(fmt.Sprintf("[%s] returned to hand", card.String()), "success", 0)
	time.Sleep(2 * time.Second)

	// Refresh the game display to show updated static field
	displayGameScreen()
}
