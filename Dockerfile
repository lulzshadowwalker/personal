FROM golang:1.24.2-bullseye AS build-base

WORKDIR /app

COPY go.mod ./

RUN [ -f "go.sum" ] && cp go.sum . || echo "no go.sum file found. skipping."

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go mod download

FROM build-base AS dev

RUN go install github.com/air-verse/air@latest && \
    go install github.com/go-delve/delve/cmd/dlv@latest

COPY . .

CMD ["air", "-c", ".air.toml"]

FROM build-base AS prod

RUN useradd -u 1001 nonroot

COPY . .

RUN go build \
  -ldflags="-linkmode external -extldflags -static" \
  -tags netgo \
  -o main

FROM scratch

WORKDIR /

COPY --from=prod /etc/passwd /etc/passwd

COPY --from=prod /app/main main

USER nonroot

EXPOSE 8272

CMD ["/main"]
