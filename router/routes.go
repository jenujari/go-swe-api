package router

import (
	"encoding/json"
	"net/http"
	"time"

	lib "github.com/jenujari/go-swe-api/lib"
	baselib "github.com/jenujari/planets-lib"
)

func SetRoutes(r *http.ServeMux) {
	r.HandleFunc("GET /ping", pong)
	r.HandleFunc("POST /api/v1/pos", pos)
	r.HandleFunc("POST /api/v1/conj", conjunction)
}

func pos(w http.ResponseWriter, r *http.Request) {
	type Payload struct {
		Time       string `json:"time"`
		PlanetName string `json:"planetName"`
	}

	defer r.Body.Close()

	var p Payload
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	t, err := time.Parse(time.RFC3339, p.Time)
	if err != nil {
		http.Error(w, "invalid time format should be RFC3339", http.StatusBadRequest)
		return
	}

	siderealTime, err := lib.UTCToSiderealTime(t)
	if err != nil {
		http.Error(w, "error converting UTC to sidereal time", http.StatusBadRequest)
		return
	}

	if p.PlanetName == "" {
		var fullResp = make(map[string]*baselib.PlanetCord)

		for planet := range baselib.PLANET_LIB_MAP {
			planetCord, err := lib.GetPlanetCalculation(siderealTime, planet)
			if err != nil {
				http.Error(w, "error calculating planet position", http.StatusBadRequest)
				return
			}

			fullResp[planet] = planetCord
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(fullResp)

	} else {
		planetCord, err := lib.GetPlanetCalculation(siderealTime, p.PlanetName)
		if err != nil {
			http.Error(w, "error calculating planet position", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(planetCord)
	}
}

func conjunction(w http.ResponseWriter, r *http.Request) {
	type Payload struct {
		StartTime string  `json:"start"`
		EndTime   string  `json:"end"`
		Planet1   string  `json:"planet1"`
		Planet2   string  `json:"planet2"`
		Orb       int     `json:"orb"`
		Step      float64 `json:"step"`
	}

	defer r.Body.Close()

	var p Payload
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	startTime, err := time.Parse(time.RFC3339, p.StartTime)
	if err != nil {
		http.Error(w, "invalid time format should be RFC3339", http.StatusBadRequest)
		return
	}

	endTime, err := time.Parse(time.RFC3339, p.EndTime)
	if err != nil {
		http.Error(w, "invalid time format should be RFC3339", http.StatusBadRequest)
		return
	}

	startConj, endConj, found := lib.FindConjunctionRange(startTime, endTime, float64(p.Orb), p.Step, baselib.PLANET_LIB_MAP[p.Planet1], baselib.PLANET_LIB_MAP[p.Planet2])

	if !found {
		http.Error(w, "no conjunction found", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(struct {
		Start time.Time `json:"start"`
		End   time.Time `json:"end"`
	}{
		Start: startConj,
		End:   endConj,
	})
}

func pong(w http.ResponseWriter, r *http.Request) {

	type pingResp struct {
		Status    string `json:"status"`
		Timestamp string `json:"timestamp"`
		Version   string `json:"version"`
	}

	w.Header().Set("Content-Type", "application/json")
	resp := pingResp{
		Status:    "ok",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Version:   lib.GetVersion(),
	}
	_ = json.NewEncoder(w).Encode(resp)
}
