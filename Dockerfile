
FROM golang:latest
WORKDIR /go/src/github.com/gempir/gempbot
COPY . .
RUN go get ./server
WORKDIR /go/src/github.com/gempir/gempbot/server
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app main.go

FROM alpine:latest  
RUN apk --no-cache add ca-certificates
COPY --from=0 /go/src/github.com/gempir/gempbot/server/app .
CMD ["./app"]