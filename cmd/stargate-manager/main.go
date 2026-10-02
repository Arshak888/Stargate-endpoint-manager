package main

import (
 "log"
 "net/http"
 "github.com/Arshak888/Stargate-endpoint-manager/internal/httpapi"
 "github.com/Arshak888/Stargate-endpoint-manager/internal/store"
)

func main() {
 store := store.NewMemoryStore()
 enrollment := &httpapi.EnrollmentHandler{Store: store}
 heartbeat := &httpapi.HeartbeatHandler{Store: store}
 mux := http.NewServeMux()
 mux.HandleFunc("/healthz", func(w http.ResponseWriter,r *http.Request) { w.Header().Set("Content-Type","application/json"); w.Write([]byte("{"status":"ok"}")) })
 mux.HandleFunc("/api/v1/enrollment/create", enrollment.Create)
 mux.HandleFunc("/api/v1/enrollment/bootstrap", enrollment.Bootstrap)
 mux.HandleFunc("/api/v1/endpoint/heartbeat", heartbeat.Receive)
 log.Println("stargate-endpoint-manager listening on :8080")
 log.Fatal(http.ListenAndServe(":8080",mux))
}
