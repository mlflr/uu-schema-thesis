package error

import (
	"fmt"
	"log/slog"
	"net/http"

	"thesis.lefler.eu/internal/util"
)

type Errors struct {
	logger *slog.Logger
}

func NewErrors(logger *slog.Logger) Errors {
	return Errors{logger}
}

func (handler *Errors) LogError(r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
	)

	handler.logger.Error(err.Error(), "method", method, "uri", uri)
}

func (handler *Errors) ErrorResponse(w http.ResponseWriter, r *http.Request, status int, message any) {
	env := util.Envelope{"error": message}

	err := util.WriteJSON(w, status, env, nil)
	if err != nil {
		handler.LogError(r, err)
		w.WriteHeader(500)
	}
}

func (handler *Errors) ServerErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	handler.LogError(r, err)

	message := "the server encountered a problem and could not process your request"
	handler.ErrorResponse(w, r, http.StatusInternalServerError, message)
}

func (handler *Errors) NotFoundResponse(w http.ResponseWriter, r *http.Request) {
	message := "the requested resource could not be found"
	handler.ErrorResponse(w, r, http.StatusNotFound, message)
}

func (handler *Errors) MethodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("the %s method is not supported for this resource", r.Method)
	handler.ErrorResponse(w, r, http.StatusMethodNotAllowed, message)
}

func (handler *Errors) BadRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	handler.ErrorResponse(w, r, http.StatusBadRequest, err.Error())
}

func (handler *Errors) FailedValidationResponse(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	handler.ErrorResponse(w, r, http.StatusUnprocessableEntity, errors)
}

func (handler *Errors) EditConflictResponse(w http.ResponseWriter, r *http.Request) {
	message := "unable to update the record due to an edit conflict, please try again"
	handler.ErrorResponse(w, r, http.StatusConflict, message)
}
