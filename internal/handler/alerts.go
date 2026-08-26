package handler

import (
	"net/http"

	"github.com/karst-pulse/karst-pulse/internal/transport"
	"github.com/karst-pulse/karst-pulse/internal/validate"
)

func (h *Handler) listAlerts(w http.ResponseWriter, r *http.Request) {
	filter, err := transport.Filter(r)
	if err == nil {
		err = validate.Filter(filter)
	}
	if err != nil {
		h.fail(w, err)
		return
	}
	result, err := h.app.ListAlerts(r.Context(), filter)
	if err != nil {
		h.fail(w, err)
		return
	}
	transport.WriteData(w, http.StatusOK, result)
}

func (h *Handler) acknowledgeAlert(w http.ResponseWriter, r *http.Request) {
	id, err := transport.ID(r, "id")
	if err != nil {
		h.fail(w, err)
		return
	}
	alert, err := h.app.AcknowledgeAlert(r.Context(), id)
	if err != nil {
		h.fail(w, err)
		return
	}
	transport.WriteData(w, http.StatusOK, alert)
}
