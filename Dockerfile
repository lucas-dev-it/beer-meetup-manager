FROM golang:1.14.4 as builder
WORKDIR /go/src/github.com/lucas-dev-it/meetup-manager
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app ./app

FROM alpine
RUN apk --no-cache add ca-certificates
COPY --from=builder /go/src/github.com/lucas-dev-it/meetup-manager/app .
ENTRYPOINT ["./app"]