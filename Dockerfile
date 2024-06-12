FROM golang:1.22

WORKDIR /app

COPY go.* ./
RUN go mod download
COPY . .

RUN make build
CMD ["./bin/complete-auth"]

