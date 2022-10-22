
FROM golang:latest
WORKDIR /go/src/github.com/gempir/gempbot
COPY . .
RUN go get
WORKDIR /go/src/github.com/gempir/gempbot
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o gempbot main.go

FROM alpine:latest  
RUN apk --no-cache add ca-certificates
COPY --from=0 /go/src/github.com/gempir/gempbot/gempbot .
CMD ["./gempbot"]
EXPOSE 3010