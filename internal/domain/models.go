package domain

import "time"

type EndpointStatus string

const (
	EndpointPending EndpointStatus = "pending"
	EndpointOnline  EndpointStatus = "online"
	EndpointStale   EndpointStatus = "stale"
	EndpointOffline EndpointStatus = "offline"
)

type Endpoint struct {
	ID                string
	Name              string
	LocationID        string
	Region            string
	Country           string
	City              string
	Status            EndpointStatus
	StargateVersion   string
	ProtocolVersion   int
	LastSeenAt        *time.Time
	Capabilities      []string
	CPUPercent        float64
	MemoryPercent     float64
	DiskPercent       float64
	ActiveSessions    int
}

type Account struct {
	ID          string
	Username    string
	DisplayName string
	Enabled     bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type AccountMembership struct {
	ID             string
	AccountID      string
	EndpointID     string
	LocalClientID  string
	Enabled        bool
	DesiredVersion int64
	ObservedState  string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type Enrollment struct {
	ID         string
	EndpointID string
	TokenHash  string
	ExpiresAt  time.Time
	UsedAt     *time.Time
	CreatedAt  time.Time
}

type Operation struct {
	ID          string
	RequestID   string
	EndpointID  string
	Operation   string
	Status      string
	ErrorCode   string
	CreatedAt   time.Time
	CompletedAt *time.Time
}

type AuditLog struct {
	ID         string
	ActorID    string
	EndpointID string
	Action     string
	CreatedAt  time.Time
}
