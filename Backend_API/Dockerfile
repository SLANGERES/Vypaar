FROM golang:1.23.5

WORKDIR /Vypaar

COPY . ./ 

RUN go mod download

RUN go build -o app ./cmd/vypaar

EXPOSE 9090

CMD ["./app"]
