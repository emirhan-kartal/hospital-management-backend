# Step 1: Use an official Golang image as the base image
FROM golang:1.20-alpine AS builder

# Step 2: Enable CGO
ENV CGO_ENABLED=1
ENV GOOS=linux

# Step 3: Install necessary packages for CGO and SQLite3
RUN apk add --no-cache gcc musl-dev sqlite-dev

# Step 4: Set the Current Working Directory inside the container
WORKDIR /app

# Step 5: Copy go.mod and go.sum files from /emir/hospital/ to the container
COPY emir/hospital/go.mod emir/hospital/go.sum emir/hospital/.env ./
RUN ls -la

# Step 6: Download all dependencies
RUN go mod download

# Step 7: Copy the entire source code from /emir/hospital/ to the container
COPY emir/hospital/ .

# Step 8: Build the Go app with CGO enabled
RUN go build -o hospital-api .

# Step 9: Use a smaller base image to reduce the size of the final image
FROM alpine:latest

# Step 10: Install Redis
RUN apk add --no-cache redis

# Step 11: Set the Current Working Directory inside the container
WORKDIR /root/

# Step 12: Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/hospital-api .

# Step 13: Expose port 3000 for the Go app and 6379 for Redis
EXPOSE 3000 6379
RUN echo "JWT_SECRET=8cbc7d4e8ded12eaacf33ba241ba7112e1287a0dcb2fbec727d7e5b622ce7ef5" > .env

# Step 14: Command to start Redis and then your Go app
CMD ["sh", "-c", "redis-server & ./hospital-api"]
