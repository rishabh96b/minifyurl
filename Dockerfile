FROM golang:alpine AS builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64


WORKDIR /build

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o minifyurl

WORKDIR /dist

# Copy minifyurl binary from build to dist dir
RUN cp /build/minifyurl .

# Using multi stage docker build process
FROM scratch

COPY --from=builder /dist/minifyurl /

ENTRYPOINT ["/minifyurl"]