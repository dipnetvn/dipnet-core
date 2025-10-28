# Support setting various labels on the final image
ARG COMMIT=""
ARG VERSION=""
ARG BUILDNUM=""

# Build DipNet in a stock Go builder container
FROM golang:1.24-alpine AS builder

RUN apk add --no-cache gcc musl-dev linux-headers git

# Get dependencies - will also be cached if we won't change go.mod/go.sum
COPY go.mod /dipnet-core/
COPY go.sum /dipnet-core/
RUN cd /dipnet-core && go mod download

ADD . /dipnet-core
RUN cd /dipnet-core && go run build/ci.go install -static ./cmd/dipnet

# Pull DipNet into a second stage deploy alpine container
FROM alpine:latest

RUN apk add --no-cache ca-certificates
COPY --from=builder /dipnet-core/build/bin/dipnet /usr/local/bin/

EXPOSE 8545 8546 30303 30303/udp
ENTRYPOINT ["dipnet"]

# Add some metadata labels to help programmatic image consumption
ARG COMMIT=""
ARG VERSION=""
ARG BUILDNUM=""

LABEL commit="$COMMIT" version="$VERSION" buildnum="$BUILDNUM"
