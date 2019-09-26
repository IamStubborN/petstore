# build stage
FROM golang:alpine as builder
RUN apk add git
RUN apk add make
WORKDIR /petstore/
COPY . .
ENV GO111MODULE=on
ENV CGO_ENABLED=0
RUN make all

# final stage
FROM alpine:latest
WORKDIR /petstore/
COPY --from=builder /petstore/db/migrations db/migrations
COPY --from=builder /petstore/petstore .
COPY --from=builder /petstore/config.yaml .

ENTRYPOINT ["./petstore"]