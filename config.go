package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// GameConfig holds configuration settings
type GameConfig struct {
	MinAnte          int    `json:"min_ante"`
	MaxAnte          int    `json:"max_ante"`
	MaxBet           int    `json:"max_bet"`
	StartingCredits  int    `json:"starting_credits"`
	ShiftProbability int    `json:"shift_probability"` // 1 in X chance
	MinRoundsToCall  int    `json:"min_rounds_to_call"`
	IdleTimeout      int    `json:"idle_timeout_seconds"`
	EnableSound      bool   `json:"enable_sound"`
	EnableStatistics bool   `json:"enable_statistics"`
	AIPersonality    string `json:"ai_personality"` // "conservative", "aggressive", "balanced"
}

// PlayerStats tracks player statistics
type PlayerStats struct {
	PlayerName      string `json:"player_name"`
	GamesPlayed     int    `json:"games_played"`
	GamesWon        int    `json:"games_won"`
	PureSabaccs     int    `json:"pure_sabaccs"`
	IdiotsArrays    int    `json:"idiots_arrays"`
	BombOuts        int    `json:"bomb_outs"`
	CreditsWon      int    `json:"credits_won"`
	CreditsLost     int    `json:"credits_lost"`
	ShiftsTriggered int    `json:"shifts_triggered"`
	FoldsCount      int    `json:"folds_count"`
	LastPlayed      string `json:"last_played"`
}

// DefaultConfig returns the default game configuration
func DefaultConfig() GameConfig {
	return GameConfig{
		MinAnte:          10,
		MaxAnte:          100,
		MaxBet:           500,
		StartingCredits:  1000,
		ShiftProbability: 6, // 1 in 6 chance
		MinRoundsToCall:  2,
		IdleTimeout:      300, // 5 minutes
		EnableSound:      true,
		EnableStatistics: true,
		AIPersonality:    "balanced",
	}
}

// LoadConfig loads configuration from file or creates default
func LoadConfig() GameConfig {
	configFile := "sabacc.conf"
	config := DefaultConfig()

	data, err := os.ReadFile(configFile)
	if err != nil {
		// Create default config file
		SaveConfig(config)
		return config
	}

	err = json.Unmarshal(data, &config)
	if err != nil {
		fmt.Printf("Error reading config file, using defaults: %v\n", err)
		return DefaultConfig()
	}

	return config
}

// SaveConfig saves configuration to file
func SaveConfig(config GameConfig) error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile("sabacc.conf", data, 0644)
}

// LoadPlayerStats loads player statistics
func LoadPlayerStats(playerName string) PlayerStats {
	statsDir := "stats"
	os.MkdirAll(statsDir, 0755)

	// Sanitize filename
	filename := strings.ReplaceAll(strings.ToLower(playerName), " ", "_")
	filename = filepath.Join(statsDir, filename+".json")

	stats := PlayerStats{
		PlayerName: playerName,
	}

	data, err := os.ReadFile(filename)
	if err != nil {
		return stats // Return empty stats if file doesn't exist
	}

	err = json.Unmarshal(data, &stats)
	if err != nil {
		fmt.Printf("Error reading stats file: %v\n", err)
		return PlayerStats{PlayerName: playerName}
	}

	return stats
}

// SavePlayerStats saves player statistics
func SavePlayerStats(stats PlayerStats) error {
	statsDir := "stats"
	os.MkdirAll(statsDir, 0755)

	filename := strings.ReplaceAll(strings.ToLower(stats.PlayerName), " ", "_")
	filename = filepath.Join(statsDir, filename+".json")

	data, err := json.MarshalIndent(stats, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filename, data, 0644)
}

// UpdateStats updates player statistics after a game
func UpdateStats(stats *PlayerStats, won bool, credits int, specialHand string, bombedOut bool, folded bool) {
	stats.GamesPlayed++

	if won {
		stats.GamesWon++
		stats.CreditsWon += credits
	} else {
		if credits < 0 {
			stats.CreditsLost += -credits
		}
	}

	switch specialHand {
	case "pure_sabacc":
		stats.PureSabaccs++
	case "idiots_array":
		stats.IdiotsArrays++
	}

	if bombedOut {
		stats.BombOuts++
	}

	if folded {
		stats.FoldsCount++
	}

	stats.LastPlayed = getCurrentDateTime()
}

// getCurrentDateTime returns current date/time as string
func getCurrentDateTime() string {
	return fmt.Sprintf("%d", getCurrentTimestamp())
}

// getCurrentTimestamp returns current Unix timestamp
func getCurrentTimestamp() int64 {
	return int64(42) // Placeholder - would use time.Now().Unix()
}

