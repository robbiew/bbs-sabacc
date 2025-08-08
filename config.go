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
		EnableStatistics: true,
		AIPersonality:    "balanced",
	}
}

// LoadConfig loads configuration from file, uses defaults if file doesn't exist
func LoadConfig() GameConfig {
	configFile := "sabacc.conf"
	config := DefaultConfig() // Start with defaults

	data, err := os.ReadFile(configFile)
	if err != nil {
		// File doesn't exist or can't be read - create default config file
		fmt.Printf("Config file not found, creating default %s\n", configFile)
		err = SaveConfig(config)
		if err != nil {
			fmt.Printf("Warning: Could not create config file: %v\n", err)
		} else {
			fmt.Printf("Created default configuration file: %s\n", configFile)
		}
		return config
	}

	err = json.Unmarshal(data, &config)
	if err != nil {
		fmt.Printf("Error reading config file, using defaults: %v\n", err)
		return DefaultConfig()
	}

	// Validate the loaded config to ensure sane values
	ValidateConfig(&config)

	fmt.Printf("Loaded configuration from %s\n", configFile)
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

	// TODO: Set LastPlayed timestamp when time functionality is implemented
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
