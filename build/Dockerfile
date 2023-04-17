FROM golang:1.19.8-alpine AS build

WORKDIR /go/src/github.com/vstdy/otus-highload
COPY . .
RUN go build -o otus-project ./cmd/project/main.go

FROM alpine:3.17.3

RUN apk add --no-cache net-tools

WORKDIR /go/src/github.com/vstdy/otus-highload
COPY --from=build /go/src/github.com/vstdy/otus-highload/otus-project .
COPY ./build/entrypoint.sh .

ENTRYPOINT ["sh", "entrypoint.sh"]
CMD ["./otus-project"]