// DisplayStats shows formatted statistics
func DisplayStats(stats PlayerStats) {
	fmt.Printf("═══════════════════════════════════════════\n")
	fmt.Printf("            PLAYER STATISTICS\n")
	fmt.Printf("═══════════════════════════════════════════\n\n")

	fmt.Printf("Player: %s\n\n", stats.PlayerName)
	fmt.Printf("Games Played: %d\n", stats.GamesPlayed)
	fmt.Printf("Games Won: %d\n", stats.GamesWon)

	if stats.GamesPlayed > 0 {
		winRate := float64(stats.GamesWon) / float64(stats.GamesPlayed) * 100
		fmt.Printf("Win Rate: %.1f%%\n", winRate)
	}

	fmt.Printf("Pure Sabaccs: %d\n", stats.PureSabaccs)
	fmt.Printf("Idiot's Arrays: %d\n", stats.IdiotsArrays)
	fmt.Printf("Bomb Outs: %d\n", stats.BombOuts)
	fmt.Printf("Folds: %d\n", stats.FoldsCount)
	fmt.Printf("Credits Won: %d\n", stats.CreditsWon)
	fmt.Printf("Credits Lost: %d\n", stats.CreditsLost)

	netCredits := stats.CreditsWon - stats.CreditsLost
	fmt.Printf("Net Credits: %d\n", netCredits)

	if stats.ShiftsTriggered > 0 {
		fmt.Printf("Shifts Triggered: %d\n", stats.ShiftsTriggered)
	}
}

// AIPersonality defines different AI behaviors
type AIPersonality struct {
	Name                 string
	DrawThreshold        int     // Will draw if hand total is below this
	StandThreshold       int     // Will stand if hand total is above this
	CallThreshold        int     // Will call if hand total is above this
	FoldThreshold        int     // Will fold if hand total is above this (bomb risk)
	StaticFieldChance    float32 // Probability of using static field for good cards
	AggressivenessFactor float32 // Multiplier for risk-taking
}

// GetAIPersonality returns AI settings based on personality type
func GetAIPersonality(personalityType string) AIPersonality {
	switch personalityType {
	case "conservative":
		return AIPersonality{
			Name:                 "Conservative",
			DrawThreshold:        15,
			StandThreshold:       20,
			CallThreshold:        18,
			FoldThreshold:        22, // Lower threshold = fold earlier
			StaticFieldChance:    0.8,
			AggressivenessFactor: 0.7,
		}
	case "aggressive":
		return AIPersonality{
			Name:                 "Aggressive",
			DrawThreshold:        18,
			StandThreshold:       22,
			CallThreshold:        20,
			FoldThreshold:        28, // Higher threshold = fold later
			StaticFieldChance:    0.3,
			AggressivenessFactor: 1.5,
		}
	default: // "balanced"
		return AIPersonality{
			Name:                 "Balanced",
			DrawThreshold:        16,
			StandThreshold:       21,
			CallThreshold:        19,
			FoldThreshold:        25, // Middle ground
			StaticFieldChance:    0.5,
			AggressivenessFactor: 1.0,
		}
	}
}

// ConfigMenu allows players to adjust game settings
func showConfigMenu() {
	// This would be called from the main menu
	// Implementation would allow players to adjust:
	// - Starting credits
	// - AI difficulty/personality
	// - Sound on/off
	// - Statistics tracking on/off
	// - Ante limits
	fmt.Println("Configuration menu not yet implemented")
}

// ValidateConfig ensures config values are within reasonable ranges
func ValidateConfig(config *GameConfig) {
	if config.MinAnte < 1 {
		config.MinAnte = 1
	}
	if config.MaxAnte < config.MinAnte {
		config.MaxAnte = config.MinAnte * 2
	}
	if config.StartingCredits < config.MaxAnte {
		config.StartingCredits = config.MaxAnte * 10
	}
	if config.ShiftProbability < 2 {
		config.ShiftProbability = 2
	}
	if config.MinRoundsToCall < 1 {
		config.MinRoundsToCall = 1
	}
	if config.IdleTimeout < 60 {
		config.IdleTimeout = 60
	}
}

// GetLeaderboard returns top players by various metrics
func GetLeaderboard(metric string, limit int) []PlayerStats {
	// This would scan all player stats files and return top players
	// Metrics could be: "wins", "credits", "sabaccs", "win_rate"
	var leaderboard []PlayerStats

	// Placeholder implementation
	// In real implementation, this would:
	// 1. Read all .json files from stats/ directory
	// 2. Parse each file into PlayerStats
	// 3. Sort by requested metric
	// 4. Return top N results

	return leaderboard
}

// ResetStats allows player to reset their statistics
func ResetPlayerStats(playerName string) error {
	filename := strings.ReplaceAll(strings.ToLower(playerName), " ", "_")
	filename = filepath.Join("stats", filename+".json")

	// Remove the stats file
	return os.Remove(filename)
}

// ExportStats exports all statistics to a readable format
func ExportStats(format string) error {
	// Could export to CSV, JSON, or text format
	// Useful for BBS sysops to track game usage
	return nil
}

// GetSystemStats returns overall system statistics
func GetSystemStats() map[string]interface{} {
	stats := make(map[string]interface{})

	// This would aggregate stats across all players
	stats["total_games"] = 0
	stats["total_players"] = 0
	stats["total_credits_won"] = 0
	stats["total_sabaccs"] = 0

	return stats
}
