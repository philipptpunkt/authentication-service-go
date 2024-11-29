FROM golang:1.23

# Set the working directory inside the container
WORKDIR /app

# Copy and download modules
COPY /apps/backend/go.mod /apps/backend/go.sum ./
RUN go mod download

# Copy the Go files
COPY /apps/backend/*.go ./

# Build the Go binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /authentication-service-backend

# Expose the port your backend will listen on
EXPOSE 8080

# Run the Go application
CMD ["/authentication-service-backend"]
