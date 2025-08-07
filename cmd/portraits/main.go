// cmd/portraits/main.go - Portrait build tool
package main

import (
	"encoding/binary"
	"fmt"
	"os"
)

// Portrait represents a single ANSI art portrait
type Portrait struct {
	Lines [6]string // 6 rows of portrait data
}

// Portrait file format constants
const (
	portraitWidth  = 9
	portraitHeight = 6
	portraitChars  = portraitWidth * portraitHeight // 54 bytes for characters
	portraitAttrs  = portraitWidth * portraitHeight // 54 bytes for attributes
	portraitSize   = portraitChars + portraitAttrs  // 108 bytes total
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage:")
		fmt.Println("  portraits extract <output.bin> <input.xb>   - Extract from XBin file")
		fmt.Println("  portraits info <input.bin>                 - Show file info")
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "extract":
		if len(os.Args) < 4 {
			fmt.Println("Usage: portraits extract <output.bin> <input.xb>")
			os.Exit(1)
		}
		outputFile := os.Args[2]
		xbinFile := os.Args[3]
		err := extractFromXBin(xbinFile, outputFile)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Successfully extracted portraits from %s to %s\n", xbinFile, outputFile)

	case "info":
		if len(os.Args) < 3 {
			fmt.Println("Usage: portraits info <input.bin>")
			os.Exit(1)
		}
		filename := os.Args[2]
		err := showPortraitInfo(filename)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

	default:
		fmt.Printf("Unknown command: %s\n", command)
		os.Exit(1)
	}
}

// showPortraitInfo displays information about a portrait file
func showPortraitInfo(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	// Read header
	var numPortraits uint32
	err = binary.Read(file, binary.LittleEndian, &numPortraits)
	if err != nil {
		return fmt.Errorf("failed to read header: %v", err)
	}

	// Get file size
	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("failed to get file info: %v", err)
	}

	expectedSize := 4 + (int(numPortraits) * portraitSize)

	fmt.Printf("Portrait file: %s\n", filename)
	fmt.Printf("Number of portraits: %d\n", numPortraits)
	fmt.Printf("Expected file size: %d bytes\n", expectedSize)
	fmt.Printf("Actual file size: %d bytes\n", fileInfo.Size())

	return nil
}

// extractFromXBin extracts portraits from an XBin file and saves them to a binary file
func extractFromXBin(xbinFile, outputFile string) error {
	// Load all portraits from the XBin file
	portraits, err := loadAllXBinPortraits(xbinFile)
	if err != nil {
		return fmt.Errorf("failed to load XBin portraits: %v", err)
	}

	if len(portraits) == 0 {
		return fmt.Errorf("no portraits found in XBin file")
	}

	// Create the output file
	file, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("failed to create output file: %v", err)
	}
	defer file.Close()

	// Write header (number of portraits)
	numPortraits := uint32(len(portraits))
	err = binary.Write(file, binary.LittleEndian, numPortraits)
	if err != nil {
		return fmt.Errorf("failed to write header: %v", err)
	}

	// Write each portrait
	for i, portraitData := range portraits {
		n, err := file.Write(portraitData)
		if err != nil {
			return fmt.Errorf("failed to write portrait %d: %v", i, err)
		}
		if n != portraitSize {
			return fmt.Errorf("wrote %d bytes for portrait %d, expected %d", n, i, portraitSize)
		}
	}

	return nil
}

// XBin file format structure
type XBinHeader struct {
	ID       [4]byte // "XBIN"
	EofChar  byte    // 0x1A
	Width    uint16  // Width in characters
	Height   uint16  // Height in characters
	FontSize byte    // Font size
	Flags    byte    // Format flags
}

