FROM golang:1.18-bullseye AS builder

ARG CGO_ENABLED=0
ARG GOOS=linux
ARG GOARCH=amd64
ARG TAG_RELEASE=dev

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download -x

COPY . ./
RUN echo "Tagging build as $TAG_RELEASE" && \
    go build -a -ldflags="-w -s -X 'main.release=$TAG_RELEASE'" -o /bin/server cmd/server/*.go

FROM gcr.io/distroless/static-debian11 as final

COPY --from=builder --chown=nonroot:nonroot /bin/server /app/server
WORKDIR /app
USER nonroot
EXPOSE 8081

CMD ["/app/server"]