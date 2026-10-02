# Stargate Endpoint Manager Architecture

## System boundary

### Stargate UI endpoint
Owns VPN/Xray/SSH runtime state, local protocol configuration, local client credentials, local accounting, nftables/kernel state and endpoint-local administration.

### Endpoint Manager
Owns global endpoint inventory, locations, global account identity, endpoint membership, desired state, fleet health history, remote jobs, audit events and manager-side permissions.

The manager must never require direct access to an endpoint SQLite database.

## Source of truth

| Data | Authority |
|---|---|
| Global account identity | Manager |
| Endpoint membership | Manager |
| Desired multi-location assignment | Manager |
| Local runtime state | Endpoint |
| Live sessions | Endpoint |
| nftables counters | Endpoint |
| Daemon health | Endpoint |
| Observed capabilities | Endpoint |
| Fleet audit/job state | Manager |

The manager stores both desired and observed state.

## Enrollment

1. Administrator creates an enrollment request.
2. Manager issues a short-lived, single-use token.
3. Endpoint presents the token over TLS.
4. Manager creates the endpoint identity.
5. Endpoint receives permanent credentials.
6. Enrollment token is invalidated.
7. Endpoint starts heartbeat and capability reporting.

The enrollment token is never used for normal operations.

## Transport

Versioned HTTPS/JSON is the first transport. Browser sessions are not involved.

Every mutating request contains request_id, endpoint_id, operation, protocol_version, timestamp and payload. Responses contain request_id, operation_id, status, observed_at and result/error.

A later mTLS transport can be introduced without changing the domain model.

## Desired versus observed

The manager never assumes that command acceptance means the endpoint reached the desired state. Partial failures remain visible.

Example:

    Manager desired: account alice, locations Germany + Finland, protocols AnyTLS + WireGuard
    Germany observed: AnyTLS active, WireGuard active
    Finland observed: AnyTLS pending, WireGuard unsupported

## Idempotency

Remote writes use unique request IDs and declarative ensure semantics where possible: ensure account, ensure membership, ensure inbound, ensure client configuration.

## Offline behavior

When an endpoint disappears, its existing runtime continues. The manager marks it stale/offline after the heartbeat timeout, performs no destructive inference, retains desired state, and resumes reconciliation after reconnect.

## Multi-location account model

A global account is independent from any single endpoint. A membership links the account to an endpoint and carries endpoint-specific identity, enabled protocols, limits, desired state and observed state.

This prevents a Stargate email/client identity from becoming the global primary key.

## Future subscription layer

The manager may later generate a global subscription containing multiple locations. It consumes endpoint availability metadata and must not proxy user traffic itself.

## Reseller scope

Scopes can be all endpoints, selected regions, selected endpoints or selected accounts. The manager resolves authorization before sending operations.

## Observability

Heartbeat data should include version, uptime, CPU, RAM, disk, load, capabilities, active sessions, traffic samples, certificate state and reconciliation status.

## Implementation phases

Phase 0: architecture and protocol contract.
Phase 1: database models, enrollment, heartbeat and capability discovery.
Phase 2: endpoint adapter and read-only fleet inventory.
Phase 3: account and membership synchronization.
Phase 4: declarative remote operations and reconciliation.
Phase 5: multi-location subscriptions, reseller scopes and fleet analytics.