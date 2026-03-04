package lib

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
	"net/url"
	"strconv"

	entities "google-flights-crawler/entities"
)

type serpLayover struct {
	Name     string `json:"name"`
	ID       string `json:"id"`
	Duration int    `json:"duration"`
}

type serpFlight struct {
	Flights []struct {
		DepartureAirport struct{ Time string `json:"time"` } `json:"departure_airport"`
		ArrivalAirport   struct{ Time string `json:"time"` } `json:"arrival_airport"`
		Duration         int    `json:"duration"`
		Airline          string `json:"airline"`
		FlightNumber     string `json:"flight_number"`
	} `json:"flights"`
	Layovers        []serpLayover `json:"layovers"`
	TotalDuration   int `json:"total_duration"`
	CarbonEmissions struct {
		ThisFlightKg int `json:"this_flight"`
	} `json:"carbon_emissions"`
	Price   int    `json:"price"`
	Airline string `json:"airline,omitempty"`
}

type serpResponse struct {
	BestFlights  []serpFlight `json:"best_flights"`
	OtherFlights []serpFlight `json:"other_flights"`
	Error        string       `json:"error,omitempty"`
}

type SearchParams struct {
	APIKey       string
	DepartureID  string // IATA code, e.g. "GRU"
	ArrivalID    string // IATA code, e.g. "JFK"
	OutboundDate string // YYYY-MM-DD
	ReturnDate   string // YYYY-MM-DD (empty = one-way)
	Adults       int
	TravelClass  int // 1=Economy 2=Premium 3=Business 4=First
	Stops        int // 0=Any 1=Nonstop 2=1stop 3=2stops
	Currency     string // e.g. "BRL", "USD"
	Language     string // e.g. "pt", "en"
	Country      string // e.g. "br", "us"
}

func BuildQuery(p SearchParams) url.Values {
	q := url.Values{}
	q.Set("engine", "google_flights")
	q.Set("api_key", p.APIKey)
	q.Set("departure_id", p.DepartureID)
	q.Set("arrival_id", p.ArrivalID)
	q.Set("outbound_date", p.OutboundDate)
	q.Set("adults", strconv.Itoa(p.Adults))
	q.Set("travel_class", strconv.Itoa(p.TravelClass))
	q.Set("stops", strconv.Itoa(p.Stops))
	q.Set("currency", p.Currency)
	q.Set("hl", p.Language)
	q.Set("gl", p.Country)

	if p.ReturnDate != "" {
		q.Set("return_date", p.ReturnDate)
		q.Set("type", "1")
	} else {
		q.Set("type", "2")
	}
	return q
}

func FetchFlights(p SearchParams) (*entities.SearchResult, error) {
	const SERPAPIBASE = "https://serpapi.com/search"
	reqURL := fmt.Sprintf("%s?%s", SERPAPIBASE, BuildQuery(p).Encode())

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Get(reqURL)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}

	var raw serpResponse
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, fmt.Errorf("parsing JSON: %w", err)
	}
	if raw.Error != "" {
		return nil, fmt.Errorf("API error: %s", raw.Error)
	}

	result := &entities.SearchResult{
		SearchedAt:  time.Now().UTC(),
		Origin:      p.DepartureID,
		Destination: p.ArrivalID,
		Date:        p.OutboundDate,
		ReturnDate:  p.ReturnDate,
		Currency:    p.Currency,
	}

	parse := func(raw []serpFlight) []entities.Flight {
		var out []entities.Flight
		for _, sf := range raw {
			airline := sf.Airline
			flightNum := ""
			depTime := ""
			arrTime := ""
			if len(sf.Flights) > 0 {
				if airline == "" {
					airline = sf.Flights[0].Airline
				}
				flightNum = sf.Flights[0].FlightNumber
				depTime = sf.Flights[0].DepartureAirport.Time
				arrTime = sf.Flights[len(sf.Flights)-1].ArrivalAirport.Time
			}
			out = append(out, entities.Flight{
				Airline:      airline,
				FlightNumber: flightNum,
				Departure:    depTime,
				Arrival:      arrTime,
				Duration:     sf.TotalDuration,
				Stops:        len(sf.Layovers),
				Price:        float64(sf.Price),
				Currency:     p.Currency,
				CarbonEmitKg: sf.CarbonEmissions.ThisFlightKg,
			})
		}
		return out
	}

	result.BestFlights = parse(raw.BestFlights)
	result.OtherFlights = parse(raw.OtherFlights)

	// Find overall best (lowest) price
	all := append(result.BestFlights, result.OtherFlights...)
	if len(all) > 0 {
		best := all[0].Price
		for _, f := range all[1:] {
			if f.Price > 0 && f.Price < best {
				best = f.Price
			}
		}
		result.BestPrice = best
	}

	return result, nil
}
