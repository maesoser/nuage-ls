FROM golang:alpine as builder

# Copy the code from the host
WORKDIR /app/
COPY . .

# Download and install dependencies
RUN apk update \
        && apk upgrade \
        && apk add --no-cache git \

# Compile it
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o /app/nsgls .

# Create docker
FROM scratch
COPY --from=builder /app/nsgls /app/
ENTRYPOINT ["/app/nsgls", "-h"]