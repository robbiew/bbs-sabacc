package main

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

// PortraitManager handles loading and managing AI player portraits
type PortraitManager struct {
	Portraits         []Portrait
	SelectedPortraits [4]int   // Indices of selected portraits for each AI player
	debugErrors       []string // Store debug error messages
}

// Portrait represents a single ANSI art portrait
type Portrait struct {
	Lines [6]string // 6 rows of portrait data
	Attrs [6]string // 6 rows of color attribute data (unused for ANSI format)
}

// Portrait dimensions
const (
	portraitWidth  = 9
	portraitHeight = 6
)

// NewPortraitManager creates a new portrait manager and loads portraits
func NewPortraitManager() *PortraitManager {
	pm := &PortraitManager{
		Portraits: make([]Portrait, 0),
	}

	// Get current working directory for debugging
	cwd, _ := os.Getwd()

	// Try to load portraits.ans file from ansi directory
	err := pm.LoadPortraits("ansi/portraits.ans")
	if err != nil {
		// Try alternative path - maybe running from different directory
		err2 := pm.LoadPortraits("./ansi/portraits.ans")
		if err2 != nil {
			// Try without ansi directory
			err3 := pm.LoadPortraits("portraits.ans")
			if err3 != nil {
				// All paths failed - store the errors for debugging
				pm.debugErrors = []string{
					fmt.Sprintf("CWD: %s", cwd),
					fmt.Sprintf("Path 1 (ansi/portraits.ans): %v", err),
					fmt.Sprintf("Path 2 (./ansi/portraits.ans): %v", err2),
					fmt.Sprintf("Path 3 (portraits.ans): %v", err3),
				}
			}
		}
		// No portraits available - system will handle gracefully with placeholders
		pm.Portraits = make([]Portrait, 0)
	}

	// Randomize portrait selection for this game session
	pm.RandomizeSelection()

	return pm
}

// LoadPortraits loads portraits from the single stacked ANSI file
func (pm *PortraitManager) LoadPortraits(filename string) error {
	// Load from single stacked ANSI file directly
	return pm.loadFromStackedANSI(filename)
}

// loadFromStackedANSI loads portraits from a single ANSI file with stacked portraits
func (pm *PortraitManager) loadFromStackedANSI(filename string) error {
	// Try to read the file
	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read stacked ANSI file: %v", err)
	}

	// Remove SAUCE metadata like godoors does
	ansiContent := pm.trimSauceData(string(data))

	fmt.Printf("DEBUG: Content length after SAUCE removal: %d\n", len(ansiContent))

	// Split by newlines - the file now has proper line breaks
	lines := strings.Split(ansiContent, "\n")

	// Filter out empty lines but keep all ANSI content
	var contentLines []string
	for _, line := range lines {
		// Keep all non-empty lines that don't start with SAUCE
		if len(strings.TrimSpace(line)) > 0 && !strings.HasPrefix(line, "SAUCE00") {
			contentLines = append(contentLines, line)
		}
	}

	fmt.Printf("DEBUG: Content lines after filtering: %d\n", len(contentLines))

	// Calculate number of portraits (each portrait is 6 rows)
	portraitsCount := len(contentLines) / portraitHeight
	if portraitsCount == 0 {
		return fmt.Errorf("no valid portraits found in stacked ANSI file")
	}

	fmt.Printf("DEBUG: Found %d portraits\n", portraitsCount)

	// Extract portraits - store the full ANSI lines for rendering
	var portraits []Portrait
	for p := 0; p < portraitsCount; p++ {
		portrait := Portrait{}

		// Get 6 rows for this portrait
		for row := 0; row < portraitHeight; row++ {
			lineIndex := (p * portraitHeight) + row
			if lineIndex < len(contentLines) {
				// Store the full ANSI line - preserve escape sequences
				portrait.Lines[row] = contentLines[lineIndex]
			} else {
				portrait.Lines[row] = ""
			}
			// For ANSI files, attributes are embedded in the line
			portrait.Attrs[row] = ""
		}

		portraits = append(portraits, portrait)
	}

	if len(portraits) == 0 {
		return fmt.Errorf("no valid portraits found in stacked ANSI file")
	}

	pm.Portraits = portraits
	fmt.Printf("DEBUG: Successfully loaded %d portraits\n", len(pm.Portraits))
	return nil
}


// trimSauceData removes SAUCE metadata from ANSI content (similar to godoors TrimStringFromSauce)
func (pm *PortraitManager) trimSauceData(content string) string {
	// Look for SAUCE00 signature and trim everything after it
	if idx := strings.Index(content, "SAUCE00"); idx != -1 {
		return content[:idx]
	}
	// Also check for COMNT (comment block)
	if idx := strings.Index(content, "COMNT"); idx != -1 {
		return content[:idx]
	}
	return content
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
	// Debug: Show which portrait index is being used for each AI player
	portraitIndex := pm.SelectedPortraits[aiPlayerIndex]

	portrait := pm.GetPortrait(aiPlayerIndex)
	if portrait == nil {
		// Fall back to simple placeholder - also show debug info if available
		MoveCursor(x, y)
		if len(pm.debugErrors) > 0 {
			fmt.Printf("DEBUG ERR")
		} else {
			fmt.Printf("  [AI %d-%d]  ", aiPlayerIndex+1, portraitIndex)
		}
		return
	}

	// Render each line of the portrait - ANSI format with embedded escape sequences
	for i, line := range portrait.Lines {
		MoveCursor(x, y+i)
		// Output the line directly since it contains ANSI escape sequences
		fmt.Printf("%s", line)
	}
}

// GetDebugErrors returns any debug errors that occurred during loading
func (pm *PortraitManager) GetDebugErrors() []string {
	return pm.debugErrors
}
