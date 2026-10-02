# Stargate Endpoint Manager

Central control plane for Stargate UI endpoints.

Goals:
- register and monitor many Stargate endpoints
- organize endpoints by location, region and tags
- maintain global account identity across endpoints
- synchronize desired state without sharing endpoint databases
- collect health, traffic and capacity telemetry
- provide auditable idempotent remote operations
- support reseller scopes across selected endpoints

Architecture: Manager is the global control plane; each Stargate installation remains authoritative for its local runtime, protocol daemons, Xray configuration, nftables state and local accounting.

Security: browser panel cookies are never reused. Endpoints enroll with short-lived tokens and receive unique server-to-server credentials.

Status: Phase 0, architecture and protocol contract.