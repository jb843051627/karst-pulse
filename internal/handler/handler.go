package handler

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/karst-pulse/karst-pulse/internal/engine"
	"github.com/karst-pulse/karst-pulse/internal/model"
	"github.com/karst-pulse/karst-pulse/internal/service"
	"github.com/karst-pulse/karst-pulse/internal/transport"
	"github.com/karst-pulse/karst-pulse/internal/validate"
)

type Handler struct {
	app       *service.App
	runtime   *engine.Runtime
	startedAt time.Time
	static    http.Handler
}

func New(app *service.App, runtime *engine.Runtime, startedAt time.Time, static http.Handler) *Handler {
	return &Handler{app: app, runtime: runtime, startedAt: startedAt, static: static}
}

func (h *Handler) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", h.health)
	mux.HandleFunc("GET /metrics", h.metrics)
	mux.HandleFunc("GET /api/v1/springs", h.listSprings)
	mux.HandleFunc("POST /api/v1/springs", h.createSpring)
	mux.HandleFunc("GET /api/v1/springs/{id}", h.getSpring)
	mux.HandleFunc("GET /api/v1/springs/{id}/dashboard", h.dashboard)
	mux.HandleFunc("GET /api/v1/sensors", h.listSensors)
	mux.HandleFunc("POST /api/v1/sensors", h.createSensor)
	mux.HandleFunc("POST /api/v1/readings", h.ingestReading)
	mux.HandleFunc("GET /api/v1/readings", h.listReadings)
	mux.HandleFunc("GET /api/v1/events", h.listEvents)
	mux.HandleFunc("GET /api/v1/batches", h.listBatches)
	mux.HandleFunc("POST /api/v1/batches", h.createBatch)
	mux.HandleFunc("GET /api/v1/batches/{id}/samples", h.listSamples)
	mux.HandleFunc("POST /api/v1/batches/{id}/samples", h.createSample)
	mux.HandleFunc("GET /api/v1/alerts", h.listAlerts)
	mux.HandleFunc("POST /api/v1/alerts/{id}/ack", h.acknowledgeAlert)
	mux.HandleFunc("GET /api/v1/maintenance", h.listMaintenance)
	mux.HandleFunc("POST /api/v1/maintenance", h.createMaintenance)
	mux.HandleFunc("POST /api/v1/maintenance/{id}/complete", h.completeMaintenance)
	mux.HandleFunc("GET /api/v1/analysis", h.analysis)
	if h.static != nil {
		mux.Handle("/", h.static)
	}
	return transport.WithRequestHeaders(mux)
}

func (h *Handler) health(w http.ResponseWriter, r *http.Request) {
	health, err := h.app.Health(r.Context(), model.Health{StartedAt: h.startedAt})
	if err != nil {
		h.fail(w, err)
		return
	}
	transport.WriteData(w, http.StatusOK, health)
}

func (h *Handler) metrics(w http.ResponseWriter, _ *http.Request) {
	transport.WriteData(w, http.StatusOK, h.app.MetricsSnapshot())
}

func (h *Handler) fail(w http.ResponseWriter, err error) {
	transport.WriteError(w, normalizeError(err))
}

func normalizeError(err error) error {
	if err == nil {
		return nil
	}
	var validationErr *validate.Errors
	if errors.As(err, &validationErr) {
		return transport.BadRequest(validationErr.Error(), validationErr)
	}
	if errors.Is(err, sql.ErrNoRows) || strings.Contains(err.Error(), "not found") || strings.Contains(err.Error(), "does not exist") {
		return transport.NotFound("resource not found")
	}
	if strings.Contains(err.Error(), "UNIQUE constraint") || strings.Contains(err.Error(), "is not open") {
		return transport.Conflict("resource state conflicts with the request")
	}
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return transport.BadRequest("request context ended", nil)
	}
	return err
}
