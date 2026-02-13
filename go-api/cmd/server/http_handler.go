package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"

	"Incident_Monitoring_Project/internal/store"
)

type Handler struct {
	repo        store.Repository
	mlService   string
	httpClient  *http.Client
}

func NewHandler(repo store.Repository, mlService string) *Handler {
	return &Handler{
		repo:      repo,
		mlService: mlService,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

type IngestLogRequest struct {
	Logs []struct {
		Timestamp *time.Time       `json:"timestamp"`
		Service   string           `json:"service"`
		Level     string           `json:"level"`
		Message   string           `json:"message"`
		Metadata  map[string]any   `json:"metadata"`
	} `json:"logs"`
}

func (h *Handler) IngestLogs(c echo.Context) error {
	var req IngestLogRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid payload"})
	}

	if len(req.Logs) == 0 {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "no logs provided"})
	}

	var logs []store.LogEntry
	now := time.Now().UTC()
	for _, l := range req.Logs {
		ts := now
		if l.Timestamp != nil {
			ts = *l.Timestamp
		}
		metaBytes, _ := json.Marshal(l.Metadata)
		logs = append(logs, store.LogEntry{
			Timestamp: ts,
			Service:   l.Service,
			Level:     l.Level,
			Message:   l.Message,
			Metadata:  string(metaBytes),
		})
	}

	ctx := c.Request().Context()
	if err := h.repo.InsertLogs(ctx, logs); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to store logs"})
	}

	return c.JSON(http.StatusAccepted, echo.Map{"status": "accepted", "count": len(logs)})
}

func (h *Handler) Health(c echo.Context) error {
	ctx, cancel := context.WithTimeout(c.Request().Context(), 3*time.Second)
	defer cancel()

	_, err := h.repo.ListRecentLogs(ctx, 1)
	dbOK := err == nil

	return c.JSON(http.StatusOK, echo.Map{
		"status": "ok",
		"checks": echo.Map{
			"db": dbOK,
		},
	})
}

func (h *Handler) ListIncidents(c echo.Context) error {
	ctx := c.Request().Context()
	incidents, err := h.repo.ListIncidents(ctx, 100)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to list incidents"})
	}
	return c.JSON(http.StatusOK, incidents)
}

func (h *Handler) GetIncidentSummary(c echo.Context) error {
	idStr := c.Param("incident_id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid incident id"})
	}

	ctx := c.Request().Context()
	incident, err := h.repo.GetIncident(ctx, id)
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "incident not found"})
	}

	if incident.Summary != nil && incident.RootCause != nil {
		return c.JSON(http.StatusOK, incident)
	}

	reqBody := map[string]any{
		"incident_id": id,
		"description": incident.Description,
	}
	bodyBytes, _ := json.Marshal(reqBody)

	url := fmt.Sprintf("%s/analyze_incident", h.mlService)
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(bodyBytes))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to create ML request"})
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := h.httpClient.Do(httpReq)
	if err != nil || resp.StatusCode >= 300 {
		return c.JSON(http.StatusBadGateway, echo.Map{"error": "ML service unavailable"})
	}
	defer resp.Body.Close()

	var mlResp struct {
		Summary   string `json:"summary"`
		RootCause string `json:"root_cause"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&mlResp); err != nil {
		return c.JSON(http.StatusBadGateway, echo.Map{"error": "invalid ML response"})
	}

	if err := h.repo.UpdateIncidentSummary(ctx, id, mlResp.Summary, mlResp.RootCause); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to save summary"})
	}

	incident.Summary = &mlResp.Summary
	incident.RootCause = &mlResp.RootCause

	return c.JSON(http.StatusOK, incident)
}

