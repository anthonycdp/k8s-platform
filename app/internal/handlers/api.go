package handlers

import (
	"net/http"
	"runtime"
	"time"

	"github.com/example/platform/api/internal/config"
)

var (
	appStartTime = time.Now()
	appVersion   string
	appConfig    *config.Config
)

func SetConfig(cfg *config.Config) {
	appConfig = cfg
}

func SetAppInfo(version, env string) {
	appVersion = version
}

type StatusResponse struct {
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
	Uptime    string `json:"uptime"`
	Version   string `json:"version"`
}

type InfoResponse struct {
	Application ApplicationInfo `json:"application"`
	Runtime     RuntimeInfo     `json:"runtime"`
	Kubernetes  KubernetesInfo  `json:"kubernetes"`
}

type ApplicationInfo struct {
	Name        string `json:"name"`
	Environment string `json:"environment"`
	Version     string `json:"version"`
}

type RuntimeInfo struct {
	GoVersion  string `json:"goVersion"`
	OS         string `json:"os"`
	Arch       string `json:"arch"`
	Goroutines int    `json:"goroutines"`
}

type KubernetesInfo struct {
	PodName      string `json:"podName"`
	PodNamespace string `json:"podNamespace"`
	NodeName     string `json:"nodeName"`
}

func StatusHandler(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, StatusResponse{
		Status:    "operational",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Uptime:    time.Since(appStartTime).String(),
		Version:   appVersion,
	})
}

func InfoHandler(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, InfoResponse{
		Application: ApplicationInfo{
			Name:        appConfig.AppName,
			Environment: appConfig.Environment,
			Version:     appVersion,
		},
		Runtime: RuntimeInfo{
			GoVersion:  runtime.Version(),
			OS:         runtime.GOOS,
			Arch:       runtime.GOARCH,
			Goroutines: runtime.NumGoroutine(),
		},
		Kubernetes: KubernetesInfo{
			PodName:      appConfig.PodName,
			PodNamespace: appConfig.PodNamespace,
			NodeName:     appConfig.NodeName,
		},
	})
}

type EchoResponse struct {
	Method  string              `json:"method"`
	Path    string              `json:"path"`
	Headers map[string][]string `json:"headers"`
	Query   map[string][]string `json:"query"`
}

func EchoHandler(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, EchoResponse{
		Method:  r.Method,
		Path:    r.URL.Path,
		Headers: r.Header,
		Query:   r.URL.Query(),
	})
}
