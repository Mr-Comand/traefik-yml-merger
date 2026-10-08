FROM golang:1.25-alpine AS build

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY main.go .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o yaml-merger .


FROM alpine:3.22

RUN apk add --no-cache iproute2

COPY --from=build /build/yaml-merger /yaml-merger

COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

EXPOSE 8080

ENTRYPOINT ["/entrypoint.sh"]