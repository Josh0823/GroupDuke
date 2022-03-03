build:
	go build -o server src/*.go

run:
	source secrets/.env && go run src/*.go

exec:
	source secrets/.env && go build -o server src/*.go && ./server
