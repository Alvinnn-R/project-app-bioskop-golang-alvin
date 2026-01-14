package entity

import "time"

// Cinema represents a cinema location
type Cinema struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Location  string    `json:"location"`
	CreatedAt time.Time `json:"created_at"`
}

// Studio represents a studio inside a cinema
type Studio struct {
	ID         int    `json:"id"`
	CinemaID   int    `json:"cinema_id"`
	Name       string `json:"name"`
	TotalSeats int    `json:"total_seats"`
}

// Seat represents a seat inside a studio
type Seat struct {
	ID       int    `json:"id"`
	StudioID int    `json:"studio_id"`
	SeatCode string `json:"seat_code"`
}

// SeatAvailability represents seat status for a showtime
type SeatAvailability struct {
	ID       int    `json:"id"`
	SeatCode string `json:"seat_code"`
	StudioID int    `json:"studio_id"`
	IsBooked bool   `json:"is_booked"`
}

// Movie represents a movie entity
type Movie struct {
	ID              int       `json:"id"`
	Title           string    `json:"title"`
	PosterURL       string    `json:"poster_url"`
	Genres          []string  `json:"genres"`
	Rating          float64   `json:"rating"`
	ReviewCount     int       `json:"review_count"`
	ReleaseDate     time.Time `json:"release_date"`
	DurationMinutes int       `json:"duration_minutes"`
	ReleaseStatus   string    `json:"release_status"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// Showtime represents a movie showtime
type Showtime struct {
	ID       int     `json:"id"`
	CinemaID int     `json:"cinema_id"`
	StudioID int     `json:"studio_id"`
	MovieID  int     `json:"movie_id"`
	ShowDate string  `json:"show_date"`
	ShowTime string  `json:"show_time"`
	Price    float64 `json:"price"`
	Movie    *Movie  `json:"movie,omitempty"`
	Studio   *Studio `json:"studio,omitempty"`
}
