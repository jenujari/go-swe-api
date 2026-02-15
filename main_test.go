package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	// "github.com/goforj/godump"
	"github.com/jenujari/go-swe-api/router"
	"github.com/stretchr/testify/assert"
)

func TestPingHandler(t *testing.T) {
	h := router.GetServer().Handler
	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}

	var pr struct {
		Status    string `json:"status"`
		Timestamp string `json:"timestamp"`
		Version   string `json:"version"`
	}
	if err := json.Unmarshal(rr.Body.Bytes(), &pr); err != nil {
		t.Fatalf("invalid json: %v", err)
	}

	// godump.Dump(pr)

	if pr.Status != "ok" {
		t.Fatalf("expected status ok, got %q", pr.Status)
	}

	if _, err := time.Parse(time.RFC3339, pr.Timestamp); err != nil {
		t.Fatalf("invalid timestamp: %v", err)
	}
}

func TestPosHandler(t *testing.T) {
	h := router.GetServer().Handler

	payload := struct {
		Time       string `json:"time"`
		PlanetName string `json:"planetName"`
	}{
		Time:       "2026-01-26T00:00:00Z",
		PlanetName: "Sun",
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("invalid json: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/v1/pos", bytes.NewBuffer(payloadBytes))
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}

	var pr struct {
		Longitude float64 `json:"longitude"`
		Latitude  float64 `json:"latitude"`
		Distance  float64 `json:"distance"`
		SpeedLong float64 `json:"speedLong"`
		SpeedLat  float64 `json:"speedLat"`
		SpeedDist float64 `json:"speedDist"`
	}

	if err := json.Unmarshal(rr.Body.Bytes(), &pr); err != nil {
		t.Fatalf("invalid json: %v", err)
	}

	assert.InDelta(t, 281.808299, pr.Longitude, 0.001)
	assert.InDelta(t, 0.000146, pr.Latitude, 0.001)
	assert.InDelta(t, 0.984579, pr.Distance, 0.001)
	assert.InDelta(t, 1.016591, pr.SpeedLong, 0.001)
	assert.InDelta(t, 0.000021, pr.SpeedLat, 0.001)
	assert.InDelta(t, 0.000104, pr.SpeedDist, 0.001)
	// godump.Dump(pr)
}
