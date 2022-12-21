FROM golang:1.19-alpine AS build_deps
WORKDIR /workspace
COPY go.mod .
COPY go.sum .
RUN go mod download

FROM build_deps AS build
COPY . .
RUN CGO_ENABLED=0 go build -o smarthub-exporter -ldflags '-w -extldflags "-static"' .

FROM alpine:3.9
COPY --from=build /workspace/smarthub-exporter /usr/local/bin/smarthub-exporter
ENTRYPOINT ["smarthub-exporter"]
