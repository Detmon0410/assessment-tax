# # Define base image
# ARG GO_VERSION=1.22.2
# FROM golang:${GO_VERSION}-alpine AS build-base

# # Set working directory
# WORKDIR /app

# # Copy Go module files and download dependencies
# COPY go.mod go.sum ./
# RUN go mod download

# # Copy the rest of the application source code
# COPY . .

# # Run tests (if applicable)
# RUN CGO_ENABLED=0 go test -v

# # Build the Go application
# RUN go build -o ./out/go-sample .

# # Define base image for the final stage
# FROM alpine:3.16.2

# # Set environment variables
# ENV APP_DIR=/app
# ENV EXECUTABLE_NAME=go-sample

# # Create directory for the application
# RUN mkdir -p ${APP_DIR}

# # Copy the compiled binary from the build stage
# COPY --from=build-base /app/out/${EXECUTABLE_NAME} ${APP_DIR}/${EXECUTABLE_NAME}

# # Set the working directory
# WORKDIR ${APP_DIR}

# # Specify the command to run the application
# CMD ["./${EXECUTABLE_NAME}"]