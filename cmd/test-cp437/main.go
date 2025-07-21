// cmd/test-cp437/main.go - Test CP437 character output
package main

import (
	"fmt"
	"os"
)

func main() {
	// Create a test file with proper CP437 characters
	testContent := generateCP437Test()

	err := os.WriteFile("cp437_test.ans", []byte(testContent), 0644)
	if err != nil {
		fmt.Printf("Error creating test file: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Created cp437_test.ans")
	fmt.Println("View this file in your ANSI viewer to verify CP437 support")
}

func generateCP437Test() string {
	// CP437 Box Drawing Characters
	var content string

	content += "\x1b[37;1mCP437 CHARACTER TEST\x1b[0m\n"
	content += "===================\n\n"

	// Box drawing characters
	content += "Box Drawing Characters:\n"
	content += "\xda\xc4\xc4\xc4\xc4\xc4\xc4\xc4\xbf  (corners and lines)\n"
	content += "\xb3       \xb3  (vertical)\n"
	content += "\xc0\xc4\xc4\xc4\xc4\xc4\xc4\xc4\xd9  (bottom)\n\n"

	// Card suit symbols
	content += "Card Suit Symbols:\n"
	content += "Spades: \x1b[34m\x06\x1b[0m  (\\x06)\n"
	content += "Hearts: \x1b[31m\x03\x1b[0m  (\\x03)\n"
	content += "Diamonds: \x1b[32m\x04\x1b[0m  (\\x04)\n"
	content += "Clubs: \x1b[33m\x05\x1b[0m  (\\x05)\n"
	content += "Star: \x1b[35m\x0f\x1b[0m  (\\x0f)\n\n"

	// Sample card using CP437 - FIXED SPACING
	content += "Sample Card (5 of Spades):\n"
	content += "\x1b[34m\xda\xc4\xc4\xc4\xc4\xc4\xc4\xc4\xbf\x1b[0m\n"   // ┌───────┐
	content += "\x1b[34m\xb3\x1b[37;1m 5     \x1b[34m\xb3\x1b[0m\n"      // │ 5     │
	content += "\x1b[34m\xb3       \xb3\x1b[0m\n"                        // │       │
	content += "\x1b[34m\xb3   \x06   \xb3\x1b[0m\n"                     // │   ♠   │
	content += "\x1b[34m\xb3       \xb3\x1b[0m\n"                        // │       │
	content += "\x1b[34m\xb3     \x1b[37;1m 5\x1b[34m\xb3\x1b[0m\n"      // │     5│
	content += "\x1b[34m\xc0\xc4\xc4\xc4\xc4\xc4\xc4\xc4\xd9\x1b[0m\n\n" // └───────┘

	// Better sample - 10 of Hearts
	content += "Sample Card (10 of Hearts):\n"
	content += "\x1b[31m\xda\xc4\xc4\xc4\xc4\xc4\xc4\xc4\xbf\x1b[0m\n"   // ┌───────┐
	content += "\x1b[31m\xb3\x1b[37;1m10     \x1b[31m\xb3\x1b[0m\n"      // │10     │
	content += "\x1b[31m\xb3       \xb3\x1b[0m\n"                        // │       │
	content += "\x1b[31m\xb3   \x03   \xb3\x1b[0m\n"                     // │   ♥   │
	content += "\x1b[31m\xb3       \xb3\x1b[0m\n"                        // │       │
	content += "\x1b[31m\xb3     \x1b[37;1m10\x1b[31m\xb3\x1b[0m\n"      // │     10│
	content += "\x1b[31m\xc0\xc4\xc4\xc4\xc4\xc4\xc4\xc4\xd9\x1b[0m\n\n" // └───────┘

	// Arcana card sample
	content += "Sample Arcana Card (Idiot):\n"
	content += "\x1b[35m\xda\xc4\xc4\xc4\xc4\xc4\xc4\xc4\xbf\x1b[0m\n"   // ┌───────┐
	content += "\x1b[35m\xb3\x1b[37;1mId     \x1b[35m\xb3\x1b[0m\n"      // │Id     │
	content += "\x1b[35m\xb3       \xb3\x1b[0m\n"                        // │       │
	content += "\x1b[35m\xb3   \x0f   \xb3\x1b[0m\n"                     // │   ☼   │
	content += "\x1b[35m\xb3       \xb3\x1b[0m\n"                        // │       │
	content += "\x1b[35m\xb3     \x1b[37;1mId\x1b[35m\xb3\x1b[0m\n"      // │     Id│
	content += "\x1b[35m\xc0\xc4\xc4\xc4\xc4\xc4\xc4\xc4\xd9\x1b[0m\n\n" // └───────┘

	// Pattern test
	content += "Pattern Test:\n"
	content += "\x1b[31m\xb0\xb0\xb0\xb0\xb0\x1b[0m  (light shade)\n"
	content += "\x1b[31m\xb1\xb1\xb1\xb1\xb1\x1b[0m  (medium shade)\n"
	content += "\x1b[31m\xb2\xb2\xb2\xb2\xb2\x1b[0m  (dark shade)\n"
	content += "\x1b[31m\xdb\xdb\xdb\xdb\xdb\x1b[0m  (full block)\n\n"

	// Color test
	content += "Color Test:\n"
	for i := 30; i <= 37; i++ {
		content += fmt.Sprintf("\x1b[%dm Color %d \x1b[0m", i, i-29)
	}
	content += "\n"
	for i := 30; i <= 37; i++ {
		content += fmt.Sprintf("\x1b[%d;1m Bright %d \x1b[0m", i, i-29)
	}
	content += "\n\n"

	content += "If you can see the cards and symbols properly,\n"
	content += "your ANSI viewer supports CP437!\n"

	return content
}

// CP437 Character Reference
const (
	// Box Drawing
	BOX_TL = "\xda" // ┌ Top Left
	BOX_TR = "\xbf" // ┐ Top Right
	BOX_BL = "\xc0" // └ Bottom Left
	BOX_BR = "\xd9" // ┘ Bottom Right
	BOX_H  = "\xc4" // ─ Horizontal
	BOX_V  = "\xb3" // │ Vertical

	// Card Suits
	HEART   = "\x03" // ♥
	DIAMOND = "\x04" // ♦
	CLUB    = "\x05" // ♣
	SPADE   = "\x06" // ♠

	// Patterns
	SHADE_LIGHT  = "\xb0" // ░
	SHADE_MEDIUM = "\xb1" // ▒
	SHADE_DARK   = "\xb2" // ▓
	BLOCK_FULL   = "\xdb" // █

	// Other useful
	STAR = "\x0f" // ☼
	NOTE = "\x0e" // ♫
)
