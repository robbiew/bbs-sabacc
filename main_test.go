package main

import (
	"testing"
)

// TestNewDeck tests deck creation
func TestNewDeck(t *testing.T) {
	deck := NewDeck()

	// Should have 76 cards total
	expectedCards := 92 // 60 numbered cards + 32 arcana cards (16 types x 2 copies)
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

	// Should have 15 cards of each regular suit
	for suit, count := range suits {
		if count != 15 {
			t.Errorf("Expected 15 cards for suit %s, got %d", suit, count)
		}
	}

	// Should have 32 arcana cards (16 types x 2 copies)
	if arcanaCount != 32 {
		t.Errorf("Expected 32 arcana cards, got %d", arcanaCount)
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

// TestIsIdiotsArray tests Idiot's Array detection
func TestIsIdiotsArray(t *testing.T) {
	tests := []struct {
		hand     []Card
		expected bool
	}{
		{
			// Valid Idiot's Array
			hand: []Card{
				{Value: -15, Suit: "Arcana", Name: "Idiot"},
				{Value: 2, Suit: SuitSabers},
				{Value: 3, Suit: SuitCoins},
			},
			expected: true,
		},
		{
			// Not Idiot's Array - missing Idiot
			hand: []Card{
				{Value: 1, Suit: SuitSabers},
				{Value: 2, Suit: SuitSabers},
				{Value: 3, Suit: SuitCoins},
			},
			expected: false,
		},
		{
			// Not Idiot's Array - wrong numbers
			hand: []Card{
				{Value: -15, Suit: "Arcana", Name: "Idiot"},
				{Value: 1, Suit: SuitSabers},
				{Value: 3, Suit: SuitCoins},
			},
			expected: false,
		},
		{
			// Not Idiot's Array - too many cards
			hand: []Card{
				{Value: -15, Suit: "Arcana", Name: "Idiot"},
				{Value: 2, Suit: SuitSabers},
				{Value: 3, Suit: SuitCoins},
				{Value: 4, Suit: SuitFlasks},
			},
			expected: false,
		},
	}

	for i, test := range tests {
		result := isIdiotsArray(test.hand)
		if result != test.expected {
			t.Errorf("Test %d: Expected %v, got %v", i+1, test.expected, result)
		}
	}
}

// TestIsPureSabacc tests Pure Sabacc detection
func TestIsPureSabacc(t *testing.T) {
	tests := []struct {
		hand     []Card
		expected bool
	}{
		{
			// Valid Pure Sabacc
			hand: []Card{
				{Value: 15, Suit: SuitSabers},
				{Value: 8, Suit: SuitCoins},
			},
			expected: true,
		},
		{
			// Not Pure Sabacc - total 22
			hand: []Card{
				{Value: 15, Suit: SuitSabers},
				{Value: 7, Suit: SuitCoins},
			},
			expected: false,
		},
		{
			// Pure Sabacc with negative card
			hand: []Card{
				{Value: -2, Suit: "Arcana", Name: "Strength"},
				{Value: 15, Suit: SuitSabers},
				{Value: 10, Suit: SuitCoins},
			},
			expected: true,
		},
	}

	for i, test := range tests {
		result := isPureSabacc(test.hand)
		if result != test.expected {
			t.Errorf("Test %d: Expected %v, got %v", i+1, test.expected, result)
		}
	}
}

// TestCheckBombOut tests bomb out detection
func TestCheckBombOut(t *testing.T) {
	tests := []struct {
		hand     []Card
		expected bool
	}{
		{
			// Bomb out - over 23
			hand: []Card{
				{Value: 15, Suit: SuitSabers},
				{Value: 10, Suit: SuitCoins},
			},
			expected: true,
		},
		{
			// Bomb out - under -23
			hand: []Card{
				{Value: -15, Suit: "Arcana", Name: "Idiot"},
				{Value: -10, Suit: "Arcana", Name: "Destruction"},
			},
			expected: true,
		},
		{
			// Bomb out - exactly 0
			hand: []Card{
				{Value: -5, Suit: "Arcana", Name: "Justice"},
				{Value: 5, Suit: SuitSabers},
			},
			expected: true,
		},
		{
			// Not bomb out - exactly 23
			hand: []Card{
				{Value: 15, Suit: SuitSabers},
				{Value: 8, Suit: SuitCoins},
			},
			expected: false,
		},
		{
			// Not bomb out - exactly -23
			hand: []Card{
				{Value: -15, Suit: "Arcana", Name: "Idiot"},
				{Value: -8, Suit: "Arcana", Name: "Balance"},
			},
			expected: false,
		},
		{
			// Not bomb out - normal positive hand
			hand: []Card{
				{Value: 10, Suit: SuitSabers},
				{Value: 5, Suit: SuitCoins},
			},
			expected: false,
		},
	}

	for i, test := range tests {
		result := checkBombOut(test.hand)
		if result != test.expected {
			t.Errorf("Test %d: Expected %v, got %v", i+1, test.expected, result)
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

	// Conservative should fold at a LOWER threshold (fold earlier/sooner)
	// Conservative: 22, Aggressive: 28 - so 22 < 28 is correct
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

// TestGameLogic tests basic game logic flows
func TestGameLogic(t *testing.T) {
	// Test hand comparison logic
	hand1 := []Card{{Value: 20, Suit: SuitSabers}, {Value: 2, Suit: SuitCoins}} // 22
	hand2 := []Card{{Value: 15, Suit: SuitSabers}, {Value: 7, Suit: SuitCoins}} // 22
	hand3 := []Card{{Value: 15, Suit: SuitSabers}, {Value: 8, Suit: SuitCoins}} // 23 (Pure Sabacc)

	total1 := calculateHandTotal(hand1)
	total2 := calculateHandTotal(hand2)
	total3 := calculateHandTotal(hand3)

	if total1 != total2 {
		t.Error("Hands with same total should be equal")
	}

	if !isPureSabacc(hand3) {
		t.Error("Hand totaling 23 should be Pure Sabacc")
	}

	if total3 != 23 {
		t.Errorf("Expected 23, got %d", total3)
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
