FROM golang:1.25-alpine AS build

WORKDIR /build

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY main.go .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o yaml-merger .


FROM scratch

COPY --from=build /build/yaml-merger /yaml-merger

EXPOSE 8080

ENTRYPOINT ["/yaml-merger"]