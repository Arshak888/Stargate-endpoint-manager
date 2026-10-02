package httpapi

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/Arshak888/Stargate-endpoint-manager/internal/auth"
	"github.com/Arshak888/Stargate-endpoint-manager/internal/domain"
	"github.com/Arshak888/Stargate-endpoint-manager/internal/store"
)

type EnrollmentHandler struct {
	Store *store.MemoryStore
}

func randomID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	return hex.EncodeToString(b)
}

func (h *EnrollmentHandler) Create(w http.ResponseWriter, r *http.Request) {
	token, hash, err := auth.GenerateToken()
	if err != nil {
		http.Error(w, "token generation failed", http.StatusInternalServerError)
		return
	}

	now := time.Now().UTC()
	e := domain.Enrollment{
		ID:         randomID(),
		EndpointID: randomID(),
		TokenHash:  hash,
		ExpiresAt:  now.Add(15 * time.Minute),
		CreatedAt:  now,
	}
	if err := h.Store.PutEnrollment(e); err != nil {
		http.Error(w, "state persistence failed", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusCreated, map[string]interface{}{
		"enrollment_id": e.ID,
		"endpoint_id":   e.EndpointID,
		"token":         token,
		"expires_at":    e.ExpiresAt,
	})
}

func (h *EnrollmentHandler) Bootstrap(w http.ResponseWriter, r *http.Request) {
	var req map[string]interface{}
	if json.NewDecoder(r.Body).Decode(&req) != nil {
		writeError(w, http.StatusBadRequest, "VALIDATION_FAILED")
		return
	}

	id, _ := req["enrollment_id"].(string)
	token, _ := req["token"].(string)
	e, ok := h.Store.GetEnrollment(id)
	if !ok || e.TokenHash != auth.HashToken(strings.TrimSpace(token)) {
		writeError(w, http.StatusUnauthorized, "AUTH_INVALID")
		return
	}

	consumed, err := h.Store.ConsumeEnrollment(id, time.Now().UTC())
	if err != nil {
		writeError(w, http.StatusUnauthorized, "ENROLLMENT_EXPIRED")
		return
	}

	credential, credentialHash, err := auth.GenerateToken()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "CREDENTIAL_GENERATION_FAILED")
		return
	}

	ep := domain.Endpoint{
		ID:              consumed.EndpointID,
		Name:            asString(req["name"]),
		Region:          asString(req["region"]),
		Country:         asString(req["country"]),
		City:            asString(req["city"]),
		Status:          domain.EndpointOnline,
		ProtocolVersion: 1,
		CredentialHash:  credentialHash,
		StargateVersion: asString(req["version"]),
	}
	if err := h.Store.PutEndpoint(ep); err != nil {
		writeError(w, http.StatusInternalServerError, "STATE_PERSISTENCE_FAILED")
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"endpoint_id":                consumed.EndpointID,
		"endpoint_credential":        credential,
		"protocol_version":            1,
		"heartbeat_interval_seconds": 30,
	})
}

func asString(v interface{}) string {
	s, _ := v.(string)
	return s
}
