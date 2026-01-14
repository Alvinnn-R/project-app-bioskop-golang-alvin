package usecase

import (
	"context"
	"project-app-bioskop/internal/data/repository"
	"project-app-bioskop/internal/dto"
)

type SeatUseCaseInterface interface {
	GetSeatAvailability(ctx context.Context, cinemaID int, date, time string) ([]dto.SeatResponse, error)
	GetShowtimesByCinema(ctx context.Context, cinemaID int) ([]dto.ShowtimeResponse, error)
}

type SeatUseCase struct {
	Repo *repository.Repository
}

func NewSeatUseCase(repo *repository.Repository) SeatUseCaseInterface {
	return &SeatUseCase{Repo: repo}
}

// GetSeatAvailability retrieves seat availability for a showtime
func (u *SeatUseCase) GetSeatAvailability(ctx context.Context, cinemaID int, date, time string) ([]dto.SeatResponse, error) {
	showtime, err := u.Repo.Seat.GetShowtimeByParams(ctx, cinemaID, date, time)
	if err != nil {
		return nil, err
	}

	seats, err := u.Repo.Seat.GetSeatsByShowtime(ctx, showtime.ID)
	if err != nil {
		return nil, err
	}

	var response []dto.SeatResponse
	for _, s := range seats {
		response = append(response, dto.SeatResponse{
			ID:       s.ID,
			SeatCode: s.SeatCode,
			IsBooked: s.IsBooked,
		})
	}

	return response, nil
}

// GetShowtimesByCinema retrieves all showtimes for a cinema
func (u *SeatUseCase) GetShowtimesByCinema(ctx context.Context, cinemaID int) ([]dto.ShowtimeResponse, error) {
	showtimes, err := u.Repo.Seat.GetShowtimesByCinema(ctx, cinemaID)
	if err != nil {
		return nil, err
	}

	var response []dto.ShowtimeResponse
	for _, st := range showtimes {
		showtimeResp := dto.ShowtimeResponse{
			ID:       st.ID,
			ShowDate: st.ShowDate,
			ShowTime: st.ShowTime,
			Price:    st.Price,
		}
		
		// Add movie info if available
		if st.Movie != nil {
			showtimeResp.Movie = dto.MovieResponse{
				ID:              st.Movie.ID,
				Title:           st.Movie.Title,
				PosterURL:       st.Movie.PosterURL,
				Genres:          st.Movie.Genres,
				Rating:          st.Movie.Rating,
				DurationMinutes: st.Movie.DurationMinutes,
			}
		}
		
		// Add studio info if available
		if st.Studio != nil {
			showtimeResp.Studio = dto.StudioResponse{
				ID:         st.Studio.ID,
				Name:       st.Studio.Name,
				TotalSeats: st.Studio.TotalSeats,
			}
		}
		
		response = append(response, showtimeResp)
	}

	return response, nil
}
