# Makefile for BBS Sabacc

.PHONY: all build cards clean test preview help install cp437-test

# Default target
all: build cards

# Build the main game
build:
	@echo "Building BBS Sabacc..."
	go build -ldflags="-s -w" -o sabacc .
	@echo "Game built successfully!"

# Create the card database
cards:
	@echo "Creating card database..."
	go run cmd/build-cards/main.go
	@echo "Card database created!"

# Test CP437 character support
cp437-test:
	@echo "Creating CP437 test file..."
	go run cmd/test-cp437/main.go
	@echo "View cp437_test.ans in your ANSI viewer to verify CP437 support"

# Test the card database
test: cards
	@echo "Testing card database..."
	go run cmd/build-cards/main.go test

# Create ANSI preview
preview: cards
	@echo "Creating ANSI preview..."
	go run cmd/build-cards/main.go preview

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -f sabacc
	rm -f sabacc_cards.bin
	rm -f card_index.txt
	rm -f card_preview.ans
	rm -f cp437_test.ans
	@echo "Clean complete!"

# Install to BBS directory (customize BBSDIR)
BBSDIR ?= /home/bbs/doors/sabacc
install: all
	@echo "Installing to $(BBSDIR)..."
	mkdir -p $(BBSDIR)
	cp sabacc $(BBSDIR)/
	cp sabacc_cards.bin $(BBSDIR)/
	chmod +x $(BBSDIR)/sabacc
	@echo "Installation complete!"
	@echo "Configure your BBS to launch: $(BBSDIR)/sabacc -path %%3"

# Show help
help:
	@echo "BBS Sabacc Build System"
	@echo "======================="
	@echo ""
	@echo "Targets:"
	@echo "  all        - Build game and create card database (default)"
	@echo "  build      - Build the main game executable"
	@echo "  cards      - Create the card database (sabacc_cards.bin)"
	@echo "  cp437-test - Create CP437 character test file"
	@echo "  test       - Test the card database"
	@echo "  preview    - Create ANSI preview file"
	@echo "  clean      - Remove build artifacts"
	@echo "  install    - Install to BBS directory (set BBSDIR=path)"
	@echo "  help       - Show this help"
	@echo ""
	@echo "Examples:"
	@echo "  make                    # Build everything"
	@echo "  make cp437-test         # Test CP437 support first"
	@echo "  make BBSDIR=/opt/bbs/doors/sabacc install"
	@echo "  make clean all          # Clean rebuild"
	@echo ""
	@echo "Development workflow:"
	@echo "  1. make cp437-test      # Verify ANSI support"
	@echo "  2. make cards           # Build card database"
	@echo "  3. make preview         # Check card appearance"
	@echo "  4. make build           # Build game"

# Development targets
dev-build:
	@echo "Building development version..."
	go build -o sabacc .

dev-run: dev-build cards
	@echo "Running development version..."
	./sabacc -path ./

dev-test: cp437-test preview
	@echo "Development test complete!"
	@echo "Check cp437_test.ans and card_preview.ans in your ANSI viewer"

# Release build
release: clean
	@echo "Building release version..."
	go build -ldflags="-s -w" -o sabacc .
	go run cmd/build-cards/main.go
	@echo "Release build complete!"
	@echo "Files to distribute:"
	@echo "  - sabacc (executable)"
	@echo "  - sabacc_cards.bin (card database)"

# Verify everything works
verify: clean cp437-test cards preview test
	@echo "Verification complete!"
	@echo "All files created and tested successfully"