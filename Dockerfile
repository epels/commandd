FROM golang:alpine AS builder
ADD . /go/src/github.com/epels/commandd
RUN go install /go/src/github.com/epels/commandd/cmd/commandd

FROM alpine
COPY --from=builder /go/bin/commandd /app/
CMD ["/app/commandd"]
