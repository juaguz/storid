# Etapa de compilación para ambos binarios
FROM golang:1.22 as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go  build -tags local -o bin/importer ./cmd/importer/cli/main.go

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go  build -tags local -o bin/sender ./cmd/sender/cli/main.go

FROM alpine:latest as importer
WORKDIR /app
COPY --from=builder /app/bin/importer /app/importer

CMD ["/app/importer"]

FROM alpine:latest as sender
WORKDIR /app
COPY --from=builder /app/bin/sender /app/sender

CMD ["/app/sender"]
