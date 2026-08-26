package handler

import (
	"net/http"

	"github.com/karst-pulse/karst-pulse/internal/model"
	"github.com/karst-pulse/karst-pulse/internal/transport"
	"github.com/karst-pulse/karst-pulse/internal/validate"
)

func (h *Handler) listSensors(w http.ResponseWriter, r *http.Request) {
	filter, err := transport.Filter(r)
	if err == nil {
		err = validate.Filter(filter)
	}
	if err != nil {
		h.fail(w, err)
		return
	}
	result, err := h.app.ListSensors(r.Context(), filter)
	if err != nil {
		h.fail(w, err)
		return
	}
	transport.WriteData(w, http.StatusOK, result)
}

func (h *Handler) createSensor(w http.ResponseWriter, r *http.Request) {
	var input model.SensorInput
	if err := transport.DecodeJSON(r, &input); err != nil {
		h.fail(w, transport.BadRequest(err.Error(), nil))
		return
	}
	sensor, err := h.app.CreateSensor(r.Context(), input)
	if err != nil {
		h.fail(w, err)
		return
	}
	transport.WriteData(w, http.StatusCreated, sensor)
}
