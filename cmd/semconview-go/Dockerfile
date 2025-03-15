FROM --platform=$BUILDPLATFORM golang:bookworm AS builder
ARG TARGETARCH
ARG VERSION=unknown

ENV GOTOOLCHAIN=auto
ENV CGO_ENABLED=0
ENV GOARCH=${TARGETARCH}

WORKDIR /opt/app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN  go build -ldflags="-s -w -X github.com/Arthur1/semconview-go.Version=${VERSION}" -o semconview-go ./cmd/semconview-go

FROM gcr.io/distroless/base-debian12:nonroot

COPY --from=builder --chown=nonroot:nonroot /opt/app/semconview-go /semconview-go

ENTRYPOINT ["/semconview-go"]
CMD ["-h"]
