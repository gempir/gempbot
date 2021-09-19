
FROM golang:latest
WORKDIR /go/src/github.com/gempir/gempbot
COPY . .
RUN go get ./cmd/bot
WORKDIR /go/src/github.com/gempir/gempbot/cmd/bot
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

FROM alpine:latest  
RUN apk --no-cache add ca-certificates
COPY --from=0 /go/src/github.com/gempir/gempbot/cmd/bot/app .
CMD ["./app"]