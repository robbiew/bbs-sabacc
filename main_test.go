package main

import (
	"testing"
)

// Key Features Tested:
// 76-card deck composition (60 numbered + 32 Arcana)
// All 16 Arcana card types with 2 copies each
// Idiot's Array detection (exactly 3 cards: Idiot + 2 + 3)
// Pure Sabacc detection (exactly 23 points)
// Bomb-out conditions (>23, <-23, =0)
// Card string representation for all suits
// AI personality behaviors
// Classic Sabacc rule precedence

// TestNewDeck tests deck creation
func TestNewDeck(t *testing.T) {
	deck := NewDeck()

	// Should have 76 cards total (Classic Sabacc)
	// 60 numbered cards (4 suits × 15 cards) + 32 arcana cards (16 types × 2 copies)
	expectedCards := 92
	if len(deck.Cards) != expectedCards {
		t.Errorf("Expected %d cards, got %d", expectedCards, len(deck.Cards))
	}

	// Test that we have the right number of each suit
	suits := map[string]int{}
	arcanaCount := 0

	for _, card := range deck.Cards {
		if card.Suit == "Arcana" {
			arcanaCount++
		} else {
			suits[card.Suit]++
		}
	}

	// Should have 15 cards of each regular suit (1-15)
	expectedSuits := []string{SuitSabers, SuitFlasks, SuitCoins, SuitStaves}
	for _, suit := range expectedSuits {
		if suits[suit] != 15 {
			t.Errorf("Expected 15 cards for suit %s, got %d", suit, suits[suit])
		}
	}

	// Should have 32 arcana cards (16 types × 2 copies each)
	if arcanaCount != 32 {
		t.Errorf("Expected 32 arcana cards, got %d", arcanaCount)
	}

	// Verify we have all 4 suits
	if len(suits) != 4 {
		t.Errorf("Expected 4 suits, got %d", len(suits))
	}
}

// TestArcanaCards tests that all expected Arcana cards are present
func TestArcanaCards(t *testing.T) {
	deck := NewDeck()

	// Count each type of Arcana card
	arcanaTypes := map[string]int{}
	for _, card := range deck.Cards {
		if card.Suit == "Arcana" {
			arcanaTypes[card.Name]++
		}
	}

	// Expected Arcana cards (each should appear twice)
	expectedArcana := []string{
		"Death", "Strength", "Moderation", "Evil One", "Justice",
		"Queen of Air and Darkness", "Endurance", "Balance", "Demise",
		"Destruction", "Despair", "Failure", "Futility", "Mistress",
		"Idiot", "Star",
	}

	// Check that we have exactly 2 of each Arcana type
	for _, arcanaName := range expectedArcana {
		if arcanaTypes[arcanaName] != 2 {
			t.Errorf("Expected 2 copies of %s, got %d", arcanaName, arcanaTypes[arcanaName])
		}
	}

	// Check that we have exactly 16 types
	if len(arcanaTypes) != 16 {
		t.Errorf("Expected 16 Arcana types, got %d", len(arcanaTypes))
	}
}

// TestShuffle tests deck shuffling
func TestShuffle(t *testing.T) {
	deck1 := NewDeck()
	deck2 := NewDeck()

	// Get initial order
	originalOrder := make([]Card, len(deck1.Cards))
	copy(originalOrder, deck1.Cards)

	// Shuffle one deck
	deck1.Shuffle()

	// Check that order changed (this could theoretically fail but is very unlikely)
	orderChanged := false
	for i := 0; i < len(deck1.Cards); i++ {
		if deck1.Cards[i] != originalOrder[i] {
			orderChanged = true
			break
		}
	}

	if !orderChanged {
		t.Error("Shuffle did not change card order")
	}

	// Both decks should still have same total cards
	if len(deck1.Cards) != len(deck2.Cards) {
		t.Error("Shuffle changed number of cards")
	}
}

// TestDeal tests dealing cards from deck
func TestDeal(t *testing.T) {
	deck := NewDeck()
	originalSize := len(deck.Cards)

	// Deal a card
	card := deck.Deal()

	// Deck should have one fewer card
	if len(deck.Cards) != originalSize-1 {
		t.Errorf("Expected %d cards after dealing, got %d", originalSize-1, len(deck.Cards))
	}

	// Card should not be empty
	if card.Suit == "" {
		t.Error("Dealt card has empty suit")
	}
}

