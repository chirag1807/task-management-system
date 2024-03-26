FROM golang:1.22.0

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o cmd/main ./cmd

WORKDIR /app/cmd

EXPOSE 9090

CMD [ "/app/cmd/main" ]