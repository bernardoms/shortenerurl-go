# build stage
FROM golang as builder

WORKDIR /app

ENV CGO_ENABLED=0
ENV GOOS=linux

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build cmd/main.go

# final stage
FROM scratch
COPY --from=builder /app/main /app/

EXPOSE 8080
ENTRYPOINT ["/app/main"]
