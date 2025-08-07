package main

import (
	"encoding/binary"
	"fmt"
	"math/rand"
	"os"
	"time"
)

// PortraitManager handles loading and managing AI player portraits
type PortraitManager struct {
	Portraits         []Portrait
	SelectedPortraits [4]int // Indices of selected portraits for each AI player
}

// Portrait represents a single ANSI art portrait
type Portrait struct {
	Lines [6]string // 6 rows of portrait data
}

// Portrait file format:
// Header: 4 bytes (uint32) - number of portraits
// Data: For each portrait, 54 bytes of raw data (9 chars × 6 lines)
const (
	portraitWidth  = 9
	portraitHeight = 6
	portraitChars  = portraitWidth * portraitHeight // 54 bytes for characters
	portraitAttrs  = portraitWidth * portraitHeight // 54 bytes for attributes
	portraitSize   = portraitChars + portraitAttrs  // 108 bytes total
)

// NewPortraitManager creates a new portrait manager and loads portraits
func NewPortraitManager() *PortraitManager {
	pm := &PortraitManager{
		Portraits: make([]Portrait, 0),
	}

	// Try to load pre-built portraits.bin file
	err := pm.LoadPortraits("portraits.bin")
	if err != nil {
		// No portraits available - system will handle gracefully with placeholders
		pm.Portraits = make([]Portrait, 0)
	}

	// Randomize portrait selection for this game session
	pm.RandomizeSelection()

	return pm
}

// LoadPortraits loads portraits from a binary file
func (pm *PortraitManager) LoadPortraits(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open portrait file: %v", err)
	}
	defer file.Close()

	// Read header (number of portraits)
	var numPortraits uint32
	err = binary.Read(file, binary.LittleEndian, &numPortraits)
	if err != nil {
		return fmt.Errorf("failed to read portrait count: %v", err)
	}

	if numPortraits == 0 || numPortraits > 100 { // Sanity check
		return fmt.Errorf("invalid portrait count: %d", numPortraits)
	}

	// Read each portrait
	pm.Portraits = make([]Portrait, numPortraits)
	for i := uint32(0); i < numPortraits; i++ {
		// Read portrait data (108 bytes: 54 chars + 54 attributes)
		portraitData := make([]byte, portraitSize)
		n, err := file.Read(portraitData)
		if err != nil {
			return fmt.Errorf("failed to read portrait %d: %v", i, err)
		}
		if n != portraitSize {
			return fmt.Errorf("incomplete portrait %d: read %d bytes, expected %d", i, n, portraitSize)
		}

		// Convert raw bytes to 6 lines of 9 characters each (only use character data for now)
		portrait := Portrait{}
		for row := 0; row < portraitHeight; row++ {
			start := row * portraitWidth
			end := start + portraitWidth
			portrait.Lines[row] = string(portraitData[start:end])
		}

		pm.Portraits[i] = portrait
	}

	return nil
}

// RandomizeSelection randomly selects portraits for each AI player
func (pm *PortraitManager) RandomizeSelection() {
	if len(pm.Portraits) == 0 {
		return
	}

	rand.Seed(time.Now().UnixNano())

	// Ensure each AI player gets a different portrait (if we have enough)
	usedIndices := make(map[int]bool)

	for i := 0; i < 4; i++ {
		var selectedIndex int
		attempts := 0

		// Try to find an unused portrait, but don't loop forever
		for attempts < 10 {
			selectedIndex = rand.Intn(len(pm.Portraits))
			if !usedIndices[selectedIndex] || len(pm.Portraits) <= 4 {
				break
			}
			attempts++
		}

		pm.SelectedPortraits[i] = selectedIndex
		usedIndices[selectedIndex] = true
	}
}

// GetPortrait returns the portrait for a specific AI player
func (pm *PortraitManager) GetPortrait(aiPlayerIndex int) *Portrait {
	if aiPlayerIndex < 0 || aiPlayerIndex >= 4 {
		return nil
	}

	portraitIndex := pm.SelectedPortraits[aiPlayerIndex]
	if portraitIndex < 0 || portraitIndex >= len(pm.Portraits) {
		return nil
	}

	return &pm.Portraits[portraitIndex]
}

// RenderPortrait renders a portrait at the specified screen position
func (pm *PortraitManager) RenderPortrait(x, y, aiPlayerIndex int) {
	portrait := pm.GetPortrait(aiPlayerIndex)
	if portrait == nil {
		// Fall back to simple placeholder
		MoveCursor(x, y)
		fmt.Printf("  [AI %d]  ", aiPlayerIndex+1)
		return
	}

	// Render each line of the portrait
	for i, line := range portrait.Lines {
		MoveCursor(x, y+i)
		fmt.Printf("%s%s%s", White, line, Reset)
	}
}
