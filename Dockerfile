FROM golang:1.19

WORKDIR /app

COPY . .

RUN go build -o os-api

EXPOSE 8081

CMD ./os-api
