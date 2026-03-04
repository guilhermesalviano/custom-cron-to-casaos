package entities

import "time"

type SearchResult struct {
	SearchedAt  time.Time `json:"searched_at"`
	Origin      string    `json:"origin"`
	Destination string    `json:"destination"`
	Date        string    `json:"outbound_date"`
	ReturnDate  string    `json:"return_date,omitempty"`
	BestFlights []Flight  `json:"best_flights"`
	OtherFlights []Flight `json:"other_flights"`
	BestPrice   float64   `json:"best_price"`
	Currency    string    `json:"currency"`
}