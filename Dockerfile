FROM golang:1.25.4-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o sse-app cmd/sse/main.go

FROM alpine:latest
WORKDIR /root/
COPY --from=build /app/sse-app .

EXPOSE 8081
CMD ["./sse-app"]
