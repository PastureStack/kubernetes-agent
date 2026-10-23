# PastureStack Kubernetes Agent

Kubernetes Agent synchronizes services, namespaces, host state, and established platform events between a Kubernetes cluster and the compatibility control plane.

PastureStack is an independent community effort to preserve, audit, and modernize the Rancher 1.6 ecosystem. It is not affiliated with or endorsed by Rancher Labs or SUSE.

**Upstream:** [`rancher/kubernetes-agent`](https://github.com/rancher/kubernetes-agent). This GitHub fork preserves the upstream Git history, authorship, dates, tags, and license notices; PastureStack maintenance is consolidated into one commit after the preserved upstream boundary.

## Project status

Published release `v0.7.3` uses Go 1.27.3, Go modules, client-go 0.36.4, and logrus 1.10.1 while
retaining the established event and Kubernetes API contracts. Kubernetes 1.36.4
is the supported and tested API line. Its runtime is a
digest-pinned distroless image that runs as UID/GID 65532. Source, runtime, and
Kubernetes 1.36.4 integration tests must all pass before a release or catalog
reference is created.

The [v0.7.3 release](https://github.com/PastureStack/kubernetes-agent/releases/tag/v0.7.3) publishes the exact checked producer `092f1ce5499c4151d0780591cf9a2b71c6d5267e` from [run 39962596272](https://github.com/PastureStack/kubernetes-agent/actions/runs/39962596272). Its registry digest is `sha256:1224ab0f4f5b621fb9f5de3ae386352e58d5d8d333a1a1bddcd9d78543fdfa57`. User-facing configuration retains the plain version tag:

```text
ghcr.io/pasturestack/kubernetes-agent:v0.7.3
```

## Configuration

Preferred settings are `PLATFORM_URL`, `PLATFORM_ACCESS_KEY`, and `PLATFORM_SECRET_KEY`. Historical `CATTLE_*` names remain accepted as temporary compatibility aliases. `KUBERNETES_URL`, `WORKER_COUNT`, and `HEALTH_CHECK_PORT` keep their established meanings.

`PLATFORM_CA_ROOT` accepts only `/var/lib/pasturestack/etc/ssl/ca.crt` or the established `/var/lib/rancher/etc/ssl/ca.crt` compatibility mount. An explicit `SSL_CERT_FILE` still takes precedence. Arbitrary environment-provided filesystem paths and symbolic links are rejected.

Set `PASTURESTACK_LOCALE=en-US` or `PASTURESTACK_LOCALE=zh-TW` for operator lifecycle messages. Event payloads, Kubernetes objects, labels, and protocol replies are never translated.

## Build and test

Run from a Docker-capable Linux host:

```sh
make test
make validate
make build
make package IMAGE_NAME=pasturestack/kubernetes-agent TAG=poc
```

The package target builds a local image only and does not push, release, or
change a catalog. See [COMPATIBILITY.md](COMPATIBILITY.md),
[SECURITY.md](SECURITY.md), and [ORIGIN.md](ORIGIN.md).

## License and attribution

The inherited project remains licensed under [Apache License 2.0](LICENSE). Copyright and attribution for inherited work and vendored dependencies remain with their respective authors and contributors. PastureStack contributors claim authorship only for their own changes.

## October 2026 security maintenance

Maintained Go builds select the patched Go 1.27.3 compiler without raising the existing Go language requirement. Ubuntu build/runtime stages apply the checksum-pinned `ca-certificates 20260816ubuntu0.26.04.1` security package after their frozen APT baseline; the final trust store and installed version are verified. Release artifacts must pass the repository's existing checks before publication.
