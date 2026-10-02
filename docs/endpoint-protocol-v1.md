# Endpoint Protocol v1

Stable server-to-server contract between **Stargate Endpoint Manager** and **Stargate Endpoint**.

## Principles

1. TLS is mandatory in production.
2. Browser session cookies are never used.
3. Every endpoint has a unique identity.
4. Enrollment credentials are short-lived and one-time.
5. After enrollment, the endpoint uses a dedicated credential.
6. Every write is idempotent or has an explicit operation ID.
7. Requests are scoped to exactly one endpoint.
8. Responses report observed state, not only command acceptance.
9. Safe unknown fields are ignored for minor-version interoperability.

## Bootstrap

Manager endpoint:

    POST /api/v1/enrollment/bootstrap

Request:

    enrollment_id
    token
    name
    region
    country
    city
    version

The enrollment token is high-entropy, expires after a short period, and can be consumed only once.

On success the Manager returns:

    endpoint_id
    endpoint_credential
    protocol_version
    heartbeat_interval_seconds

The endpoint credential is returned only during bootstrap. The Manager stores only its SHA-256 hash.

## Endpoint authentication

Authenticated Endpoint requests use:

    Authorization: Bearer <endpoint_credential>

The enrollment token must never be reused for heartbeat or operational APIs.

An invalid or missing endpoint credential returns:

    ENDPOINT_AUTH_INVALID

## Heartbeat

Endpoint:

    POST /api/v1/endpoint/heartbeat

The request includes:

    endpoint_id
    version
    capabilities
    cpu_percent
    memory_percent
    disk_percent
    active_sessions

The Manager records the latest observed state and marks the endpoint online.

The current heartbeat interval is 30 seconds.

## Capability identifiers

Examples:

    protocol.xray.anytls
    protocol.xray.tuic5
    protocol.xray.naiveproxy
    protocol.wireguard
    protocol.amneziawg
    feature.real_ip_ssl
    feature.telegram_automation

Operations must be gated by capabilities instead of assuming every endpoint is identical.

## Initial read operations

- endpoint.info.get
- endpoint.capabilities.get
- endpoint.health.get
- accounts.list
- inbounds.list
- account.traffic.get

## Initial write operations

- account.ensure
- account.disable
- account.delete
- membership.ensure
- membership.remove
- inbound.ensure
- inbound.remove

## Stable error codes

AUTH_REQUIRED, AUTH_INVALID, ENDPOINT_AUTH_INVALID, ENDPOINT_NOT_FOUND, UNSUPPORTED_CAPABILITY, VALIDATION_FAILED, RESOURCE_NOT_FOUND, CONFLICT, RATE_LIMITED, INTERNAL_ERROR, TEMPORARY_UNAVAILABLE, STATE_PERSISTENCE_FAILED.

Human-readable messages are supplemental and must not be used for program logic.

## Compatibility

A v1 endpoint may add fields without breaking a v1 manager. The manager must not send an operation unless the endpoint advertises the required capability.
