FROM golang:1.23.5

WORKDIR /Vypaar-SendMailService

COPY . ./ 

RUN go mod download

RUN go build -o app ./


CMD ["./app"]
