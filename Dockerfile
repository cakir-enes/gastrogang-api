FROM golang:1.12-alpine as builder

# To fix go get and build with cgo
RUN apk add --no-cache --virtual .build-deps \
    bash \
    gcc \
    git \
    musl-dev

RUN mkdir build
COPY . /build
WORKDIR /build

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o webserver ./cmd
RUN adduser -S -D -H -h /build webserver
USER webserver

FROM scratch
COPY --from=builder /build/webserver /app/
WORKDIR /app
EXPOSE 8080
CMD ["./webserver"]