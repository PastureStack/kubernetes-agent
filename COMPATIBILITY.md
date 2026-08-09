# Compatibility Contract

The migration preserves the established event names and reply shapes, Kubernetes object mappings, source metadata schema, host status flow, and required legacy labels consumed by existing deployments.

Preferred settings use `PLATFORM_URL`, `PLATFORM_ACCESS_KEY`, `PLATFORM_SECRET_KEY`, and `PLATFORM_METADATA_ADDRESS`. Historical environment names and `io.rancher.*` labels remain temporary protocol and data-contract aliases. They are not PastureStack branding and cannot be removed until every recorded producer and consumer has migrated.

The `github.com/rancher/*` imports and generated `RancherClient` type names are inherited dependency and wire-schema contracts. Two minimal locally patched compatibility modules retain exact upstream provenance and licenses; all unused historical vendoring has been removed.

The candidate uses client-go 0.36.3 with JSON negotiation explicitly selected.
Only stable Core/v1 resources used by this agent are in scope. A real
Kubernetes 1.12.10 API gate must prove service CRUD, namespace and service
watches, node label/status updates, and clean shutdown before release.
The agent explicitly uses list-then-watch startup because Kubernetes 1.12 does
not implement the streaming initial-list semantics enabled by current client-go
releases. Startup waits for each watch cache to complete its first list so that
no events are lost during initialization.

Operator messages support `en-US` and `zh-TW`. Event payloads, Kubernetes resources, labels, identifiers, and API replies remain language-neutral.

Before release, validate event reconnect and reply behavior, namespace and service synchronization, host labels and status, health isolation, redaction, CA bootstrap, and an old/new agent rollout and rollback in the isolated VM environment. No production endpoint is part of the automated gate.
