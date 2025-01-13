FROM golang:1.22-alpine

WORKDIR /app

# Copy the entire project directory
COPY . .

# Download dependencies
RUN go mod tidy

# Build the application
RUN go build -o app .

EXPOSE 8080

CMD ["./app"]