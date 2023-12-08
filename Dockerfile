FROM golang:1.21 as build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./

RUN go build
EXPOSE 8080

FROM debian:bookworm

WORKDIR /app

COPY --from=build /app/reversocks .

# Run
CMD ["./reversocks"]
