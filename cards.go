package main

import (
	"fmt"
	"math/rand"
	"time"

	gd "github.com/robbiew/godoors"
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

// String returns a string representation of the card
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

// NewDeck creates a new 76-card Sabacc deck
func NewDeck() Deck {
	var cards []Card

	// Add numbered cards (1-15) for each suit (positive values)
	suits := []string{SuitSabers, SuitFlasks, SuitCoins, SuitStaves}
	for _, suit := range suits {
		for value := 1; value <= 15; value++ {
			cards = append(cards, Card{Value: value, Suit: suit, Name: fmt.Sprintf("%d", value)})
		}
	}

	// Add two copies of each Arcana card
	for _, arcana := range ArcanaCards {
		cards = append(cards, arcana)
		cards = append(cards, arcana) // Second copy
	}

	return Deck{Cards: cards}
}

// Shuffle randomizes the deck
func (d *Deck) Shuffle() {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(d.Cards), func(i, j int) {
		d.Cards[i], d.Cards[j] = d.Cards[j], d.Cards[i]
	})
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

func getCardColor(card Card) string {
	switch card.Suit {
	case SuitSabers:
		return gd.Blue
	case SuitFlasks:
		return gd.Green
	case SuitCoins:
		return gd.Yellow
	case SuitStaves:
		return gd.Red
	case "Arcana":
		return gd.Magenta
	default:
		return gd.White
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

func isPureSabacc(hand []Card) bool {
	return calculateHandTotal(hand) == 23
}

func handleTradeCard() {
	player := &game.Players[0]

	if len(player.Hand) < 2 {
		fmt.Printf("\n%sYou need at least 2 cards to trade!%s\n", gd.Red, gd.Reset)
		time.Sleep(1 * time.Second)
		return
	}

	gd.ClearScreen()
	fmt.Print(gd.CyanHi + "═══════════════════════════════════════════\n" + gd.Reset)
	fmt.Print(gd.CyanHi + "              TRADE A CARD\n" + gd.Reset)
	fmt.Print(gd.CyanHi + "═══════════════════════════════════════════\n\n" + gd.Reset)

	fmt.Printf("%sSelect a card to trade:%s\n\n", gd.Yellow, gd.Reset)

	for i, card := range player.Hand {
		fmt.Printf("%s[%d]%s %s[%s]%s\n",
			gd.Green, i+1, gd.Reset,
			getCardColor(card), card.String(), gd.Reset)
	}

	fmt.Printf("\n%s[0] Cancel%s\n\n", gd.Red, gd.Reset)
	fmt.Print(gd.Green + "Choice: " + gd.Reset)

	char, _, err := getKeyWithTimeout()
	if err != nil {
		return
	}

	choice := int(char - '0')
	if choice == 0 {
		return
	}

	if choice < 1 || choice > len(player.Hand) {
		fmt.Printf("\n%sInvalid choice!%s\n", gd.Red, gd.Reset)
		time.Sleep(1 * time.Second)
		return
	}

	// Remove the selected card
	tradedCard := player.Hand[choice-1]
	player.Hand = append(player.Hand[:choice-1], player.Hand[choice:]...)

	// Draw a new card
	if len(game.Deck.Cards) > 0 {
		newCard := game.Deck.Deal()
		player.Hand = append(player.Hand, newCard)

		fmt.Printf("\n%sYou traded %s[%s]%s for %s[%s]%s\n",
			gd.Green, getCardColor(tradedCard), tradedCard.String(), gd.Reset,
			getCardColor(newCard), newCard.String(), gd.Reset)
	} else {
		fmt.Printf("\n%sNo more cards in deck!%s\n", gd.Red, gd.Reset)
	}

	time.Sleep(2 * time.Second)
}

func handleStaticField() {
	player := &game.Players[0]

	gd.ClearScreen()
	fmt.Print(gd.CyanHi + "═══════════════════════════════════════════\n" + gd.Reset)
	fmt.Print(gd.CyanHi + "             STATIC FIELD\n" + gd.Reset)
	fmt.Print(gd.CyanHi + "═══════════════════════════════════════════\n\n" + gd.Reset)

	fmt.Printf("%s[1]%s Place card in Static Field\n", gd.Green, gd.Reset)
	fmt.Printf("%s[2]%s Remove card from Static Field\n", gd.Green, gd.Reset)
	fmt.Printf("%s[0]%s Cancel\n\n", gd.Red, gd.Reset)
	fmt.Print(gd.Green + "Choice: " + gd.Reset)

	char, _, err := getKeyWithTimeout()
	if err != nil {
		return
	}

	switch char {
	case '1':
		placeInStaticField()
	case '2':
		removeFromStaticField()
	case '0':
		return
	}
}

func placeInStaticField() {
	player := &game.Players[0]

	if len(player.Hand) == 0 {
		fmt.Printf("\n%sNo cards in hand!%s\n", gd.Red, gd.Reset)
		time.Sleep(1 * time.Second)
		return
	}

	fmt.Printf("\n%sSelect a card to place in Static Field:%s\n\n", gd.Yellow, gd.Reset)

	for i, card := range player.Hand {
		fmt.Printf("%s[%d]%s %s[%s]%s\n",
			gd.Green, i+1, gd.Reset,
			getCardColor(card), card.String(), gd.Reset)
	}

	fmt.Printf("\n%s[0] Cancel%s\n\n", gd.Red, gd.Reset)
	fmt.Print(gd.Green + "Choice: " + gd.Reset)

	char, _, err := getKeyWithTimeout()
	if err != nil {
		return
	}

	choice := int(char - '0')
	if choice == 0 {
		return
	}

	if choice < 1 || choice > len(player.Hand) {
		fmt.Printf("\n%sInvalid choice!%s\n", gd.Red, gd.Reset)
		time.Sleep(1 * time.Second)
		return
	}

	// Move card to static field
	card := player.Hand[choice-1]
	player.Hand = append(player.Hand[:choice-1], player.Hand[choice:]...)
	player.StaticField = append(player.StaticField, card)

	fmt.Printf("\n%s%s[%s]%s placed in Static Field (protected from shifts)%s\n",
		gd.Green, getCardColor(card), card.String(), gd.Reset, gd.Reset)
	time.Sleep(2 * time.Second)
}

func removeFromStaticField() {
	player := &game.Players[0]

	if len(player.StaticField) == 0 {
		fmt.Printf("\n%sNo cards in Static Field!%s\n", gd.Red, gd.Reset)
		time.Sleep(1 * time.Second)
		return
	}

	fmt.Printf("\n%sSelect a card to remove from Static Field:%s\n\n", gd.Yellow, gd.Reset)

	for i, card := range player.StaticField {
		fmt.Printf("%s[%d]%s %s[%s]%s\n",
			gd.Green, i+1, gd.Reset,
			getCardColor(card), card.String(), gd.Reset)
	}

	fmt.Printf("\n%s[0] Cancel%s\n\n", gd.Red, gd.Reset)
	fmt.Print(gd.Green + "Choice: " + gd.Reset)

	char, _, err := getKeyWithTimeout()
	if err != nil {
		return
	}

	choice := int(char - '0')
	if choice == 0 {
		return
	}

	if choice < 1 || choice > len(player.StaticField) {
		fmt.Printf("\n%sInvalid choice!%s\n", gd.Red, gd.Reset)
		time.Sleep(1 * time.Second)
		return
	}

	// Move card back to hand
	card := player.StaticField[choice-1]
	player.StaticField = append(player.StaticField[:choice-1], player.StaticField[choice:]...)
	player.Hand = append(player.Hand, card)

	fmt.Printf("\n%s%s[%s]%s returned to hand%s\n",
		gd.Green, getCardColor(card), card.String(), gd.Reset, gd.Reset)
	time.Sleep(2 * time.Second)
}
