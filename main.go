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

	utils "google-flights-crawler/utils"
	entities "google-flights-crawler/entities"
	lib "google-flights-crawler/lib"
	_ "github.com/go-sql-driver/mysql"
)

// ─── Display ─────────────────────────────────────────────────────────────────

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

// ─── Main ─────────────────────────────────────────────────────────────────────

func main() {
	err := godotenv.Load()
	if err != nil { log.Println("Warning: .env file not found, relying on environment variables") }

	apiKey     := flag.String("key", os.Getenv("SERPAPI_KEY"), "SerpApi API key (or set SERPAPI_KEY env var)")
	from       := flag.String("from", "GRU", "Departure IATA code (e.g. GRU, JFK, LHR)")
	to         := flag.String("to", "JFK", "Arrival IATA code (e.g. JFK, GRU, CDG)")
	outbound   := flag.String("date", time.Now().AddDate(0, 1, 0).Format("2006-01-02"), "Outbound date YYYY-MM-DD")
	returnDate := flag.String("return", "", "Return date YYYY-MM-DD (empty = one-way)")
	adults     := flag.Int("adults", 1, "Number of adult passengers")
	class      := flag.Int("class", 1, "Travel class: 1=Economy 2=Premium Economy 3=Business 4=First")
	stops      := flag.Int("stops", 0, "Max stops: 0=Any 1=Nonstop 2=1stop 3=2stops")
	currency   := flag.String("currency", "BRL", "Currency code (e.g. BRL, USD, EUR)")
	lang       := flag.String("lang", "pt", "Language code (e.g. pt, en)")
	country    := flag.String("country", "br", "Country code (e.g. br, us)")
	output     := flag.String("output", "", "Save results to JSON file (optional)")
	flag.Parse()

	if *apiKey == "" {
		fmt.Fprintln(os.Stderr, "❌  API key required. Use -key flag or set SERPAPI_KEY env var.")
		fmt.Fprintln(os.Stderr, "    Get a free key at https://serpapi.com/")
		os.Exit(1)
	}

	params := lib.SearchParams{
		APIKey:       *apiKey,
		DepartureID:  *from,
		ArrivalID:    *to,
		OutboundDate: *outbound,
		ReturnDate:   *returnDate,
		Adults:       *adults,
		TravelClass:  *class,
		Stops:        *stops,
		Currency:     *currency,
		Language:     *lang,
		Country:      *country,
	}

	fmt.Printf("🔍 Searching flights %s → %s on %s...\n", params.DepartureID, params.ArrivalID, params.OutboundDate)

	result, err := lib.FetchFlights(params)
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌  Error: %v\n", err)
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

	if *output != "" {
		data, _ := json.MarshalIndent(result, "", "  ")
		if err := os.WriteFile(*output, data, 0644); err != nil {
			fmt.Fprintf(os.Stderr, "⚠️  Could not write output file: %v\n", err)
		} else {
			fmt.Printf("💾 Results saved to %s\n", *output)
		}
	}
}
