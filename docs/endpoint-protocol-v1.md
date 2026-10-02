# Endpoint Protocol v1

First stable contract between Stargate Endpoint Manager and a Stargate UI endpoint.

## Principles

1. TLS is mandatory.
2. Browser session cookies are never used.
3. Every endpoint has a unique identity.
4. Every write is idempotent or has an explicit operation ID.
5. Requests are scoped to exactly one endpoint.
6. Responses report observed state, not only command acceptance.
7. Safe unknown fields are ignored for minor-version interoperability.

## Request envelope

    protocol_version
    request_id
    endpoint_id
    operation
    sent_at
    payload

## Response envelope

    protocol_version
    request_id
    operation_id
    accepted
    status
    observed_at
    result or error

## Enrollment

Manager creates enrollment_id, endpoint_id, expiry and a high-entropy one-time token. After successful enrollment the token is invalidated and the endpoint receives permanent credentials and heartbeat settings.

## Heartbeat

Heartbeat reports Stargate version, protocol version, uptime, CPU, memory, disk, load, capabilities and active sessions.

## Capability identifiers

Examples: protocol.xray.anytls, protocol.xray.tuic5, protocol.xray.naiveproxy, protocol.wireguard, protocol.amneziawg, feature.real_ip_ssl, feature.telegram_automation.

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

AUTH_REQUIRED, AUTH_INVALID, ENDPOINT_NOT_FOUND, UNSUPPORTED_CAPABILITY, VALIDATION_FAILED, RESOURCE_NOT_FOUND, CONFLICT, RATE_LIMITED, INTERNAL_ERROR, TEMPORARY_UNAVAILABLE.

Human-readable messages are supplemental and must not be used for program logic.

## Compatibility

A v1 endpoint may add fields without breaking a v1 manager. The manager must not send an operation unless the endpoint advertises the required capability.