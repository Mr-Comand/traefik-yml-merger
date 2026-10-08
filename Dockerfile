FROM golang:1.24-alpine AS build

WORKDIR /build

COPY go.mod .
RUN go mod download

COPY main.go .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o yaml-merger .


FROM scratch

COPY --from=build /build/yaml-merger /yaml-merger
COPY config.yml /config.yml

EXPOSE 8080

ENTRYPOINT ["/yaml-merger"]