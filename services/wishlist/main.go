package main

import (
	"log"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/joho/godotenv"

	"github.com/guilhermesalviano/casaos-cron/services/wishlist/internal/config"
	"github.com/guilhermesalviano/casaos-cron/services/wishlist/internal/notify"
	"github.com/guilhermesalviano/casaos-cron/services/wishlist/internal/wishlist"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, relying on environment variables")
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Invalid wishlist configuration: %v", err)
	}

	location, err := time.LoadLocation(cfg.Timezone)
	if err != nil {
		log.Fatalf("Invalid wishlist timezone %q: %v", cfg.Timezone, err)
	}

	scheduler := gocron.NewScheduler(location)
	jobScheduler, err := scheduleOnDay(scheduler, cfg.Day, cfg.Time)
	if err != nil {
		log.Fatalf("Could not schedule wishlist crawler: %v", err)
	}

	_, err = jobScheduler.Do(func() {
		wishlist.CrawlAndStore()
	})
	if err != nil {
		log.Fatalf("Could not register wishlist crawler: %v", err)
	}

	notify.Notify("Wishlist service started")
	log.Printf("Scheduled Amazon wishlist crawler every %s at %s (%s)", cfg.Day, cfg.Time, cfg.Timezone)
	scheduler.StartBlocking()
}

func scheduleOnDay(scheduler *gocron.Scheduler, day, at string) (*gocron.Scheduler, error) {
	job := scheduler.Every(1)

	switch day {
	case "monday":
		job = job.Monday()
	case "tuesday":
		job = job.Tuesday()
	case "wednesday":
		job = job.Wednesday()
	case "thursday":
		job = job.Thursday()
	case "friday":
		job = job.Friday()
	case "saturday":
		job = job.Saturday()
	case "sunday":
		job = job.Sunday()
	default:
		return nil, &config.InvalidDayError{Day: day}
	}

	return job.At(at), nil
}