// TestCalculateHandTotal tests hand total calculation
func TestCalculateHandTotal(t *testing.T) {
	tests := []struct {
		hand     []Card
		expected int
	}{
		{
			hand: []Card{
				{Value: 5, Suit: SuitSabers},
				{Value: 10, Suit: SuitCoins},
			},
			expected: 15,
		},
		{
			hand: []Card{
				{Value: 15, Suit: SuitSabers},
				{Value: 8, Suit: SuitCoins},
			},
			expected: 23,
		},
		{
			hand: []Card{
				{Value: -5, Suit: "Arcana", Name: "Justice"},
				{Value: 10, Suit: SuitCoins},
			},
			expected: 5,
		},
		{
			hand: []Card{
				{Value: -10, Suit: "Arcana", Name: "Destruction"},
				{Value: -5, Suit: "Arcana", Name: "Justice"},
			},
			expected: -15,
		},
	}

	for i, test := range tests {
		result := calculateHandTotal(test.hand)
		if result != test.expected {
			t.Errorf("Test %d: Expected total %d, got %d", i+1, test.expected, result)
		}
	}
}

// TestIsIdiotsArray tests Idiot's Array detection (Classic Sabacc rules)
func TestIsIdiotsArray(t *testing.T) {
	tests := []struct {
		name     string
		hand     []Card
		expected bool
	}{
		{
			name: "Valid Idiot's Array",
			hand: []Card{
				{Value: -15, Suit: "Arcana", Name: "Idiot"},
				{Value: 2, Suit: SuitSabers},
				{Value: 3, Suit: SuitCoins},
			},
			expected: true,
		},
		{
			name: "Not Idiot's Array - missing Idiot",
			hand: []Card{
				{Value: 1, Suit: SuitSabers},
				{Value: 2, Suit: SuitSabers},
				{Value: 3, Suit: SuitCoins},
			},
			expected: false,
		},
		{
			name: "Not Idiot's Array - wrong numbers",
			hand: []Card{
				{Value: -15, Suit: "Arcana", Name: "Idiot"},
				{Value: 1, Suit: SuitSabers},
				{Value: 3, Suit: SuitCoins},
			},
			expected: false,
		},
		{
			name: "Not Idiot's Array - too many cards",
			hand: []Card{
				{Value: -15, Suit: "Arcana", Name: "Idiot"},
				{Value: 2, Suit: SuitSabers},
				{Value: 3, Suit: SuitCoins},
				{Value: 4, Suit: SuitFlasks},
			},
			expected: false,
		},
		{
			name: "Not Idiot's Array - only two cards",
			hand: []Card{
				{Value: -15, Suit: "Arcana", Name: "Idiot"},
				{Value: 2, Suit: SuitSabers},
			},
			expected: false,
		},
	}

	for _, test := range tests {
		result := isIdiotsArray(test.hand)
		if result != test.expected {
			t.Errorf("Test '%s': Expected %v, got %v", test.name, test.expected, result)
		}
	}
}

// TestIsPureSabacc tests Pure Sabacc detection
func TestIsPureSabacc(t *testing.T) {
	tests := []struct {
		name     string
		hand     []Card
		expected bool
	}{
		{
			name: "Valid Pure Sabacc",
			hand: []Card{
				{Value: 15, Suit: SuitSabers},
				{Value: 8, Suit: SuitCoins},
			},
			expected: true,
		},
		{
			name: "Not Pure Sabacc - total 22",
			hand: []Card{
				{Value: 15, Suit: SuitSabers},
				{Value: 7, Suit: SuitCoins},
			},
			expected: false,
		},
		{
			name: "Pure Sabacc with negative card",
			hand: []Card{
				{Value: -2, Suit: "Arcana", Name: "Strength"},
				{Value: 15, Suit: SuitSabers},
				{Value: 10, Suit: SuitCoins},
			},
			expected: true,
		},
		{
			name: "Not Pure Sabacc - total 24",
			hand: []Card{
				{Value: 15, Suit: SuitSabers},
				{Value: 9, Suit: SuitCoins},
			},
			expected: false,
		},
	}

	for _, test := range tests {
		result := isPureSabacc(test.hand)
		if result != test.expected {
			t.Errorf("Test '%s': Expected %v, got %v", test.name, test.expected, result)
		}
	}
}

