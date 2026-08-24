# Source provenance

This maintained compatibility module started from the Apache-2.0 licensed
`github.com/rancher/go-rancher/v3` source at commit
`463b431ca017` (2017-09-28). PastureStack owns subsequent maintenance and the
module path; builds no longer download or import the unmaintained upstream
module.

Only the v3 control-plane wire client required by Kubernetes Agent is retained.
The API schema and wire field names remain protocol contracts and do not imply
an active dependency on Rancher software.
