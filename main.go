package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"sort"
	"time"
	"github.com/joho/godotenv"
	"github.com/go-co-op/gocron"

	notifier "google-flights-crawler/notifier"
	utils "google-flights-crawler/utils"
	entities "google-flights-crawler/entities"
	lib "google-flights-crawler/lib"
	_ "github.com/go-sql-driver/mysql"
)

type Scheduler struct {
	*gocron.Scheduler
}

func main() {
	err := godotenv.Load()
	if err != nil { notifier.Notify("Warning: .env file not found, relying on environment variables") }

	apiKey := getApiKey()
	if *apiKey == "" {
		notifier.Notify("Error: API key required. Use -key flag or set SERPAPI_KEY env var.")
		os.Exit(1)
	}
	notifier.Notify("Api key loaded successfully")

	path := "/media/ShareDatabase/flightsToFollow.csv"
	flights, err := utils.LoadSearchParams(path)
	if err != nil {
		flights, err = utils.LoadSearchParams("./flightsToFollow.csv")
		if err != nil {
			notifier.Notify("Error loading search params: " + err.Error())
			os.Exit(1)
		}
	}
	notifier.Notify(fmt.Sprintf("Loaded %d flight search parameters", len(flights)))

	local, _ := time.LoadLocation("America/Sao_Paulo")
	scheduler := &Scheduler{ gocron.NewScheduler(local) }

	scheduler.scheduleFlightsCrawler(flights, apiKey)

	scheduler.StartBlocking()
}

func (scheduler *Scheduler) scheduleFlightsCrawler(flights []utils.FlightCsv, apiKey *string) {
	for _, flight := range flights {
		notifier.Notify(fmt.Sprintf("📅 Schedule: %s → %s on %s (every %s at %s)",
			flight.DepartureID, flight.ArrivalID, flight.OutboundDate, flight.Day, flight.Time))

		params := lib.SearchParams{
			APIKey:       *apiKey,
			DepartureID:  flight.DepartureID,
			ArrivalID:    flight.ArrivalID,
			OutboundDate: flight.OutboundDate,
			ReturnDate:   flight.ReturnDate,
			Adults:       flight.Adults,
			TravelClass:  flight.TravelClass,
			Stops:        flight.Stops,
			Currency:     flight.Currency,
			Language:     flight.Language,
			Country:      flight.Country,
		}

		_, err := utils.ScheduleOnDay(scheduler.Scheduler, flight.Day).At(flight.Time).Do(func() {
			startGoogleFlightsCrawler(params, nil)
		})

		if err != nil { 
			notifier.Notify(fmt.Sprintf("Error scheduling job: %s", err))
			os.Exit(1)
		}
	}
}

func getApiKey() *string {
    if key := os.Getenv("SERPAPI_KEY"); key != "" {
        return &key
    }

    key := flag.String("key", "", "SerpApi API key")
    flag.Parse()

    if *key == "" {
        notifier.Notify("SERPAPI_KEY não definida. Use a env var ou o flag -key")
    }

    return key
}

func getFlagsValuesOld() (lib.SearchParams, *string) {
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

	return params, output
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

func startGoogleFlightsCrawler(params lib.SearchParams, output *string) {
	notifier.Notify(fmt.Sprintf("🔍 Searching flights %s → %s on %s...\n", params.DepartureID, params.ArrivalID, params.OutboundDate))

	result, err := lib.ScrapeFlights(params)
	if err != nil {
		notifier.Notify(fmt.Sprintf("❌  Error: %v\n", err))
		os.Exit(1)
	}

	db := lib.CreateDatabaseConnection()
	defer db.Close()

	printResults(result)
	
	notifier.Notify(fmt.Sprintf("✅ Search completed: %s → %s on %s. Best price: %.0f %s", result.Origin, result.Destination, result.Date, result.BestPrice, result.Currency))

	er := lib.InsertIntoDB(db, result)
	if er != nil {
		notifier.Notify(fmt.Sprintf("⚠️ Could not insert into database: %v\n", er))
	} else {
		notifier.Notify("💾 Search saved to database")
	}

	if output != nil && *output != "" {
		data, _ := json.MarshalIndent(result, "", "  ")
		if err := os.WriteFile(*output, data, 0644); err != nil {
			notifier.Notify(fmt.Sprintf("⚠️  Could not write output file: %v\n", err))
		} else {
			notifier.Notify(fmt.Sprintf("💾 Results saved to %s", *output))
		}
	}
}