# Source provenance

This PastureStack-maintained compatibility module started from the Apache-2.0
licensed `rancher/go-rancher-metadata` source at commit
`d2103caca5873119ff423d29cba09b4d03cd69b8`. PastureStack owns subsequent
maintenance and the module path; builds no longer download or import the
unmaintained upstream module.

Only the metadata protocol package used by Kubernetes Agent is retained. The
wire schema, polling behaviour, and callback behaviour remain compatible.
