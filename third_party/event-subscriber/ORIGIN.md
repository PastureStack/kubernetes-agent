# Source provenance

This PastureStack-maintained compatibility module started from the Apache-2.0
licensed `rancher/event-subscriber` source at commit
`dddf44b13e15ae38f8bb70a6de0a095a45da91ad`. PastureStack owns subsequent
maintenance and the module path; builds no longer download or import the
unmaintained upstream module.

Only the event and lock packages used by Kubernetes Agent are retained. Event
names, payloads, locking behaviour, and retry behaviour remain wire compatible;
logging now uses bounded fields and standard-library error wrapping.
