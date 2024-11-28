
build:
	mkdir -p ./bin
	go build -o ./bin/fileserver cmd/server/main.go

build-client:
	mkdir -p ./bin
	go build -o ./bin/fileagent cmd/client/main.go

run: build
	chmod +x ./bin/fileserver
	./bin/fileserver

run-client: build-client
	chmod +x ./bin/fileagent
	./bin/fileagent

