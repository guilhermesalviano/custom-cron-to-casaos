package utils

import (
	"encoding/csv"
	"os"
	"strconv"
)

type FlightCsv struct {
	APIKey       string
	DepartureID  string
	ArrivalID    string
	OutboundDate string
	ReturnDate   string
	Adults       int
	TravelClass  int
	Stops        int
	Currency     string
	Language     string
	Country      string
	Day          string
	Time         string
}

func LoadSearchParams(filePath string) ([]FlightCsv, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	if _, err := reader.Read(); err != nil {
		return nil, err
	}

	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var allParams []FlightCsv

	for _, row := range records {
		adults, _ := strconv.Atoi(row[4])
		class, _  := strconv.Atoi(row[5])
		stops, _  := strconv.Atoi(row[6])

		p := FlightCsv{
			DepartureID:  row[0],
			ArrivalID:    row[1],
			OutboundDate: row[2],
			ReturnDate:   row[3],
			Adults:       adults,
			TravelClass:  class,
			Stops:        stops,
			Currency:     row[7],
			Language:     row[8],
			Country:      row[9],
			Day:      		row[10],
			Time:      		row[11],
		}
		allParams = append(allParams, p)
	}

	return allParams, nil
}