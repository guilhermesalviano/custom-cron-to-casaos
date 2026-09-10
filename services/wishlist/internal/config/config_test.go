package config

import "testing"

func TestLoadDefaultsTimezoneAndNormalizesDay(t *testing.T) {
	t.Setenv("WISHLIST_ID", "ABC123")
	t.Setenv("WISHLIST_DAY", "Saturday")
	t.Setenv("WISHLIST_TIME", "17:00")
	t.Setenv("WISHLIST_TIMEZONE", "")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if cfg.Day != "saturday" {
		t.Fatalf("Day = %q, want saturday", cfg.Day)
	}
	if cfg.Timezone != "America/Sao_Paulo" {
		t.Fatalf("Timezone = %q, want America/Sao_Paulo", cfg.Timezone)
	}
}

func TestLoadRequiresConfiguration(t *testing.T) {
	for _, key := range []string{"WISHLIST_ID", "WISHLIST_DAY", "WISHLIST_TIME"} {
		t.Run(key, func(t *testing.T) {
			t.Setenv("WISHLIST_ID", "ABC123")
			t.Setenv("WISHLIST_DAY", "Saturday")
			t.Setenv("WISHLIST_TIME", "17:00")
			t.Setenv(key, "")

			if _, err := Load(); err == nil {
				t.Fatalf("Load() error = nil, want error for %s", key)
			}
		})
	}
}

func TestLoadRejectsInvalidDayAndTime(t *testing.T) {
	t.Setenv("WISHLIST_ID", "ABC123")
	t.Setenv("WISHLIST_DAY", "Someday")
	t.Setenv("WISHLIST_TIME", "not-a-time")

	if _, err := Load(); err == nil {
		t.Fatal("Load() error = nil, want validation error")
	}

	t.Setenv("WISHLIST_DAY", "Saturday")
	if _, err := Load(); err == nil {
		t.Fatal("Load() error = nil, want invalid time error")
	}
}
