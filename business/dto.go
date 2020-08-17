package business

import "time"

type MeetupBeersData struct {
	BeerPacks      *float64        `json:"beer_packs,omitempty"`
	MaxTemperature float64         `json:"max_temperature"`
	MinTemperature float64         `json:"min_temperature"`
	MeetupMetadata *MeetupMetadata `json:"meetup_metadata"`
}

type MeetupMetadata struct {
	ID             uint       `json:"id"`
	Name           string     `json:"name"`
	StartDate      *time.Time `json:"start_date"`
	EndDate        *time.Time `json:"end_date"`
	Country        string     `json:"country"`
	State          string     `json:"state"`
	City           string     `json:"city"`
	AttendeesCount *int       `json:"attendees_count,omitempty"`
}

type TokenIssue struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type ClaimSet struct {
	AccessToken map[string]interface{} `json:"access_token"`
	ExpiresAt   int64                  `json:"expires_at"`
}
