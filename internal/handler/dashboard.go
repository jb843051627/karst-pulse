package handler

import (
	"net/http"

	"github.com/karst-pulse/karst-pulse/internal/transport"
)

func (h *Handler) dashboard(w http.ResponseWriter, r *http.Request) {
	springID, err := transport.ID(r, "id")
	if err != nil {
		h.fail(w, err)
		return
	}
	result, err := h.app.RecentDashboard(r.Context(), springID)
	if err != nil {
		h.fail(w, err)
		return
	}
	transport.WriteData(w, http.StatusOK, result)
}
