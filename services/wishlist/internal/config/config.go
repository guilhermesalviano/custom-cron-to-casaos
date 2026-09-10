package config

import (
	"fmt"
	"os"
	"strings"
	"time"
)

type Config struct {
	WishlistID string
	Day        string
	Time       string
	Timezone   string
}

type InvalidDayError struct {
	Day string
}

func (e *InvalidDayError) Error() string {
	return fmt.Sprintf("WISHLIST_DAY must be a weekday, got %q", e.Day)
}

func Load() (Config, error) {
	cfg := Config{
		WishlistID: strings.TrimSpace(os.Getenv("WISHLIST_ID")),
		Day:        strings.ToLower(strings.TrimSpace(os.Getenv("WISHLIST_DAY"))),
		Time:       strings.TrimSpace(os.Getenv("WISHLIST_TIME")),
		Timezone:   strings.TrimSpace(os.Getenv("WISHLIST_TIMEZONE")),
	}

	if cfg.WishlistID == "" {
		return Config{}, fmt.Errorf("WISHLIST_ID is required")
	}
	if cfg.Day == "" {
		return Config{}, fmt.Errorf("WISHLIST_DAY is required")
	}
	if cfg.Time == "" {
		return Config{}, fmt.Errorf("WISHLIST_TIME is required")
	}
	if _, err := time.Parse("15:04", cfg.Time); err != nil {
		return Config{}, fmt.Errorf("WISHLIST_TIME must use HH:MM format: %w", err)
	}
	if cfg.Timezone == "" {
		cfg.Timezone = "America/Sao_Paulo"
	}

	if !validDay(cfg.Day) {
		return Config{}, &InvalidDayError{Day: cfg.Day}
	}

	return cfg, nil
}

func validDay(day string) bool {
	switch day {
	case "monday", "tuesday", "wednesday", "thursday", "friday", "saturday", "sunday":
		return true
	default:
		return false
	}
}
