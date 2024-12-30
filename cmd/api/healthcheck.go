package main

import (
	"net/http"

	"thesis.lefler.eu/internal/util"
)

func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	env := util.Envelope{
		"status": "available",
		"system_info": map[string]string{
			"environment": app.config.env,
			"version":     version,
		},
	}

	err := util.WriteJSON(w, http.StatusOK, env, nil)
	if err != nil {
		app.errors.ServerErrorResponse(w, r, err)
	}
}
