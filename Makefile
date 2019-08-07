.PHONY: build clean run test-api

build:
	GOOS=linux GOARCH=amd64 go build -o auth main.go

clean: 
	rm -rf ./auth

run:
	go run main.go

test-api:
	GOOS=linux GOARCH=amd64 go build -o auth main.go
	open http://127.0.0.1:8080/healthcheck
	sam local start-api --template template.json --port 8080