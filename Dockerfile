FROM golang:latest

LABEL maintainer="Bryan Porter <bryan@bryanporter.com>"

WORKDIR "/app"

COPY go.mod go.sum storage.cfg ./

RUN go mod download

COPY . .

RUN go build -o main .

EXPOSE 8080

CMD ["./main"]