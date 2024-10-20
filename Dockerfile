FROM docker.io/library/golang:1.23 AS builder


WORKDIR /walletd

# Install dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Enable CGO for sqlite3 support
ENV CGO_ENABLED=1 

RUN go generate ./...
RUN go mod tidy
RUN go build -o bin/ -tags='netgo timetzdata' -trimpath -a -ldflags '-s -w -linkmode external -extldflags "-static"'  ./cmd/walletd

FROM docker.io/library/alpine:3

ENV PUID=0
ENV PGID=0
# This is intended to be overridden at runtime via `-e WALLETD_API_PASSWORD=...`
ENV WALLETD_API_PASSWORD=password

# copy binary and prepare data dir.
VOLUME [ "/data" ]
COPY --from=builder /walletd/bin/* /usr/bin/
# API port
EXPOSE 9980/tcp
# RPC port
EXPOSE 9981/tcp

USER ${PUID}:${PGID}

# Dockerfile
ENTRYPOINT ["walletd"]
# Default arguments; will be overridden at runtime if any params are passed.
# note: --dir is required to be set to /data or /tmp
CMD ["--dir", "/data", "--http", ":9980", "-network=komodo", "-index.mode=full", "-debug"]