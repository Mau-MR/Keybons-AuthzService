FROM golang:1.14-alpine AS build

RUN apk add --no-cache git

WORKDIR /keybons

#This has to be run with the context of the parent root directory
COPY . .

#Dependency download
RUN go mod download

#Unit test
#RUN CGO_ENABLED=0 go test -v

RUN go build -o ./temp/authzserver cmd/server/main.go

#new smaller container image
FROM alpine:3.9 AS bin

COPY --from=build /keybons/temp/authzserver /app/authzserver

ENTRYPOINT ["/app/authzserver","-port=9090"]

#The port on which the server runs
EXPOSE 9090
