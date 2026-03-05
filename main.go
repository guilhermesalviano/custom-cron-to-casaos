package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"sort"
	"time"
	"github.com/joho/godotenv"
	"log"
	"database/sql"
	"github.com/go-co-op/gocron"

	notifier "google-flights-crawler/notifier"
	utils "google-flights-crawler/utils"
	entities "google-flights-crawler/entities"
	lib "google-flights-crawler/lib"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	err := godotenv.Load()
	if err != nil { log.Println("Warning: .env file not found, relying on environment variables") }

	output, params := getFlagsValues()
	if params.APIKey == "" {
		notifier.Notify("API key required. Use -key flag or set SERPAPI_KEY env var.\n Get a free key at https://serpapi.com/")
		os.Exit(1)
	}

	local, _ := time.LoadLocation("America/Sao_Paulo")
	scheduler := gocron.NewScheduler(local)

	scheduler.Every(1).Saturday().At("08:00").Do(func() {
		search(params, output)
	})

	notifier.Notify("📅 Agendador iniciado. Aguardando próximo sábado...")

	scheduler.StartBlocking()
}

func getFlagsValues() (*string, lib.SearchParams) {
	var p lib.SearchParams
	var output string

	flag.StringVar(&p.APIKey,       "key",      os.Getenv("SERPAPI_KEY"), "SerpApi API key")
	flag.StringVar(&p.DepartureID,  "from",     "GYN",        "Departure IATA code")
	flag.StringVar(&p.ArrivalID,    "to",       "GRU",        "Arrival IATA code")
	flag.StringVar(&p.OutboundDate, "date",     "2026-05-02", "Outbound date YYYY-MM-DD")
	flag.StringVar(&p.ReturnDate,   "return",   "",           "Return date YYYY-MM-DD")
	flag.IntVar(&p.Adults,          "adults",   1,            "Number of adult passengers")
	flag.IntVar(&p.TravelClass,     "class",    1,            "Travel class")
	flag.IntVar(&p.Stops,           "stops",    0,            "Max stops")
	flag.StringVar(&p.Currency,     "currency", "BRL",        "Currency code")
	flag.StringVar(&p.Language,     "lang",     "pt",         "Language code")
	flag.StringVar(&p.Country,      "country",  "br",         "Country code")
	flag.StringVar(&output,         "output",   "",           "Save results to JSON file")
	flag.Parse()

	return &output, p
}

func printResults(r *entities.SearchResult) {
	fmt.Printf("\n╔══════════════════════════════════════════════════════╗\n")
	fmt.Printf("║         Google Flights — Best Price Crawler          ║\n")
	fmt.Printf("╚══════════════════════════════════════════════════════╝\n")
	fmt.Printf("  Route     : %s → %s\n", r.Origin, r.Destination)
	fmt.Printf("  Outbound  : %s\n", r.Date)
	if r.ReturnDate != "" {
		fmt.Printf("  Return    : %s\n", r.ReturnDate)
	}
	fmt.Printf("  Searched  : %s\n", r.SearchedAt.Format("2006-01-02 15:04:05 UTC"))
	fmt.Printf("  ★ Best price: %.0f %s\n\n", r.BestPrice, r.Currency)

	printSection := func(label string, flights []entities.Flight) {
		if len(flights) == 0 {
			return
		}
		// Sort by price
		sort.Slice(flights, func(i, j int) bool {
			return flights[i].Price < flights[j].Price
		})
		fmt.Printf("── %s (%d results) ──────────────────────────────\n", label, len(flights))
		fmt.Printf("  %-22s %-8s %-8s %-10s %-6s %s\n",
			"Airline", "Dep", "Arr", "Duration", "Stops", "Price")
		fmt.Printf("  %s\n", "─────────────────────────────────────────────────────────")
		for _, f := range flights {
			dur := fmt.Sprintf("%dh%02dm", f.Duration/60, f.Duration%60)
			stops := "nonstop"
			if f.Stops == 1 {
				stops = "1 stop"
			} else if f.Stops > 1 {
				stops = fmt.Sprintf("%d stops", f.Stops)
			}
			fmt.Printf("  %-22s %-8s %-8s %-10s %-6s %.0f %s\n",
				utils.Truncate(f.Airline, 22),
				utils.TimeOnly(f.Departure),
				utils.TimeOnly(f.Arrival),
				dur, stops,
				f.Price, f.Currency,
			)
		}
		fmt.Println()
	}

	printSection("Best Flights", r.BestFlights)
	printSection("Other Flights", r.OtherFlights)
}

func search(params lib.SearchParams, output *string) {
	notifier.Notify(fmt.Sprintf("🔍 Searching flights %s → %s on %s...\n", params.DepartureID, params.ArrivalID, params.OutboundDate))

	result, err := lib.FetchFlights(params)
	if err != nil {
		notifier.Notify(fmt.Sprintf("❌  Error: %v\n", err))
		os.Exit(1)
	}

	dbConfig := lib.DBConfig{
		Username: os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Host: os.Getenv("DB_HOST"),
		Port: os.Getenv("DB_PORT"),
		Database: os.Getenv("DB_NAME"),
	}

	db := lib.DbConnection(dbConfig)
	defer db.Close()

	printResults(result)
	
	notifier.Notify(fmt.Sprintf("✅ Search completed: %s → %s on %s. Best price: %.0f %s", result.Origin, result.Destination, result.Date, result.BestPrice, result.Currency))

	er := insertToDB(db, result)
	if er != nil {
		notifier.Notify(fmt.Sprintf("⚠️ Could not insert into database: %v\n", er))
	} else {
		notifier.Notify("💾 Search result saved to database")
	}

	if *output != "" {
		data, _ := json.MarshalIndent(result, "", "  ")
		if err := os.WriteFile(*output, data, 0644); err != nil {
			notifier.Notify(fmt.Sprintf("⚠️  Could not write output file: %v\n", err))
		} else {
			notifier.Notify(fmt.Sprintf("💾 Results saved to %s", *output))
		}
	}
}

func insertToDB(db *sql.DB, r *entities.SearchResult) error {
	_, err := db.Exec(
		"INSERT INTO flight_crawled (origin, destination, airline, stops, price, flightDate, searchDate) VALUES (?, ?, ?, ?, ?, ?, ?)",
		r.Origin,
		r.Destination,
		r.BestFlights[0].Airline, 
		r.BestFlights[0].Stops, 
		r.BestFlights[0].Departure, 
		r.BestFlights[0].Price,
		time.Now(),
	)
	return err
}