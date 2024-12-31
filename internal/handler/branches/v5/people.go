package v5

import (
	"errors"
	"fmt"
	"net/http"

	"cloud.google.com/go/civil"
	data "thesis.lefler.eu/internal/data/branches/v5"
	e "thesis.lefler.eu/internal/error"
	"thesis.lefler.eu/internal/util"
	"thesis.lefler.eu/internal/validator"
)

type PersonHandler struct {
	errors *e.Errors
	models *data.Models
}

func (handler *PersonHandler) CreateHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name      string     `json:"name"`      // Person name
		Birthdate civil.Date `json:"birthdate"` // Person birthdate
	}

	err := util.ReadJSON(w, r, &input)
	if err != nil {
		handler.errors.BadRequestResponse(w, r, err)
		return
	}

	person := &data.Person{
		Name:      input.Name,
		Birthdate: &input.Birthdate,
	}

	v := validator.New()

	if data.ValidatePerson(v, person); !v.Valid() {
		handler.errors.FailedValidationResponse(w, r, v.Errors)
		return
	}

	err = handler.models.People.Insert(person)
	if err != nil {
		handler.errors.ServerErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v5/people/%d", person.ID))

	err = util.WriteJSON(w, http.StatusCreated, util.Envelope{"person": person}, headers)
	if err != nil {
		handler.errors.ServerErrorResponse(w, r, err)
	}
}

func (handler *PersonHandler) GetHandler(w http.ResponseWriter, r *http.Request) {
	id, err := util.ReadIDParam(r)
	if err != nil {
		handler.errors.NotFoundResponse(w, r)
		return
	}

	person, err := handler.models.People.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			handler.errors.NotFoundResponse(w, r)
		default:
			handler.errors.ServerErrorResponse(w, r, err)
		}
		return
	}

	err = util.WriteJSON(w, http.StatusOK, util.Envelope{"person": person}, nil)
	if err != nil {
		handler.errors.ServerErrorResponse(w, r, err)
	}
}

func (handler *PersonHandler) UpdateHandler(w http.ResponseWriter, r *http.Request) {
	id, err := util.ReadIDParam(r)
	if err != nil {
		handler.errors.NotFoundResponse(w, r)
		return
	}

	person, err := handler.models.People.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			handler.errors.NotFoundResponse(w, r)
		default:
			handler.errors.ServerErrorResponse(w, r, err)
		}
		return
	}

	var input struct {
		Name      *string     `json:"name"`      // Person name
		Birthdate *civil.Date `json:"birthdate"` // Person birthdate
	}

	err = util.ReadJSON(w, r, &input)
	if err != nil {
		handler.errors.BadRequestResponse(w, r, err)
		return
	}

	if input.Name != nil {
		person.Name = *input.Name
	}

	if input.Birthdate != nil {
		person.Birthdate = input.Birthdate
	}

	v := validator.New()

	if data.ValidatePerson(v, person); !v.Valid() {
		handler.errors.FailedValidationResponse(w, r, v.Errors)
		return
	}

	err = handler.models.People.Update(person)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			handler.errors.EditConflictResponse(w, r)
		default:
			handler.errors.ServerErrorResponse(w, r, err)
		}
		return
	}

	err = util.WriteJSON(w, http.StatusOK, util.Envelope{"person": person}, nil)
	if err != nil {
		handler.errors.ServerErrorResponse(w, r, err)
	}
}

func (handler *PersonHandler) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	id, err := util.ReadIDParam(r)
	if err != nil {
		handler.errors.NotFoundResponse(w, r)
		return
	}

	err = handler.models.People.Delete(id)
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

func (handler *PersonHandler) ListHandler(w http.ResponseWriter, r *http.Request) {
	people, err := handler.models.People.GetAll()
	if err != nil {
		handler.errors.ServerErrorResponse(w, r, err)
		return
	}

	err = util.WriteJSON(w, http.StatusOK, util.Envelope{"people": people}, nil)
	if err != nil {
		handler.errors.ServerErrorResponse(w, r, err)
	}
}
