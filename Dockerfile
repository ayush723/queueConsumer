FROM golang:1.16 as builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o consumer .

FROM golang:1.16
WORKDIR /app
COPY --from=builder /app/consumer /app
CMD ["/app/consumer"]