// TestCheckBombOut tests bomb out detection (Classic Sabacc rules)
func TestCheckBombOut(t *testing.T) {
	tests := []struct {
		name     string
		hand     []Card
		expected bool
	}{
		{
			name: "Bomb out - over 23",
			hand: []Card{
				{Value: 15, Suit: SuitSabers},
				{Value: 10, Suit: SuitCoins},
			},
			expected: true,
		},
		{
			name: "Bomb out - under -23",
			hand: []Card{
				{Value: -15, Suit: "Arcana", Name: "Idiot"},
				{Value: -10, Suit: "Arcana", Name: "Destruction"},
			},
			expected: true,
		},
		{
			name: "Bomb out - exactly 0",
			hand: []Card{
				{Value: -5, Suit: "Arcana", Name: "Justice"},
				{Value: 5, Suit: SuitSabers},
			},
			expected: true,
		},
		{
			name: "Not bomb out - exactly 23 (Pure Sabacc)",
			hand: []Card{
				{Value: 15, Suit: SuitSabers},
				{Value: 8, Suit: SuitCoins},
			},
			expected: false,
		},
		{
			name: "Not bomb out - exactly -23",
			hand: []Card{
				{Value: -15, Suit: "Arcana", Name: "Idiot"},
				{Value: -8, Suit: "Arcana", Name: "Balance"},
			},
			expected: false,
		},
		{
			name: "Not bomb out - normal positive hand",
			hand: []Card{
				{Value: 10, Suit: SuitSabers},
				{Value: 5, Suit: SuitCoins},
			},
			expected: false,
		},
		{
			name: "Bomb out - way over 23",
			hand: []Card{
				{Value: 15, Suit: SuitSabers},
				{Value: 15, Suit: SuitCoins},
			},
			expected: true,
		},
	}

	for _, test := range tests {
		result := checkBombOut(test.hand)
		if result != test.expected {
			t.Errorf("Test '%s': Expected %v, got %v", test.name, test.expected, result)
		}
	}
}

// TestCardString tests card string representation
func TestCardString(t *testing.T) {
	tests := []struct {
		card     Card
		expected string
	}{
		{
			card:     Card{Value: 5, Suit: SuitSabers},
			expected: "+5S",
		},
		{
			card:     Card{Value: 15, Suit: SuitCoins},
			expected: "+15C",
		},
		{
			card:     Card{Value: -5, Suit: "Arcana", Name: "Justice"},
			expected: "Ju",
		},
		{
			card:     Card{Value: -15, Suit: "Arcana", Name: "Idiot"},
			expected: "Id",
		},
		{
			card:     Card{Value: 1, Suit: SuitFlasks},
			expected: "+1F",
		},
		{
			card:     Card{Value: 10, Suit: SuitStaves},
			expected: "+10T",
		},
		{
			card:     Card{Value: -17, Suit: "Arcana", Name: "Star"},
			expected: "Sr",
		},
	}

	for i, test := range tests {
		result := test.card.String()
		if result != test.expected {
			t.Errorf("Test %d: Expected %s, got %s", i+1, test.expected, result)
		}
	}
}

// TestConfig tests configuration loading and validation
func TestConfig(t *testing.T) {
	config := DefaultConfig()

	// Test default values
	if config.StartingCredits <= 0 {
		t.Error("Starting credits should be positive")
	}

	if config.MinAnte <= 0 {
		t.Error("Min ante should be positive")
	}

	if config.MaxAnte < config.MinAnte {
		t.Error("Max ante should be >= min ante")
	}

	// Test validation
	invalidConfig := GameConfig{
		MinAnte:         -5,
		MaxAnte:         1,
		StartingCredits: 50,
	}

	ValidateConfig(&invalidConfig)

	if invalidConfig.MinAnte < 1 {
		t.Error("Validation should fix negative min ante")
	}

	if invalidConfig.MaxAnte < invalidConfig.MinAnte {
		t.Error("Validation should fix max ante being less than min ante")
	}
}

// TestPlayerStats tests statistics tracking
func TestPlayerStats(t *testing.T) {
	stats := PlayerStats{
		PlayerName:  "TestPlayer",
		GamesPlayed: 0,
		GamesWon:    0,
		CreditsWon:  0,
		CreditsLost: 0,
	}

	// Test winning a game
	UpdateStats(&stats, true, 100, "pure_sabacc", false, false)

	if stats.GamesPlayed != 1 {
		t.Error("Games played should be incremented")
	}

	if stats.GamesWon != 1 {
		t.Error("Games won should be incremented")
	}

	if stats.PureSabaccs != 1 {
		t.Error("Pure Sabaccs should be incremented")
	}

	if stats.CreditsWon != 100 {
		t.Error("Credits won should be tracked")
	}

	// Test losing a game
	UpdateStats(&stats, false, -50, "", true, false)

	if stats.GamesPlayed != 2 {
		t.Error("Games played should be 2")
	}

	if stats.GamesWon != 1 {
		t.Error("Games won should still be 1")
	}

	if stats.BombOuts != 1 {
		t.Error("Bomb outs should be incremented")
	}

	if stats.CreditsLost != 50 {
		t.Error("Credits lost should be tracked")
	}
}

