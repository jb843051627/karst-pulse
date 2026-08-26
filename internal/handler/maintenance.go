package handler

import (
	"net/http"

	"github.com/karst-pulse/karst-pulse/internal/model"
	"github.com/karst-pulse/karst-pulse/internal/transport"
	"github.com/karst-pulse/karst-pulse/internal/validate"
)

func (h *Handler) listMaintenance(w http.ResponseWriter, r *http.Request) {
	filter, err := transport.Filter(r)
	if err == nil {
		err = validate.Filter(filter)
	}
	if err != nil {
		h.fail(w, err)
		return
	}
	result, err := h.app.ListMaintenance(r.Context(), filter)
	if err != nil {
		h.fail(w, err)
		return
	}
	transport.WriteData(w, http.StatusOK, result)
}

func (h *Handler) createMaintenance(w http.ResponseWriter, r *http.Request) {
	var input model.MaintenanceInput
	if err := transport.DecodeJSON(r, &input); err != nil {
		h.fail(w, transport.BadRequest(err.Error(), nil))
		return
	}
	task, err := h.app.CreateMaintenance(r.Context(), input)
	if err != nil {
		h.fail(w, err)
		return
	}
	transport.WriteData(w, http.StatusCreated, task)
}

func (h *Handler) completeMaintenance(w http.ResponseWriter, r *http.Request) {
	id, err := transport.ID(r, "id")
	if err != nil {
		h.fail(w, err)
		return
	}
	task, err := h.app.CompleteMaintenance(r.Context(), id)
	if err != nil {
		h.fail(w, err)
		return
	}
	transport.WriteData(w, http.StatusOK, task)
}
