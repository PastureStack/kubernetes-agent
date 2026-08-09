FROM gcr.io/distroless/static-debian12@sha256:4b2a093ef4649bccd586625090a3c668b254cfe180dee54f4c94f3e9bd7e381e

ARG IMAGE_VERSION=dev
ARG SOURCE_REVISION=unknown

COPY ./kubernetes-agent /usr/local/bin/kubernetes-agent
COPY ./LICENSE ./SECURITY.md ./ORIGIN.md ./COMPATIBILITY.md /licenses/
COPY ./third-party-licenses/ /licenses/go-modules/

LABEL org.opencontainers.image.title="pasturestack-kubernetes-agent" \
      org.opencontainers.image.description="PastureStack Kubernetes compatibility agent" \
      org.opencontainers.image.version="${IMAGE_VERSION}" \
      org.opencontainers.image.revision="${SOURCE_REVISION}" \
      org.opencontainers.image.licenses="Apache-2.0" \
      org.opencontainers.image.source="https://github.com/PastureStack/kubernetes-agent" \
      io.pasturestack.compatibility.contract="legacy-control-plane-kubernetes-v1"

USER 65532:65532
ENTRYPOINT ["/usr/local/bin/kubernetes-agent"]
