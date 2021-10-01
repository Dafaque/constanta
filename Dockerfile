FROM golang:1.16-alpine as builder
WORKDIR /build
COPY . .
RUN go build .

FROM alpine:latest
WORKDIR /app
COPY --from=builder /build/constanta constanta
ENTRYPOINT [ "constanta" ]