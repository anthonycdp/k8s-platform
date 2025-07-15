package handlers

import (
	"net/http"
	"runtime"
	"time"
)

type HealthStatus struct {
	Status  string `json:"status"`
	Time    string `json:"time"`
	Version string `json:"version,omitempty"`
}

type ReadyStatus struct {
	Status     string `json:"status"`
	Time       string `json:"time"`
	Goroutines int    `json:"goroutines"`
	GoVersion  string `json:"goVersion"`
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, HealthStatus{
		Status: "healthy",
		Time:   time.Now().UTC().Format(time.RFC3339),
	})
}

func ReadyHandler(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, ReadyStatus{
		Status:     "ready",
		Time:       time.Now().UTC().Format(time.RFC3339),
		Goroutines: runtime.NumGoroutine(),
		GoVersion:  runtime.Version(),
	})
}
