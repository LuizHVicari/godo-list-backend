FROM golang:1.25-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /api cmd/api/main.go

FROM gcr.io/distroless/static-debian12

COPY --from=build /api /api

EXPOSE 8080

ENTRYPOINT ["/api"]
