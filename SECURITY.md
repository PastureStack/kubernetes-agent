# Security Policy

## Supported state

There is no supported PastureStack release yet. `0.7.3` is an unpublished
candidate and must not enter a catalog until its source, image, and Kubernetes
1.12.10 integration gates have passed at the exact candidate commit.

## Security boundaries

- Platform events and Kubernetes objects can contain sensitive topology; do not log complete payloads.
- API keys, kubeconfigs, service-account tokens, certificates, private endpoints, and live metadata must never be committed.
- Preserve bounded health listeners and avoid sharing the default HTTP mux.
- A mounted CA is parsed and selected without modifying the system trust store; an explicitly configured `SSL_CERT_FILE` always wins. Only the documented platform and compatibility mount paths are accepted, and symbolic-link entries are rejected.
- Untrusted event payloads and reply destinations are never logged. Operational identifiers are length-bounded and encoded onto one printable line before reaching the logger.
- The runtime is distroless, contains no shell or package manager, and runs as UID/GID 65532.
- The build fails on reachable Go vulnerabilities and applicable Critical/High source, runtime, and build-environment findings.
- Raw and applicable build-environment findings remain in the evidence bundle. The
  build is intentionally pure Go (`CGO_ENABLED=0`), so compiler and kernel-header
  packages are not installed. OpenVEX is retained as an empty reviewed assertion
  set rather than being used to hide scanner findings.

## Reporting

Report suspected vulnerabilities through this repository's private security advisory channel. Do not include live credentials or cluster data in a public issue.
