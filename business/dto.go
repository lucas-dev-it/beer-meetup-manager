package business

type MeetupBeersData struct {
	BeerPacks      float64 `json:"beer_packs"`
	MaxTemperature float64 `json:"max_temperature"`
	MinTemperature float64 `json:"min_temperature"`
	AttendeesCount int     `json:"attendees_count"`
}
