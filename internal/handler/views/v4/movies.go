package v4

import (
	"errors"
	"fmt"
	"net/http"

	data "thesis.lefler.eu/internal/data/views/v4"
	e "thesis.lefler.eu/internal/error"
	util "thesis.lefler.eu/internal/util"
	"thesis.lefler.eu/internal/validator"
)

type MovieHandler struct {
	errors *e.Errors
	models *data.Models
}

func (handler *MovieHandler) CreateMovieHandler(w http.ResponseWriter, r *http.Request) {
	type actor struct {
		ActorID int64  `json:"actor_id"`
		Role    string `json:"role"`
	}
	var input struct {
		Title    string   `json:"title"`
		Year     int32    `json:"year"`
		Genres   []string `json:"genres"`
		Director string   `json:"director"`
		Runtime  int32    `json:"runtime"`
		Language string   `json:"language"`
		Actors   []actor  `json:"actors"`
	}

	err := util.ReadJSON(w, r, &input)
	if err != nil {
		handler.errors.BadRequestResponse(w, r, err)
		return
	}

	movie := &data.Movie{
		Title:    input.Title,
		Year:     input.Year,
		Genres:   input.Genres,
		Director: &input.Director,
		Runtime:  &input.Runtime,
		Language: &input.Language,
	}

	v := validator.New()

	if data.ValidateMovie(v, movie); !v.Valid() {
		handler.errors.FailedValidationResponse(w, r, v.Errors)
		return
	}

	err = handler.models.Movies.Insert(movie)
	if err != nil {
		handler.errors.ServerErrorResponse(w, r, err)
		return
	}

	for _, a := range input.Actors {
		if a.ActorID < 1 {
			handler.errors.BadRequestResponse(w, r, errors.New("invalid actor_id"))
			return
		}

		movieActor := data.MovieActor{
			MovieID: movie.ID,
			ActorID: a.ActorID,
			Role:    a.Role,
		}

		err = handler.models.MovieActors.Insert(&movieActor)
		if err != nil {
			handler.errors.ServerErrorResponse(w, r, err)
			return
		}
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v4/movie/%d", movie.ID))

	err = util.WriteJSON(w, http.StatusCreated, util.Envelope{"movie": movie}, headers)
	if err != nil {
		handler.errors.ServerErrorResponse(w, r, err)
	}
}

func (handler *MovieHandler) GetMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := util.ReadIDParam(r)
	if err != nil {
		handler.errors.NotFoundResponse(w, r)
		return
	}

	movie, err := handler.models.Movies.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			handler.errors.NotFoundResponse(w, r)
		default:
			handler.errors.ServerErrorResponse(w, r, err)
		}
		return
	}

	movie.Actors, err = handler.models.MovieActors.GetForMovie(movie.ID)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			handler.errors.NotFoundResponse(w, r)
		default:
			handler.errors.ServerErrorResponse(w, r, err)
		}
		return
	}

	err = util.WriteJSON(w, http.StatusOK, util.Envelope{"movie": movie}, nil)
	if err != nil {
		handler.errors.ServerErrorResponse(w, r, err)
	}
}

func (handler *MovieHandler) UpdateMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := util.ReadIDParam(r)
	if err != nil {
		handler.errors.NotFoundResponse(w, r)
		return
	}

	movie, err := handler.models.Movies.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			handler.errors.NotFoundResponse(w, r)
		default:
			handler.errors.ServerErrorResponse(w, r, err)
		}
		return
	}

	type actor struct {
		ActorID int64  `json:"actor_id"`
		Role    string `json:"role"`
	}
	var input struct {
		Title    *string   `json:"title"`
		Year     *int32    `json:"year"`
		Genres   *[]string `json:"genres"`
		Director *string   `json:"director"`
		Runtime  *int32    `json:"runtime"`
		Language *string   `json:"language"`
		Actors   []actor   `json:"actors"`
	}

	err = util.ReadJSON(w, r, &input)
	if err != nil {
		handler.errors.BadRequestResponse(w, r, err)
		return
	}

	if input.Title != nil {
		movie.Title = *input.Title
	}
	if input.Year != nil {
		movie.Year = *input.Year
	}
	if input.Genres != nil {
		movie.Genres = *input.Genres
	}
	if input.Director != nil {
		movie.Director = input.Director
	}
	if input.Runtime != nil {
		movie.Runtime = input.Runtime
	}
	if input.Language != nil {
		movie.Language = input.Language
	}

	v := validator.New()

	if data.ValidateMovie(v, movie); !v.Valid() {
		handler.errors.FailedValidationResponse(w, r, v.Errors)
		return
	}

	err = handler.models.Movies.Update(movie)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			handler.errors.EditConflictResponse(w, r)
		default:
			handler.errors.ServerErrorResponse(w, r, err)
		}
		return
	}

	if input.Actors != nil {
		err = handler.models.MovieActors.DeleteForMovie(movie.ID)
		if err != nil {
			handler.errors.ServerErrorResponse(w, r, err)
			return
		}

		for _, a := range input.Actors {
			if a.ActorID < 1 {
				handler.errors.BadRequestResponse(w, r, errors.New("invalid actor_id"))
				return
			}

			actor, err := handler.models.Actors.Get(a.ActorID)
			if err != nil {
				switch {
				case errors.Is(err, data.ErrRecordNotFound):
					handler.errors.NotFoundResponse(w, r)
				default:
					handler.errors.ServerErrorResponse(w, r, err)
				}
				return
			}

			movieActor := data.MovieActor{
				MovieID:   movie.ID,
				ActorID:   a.ActorID,
				ActorName: actor.Name,
				Role:      a.Role,
			}

			err = handler.models.MovieActors.Insert(&movieActor)
			if err != nil {
				handler.errors.ServerErrorResponse(w, r, err)
				return
			}
			movie.Actors = append(movie.Actors, &movieActor)
		}
	}

	err = util.WriteJSON(w, http.StatusOK, util.Envelope{"movie": movie}, nil)
	if err != nil {
		handler.errors.ServerErrorResponse(w, r, err)
	}
}

func (handler *MovieHandler) DeleteMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := util.ReadIDParam(r)
	if err != nil {
		handler.errors.NotFoundResponse(w, r)
		return
	}

	err = handler.models.Movies.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			handler.errors.NotFoundResponse(w, r)
		default:
			handler.errors.ServerErrorResponse(w, r, err)
		}
		return
	}

	err = util.WriteJSON(w, http.StatusOK, util.Envelope{"message": "movie successfully deleted"}, nil)
	if err != nil {
		handler.errors.ServerErrorResponse(w, r, err)
	}
}

func (handler *MovieHandler) ListMoviesHandler(w http.ResponseWriter, r *http.Request) {
	movies, err := handler.models.Movies.GetAll()
	if err != nil {
		handler.errors.ServerErrorResponse(w, r, err)
		return
	}

	err = util.WriteJSON(w, http.StatusOK, util.Envelope{"movies": movies}, nil)
	if err != nil {
		handler.errors.ServerErrorResponse(w, r, err)
	}
}
