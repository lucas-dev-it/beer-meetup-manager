package business

import "time"

type MeetupBeersData struct {
	BeerPacks      float64         `json:"beer_packs"`
	MaxTemperature float64         `json:"max_temperature"`
	MinTemperature float64         `json:"min_temperature"`
	MeetupMetadata *MeetupMetadata `json:"meetup_metadata"`
}

type MeetupMetadata struct {
	Name           string     `json:"name"`
	StartDate      *time.Time `json:"start_date"`
	EndDate        *time.Time `json:"end_date"`
	Country        string     `json:"country"`
	State          string     `json:"state"`
	City           string     `json:"city"`
	AttendeesCount int        `json:"attendees_count"`
}
