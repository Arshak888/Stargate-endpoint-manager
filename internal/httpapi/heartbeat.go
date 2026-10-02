package httpapi

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/Arshak888/Stargate-endpoint-manager/internal/domain"
	"github.com/Arshak888/Stargate-endpoint-manager/internal/store"
)

type HeartbeatHandler struct {
	Store *store.MemoryStore
}

func (h *HeartbeatHandler) Receive(w http.ResponseWriter, r *http.Request) {
	var req map[string]interface{}
	if json.NewDecoder(r.Body).Decode(&req) != nil {
		writeError(w, http.StatusBadRequest, "VALIDATION_FAILED")
		return
	}

	id, _ := req["endpoint_id"].(string)
	credential := bearerCredential(r.Header.Get("Authorization"))
	if !h.Store.AuthenticateEndpoint(id, credential) {
		writeError(w, http.StatusUnauthorized, "ENDPOINT_AUTH_INVALID")
		return
	}

	ep, ok := h.Store.GetEndpoint(id)
	if !ok {
		writeError(w, http.StatusNotFound, "ENDPOINT_NOT_FOUND")
		return
	}

	now := time.Now().UTC()
	ep.Status = domain.EndpointOnline
	ep.LastSeenAt = &now

	if v, ok := req["version"].(string); ok {
		ep.StargateVersion = v
	}
	if v, ok := req["capabilities"].([]interface{}); ok {
		ep.Capabilities = nil
		for _, x := range v {
			if s, ok := x.(string); ok {
				ep.Capabilities = append(ep.Capabilities, s)
			}
		}
	}
	if v, ok := req["cpu_percent"].(float64); ok {
		ep.CPUPercent = v
	}
	if v, ok := req["memory_percent"].(float64); ok {
		ep.MemoryPercent = v
	}
	if v, ok := req["disk_percent"].(float64); ok {
		ep.DiskPercent = v
	}
	if v, ok := req["active_sessions"].(float64); ok {
		ep.ActiveSessions = int(v)
	}

	h.Store.PutEndpoint(ep)
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"status":                     "ok",
		"observed_at":                now,
		"heartbeat_interval_seconds": 30,
	})
}

func bearerCredential(header string) string {
	const prefix = "Bearer "
	if !strings.HasPrefix(header, prefix) {
		return ""
	}
	return strings.TrimSpace(strings.TrimPrefix(header, prefix))
}
