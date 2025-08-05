package main

import (
	"fmt"
	"regexp"
	"strings"
)

// ScreenLayout manages the persistent UI layout for the Sabacc game
type ScreenLayout struct {
	// Terminal dimensions
	TerminalW int
	TerminalH int

	// Fixed screen positions for 4 AI players + human player layout
	// AI Player positions (corners of screen)
	AIPlayer1X int // Top-left AI player (Phoo_ja)
	AIPlayer1Y int
	AIPlayer2X int // Top-right AI player (Rsh-Taac)
	AIPlayer2Y int
	AIPlayer3X int // Bottom-left AI player (Soladi)
	AIPlayer3Y int
	AIPlayer4X int // Bottom-right AI player (Ky'Ola)
	AIPlayer4Y int

	// Central game log area
	GameLogX int // Center area for scrolling messages
	GameLogY int
	GameLogW int // Width of game log area
	GameLogH int // Height of game log area

	// Human player area (bottom center)
	HumanPlayerX int // Player cards area
	HumanPlayerY int
	HumanNameY   int // Player name/info line

	// Turn indicator area
	TurnIndicatorX int
	TurnIndicatorY int

	// Bottom status and menu
	StatusY int // Bottom status line (pots, credits)
	MenuY   int // Action buttons line

	// Card display settings
	CardWidth   int // Width of each card display
	CardHeight  int // Height of each card display
	MaxCardsRow int // Maximum cards per row

	// Face art dimensions
	FaceWidth  int // Width of character face art
	FaceHeight int // Height of character face art

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

// Add a helper to strip ANSI escape codes
var ansiRegexp = regexp.MustCompile(`\x1b\[[0-9;]*m`)

func stripANSI(s string) string {
	return ansiRegexp.ReplaceAllString(s, "")
}

// NewScreenLayout creates a new screen layout for the given terminal size
func NewScreenLayout(termW, termH int) *ScreenLayout {
	layout := &ScreenLayout{
		TerminalW: termW,
		TerminalH: termH,

		// AI Player positions (4 corners) - adjusted for 24-row terminal
		AIPlayer1X: 2, // Top-left (Phoo_ja)
		AIPlayer1Y: 1,
		AIPlayer2X: termW - 15, // Top-right (Rsh-Taac)
		AIPlayer2Y: 1,
		AIPlayer3X: 2, // Bottom-left (Soladi)
		AIPlayer3Y: 9,
		AIPlayer4X: termW - 15, // Bottom-right (Ky'Ola)
		AIPlayer4Y: 9,

		// Central game log area (blue bordered box from reference)
		GameLogX: 18,
		GameLogY: 6,
		GameLogW: 42, // Fixed width to match reference file
		GameLogH: 6,  // Reduced height for 24-row terminal

		// Human player area (bottom center - adjusted for 24 rows)
		HumanPlayerX: termW/2 - 15,
		HumanPlayerY: 16,
		HumanNameY:   15,

		// Turn indicator (above game log)
		TurnIndicatorX: termW/2 - 8,
		TurnIndicatorY: 4,

		// Bottom UI elements (properly positioned for 24 rows)
		StatusY: 24, // Status line at row 23 (pots, credits)
		MenuY:   20, // Action buttons at row 20-22

		// Card and face display settings
		CardWidth:   6,  // Diamond card width (matches card database)
		CardHeight:  5,  // Diamond card height (matches new diamond cards)
		MaxCardsRow: 8,  // More cards fit with smaller diamond cards
		FaceWidth:   10, // Character face art width
		FaceHeight:  6,  // Character face art height
	}

	// Initialize game log with proper dimensions
	layout.GameLog = &GameLog{
		Messages: make([]string, 0),
		MaxLines: layout.GameLogH - 2, // Account for borders
		StartY:   layout.GameLogY + 2, // Inside border
		Width:    layout.GameLogW - 2, // Inside border with padding
	}

	return layout
}

// InitializeScreen draws the static UI elements once
func (sl *ScreenLayout) InitializeScreen() {
	ClearScreen()

	// Draw the 4 AI player areas with face art
	sl.drawAIPlayerAreas()

	// Draw central game log with blue border
	sl.drawGameLogBorder()

	// Draw human player area
	sl.drawHumanPlayerArea()

	// Initialize game log
	sl.GameLog.Clear()
}

// Helper: Complete CP437 character set (IBM PC character set)
var (
	// Control characters (0x00-0x1F)
	cp437BEL = "\x07" // Bell
	cp437BS  = "\x08" // Backspace

	// Extended characters (0x7F-0xFF) - CP437 specific
	cp437DEL  = "\x7f" // █ (solid block)
	cp437C128 = "\x80" // Ç (C with cedilla)
	cp437C129 = "\x81" // ü (u with umlaut)
	cp437C130 = "\x82" // é (e with acute)
	cp437C131 = "\x83" // â (a with circumflex)
	cp437C132 = "\x84" // ä (a with umlaut)
	cp437C133 = "\x85" // à (a with grave)
	cp437C134 = "\x86" // å (a with ring)
	cp437C135 = "\x87" // ç (c with cedilla)
	cp437C136 = "\x88" // ê (e with circumflex)
	cp437C137 = "\x89" // ë (e with umlaut)
	cp437C138 = "\x8a" // è (e with grave)
	cp437C139 = "\x8b" // ï (i with umlaut)
	cp437C140 = "\x8c" // î (i with circumflex)
	cp437C141 = "\x8d" // ì (i with grave)
	cp437C142 = "\x8e" // Ä (A with umlaut)
	cp437C143 = "\x8f" // Å (A with ring)
	cp437C144 = "\x90" // É (E with acute)
	cp437C145 = "\x91" // æ (ae ligature)
	cp437C146 = "\x92" // Æ (AE ligature)
	cp437C147 = "\x93" // ô (o with circumflex)
	cp437C148 = "\x94" // ö (o with umlaut)
	cp437C149 = "\x95" // ò (o with grave)
	cp437C150 = "\x96" // û (u with circumflex)
	cp437C151 = "\x97" // ù (u with grave)
	cp437C152 = "\x98" // ÿ (y with umlaut)
	cp437C153 = "\x99" // Ö (O with umlaut)
	cp437C154 = "\x9a" // Ü (U with umlaut)
	cp437C155 = "\x9b" // ¢ (cent sign)
	cp437C156 = "\x9c" // £ (pound sign)
	cp437C157 = "\x9d" // ¥ (yen sign)
	cp437C158 = "\x9e" // ₧ (peseta sign)
	cp437C159 = "\x9f" // ƒ (function sign)
	cp437C160 = "\xa0" // á (a with acute)
	cp437C161 = "\xa1" // í (i with acute)
	cp437C162 = "\xa2" // ó (o with acute)
	cp437C163 = "\xa3" // ú (u with acute)
	cp437C164 = "\xa4" // ñ (n with tilde)
	cp437C165 = "\xa5" // Ñ (N with tilde)
	cp437C166 = "\xa6" // ª (feminine ordinal)
	cp437C167 = "\xa7" // º (masculine ordinal)
	cp437C168 = "\xa8" // ¿ (inverted question mark)
	cp437C169 = "\xa9" // ⌐ (reversed not sign)
	cp437C170 = "\xaa" // ¬ (not sign)
	cp437C171 = "\xab" // ½ (one half)
	cp437C172 = "\xac" // ¼ (one quarter)
	cp437C173 = "\xad" // ¡ (inverted exclamation mark)
	cp437C174 = "\xae" // « (left double angle quote)
	cp437C175 = "\xaf" // » (right double angle quote)
	cp437C176 = "\xb0" // ░ (light shade)
	cp437C177 = "\xb1" // ▒ (medium shade)
	cp437C178 = "\xb2" // ▓ (dark shade)
	cp437C179 = "\xb3" // │ (vertical line)
	cp437C180 = "\xb4" // ┤ (right T)
	cp437C181 = "\xb5" // ╡ (right double T)
	cp437C182 = "\xb6" // ╢ (double vertical)
	cp437C183 = "\xb7" // ╖ (double top-right)
	cp437C184 = "\xb8" // ╕ (double bottom-right)
	cp437C185 = "\xb9" // ╣ (double right T)
	cp437C186 = "\xba" // ║ (double vertical)
	cp437C187 = "\xbb" // ╗ (double top-right)
	cp437C188 = "\xbc" // ╝ (double bottom-right)
	cp437C189 = "\xbd" // ╜ (double bottom-left)
	cp437C190 = "\xbe" // ╛ (double top-left)
	cp437C191 = "\xbf" // ┐ (top-right corner)
	cp437C192 = "\xc0" // └ (bottom-left corner)
	cp437C193 = "\xc1" // ┴ (bottom T)
	cp437C194 = "\xc2" // ┬ (top T)
	cp437C195 = "\xc3" // ├ (left T)
	cp437C196 = "\xc4" // ─ (horizontal line)
	cp437C197 = "\xc5" // ┼ (cross)
	cp437C198 = "\xc6" // ╞ (left double T)
	cp437C199 = "\xc7" // ╟ (double left T)
	cp437C200 = "\xc8" // ╚ (double bottom-left)
	cp437C201 = "\xc9" // ╔ (double top-left)
	cp437C202 = "\xca" // ╩ (double bottom T)
	cp437C203 = "\xcb" // ╦ (double top T)
	cp437C204 = "\xcc" // ╠ (double left T)
	cp437C205 = "\xcd" // ═ (double horizontal)
	cp437C206 = "\xce" // ╬ (double cross)
	cp437C207 = "\xcf" // ╧ (double bottom T)
	cp437C208 = "\xd0" // ╨ (double top T)
	cp437C209 = "\xd1" // ╤ (double top T)
	cp437C210 = "\xd2" // ╥ (double top T)
	cp437C211 = "\xd3" // ╙ (double bottom-left)
	cp437C212 = "\xd4" // ╘ (double bottom-right)
	cp437C213 = "\xd5" // ╒ (double top-left)
	cp437C214 = "\xd6" // ╓ (double top-left)
	cp437C215 = "\xd7" // ╫ (double vertical)
	cp437C216 = "\xd8" // ╪ (double vertical)
	cp437C217 = "\xd9" // ┘ (bottom-right corner)
	cp437C218 = "\xda" // ┌ (top-left corner)
	cp437C219 = "\xdb" // █ (solid block)
	cp437C220 = "\xdc" // ▄ (bottom half block)
	cp437C221 = "\xdd" // ▌ (left half block)
	cp437C222 = "\xde" // ▐ (right half block)
	cp437C223 = "\xdf" // ▀ (top half block)
	cp437C224 = "\xe0" // α (alpha)
	cp437C225 = "\xe1" // ß (beta)
	cp437C226 = "\xe2" // Γ (gamma)
	cp437C227 = "\xe3" // π (pi)
	cp437C228 = "\xe4" // Σ (sigma)
	cp437C229 = "\xe5" // σ (sigma)
	cp437C230 = "\xe6" // µ (mu)
	cp437C231 = "\xe7" // τ (tau)
	cp437C232 = "\xe8" // Φ (phi)
	cp437C233 = "\xe9" // Θ (theta)
	cp437C234 = "\xea" // Ω (omega)
	cp437C235 = "\xeb" // δ (delta)
	cp437C236 = "\xec" // ∞ (infinity)
	cp437C237 = "\xed" // φ (phi)
	cp437C238 = "\xee" // ε (epsilon)
	cp437C239 = "\xef" // ∩ (intersection)
	cp437C240 = "\xf0" // ≡ (equivalence)
	cp437C241 = "\xf1" // ± (plus-minus)
	cp437C242 = "\xf2" // ≥ (greater than or equal)
	cp437C243 = "\xf3" // ≤ (less than or equal)
	cp437C244 = "\xf4" // ⌠ (top half integral)
	cp437C245 = "\xf5" // ⌡ (bottom half integral)
	cp437C246 = "\xf6" // ÷ (division sign)
	cp437C247 = "\xf7" // ≈ (approximately equal)
	cp437C248 = "\xf8" // ° (degree)
	cp437C249 = "\xf9" // ∙ (bullet)
	cp437C250 = "\xfa" // · (middle dot)
	cp437C251 = "\xfb" // √ (square root)
	cp437C252 = "\xfc" // ⁿ (superscript n)
	cp437C253 = "\xfd" // ² (superscript 2)
	cp437C254 = "\xfe" // ■ (black square)
	cp437C255 = "\xff" // (no break space)

	// Box drawing shortcuts (descriptive names)
	HORIZONTAL_LINE            = cp437C196 // ─ (horizontal line)
	VERTICAL_LINE              = cp437C179 // │ (vertical line)
	TOP_LEFT_CORNER            = cp437C218 // ┌ (top-left corner)
	TOP_RIGHT_CORNER           = cp437C191 // ┐ (top-right corner)
	BOTTOM_LEFT_CORNER         = cp437C192 // └ (bottom-left corner)
	BOTTOM_RIGHT_CORNER        = cp437C217 // ┘ (bottom-right corner)
	LEFT_T_JUNCTION            = cp437C195 // ├ (left T)
	RIGHT_T_JUNCTION           = cp437C180 // ┤ (right T)
	DOUBLE_HORIZONTAL_LINE     = cp437C205 // ═ (double horizontal)
	DOUBLE_TOP_LEFT_CORNER     = cp437C201 // ╔ (double top-left)
	DOUBLE_TOP_RIGHT_CORNER    = cp437C187 // ╗ (double top-right)
	DOUBLE_BOTTOM_LEFT_CORNER  = cp437C200 // ╚ (double bottom-left)
	DOUBLE_BOTTOM_RIGHT_CORNER = cp437C188 // ╝ (double bottom-right)
	DOUBLE_VERTICAL_LINE       = cp437C186 // ║ (double vertical)
	SOLID_BLOCK                = cp437C219 // █ (solid block)
	LIGHT_SHADE                = cp437C176 // ░ (light shade)
	MEDIUM_SHADE               = cp437C177 // ▒ (medium shade)
	DARK_SHADE                 = cp437C178 // ▓ (dark shade)

)

// drawAIPlayerAreas draws the 4 AI player areas with face art and card blocks
func (sl *ScreenLayout) drawAIPlayerAreas() {
	// AI player names and their positions
	aiPlayers := []struct {
		name string
		x, y int
	}{
		{"Phoo_ja", sl.AIPlayer1X, sl.AIPlayer1Y},  // Top-left
		{"Rsh-Taac", sl.AIPlayer2X, sl.AIPlayer2Y}, // Top-right
		{"Soladi", sl.AIPlayer3X, sl.AIPlayer3Y},   // Bottom-left
		{"Ky'Ola", sl.AIPlayer4X, sl.AIPlayer4Y},   // Bottom-right
	}

	for i, player := range aiPlayers {
		// Draw player name
		MoveCursor(player.x, player.y)
		fmt.Printf("%s%s%s", CyanHi, player.name, Reset)

		// Draw simple face art placeholder (will be enhanced with actual art)
		sl.drawSimpleFaceArt(player.x, player.y+1, i)

		// Draw CP437 card blocks (placeholder for now)
		sl.drawAICardBlocks(player.x, player.y+sl.FaceHeight+1, 5) // 5 cards max
	}
}

// drawHumanPlayerArea draws the human player area at bottom center
func (sl *ScreenLayout) drawHumanPlayerArea() {
	// Draw player name area
	MoveCursor(sl.HumanPlayerX, sl.HumanNameY)
	fmt.Printf("%s[Player Area]%s", GreenHi, Reset)
}

// drawSimpleFaceArt draws character faces using exact CP437 from reference file
func (sl *ScreenLayout) drawSimpleFaceArt(x, y, playerIndex int) {
	// Exact CP437 face art from reference .ans file
	faces := [][]string{
		{ // Phoo_ja (robot/droid) - top-left
			"\xda" + "pppppppp" + "\xbf",
			"\xb3" + " \x09\xdb\xdb\xdb\x09 " + "\xb3",
			"\xb3" + " \x09\xdb\xdb\xdb\x09 " + "\xb3",
			"\xb3" + "  \xc4--\xc4  " + "\xb3",
			"\xb3" + "   \xc4\xc4   " + "\xb3",
			"\xc0" + "\xcd\xcd\xcd\xcd\xcd\xcd\xcd\xcd" + "\xd9",
		},
		{ // Rsh-Taac (alien) - top-right
			"\xda" + "pppppppp" + "\xbf",
			"\xb3" + "  \x07\xdb\xdb\x07  " + "\xb3",
			"\xb3" + " o\xdb\xdb\xdb\xdbo " + "\xb3",
			"\xb3" + "   \xdb\xdb   " + "\xb3",
			"\xb3" + "   \xdb\xdb   " + "\xb3",
			"\xc0" + "\xcd\xcd\xcd\xcd\xcd\xcd\xcd\xcd" + "\xd9",
		},
		{ // Soladi (human) - bottom-left
			"\xda" + "pppppppp" + "\xbf",
			"\xb3" + " \x09    \x09 " + "\xb3",
			"\xb3" + "  \x07\x07\x07\x07  " + "\xb3",
			"\xb3" + "  \x07\x07\x07\x07  " + "\xb3",
			"\xb3" + "   \xdb\xdb   " + "\xb3",
			"\xc0" + "\xcd\xcd\xcd\xcd\xcd\xcd\xcd\xcd" + "\xd9",
		},
		{ // Ky'Ola (alien) - bottom-right
			"\xda" + "pppppppp" + "\xbf",
			"\xb3" + "      " + "\xb3",
			"\xb3" + " \xcd\xcd    " + "\xb3",
			"\xb3" + " \xcd\xcd    " + "\xb3",
			"\xb3" + "   \xcd\xcd   " + "\xb3",
			"\xc0" + "\xcd\xcd\xcd\xcd\xcd\xcd\xcd\xcd" + "\xd9",
		},
	}

	if playerIndex < len(faces) {
		for i, line := range faces[playerIndex] {
			MoveCursor(x, y+i)
			fmt.Printf("%s%s%s", White, line, Reset)
		}
	}
}

// drawAICardBlocks draws simple face-down card blocks for AI players
func (sl *ScreenLayout) drawAICardBlocks(x, y, cardCount int) {
	// Simple face-down card blocks (from reference file pattern)
	for i := 0; i < cardCount && i < 6; i++ { // Max 6 cards
		cardX := x + (i * 2) // Tight spacing for small blocks

		// Small face-down card block - just a simple colored block
		MoveCursor(cardX, y)
		fmt.Printf("%s\xdb\xdb%s", Yellow, Reset) // ██ (solid blocks)
	}
}

// drawGameLogBorder draws the border around the game log area (matching reference file)
func (sl *ScreenLayout) drawGameLogBorder() {
	// Use exact positioning from reference file
	logX := sl.GameLogX
	logY := sl.GameLogY
	logW := sl.GameLogW

	// Create horizontal line
	borderLine := strings.Repeat(DOUBLE_HORIZONTAL_LINE, logW-2)

	// Top border with proper blue background like reference
	MoveCursor(logX, logY)
	fmt.Printf("\x1b[30;44m%s%s%s\x1b[0m", DOUBLE_TOP_LEFT_CORNER, borderLine, DOUBLE_TOP_RIGHT_CORNER)

	// Side borders for each line with blue background
	for i := 1; i <= sl.GameLogH; i++ {
		MoveCursor(logX, logY+i)
		fmt.Printf("\x1b[30;44m%s\x1b[0m", DOUBLE_VERTICAL_LINE)
		MoveCursor(logX+logW-1, logY+i)
		fmt.Printf("\x1b[30;44m%s\x1b[0m", DOUBLE_VERTICAL_LINE)
	}

	// Bottom border with blue background
	MoveCursor(logX, logY+sl.GameLogH+1)
	fmt.Printf("\x1b[30;44m%s%s%s\x1b[0m", DOUBLE_BOTTOM_LEFT_CORNER, borderLine, DOUBLE_BOTTOM_RIGHT_CORNER)
}

// UpdateHeader updates the turn indicator and game status
func (sl *ScreenLayout) UpdateHeader(round, handPot, sabaccPot, deckSize int, currentPlayer string) {
	// Update turn indicator above game log using proper CP437 arrows
	MoveCursor(sl.TurnIndicatorX, sl.TurnIndicatorY)
	fmt.Print(EraseLine)
	fmt.Printf("%s\x10 \x10 %s \x11 \x11%s", YellowHi, currentPlayer, Reset)

	// Update bottom status line with pot information
	sl.UpdateStatusLine(round, handPot, sabaccPot, deckSize)
}

// UpdateStatusLine updates the bottom status line with game info (CP437 style)
func (sl *ScreenLayout) UpdateStatusLine(round, handPot, sabaccPot, deckSize int) {
	MoveCursor(1, sl.StatusY)
	fmt.Print(EraseLine)

	// Create status line with cyan background like in reference file
	statusLeft := fmt.Sprintf("\x1b[36;46m     \x1b[0;30;46m \x1b[1;37mGame Pot: \x1b[33m%d\x1b[36m \x1b[37m  Sabaac Pot: \x1b[33m%d\x1b[37m \x1b[36m \x1b[37m  Side Pot: \x1b[33m0\x1b[37m \x1b[36m \x1b[37m  Credits:\x1b[33m 22\x1b[37m \x1b[0;30;46m        \x1b[36;40m", handPot, sabaccPot)

	fmt.Print(statusLeft)

	// Fill rest of line to terminal width
	statusLen := len(stripANSI(statusLeft))
	if statusLen < sl.TerminalW {
		remaining := sl.TerminalW - statusLen
		fmt.Printf("\x1b[36;40m%s\x1b[37m", strings.Repeat(" ", remaining))
	}

	fmt.Print(Reset)
}

// UpdatePlayerInfo updates a player's name, credits, and total for new layout
func (sl *ScreenLayout) UpdatePlayerInfo(playerIndex int, name string, credits, total int, showTotal bool) {
	switch playerIndex {
	case 0: // Human player (bottom center)
		MoveCursor(sl.HumanPlayerX, sl.HumanNameY)
		fmt.Print(EraseLine)
		if showTotal {
			fmt.Printf("%s%s%s - Total: %s Credits: %d",
				GreenHi, name, Reset, displayHandValue(total), credits)
		} else {
			fmt.Printf("%s%s%s - Credits: %d", GreenHi, name, Reset, credits)
		}
	case 1, 2, 3, 4: // AI players (corners) - don't overwrite, just update credits below face art
		var x, y int
		switch playerIndex {
		case 1:
			x, y = sl.AIPlayer1X, sl.AIPlayer1Y+sl.FaceHeight+2 // Position below cards
		case 2:
			x, y = sl.AIPlayer2X, sl.AIPlayer2Y+sl.FaceHeight+2
		case 3:
			x, y = sl.AIPlayer3X, sl.AIPlayer3Y+sl.FaceHeight+2
		case 4:
			x, y = sl.AIPlayer4X, sl.AIPlayer4Y+sl.FaceHeight+2
		}
		// Clear the line first, then write credits
		MoveCursor(x, y)
		fmt.Print(strings.Repeat(" ", 10)) // Clear previous text
		MoveCursor(x, y)
		fmt.Printf("%s$%d%s", Yellow, credits, Reset)
	}
}

// ClearPlayerArea clears a player's card display area
func (sl *ScreenLayout) ClearPlayerArea(playerIndex int) {
	switch playerIndex {
	case 0: // Human player (bottom center)
		for i := 0; i < sl.CardHeight+1; i++ {
			MoveCursor(sl.HumanPlayerX, sl.HumanPlayerY+i)
			fmt.Print(EraseLine)
		}
	case 1, 2, 3, 4: // AI players (corners)
		var x, y int
		switch playerIndex {
		case 1:
			x, y = sl.AIPlayer1X, sl.AIPlayer1Y+sl.FaceHeight+1
		case 2:
			x, y = sl.AIPlayer2X, sl.AIPlayer2Y+sl.FaceHeight+1
		case 3:
			x, y = sl.AIPlayer3X, sl.AIPlayer3Y+sl.FaceHeight+1
		case 4:
			x, y = sl.AIPlayer4X, sl.AIPlayer4Y+sl.FaceHeight+1
		}
		// Clear AI card area
		for i := 0; i < sl.CardHeight+1; i++ {
			MoveCursor(x, y+i)
			fmt.Print(strings.Repeat(" ", sl.FaceWidth))
		}
	}
}

// ClearMenuArea clears the menu/command area
func (sl *ScreenLayout) ClearMenuArea() {
	for i := 0; i < 3; i++ { // Menu area is 3 lines
		MoveCursor(1, sl.MenuY+i)
		fmt.Print(EraseLine)
		fmt.Print(VERTICAL_LINE + strings.Repeat(" ", sl.TerminalW-2) + VERTICAL_LINE)
	}
}

// ShowMenu displays the command menu
func (sl *ScreenLayout) ShowMenu(title string, options []MenuOption, prompt string) {
	sl.ClearMenuArea()

	// Show title
	MoveCursor(2, sl.MenuY)
	fmt.Printf("%s%s%s", GreenHi, title, Reset)

	// Show options
	MoveCursor(2, sl.MenuY+1)
	optionTexts := make([]string, 0)
	for _, opt := range options {
		if opt.Enabled {
			text := fmt.Sprintf("[%s%c%s] %s", YellowHi, opt.Key, Yellow, opt.Description)
			optionTexts = append(optionTexts, text)
		}
	}
	fmt.Printf("%s%s%s", White, strings.Join(optionTexts, "  "), Reset)

	// Show prompt
	MoveCursor(2, sl.MenuY+2)
	fmt.Printf("%s%s%s", Green, prompt, Reset)
}

// ShowPlayerTurnMenu displays the player turn menu
func (sl *ScreenLayout) ShowPlayerTurnMenu(round int) {
	options := []MenuOption{
		{'D', "Draw card", true},
		{'T', "Trade card", true},
		{'S', "Stand", true},
		{'F', "Static Field", true},
		{'C', "Call", round >= 2},
		{'Q', "Fold", true},
	}
	sl.ShowMenu("Your Turn", options, "Choice: ")
}

// LogMessage adds a message to the game log
func (sl *ScreenLayout) LogMessage(message, msgType string) {
	sl.GameLog.AddMessage(message, msgType)
	sl.RefreshGameLog()
}

// DisplayMessage shows a temporary message
func (sl *ScreenLayout) DisplayMessage(message, msgType string, duration int) {
	// Simple implementation - just add to log for now
	sl.LogMessage(message, msgType)
}

// RefreshGameLog redraws the game log area
func (sl *ScreenLayout) RefreshGameLog() {
	// Clear the log area (inside borders)
	for i := 0; i < sl.GameLog.MaxLines; i++ {
		MoveCursor(sl.GameLogX+1, sl.GameLogY+1+i)
		fmt.Print(strings.Repeat(" ", sl.GameLog.Width))
	}

	// Display recent messages
	messages := sl.GameLog.GetRecentMessages()
	for i, msg := range messages {
		if i < sl.GameLog.MaxLines {
			MoveCursor(sl.GameLogX+2, sl.GameLogY+1+i)
			// Truncate message if too long
			if len(msg) > sl.GameLog.Width-2 {
				msg = msg[:sl.GameLog.Width-5] + "..."
			}
			fmt.Print(msg)
		}
	}
}

// RenderPlayerCards renders cards for a player
func (sl *ScreenLayout) RenderPlayerCards(playerIndex int, cards []Card, faceDown bool, renderer *CardRenderer) {
	var x, y int

	switch playerIndex {
	case 0: // Human player
		x, y = sl.HumanPlayerX, sl.HumanPlayerY
	case 1: // AI 1
		x, y = sl.AIPlayer1X, sl.AIPlayer1Y+sl.FaceHeight+1
	case 2: // AI 2
		x, y = sl.AIPlayer2X, sl.AIPlayer2Y+sl.FaceHeight+1
	case 3: // AI 3
		x, y = sl.AIPlayer3X, sl.AIPlayer3Y+sl.FaceHeight+1
	case 4: // AI 4
		x, y = sl.AIPlayer4X, sl.AIPlayer4Y+sl.FaceHeight+1
	default:
		return
	}

	// Clear the card area first
	sl.ClearPlayerArea(playerIndex)

	// Render cards
	for i, card := range cards {
		cardX := x + (i * (sl.CardWidth + 1))
		if faceDown {
			// Show face-down cards as simple blocks
			MoveCursor(cardX, y)
			fmt.Printf("%s██%s", Yellow, Reset)
		} else {
			// Show actual card
			if renderer != nil {
				MoveCursor(cardX, y)
				renderer.RenderCard(card)
			} else {
				// Fallback text representation
				MoveCursor(cardX, y)
				fmt.Printf("[%s]", card.String())
			}
		}
	}
}

// RenderStaticField renders the static field cards
func (sl *ScreenLayout) RenderStaticField(cards []Card, renderer *CardRenderer) {
	// Position static field below human player cards
	staticY := sl.HumanPlayerY + sl.CardHeight + 2

	// Clear static field area
	MoveCursor(sl.HumanPlayerX, staticY-1)
	fmt.Print(EraseLine)
	fmt.Printf("%sStatic Field:%s", Cyan, Reset)

	for i := 0; i < 2; i++ { // Clear 2 lines for static field
		MoveCursor(sl.HumanPlayerX, staticY+i)
		fmt.Print(EraseLine)
	}

	// Render static field cards
	for i, card := range cards {
		cardX := sl.HumanPlayerX + (i * (sl.CardWidth + 1))
		if renderer != nil {
			MoveCursor(cardX, staticY)
			renderer.RenderCard(card)
		} else {
			// Fallback text representation
			MoveCursor(cardX, staticY)
			fmt.Printf("[%s]", card.String())
		}
	}
}

// GameLog methods
func (gl *GameLog) AddMessage(message, msgType string) {
	// Add color coding based on message type
	coloredMsg := message
	switch msgType {
	case "error":
		coloredMsg = Red + message + Reset
	case "success":
		coloredMsg = Green + message + Reset
	case "warning":
		coloredMsg = Yellow + message + Reset
	case "important":
		coloredMsg = RedHi + message + Reset
	case "action":
		coloredMsg = Cyan + message + Reset
	case "info":
		coloredMsg = White + message + Reset
	}

	gl.Messages = append(gl.Messages, coloredMsg)

	// Keep only recent messages (circular buffer)
	if len(gl.Messages) > gl.MaxLines*3 {
		gl.Messages = gl.Messages[len(gl.Messages)-gl.MaxLines:]
	}
}

func (gl *GameLog) GetRecentMessages() []string {
	if len(gl.Messages) <= gl.MaxLines {
		return gl.Messages
	}
	return gl.Messages[len(gl.Messages)-gl.MaxLines:]
}

func (gl *GameLog) Clear() {
	gl.Messages = make([]string, 0)
}
