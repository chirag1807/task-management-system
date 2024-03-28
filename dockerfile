FROM golang:1.22-alpine

RUN wget -O /usr/local/bin/dbmate https://github.com/amacneil/dbmate/releases/download/v1.12.0/dbmate-linux-amd64
RUN chmod +x /usr/local/bin/dbmate

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o cmd/main ./cmd

WORKDIR /app/cmd

ADD https://github.com/ufoscout/docker-compose-wait/releases/download/2.5.0/wait /wait
RUN chmod +x /wait

EXPOSE 9090

CMD /wait && dbmate -d /app/db/migrations up && "/app/cmd/main"
