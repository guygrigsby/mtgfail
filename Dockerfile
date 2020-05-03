FROM golang:1.13-alpine as build

RUN apk update && apk --no-cache add ca-certificates git && update-ca-certificates

RUN mkdir /go/src/app 
ADD . /go/src/app/
WORKDIR /go/src/app

RUN go build -o app cmd/server/main.go

FROM alpine AS runtime
RUN apk update && apk --no-cache add ca-certificates git && update-ca-certificates

WORKDIR /
RUN wget https://archive.scryfall.com/json/scryfall-default-cards.json 
RUN adduser -S -D -H -h /app appuser
USER appuser
COPY --from=build /go/src/app/* /
EXPOSE 8080
ENTRYPOINT ["/app"]
