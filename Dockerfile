
FROM golang:latest
WORKDIR /go/src/github.com/gempir/gempbot
COPY . .
RUN go get ./cmd/server
WORKDIR /go/src/github.com/gempir/gempbot
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app cmd/server/main.go

FROM alpine:latest  
RUN apk --no-cache add ca-certificates
COPY --from=0 /go/src/github.com/gempir/gempbot/app .
CMD ["./app"]
EXPOSE 3010