FROM golang:alpine AS builder
ADD . /go/src/github.com/epels/uptimed
RUN go install /go/src/github.com/epels/uptimed/cmd/uptimed

FROM alpine
COPY --from=builder /go/bin/uptimed /app/
CMD ["/app/uptimed"]
