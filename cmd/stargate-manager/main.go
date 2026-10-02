package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Arshak888/Stargate-endpoint-manager/internal/httpapi"
	"github.com/Arshak888/Stargate-endpoint-manager/internal/store"
)

func main() {
	statePath := os.Getenv("STARGATE_MANAGER_STATE")
	if statePath == "" {
		statePath = "/var/lib/stargate-endpoint-manager/state.json"
	}

	store, err := store.NewPersistentStore(statePath)
	if err != nil {
		log.Fatalf("load manager state: %v", err)
	}

	enrollment := &httpapi.EnrollmentHandler{Store: store}
	heartbeat := &httpapi.HeartbeatHandler{Store: store}

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})
	mux.HandleFunc("/api/v1/enrollment/create", enrollment.Create)
	mux.HandleFunc("/api/v1/enrollment/bootstrap", enrollment.Bootstrap)
	mux.HandleFunc("/api/v1/endpoint/heartbeat", heartbeat.Receive)

	log.Printf("stargate-endpoint-manager listening on :8080 (state: %s)", statePath)
	log.Fatal(http.ListenAndServe(":8080", mux))
}
