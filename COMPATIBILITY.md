# Compatibility Contract

The migration preserves the established event names and reply shapes, Kubernetes object mappings, source metadata schema, host status flow, and required legacy labels consumed by existing deployments.

Preferred settings use `PLATFORM_URL`, `PLATFORM_ACCESS_KEY`, `PLATFORM_SECRET_KEY`, and `PLATFORM_METADATA_ADDRESS`. Historical environment names and `io.rancher.*` labels remain temporary protocol and data-contract aliases. They are not PastureStack branding and cannot be removed until every recorded producer and consumer has migrated.

PastureStack-owned control API, event, and metadata modules replace all active
`github.com/rancher/*` dependencies. Generated `RancherClient` type names and
`io.rancher.*` data fields remain wire-schema contracts only. The local modules
retain exact upstream provenance and licenses; builds do not download the
unmaintained projects.

The candidate uses client-go 0.36.4 with JSON negotiation explicitly selected.
Only stable Core/v1 resources used by this agent are in scope. A real
Kubernetes 1.36.4 API gate must prove service CRUD, namespace and service
watches, and clean shutdown before release. Focused host-watch tests cover node
label and status reconciliation. Current client-go streaming initial-list
semantics remain enabled; startup waits for each watch cache to synchronize so
that no events are lost during initialization.

Operator messages support `en-US` and `zh-TW`. Event payloads, Kubernetes resources, labels, identifiers, and API replies remain language-neutral.

Release validation is proportional to the changed behavior: the automated gates
cover control-event handling, namespace and service synchronization, host labels
and status, health isolation, redaction, CA bootstrap, and the supported
Kubernetes API. A rollout or rollback exercise is required only when deployment
behavior changes.
