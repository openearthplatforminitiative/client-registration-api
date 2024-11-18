FROM golang:1.23.2
WORKDIR /go/src/app
COPY go.mod .
COPY go.sum .
COPY cmd cmd
COPY config config
COPY handlers handlers
COPY keycloak keycloak
COPY middleware middleware
COPY models models
COPY routes routes

RUN go build -o main cmd/app/main.go

CMD ["./main"]