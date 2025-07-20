FROM golang:1.24-alpine AS build-stage

WORKDIR /app

COPY cmd/ ./cmd/
COPY internal/ ./internal/

COPY go.mod \ 
	go.sum \ 
	index.html \
	./

RUN ls -laR

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server ./cmd/ride-location-tracker/main.go

FROM gcr.io/distroless/base-debian11 AS release-stage

COPY --from=build-stage /app/server \
	/app/index.html \
	.

EXPOSE 8080

CMD ["./server"]
