FROM golang:1.24-alpine AS builder

RUN apk add --no-cache git

COPY ./go-service /usr/local/src/go-service
WORKDIR /usr/local/src/go-service
COPY go.mod .
COPY go.sum .

# The CGO_ENABLED will allow the code to run on 'scratch' image.
RUN go mod download && \
  CGO_ENABLED=0 go build -o /usr/local/bin/go-service .


FROM gcr.io/distroless/static-debian12:nonroot
LABEL maintainer="neverping@gmail.com"

USER nonroot

COPY --from=builder /usr/local/bin/go-service /usr/local/bin/

ENV SERVICE_ENDPOINT=localhost
ENV REQUEST_TIMEOUT=5
ENV IDLE_TIMEOUT=10
# The app listens on port 4458 by default via code.
ENV PORT=4458
EXPOSE 4458/tcp

CMD ["/usr/local/bin/go-service"]
