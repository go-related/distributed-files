
build:
	mkdir -p ./bin
	go build -o ./bin/fileserver cmd/server/main.go

run: build
	./bin/fileserver

