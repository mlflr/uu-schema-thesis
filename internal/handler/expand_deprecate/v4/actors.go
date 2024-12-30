package v4

import (
	"errors"
	"fmt"
	"net/http"

	"cloud.google.com/go/civil"
	data "thesis.lefler.eu/internal/data/expand_deprecate/v4"
	e "thesis.lefler.eu/internal/error"
	"thesis.lefler.eu/internal/util"
	"thesis.lefler.eu/internal/validator"
)

type ActorHandler struct {
	errors *e.Errors
	models *data.Models
}

func (handler *ActorHandler) CreateActorHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name      string     `json:"name"`      // Actor name
		Birthdate civil.Date `json:"birthdate"` // Actor birthdate
	}

	err := util.ReadJSON(w, r, &input)
	if err != nil {
		handler.errors.BadRequestResponse(w, r, err)
		return
	}

	actor := &data.Actor{
		Name:      input.Name,
		Birthdate: &input.Birthdate,
	}

	v := validator.New()

	if data.ValidateActor(v, actor); !v.Valid() {
		handler.errors.FailedValidationResponse(w, r, v.Errors)
		return
	}

	err = handler.models.Actors.Insert(actor)
	if err != nil {
		handler.errors.ServerErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v4/actor/%d", actor.ID))

	err = util.WriteJSON(w, http.StatusCreated, util.Envelope{"actor": actor}, headers)
	if err != nil {
		handler.errors.ServerErrorResponse(w, r, err)
	}
}

func (handler *ActorHandler) GetActorHandler(w http.ResponseWriter, r *http.Request) {
	id, err := util.ReadIDParam(r)
	if err != nil {
		handler.errors.NotFoundResponse(w, r)
		return
	}

	actor, err := handler.models.Actors.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			handler.errors.NotFoundResponse(w, r)
		default:
			handler.errors.ServerErrorResponse(w, r, err)
		}
		return
	}

	err = util.WriteJSON(w, http.StatusOK, util.Envelope{"actor": actor}, nil)
	if err != nil {
		handler.errors.ServerErrorResponse(w, r, err)
	}
}

func (handler *ActorHandler) UpdateActorHandler(w http.ResponseWriter, r *http.Request) {
	id, err := util.ReadIDParam(r)
	if err != nil {
		handler.errors.NotFoundResponse(w, r)
		return
	}

	actor, err := handler.models.Actors.Get(id)
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
		Name      *string     `json:"name"`      // Actor name
		Birthdate *civil.Date `json:"birthdate"` // Actor birthdate
	}

	err = util.ReadJSON(w, r, &input)
	if err != nil {
		handler.errors.BadRequestResponse(w, r, err)
		return
	}

	if input.Name != nil {
		actor.Name = *input.Name
	}

	if input.Birthdate != nil {
		actor.Birthdate = input.Birthdate
	}

	v := validator.New()

	if data.ValidateActor(v, actor); !v.Valid() {
		handler.errors.FailedValidationResponse(w, r, v.Errors)
		return
	}

	err = handler.models.Actors.Update(actor)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			handler.errors.EditConflictResponse(w, r)
		default:
			handler.errors.ServerErrorResponse(w, r, err)
		}
		return
	}

	err = util.WriteJSON(w, http.StatusOK, util.Envelope{"actor": actor}, nil)
	if err != nil {
		handler.errors.ServerErrorResponse(w, r, err)
	}
}

func (handler *ActorHandler) DeleteActorHandler(w http.ResponseWriter, r *http.Request) {
	id, err := util.ReadIDParam(r)
	if err != nil {
		handler.errors.NotFoundResponse(w, r)
		return
	}

	err = handler.models.Actors.Delete(id)
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

func (handler *ActorHandler) ListActorsHandler(w http.ResponseWriter, r *http.Request) {
	actors, err := handler.models.Actors.GetAll()
	if err != nil {
		handler.errors.ServerErrorResponse(w, r, err)
		return
	}

	err = util.WriteJSON(w, http.StatusOK, util.Envelope{"actors": actors}, nil)
	if err != nil {
		handler.errors.ServerErrorResponse(w, r, err)
	}
}
