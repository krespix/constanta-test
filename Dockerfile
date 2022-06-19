FROM golang:1.18-alpine

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY ./ ./

# Build
RUN go build -o /constanta-test ./cmd/constanta-test/main.go

# Run
CMD [ "/constanta-test" ]
