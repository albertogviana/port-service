package entity

import (
	"encoding/json"
	"log"
)

type Port struct {
	ID        int64
	Name      string
	City      string
	Country   string
	Alias     []string
	Regions   []string
	Latitude  float64
	Longitude float64
	Province  string
	Timezone  string
	Unloc     string
	Code      string
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (p *Port) UnmarshalJSON(data []byte) error {
	var pr PortRequest

	if err := json.Unmarshal(data, &pr); err != nil {
		log.Fatal(err)
		return err
	}

	p.Name = pr.Name
	p.City = pr.City
	p.Country = pr.Country
	p.Alias = pr.Alias
	p.Regions = pr.Regions

	if len(pr.Coordinates) > 0 {
		p.Latitude = pr.Coordinates[0]
		p.Longitude = pr.Coordinates[1]
	}

	p.Province = pr.Province
	p.Timezone = pr.Timezone
	p.Unloc = pr.Unlocs[0]
	p.Code = pr.Code

	return nil
}

type PortRequest struct {
	Name        string    `json:"name"`
	City        string    `json:"city"`
	Country     string    `json:"country"`
	Alias       []string  `json:"alias"`
	Regions     []string  `json:"regions"`
	Coordinates []float64 `json:"coordinates"`
	Province    string    `json:"province"`
	Timezone    string    `json:"timezone"`
	Unlocs      []string  `json:"unlocs"`
	Code        string    `json:"code"`
}
