FROM golang:1.21

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -o sample-app cmd/app/main.go

EXPOSE 8080

CMD ["./sample-app"]