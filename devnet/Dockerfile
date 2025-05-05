# devnet/Dockerfile

FROM alpine:latest

# Install runtime deps (for Go binary)
RUN apk add --no-cache ca-certificates tzdata

# Create a user and working dir
RUN addgroup -S flatgas && adduser -S flatgas -G flatgas
WORKDIR /app

# Copy the built Flatgas binary
COPY devnet/build/inso /usr/local/bin/inso

# Expose RPC and P2P ports
EXPOSE 8545 8546 30303 30303/udp

# Default entrypoint (overridable in docker-compose)
ENTRYPOINT ["inso"]
