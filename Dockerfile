# Build the manager binary
FROM golang:1.15 as builder

WORKDIR /workspace
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

# Copy the project source
COPY main.go main.go
COPY Makefile Makefile
COPY hack/ hack/
COPY api/ api/
COPY controllers/ controllers/
COPY pkg/ pkg/
COPY webhook/ webhook/
COPY config/ config/
COPY metrics/ metrics/
COPY console/ console/

# Run tests and linting
RUN make test

# Build
RUN make go-build

# Use distroless as minimal base image to package the manager binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=builder /workspace/bin/manager .
USER 65532:65532

ENTRYPOINT ["/manager"]
