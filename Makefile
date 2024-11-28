
build:
	mkdir -p ./bin
	go build -o ./bin/fileserver cmd/server/main.go

build-client:
	mkdir -p ./bin
	go build -o ./bin/fileagent cmd/client/main.go

run: build
	./bin/fileserver

run-client: build-client
	./bin/fileagent

