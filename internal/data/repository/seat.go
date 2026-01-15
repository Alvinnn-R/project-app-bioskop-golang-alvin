package repository

import (
	"context"
	"project-app-bioskop/internal/data/entity"
)

type SeatRepoInterface interface {
	GetSeatsByShowtime(ctx context.Context, showtimeID int) ([]entity.SeatAvailability, error)
	GetShowtimeByParams(ctx context.Context, cinemaID int, date, time string) (entity.Showtime, error)
	GetShowtimeByID(ctx context.Context, id int) (entity.Showtime, error)
	GetSeatsByIDs(ctx context.Context, seatIDs []int) ([]entity.Seat, error)
	CheckSeatsAvailable(ctx context.Context, showtimeID int, seatIDs []int) (bool, error)
	GetShowtimesByCinema(ctx context.Context, cinemaID int) ([]entity.Showtime, error)
}

type SeatRepo struct {
	DB DBPool
}

func NewSeatRepo(db DBPool) SeatRepoInterface {
	return &SeatRepo{DB: db}
}

// GetSeatsByShowtime retrieves all seats with availability status for a showtime
func (r *SeatRepo) GetSeatsByShowtime(ctx context.Context, showtimeID int) ([]entity.SeatAvailability, error) {
	query := `
		SELECT DISTINCT s.id, s.seat_code, s.studio_id,
			   CASE WHEN EXISTS (
			   		SELECT 1 FROM booking_seats bs
			   		INNER JOIN bookings b ON b.id = bs.booking_id
			   		WHERE bs.seat_id = s.id 
			   		AND b.showtime_id = $1 
			   		AND b.status != 'cancelled'
			   ) THEN true ELSE false END as is_booked
		FROM seats s
		INNER JOIN showtimes st ON st.studio_id = s.studio_id
		WHERE st.id = $1
		ORDER BY s.seat_code`

	rows, err := r.DB.Query(ctx, query, showtimeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var seats []entity.SeatAvailability
	for rows.Next() {
		var s entity.SeatAvailability
		if err := rows.Scan(&s.ID, &s.SeatCode, &s.StudioID, &s.IsBooked); err != nil {
			return nil, err
		}
		seats = append(seats, s)
	}
	return seats, nil
}

// GetShowtimeByParams retrieves showtime by cinema, date, and time
func (r *SeatRepo) GetShowtimeByParams(ctx context.Context, cinemaID int, date, time string) (entity.Showtime, error) {
	query := `SELECT id, cinema_id, studio_id, movie_id, 
			  show_date::text as show_date, 
			  show_time::text as show_time, 
			  price 
			  FROM showtimes 
			  WHERE cinema_id = $1 AND show_date::text = $2 AND show_time::text = $3`
	var st entity.Showtime
	err := r.DB.QueryRow(ctx, query, cinemaID, date, time).Scan(
		&st.ID, &st.CinemaID, &st.StudioID, &st.MovieID,
		&st.ShowDate, &st.ShowTime, &st.Price,
	)
	if err != nil {
		return st, err
	}
	return st, nil
}

// GetShowtimeByID retrieves showtime by ID with movie and studio details
func (r *SeatRepo) GetShowtimeByID(ctx context.Context, id int) (entity.Showtime, error) {
	query := `SELECT 
				st.id, st.cinema_id, st.studio_id, st.movie_id, 
			  	st.show_date::text as show_date, 
			  	st.show_time::text as show_time, 
			  	st.price,
				m.id, m.title, m.poster_url, m.genres, m.rating, m.duration_in_minutes,
				s.id, s.name, s.total_seats
			  FROM showtimes st
			  INNER JOIN movies m ON m.id = st.movie_id
			  INNER JOIN studios s ON s.id = st.studio_id
			  WHERE st.id = $1`

	var st entity.Showtime
	var movie entity.Movie
	var studio entity.Studio

	err := r.DB.QueryRow(ctx, query, id).Scan(
		&st.ID, &st.CinemaID, &st.StudioID, &st.MovieID,
		&st.ShowDate, &st.ShowTime, &st.Price,
		&movie.ID, &movie.Title, &movie.PosterURL, &movie.Genres, &movie.Rating, &movie.DurationMinutes,
		&studio.ID, &studio.Name, &studio.TotalSeats,
	)
	if err != nil {
		return st, err
	}

	st.Movie = &movie
	st.Studio = &studio
	return st, nil
}

// GetSeatsByIDs retrieves seats by their IDs
func (r *SeatRepo) GetSeatsByIDs(ctx context.Context, seatIDs []int) ([]entity.Seat, error) {
	query := `SELECT id, studio_id, seat_code FROM seats WHERE id = ANY($1)`
	rows, err := r.DB.Query(ctx, query, seatIDs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var seats []entity.Seat
	for rows.Next() {
		var s entity.Seat
		if err := rows.Scan(&s.ID, &s.StudioID, &s.SeatCode); err != nil {
			return nil, err
		}
		seats = append(seats, s)
	}
	return seats, nil
}

// CheckSeatsAvailable checks if all seats are available for booking
func (r *SeatRepo) CheckSeatsAvailable(ctx context.Context, showtimeID int, seatIDs []int) (bool, error) {
	query := `
		SELECT COUNT(*) FROM booking_seats bs
		INNER JOIN bookings b ON b.id = bs.booking_id
		WHERE b.showtime_id = $1 AND bs.seat_id = ANY($2) AND b.status != 'cancelled'`
	var count int
	err := r.DB.QueryRow(ctx, query, showtimeID, seatIDs).Scan(&count)
	if err != nil {
		return false, err
	}
	return count == 0, nil
}

// GetShowtimesByCinema retrieves all showtimes for a cinema
func (r *SeatRepo) GetShowtimesByCinema(ctx context.Context, cinemaID int) ([]entity.Showtime, error) {
	query := `SELECT 
				st.id, st.cinema_id, st.studio_id, st.movie_id, 
			  	st.show_date::text as show_date, 
			  	st.show_time::text as show_time, 
			  	st.price,
				m.id, m.title, m.poster_url, m.genres, m.rating, m.duration_in_minutes,
				s.id, s.name, s.total_seats
			  FROM showtimes st
			  INNER JOIN movies m ON m.id = st.movie_id
			  INNER JOIN studios s ON s.id = st.studio_id
			  WHERE st.cinema_id = $1 
			  ORDER BY st.show_date, st.show_time`
	rows, err := r.DB.Query(ctx, query, cinemaID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var showtimes []entity.Showtime
	for rows.Next() {
		var st entity.Showtime
		var movie entity.Movie
		var studio entity.Studio

		if err := rows.Scan(
			&st.ID, &st.CinemaID, &st.StudioID, &st.MovieID,
			&st.ShowDate, &st.ShowTime, &st.Price,
			&movie.ID, &movie.Title, &movie.PosterURL, &movie.Genres, &movie.Rating, &movie.DurationMinutes,
			&studio.ID, &studio.Name, &studio.TotalSeats,
		); err != nil {
			return nil, err
		}

		st.Movie = &movie
		st.Studio = &studio
		showtimes = append(showtimes, st)
	}
	return showtimes, nil
}
