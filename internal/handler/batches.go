package handler

import (
	"net/http"

	"github.com/karst-pulse/karst-pulse/internal/model"
	"github.com/karst-pulse/karst-pulse/internal/transport"
	"github.com/karst-pulse/karst-pulse/internal/validate"
)

func (h *Handler) listBatches(w http.ResponseWriter, r *http.Request) {
	filter, err := transport.Filter(r)
	if err == nil {
		err = validate.Filter(filter)
	}
	if err != nil {
		h.fail(w, err)
		return
	}
	result, err := h.app.ListBatches(r.Context(), filter)
	if err != nil {
		h.fail(w, err)
		return
	}
	transport.WriteData(w, http.StatusOK, result)
}

func (h *Handler) createBatch(w http.ResponseWriter, r *http.Request) {
	var input model.BatchInput
	if err := transport.DecodeJSON(r, &input); err != nil {
		h.fail(w, transport.BadRequest(err.Error(), nil))
		return
	}
	batch, err := h.app.CreateBatch(r.Context(), input)
	if err != nil {
		h.fail(w, err)
		return
	}
	transport.WriteData(w, http.StatusCreated, batch)
}

func (h *Handler) listSamples(w http.ResponseWriter, r *http.Request) {
	batchID, err := transport.ID(r, "id")
	if err != nil {
		h.fail(w, err)
		return
	}
	items, err := h.app.ListSamples(r.Context(), batchID)
	if err != nil {
		h.fail(w, err)
		return
	}
	transport.WriteData(w, http.StatusOK, items)
}

func (h *Handler) createSample(w http.ResponseWriter, r *http.Request) {
	batchID, err := transport.ID(r, "id")
	if err != nil {
		h.fail(w, err)
		return
	}
	var input model.SampleInput
	if err := transport.DecodeJSON(r, &input); err != nil {
		h.fail(w, transport.BadRequest(err.Error(), nil))
		return
	}
	sample, err := h.app.AddSample(r.Context(), batchID, input)
	if err != nil {
		h.fail(w, err)
		return
	}
	transport.WriteData(w, http.StatusCreated, sample)
}
