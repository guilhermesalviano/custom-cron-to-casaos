package entities

type Flight struct {
	Airline       string  `json:"airline"`
	FlightNumber  string  `json:"flight_number"`
	Departure     string  `json:"departure_time"`
	Arrival       string  `json:"arrival_time"`
	Duration      int     `json:"duration_minutes"`
	Stops         int     `json:"stops"`
	Price         float64 `json:"price"`
	Currency      string  `json:"currency"`
	CarbonEmitKg  int     `json:"carbon_emissions_kg,omitempty"`
}