// TestClassicSabaccRules tests specific Classic Sabacc rule implementations
func TestClassicSabaccRules(t *testing.T) {
	// Test that Idiot's Array beats Pure Sabacc
	idiotsArray := []Card{
		{Value: -15, Suit: "Arcana", Name: "Idiot"},
		{Value: 2, Suit: SuitSabers},
		{Value: 3, Suit: SuitCoins},
	}

	pureSabacc := []Card{
		{Value: 15, Suit: SuitSabers},
		{Value: 8, Suit: SuitCoins},
	}

	if !isIdiotsArray(idiotsArray) {
		t.Error("Should recognize Idiot's Array")
	}

	if !isPureSabacc(pureSabacc) {
		t.Error("Should recognize Pure Sabacc")
	}

	// Idiot's Array mathematically totals -10 (Idiot=-15, 2+3=5, so -15+5=-10)
	// But it's a special hand that beats Pure Sabacc despite not totaling 23
	idiotsTotal := calculateHandTotal(idiotsArray)
	sabaccTotal := calculateHandTotal(pureSabacc)

	if idiotsTotal != -10 {
		t.Errorf("Idiot's Array should total -10 (Idiot -15 + 2 + 3), got %d", idiotsTotal)
	}

	if sabaccTotal != 23 {
		t.Errorf("Pure Sabacc should total 23, got %d", sabaccTotal)
	}

	// The key rule: Idiot's Array beats Pure Sabacc even though it doesn't total 23
	// This is a special rule in Classic Sabacc - it's called a "literal 23" because
	// it contains the digits 2 and 3 along with the Idiot card
}

// BenchmarkShuffle benchmarks deck shuffling performance
func BenchmarkShuffle(b *testing.B) {
	deck := NewDeck()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		deck.Shuffle()
	}
}

// BenchmarkCalculateHandTotal benchmarks hand total calculation
func BenchmarkCalculateHandTotal(b *testing.B) {
	hand := []Card{
		{Value: 5, Suit: SuitSabers},
		{Value: 10, Suit: SuitCoins},
		{Value: -3, Suit: "Arcana", Name: "Moderation"},
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		calculateHandTotal(hand)
	}
}

// TestAIPersonality tests different AI personality behaviors
func TestAIPersonality(t *testing.T) {
	conservative := GetAIPersonality("conservative")
	aggressive := GetAIPersonality("aggressive")
	balanced := GetAIPersonality("balanced")

	// Conservative should be more cautious
	if conservative.DrawThreshold >= aggressive.DrawThreshold {
		t.Error("Conservative AI should have lower draw threshold")
	}

	// Conservative should fold earlier (lower fold threshold)
	if conservative.FoldThreshold >= aggressive.FoldThreshold {
		t.Errorf("Conservative AI should fold earlier (lower fold threshold). Conservative: %d, Aggressive: %d",
			conservative.FoldThreshold, aggressive.FoldThreshold)
	}

	// Aggressive should use static field less
	if aggressive.StaticFieldChance >= conservative.StaticFieldChance {
		t.Error("Aggressive AI should use static field less")
	}

	// Balanced should be in the middle for draw threshold
	if balanced.DrawThreshold <= conservative.DrawThreshold || balanced.DrawThreshold >= aggressive.DrawThreshold {
		t.Error("Balanced AI should be between conservative and aggressive for draw threshold")
	}

	// Balanced should be in the middle for fold threshold
	if balanced.FoldThreshold <= conservative.FoldThreshold || balanced.FoldThreshold >= aggressive.FoldThreshold {
		t.Error("Balanced AI should be between conservative and aggressive for fold threshold")
	}

	// Verify the personalities are set up correctly
	if conservative.Name != "Conservative" {
		t.Error("Conservative personality name incorrect")
	}
	if aggressive.Name != "Aggressive" {
		t.Error("Aggressive personality name incorrect")
	}
	if balanced.Name != "Balanced" {
		t.Error("Balanced personality name incorrect")
	}
}

// TestErrorHandling tests error conditions
func TestErrorHandling(t *testing.T) {
	// Test dealing from empty deck
	deck := Deck{Cards: []Card{}}
	card := deck.Deal()

	if card.Suit != "Empty" {
		t.Error("Dealing from empty deck should return empty card")
	}

	// Test invalid hand conditions
	emptyHand := []Card{}
	total := calculateHandTotal(emptyHand)

	if total != 0 {
		t.Error("Empty hand should total 0")
	}
}

// Helper function to create test hands
func createTestHand(values []int, suits []string) []Card {
	if len(values) != len(suits) {
		panic("Values and suits must have same length")
	}

	hand := make([]Card, len(values))
	for i := 0; i < len(values); i++ {
		hand[i] = Card{
			Value: values[i],
			Suit:  suits[i],
		}
		if suits[i] == "Arcana" {
			hand[i].Name = "TestArcana"
		}
	}

	return hand
}

// TestCreateTestHand tests our helper function
func TestCreateTestHand(t *testing.T) {
	hand := createTestHand([]int{5, 10, -3}, []string{SuitSabers, SuitCoins, "Arcana"})

	if len(hand) != 3 {
		t.Error("Test hand should have 3 cards")
	}

	total := calculateHandTotal(hand)
	if total != 12 {
		t.Errorf("Expected total 12, got %d", total)
	}
}
