ARG NAME=task-manager
ARG SOURCEROOT=/go/src/${NAME}

# Builder Image
FROM golang:1.22-alpine as builder

ARG NAME
ARG SOURCEROOT

COPY . ${SOURCEROOT}
WORKDIR ${SOURCEROOT}

RUN wget -O /usr/local/bin/dbmate https://github.com/amacneil/dbmate/releases/download/v1.12.0/dbmate-linux-amd64
RUN chmod +x /usr/local/bin/dbmate
RUN go mod download
RUN GOOS=linux go build -o bin/${NAME} cmd/main.go

# Runner Image
FROM alpine:latest
ARG NAME
ARG SOURCEROOT
WORKDIR /usr/bin

ADD https://github.com/ufoscout/docker-compose-wait/releases/download/2.5.0/wait /wait
RUN chmod +x /wait
RUN apk update && apk add bash && apk --no-cache add ca-certificates && rm -rf /var/cache/apk/*

COPY --from=builder /usr/local/bin/dbmate /usr/local/bin/dbmate
COPY --from=builder ${SOURCEROOT}/bin/${NAME} /usr/bin/
COPY --from=builder ${SOURCEROOT}/.config /usr/.config
COPY --from=builder ${SOURCEROOT}/db/migrations /usr/migrations

CMD /wait && dbmate -d /usr/migrations up && task-manager