package handler

import (
	"net/http"

	"github.com/karst-pulse/karst-pulse/internal/model"
	"github.com/karst-pulse/karst-pulse/internal/transport"
	"github.com/karst-pulse/karst-pulse/internal/validate"
)

func (h *Handler) listSprings(w http.ResponseWriter, r *http.Request) {
	filter, err := transport.Filter(r)
	if err == nil {
		err = validate.Filter(filter)
	}
	if err != nil {
		h.fail(w, err)
		return
	}
	result, err := h.app.ListSprings(r.Context(), filter)
	if err != nil {
		h.fail(w, err)
		return
	}
	transport.WriteData(w, http.StatusOK, result)
}

func (h *Handler) createSpring(w http.ResponseWriter, r *http.Request) {
	var input model.SpringInput
	if err := transport.DecodeJSON(r, &input); err != nil {
		h.fail(w, transport.BadRequest(err.Error(), nil))
		return
	}
	spring, err := h.app.CreateSpring(r.Context(), input)
	if err != nil {
		h.fail(w, err)
		return
	}
	transport.WriteData(w, http.StatusCreated, spring)
}

func (h *Handler) getSpring(w http.ResponseWriter, r *http.Request) {
	id, err := transport.ID(r, "id")
	if err != nil {
		h.fail(w, err)
		return
	}
	spring, err := h.app.GetSpring(r.Context(), id)
	if err != nil {
		h.fail(w, err)
		return
	}
	transport.WriteData(w, http.StatusOK, spring)
}
