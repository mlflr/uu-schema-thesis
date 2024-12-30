package v5

import (
	"errors"
	"fmt"
	"net/http"

	data "thesis.lefler.eu/internal/data/branches/v5"
	e "thesis.lefler.eu/internal/error"
	util "thesis.lefler.eu/internal/util"
	"thesis.lefler.eu/internal/validator"
)

type MovieHandler struct {
	errors *e.Errors
	models *data.Models
}

func (handler *MovieHandler) CreateMovieHandler(w http.ResponseWriter, r *http.Request) {
	type crew struct {
		PersonID int64  `json:"person_id"`
		CrewType string `json:"crew_type"`
		Role     string `json:"role,omitempty"`
	}
	var input struct {
		Title    string   `json:"title"`
		Year     int32    `json:"year"`
		Genres   []string `json:"genres"`
		Runtime  int32    `json:"runtime"`
		Language string   `json:"language"`
		Crew     []crew   `json:"crew"`
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
		Runtime:  &input.Runtime,
		Language: &input.Language,
	}

	v := validator.New()

	data.ValidateMovie(v, movie)

	var director string

	for _, a := range input.Crew {

		person, err := handler.models.People.Get(a.PersonID)
		if err != nil {
			switch {
			case errors.Is(err, data.ErrRecordNotFound):
				handler.errors.NotFoundResponse(w, r)
			default:
				handler.errors.ServerErrorResponse(w, r, err)
			}
			return
		}

		crewMember := &data.Crew{
			PersonID:   a.PersonID,
			PersonName: person.Name,
			CrewType:   a.CrewType,
			Role:       a.Role,
		}
		data.ValidateCrew(v, crewMember)
		movie.Crew = append(movie.Crew, crewMember)

		if crewMember.CrewType == "Director" && director == "" {
			director = crewMember.PersonName
		}
	}

	if !v.Valid() {
		handler.errors.FailedValidationResponse(w, r, v.Errors)
		return
	}

	err = handler.models.Movies.Insert(movie, director)
	if err != nil {
		handler.errors.ServerErrorResponse(w, r, err)
		return
	}

	for _, crewMember := range movie.Crew {
		crewMember.MovieID = movie.ID

		err = handler.models.Crew.Insert(crewMember)
		if err != nil {
			handler.errors.ServerErrorResponse(w, r, err)
			return
		}
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v5/movie/%d", movie.ID))

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

	movie.Crew, err = handler.models.Crew.GetForMovie(movie.ID)
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

	type crew struct {
		PersonID int64  `json:"person_id"`
		CrewType string `json:"crew_type"`
		Role     string `json:"role,omitempty"`
	}
	var input struct {
		Title    *string   `json:"title"`
		Year     *int32    `json:"year"`
		Genres   *[]string `json:"genres"`
		Runtime  *int32    `json:"runtime"`
		Language *string   `json:"language"`
		Crew     []crew    `json:"crew"`
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

	var director string

	if input.Crew != nil {
		err = handler.models.Crew.DeleteForMovie(movie.ID)
		if err != nil {
			handler.errors.ServerErrorResponse(w, r, err)
			return
		}

		for _, a := range input.Crew {
			if a.PersonID < 1 {
				handler.errors.BadRequestResponse(w, r, errors.New("invalid person_id"))
				return
			}

			person, err := handler.models.People.Get(a.PersonID)
			if err != nil {
				switch {
				case errors.Is(err, data.ErrRecordNotFound):
					handler.errors.NotFoundResponse(w, r)
				default:
					handler.errors.ServerErrorResponse(w, r, err)
				}
				return
			}

			crewMember := &data.Crew{
				MovieID:    movie.ID,
				PersonID:   a.PersonID,
				PersonName: person.Name,
				CrewType:   a.CrewType,
				Role:       a.Role,
			}
			data.ValidateCrew(v, crewMember)

			err = handler.models.Crew.Insert(crewMember)
			if err != nil {
				handler.errors.ServerErrorResponse(w, r, err)
				return
			}
			movie.Crew = append(movie.Crew, crewMember)

			if crewMember.CrewType == "Director" && director == "" {
				director = crewMember.PersonName
			}
		}
	}

	err = handler.models.Movies.Update(movie, director)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			handler.errors.EditConflictResponse(w, r)
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
