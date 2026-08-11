---
tldr: Design local services around typed shared operations, explicit trust boundaries, and tested lifecycle behavior.
---
# Local Service Design

## Public Local Boundary

Treat a local socket, IPC endpoint, command-line client, GUI integration, state
file, executable name, or service unit as a public interface when another
local program, user workflow, or deployment configuration depends on it. Name
the trust boundary: who may connect, what may be read or changed, where endpoint
and state paths come from, and which actions need user approval.

## Cross-Surface Operations

When an application intentionally offers GUI, CLI, and local API access, map
each supported user capability to one typed domain operation and state which
surfaces expose it. Do not let three independently implemented command paths
drift in names, validation, state changes, or error meaning.

## Protocol And Evolution Decisions

Before implementing a durable protocol, record decisions for framing, request
and operation names, typed request/response shape, error shape, versioning and
compatibility negotiation, malformed input, request size limits, and logging.
For a legacy protocol, decide explicitly whether it remains unchanged, is
adapted at a boundary, or runs beside a structured versioned protocol. Do not
silently replace an interface merely because a new format is easier to parse.

## State, Concurrency, And Cancellation

Name one owner for each mutable runtime state and document temporary exceptions
during a migration. Define concurrent-client behavior: serialization, queue
order and capacity, idempotence, cancellation, duplicate requests, timeouts,
and what a client can infer after disconnecting. State transitions should be
observable through a documented interface rather than reconstructed from
unrelated files or UI state.

## Process Lifecycle

Define singleton ownership, endpoint creation and stale-endpoint recovery,
startup ordering, crash and restart semantics, child-process supervision and
cleanup, shutdown deadlines, and version-skew behavior between clients and a
service. Test malformed requests and restart paths as first-class behavior; a
daemon that only works on the happy path is not a compatible service.

## Evidence And Handoff

For each supported operation, retain deterministic fixtures for valid and
malformed requests, response/error shape, and state change. Report decisions,
open questions, public and runtime paths touched, compatibility effects, and
the checks that exercised lifecycle and recovery behavior.

## See Also

Select `personal:engineering/staged-migration` for Python-to-Rust or other
successor work, `cdint:process/decision-first` for unresolved architecture,
and `personal:lang/rust` for Rust implementation practices once the first slice
is defined.
