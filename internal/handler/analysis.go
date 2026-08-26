package handler

import (
	"net/http"
	"time"

	"github.com/karst-pulse/karst-pulse/internal/transport"
)

func (h *Handler) analysis(w http.ResponseWriter, r *http.Request) {
	springID, err := transport.OptionalID(r, "spring_id")
	if err != nil || springID <= 0 {
		if err == nil {
			err = transport.BadRequest("spring_id is required", nil)
		}
		h.fail(w, err)
		return
	}
	filter, err := transport.Filter(r)
	if err != nil {
		h.fail(w, err)
		return
	}
	now := h.app.Now()
	from, to := filter.From, filter.To
	if from.IsZero() {
		from = now.Add(-24 * time.Hour)
	}
	if to.IsZero() {
		to = now
	}
	result, err := h.app.Analyze(r.Context(), springID, from, to)
	if err != nil {
		h.fail(w, err)
		return
	}
	transport.WriteData(w, http.StatusOK, result)
}
