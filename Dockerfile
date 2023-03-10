###########################
# STEP 1 build executable binary
############################
FROM golang:alpine AS builder
# Install git.
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git ca-certificates

WORKDIR /app/src
COPY src /app/src
RUN go build -o /app/main


############################
# STEP 2 build a small image
############################
FROM scratch
# Copy our static executable.
COPY --from=builder /app/main /app/main
# Copy certs
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
# Run the binary
ENTRYPOINT ["/app/main"]