// loadAllXBinPortraits loads multiple portraits from an XBin file
func loadAllXBinPortraits(filename string) ([][]byte, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open XBin file: %v", err)
	}
	defer file.Close()

	// Read XBin header
	var header XBinHeader
	err = binary.Read(file, binary.LittleEndian, &header)
	if err != nil {
		return nil, fmt.Errorf("failed to read XBin header: %v", err)
	}

	// Verify XBin signature
	if string(header.ID[:]) != "XBIN" {
		return nil, fmt.Errorf("invalid XBin file: missing XBIN signature")
	}

	if header.EofChar != 0x1A {
		return nil, fmt.Errorf("invalid XBin file: incorrect EOF character")
	}

	fmt.Printf("XBin file: %dx%d, font size: %d, estimated portraits: %d\n",
		header.Width, header.Height, header.FontSize, header.Height/portraitHeight)

	// Calculate number of portraits that can fit
	portraitsPerFile := int(header.Height) / portraitHeight
	if portraitsPerFile == 0 {
		portraitsPerFile = 1 // At least try to extract one portrait
	}

	// Skip palette data if present (flag bit 0)
	if header.Flags&0x01 != 0 {
		// Skip 48 bytes of palette data (16 colors * 3 bytes RGB)
		_, err = file.Seek(48, 1)
		if err != nil {
			return nil, fmt.Errorf("failed to skip palette data: %v", err)
		}
	}

	// Skip font data if present (flag bit 1)
	if header.Flags&0x02 != 0 {
		fontDataSize := int(header.FontSize) * 256
		_, err = file.Seek(int64(fontDataSize), 1)
		if err != nil {
			return nil, fmt.Errorf("failed to skip font data: %v", err)
		}
	}

	// Read character data
	charDataSize := int(header.Width) * int(header.Height)
	charData := make([]byte, charDataSize)
	n, err := file.Read(charData)
	if err != nil {
		return nil, fmt.Errorf("failed to read character data: %v", err)
	}
	if n != charDataSize {
		return nil, fmt.Errorf("incomplete character data: read %d, expected %d", n, charDataSize)
	}

	// Read attribute data (colors) - XBin stores this after character data
	attrDataSize := int(header.Width) * int(header.Height)
	attrData := make([]byte, attrDataSize)
	n, err = file.Read(attrData)
	if err != nil {
		return nil, fmt.Errorf("failed to read attribute data: %v", err)
	}
	if n != attrDataSize {
		return nil, fmt.Errorf("incomplete attribute data: read %d, expected %d", n, attrDataSize)
	}

	fmt.Printf("Successfully read %d chars + %d attributes from XBin\n", len(charData), len(attrData))

	// Extract multiple portraits from the XBin data
	var portraits [][]byte

	for p := 0; p < portraitsPerFile; p++ {
		portraitData := make([]byte, portraitSize)

		// Extract 9x6 characters + attributes for this portrait
		for row := 0; row < portraitHeight; row++ {
			for col := 0; col < portraitWidth; col++ {
				srcRow := (p * portraitHeight) + row
				if srcRow < int(header.Height) && col < int(header.Width) {
					srcIndex := srcRow*int(header.Width) + col
					charDstIndex := row*portraitWidth + col
					attrDstIndex := portraitChars + row*portraitWidth + col

					if srcIndex < len(charData) && charDstIndex < portraitChars {
						// Store character
						portraitData[charDstIndex] = charData[srcIndex]
						// Store attribute (color)
						if srcIndex < len(attrData) && attrDstIndex < portraitSize {
							portraitData[attrDstIndex] = attrData[srcIndex]
						} else {
							portraitData[attrDstIndex] = 0x07 // Default white on black
						}
					} else {
						portraitData[charDstIndex] = ' '  // Space if out of bounds
						portraitData[attrDstIndex] = 0x07 // Default white on black
					}
				} else {
					// Pad with spaces if XBin is smaller than expected
					charDstIndex := row*portraitWidth + col
					attrDstIndex := portraitChars + row*portraitWidth + col
					portraitData[charDstIndex] = ' '
					portraitData[attrDstIndex] = 0x07 // Default white on black
				}
			}
		}

		portraits = append(portraits, portraitData)
	}

	return portraits, nil
}
