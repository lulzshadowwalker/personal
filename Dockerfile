FROM golang:1.24.2-bullseye AS build-base

WORKDIR /app

COPY go.mod ./

RUN [ -f "go.sum" ] && cp go.sum . || echo "no go.sum file found. skipping."

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go mod download

FROM build-base AS dev

RUN curl -sLo /usr/local/bin/tailwindcss https://github.com/tailwindlabs/tailwindcss/releases/latest/download/tailwindcss-linux-x64 && \
    chmod +x /usr/local/bin/tailwindcss

RUN go install github.com/air-verse/air@latest && \
    go install github.com/go-delve/delve/cmd/dlv@latest && \
    go install github.com/a-h/templ/cmd/templ@latest

COPY . .

CMD ["make", "dev"]

FROM ghcr.io/a-h/templ:latest AS templ-gen

COPY --chown=65532:65532 . /app

WORKDIR /app

RUN ["templ", "generate"]

FROM dev AS tailwindcss-gen

COPY --chown=65532:65532 . /app

WORKDIR /app

RUN ["make", "tailwind"]

FROM build-base AS prod

RUN useradd -u 1001 nonroot

COPY --from=templ-gen /app /app

COPY --from=tailwindcss-gen /app/cmd/http/public/ /app/cmd/http/public/

RUN go build \
  -ldflags="-linkmode external -extldflags -static" \
  -tags netgo \
  -o http \
  ./cmd/http/main.go

FROM build-base AS test

RUN go test -v ./...

FROM scratch

WORKDIR /

COPY --from=prod /etc/passwd /etc/passwd

COPY --from=prod /app/http http

COPY --from=prod /app/cmd/http/public/ /cmd/http/public/

USER nonroot

EXPOSE 3000

CMD ["/http"]

