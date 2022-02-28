FROM golang:1.17

WORKDIR /usr/src/app

RUN apt-get update && apt-get install sqlite3 libsqlite3-dev

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o /usr/local/bin/server ./...
RUN . ./.env

ENV DB_STRING=./courses.sqlite

ENV PORT=4200
EXPOSE 4200

CMD ["server"]
