package main

import (
	"github.com/juliant/distributed_file_server/internal/server"
	"github.com/sirupsen/logrus"
)

func main() {
	port := 8080
	runServer(port)

	logrus.Infof("Server stopped")
}

func runServer(port int) {
	// dependencies
	storage := server.NewInMemoryFileStorage()

	// run the server
	srv := server.New(port, storage)
	err := srv.StartServer()

	if err != nil {
		logrus.WithError(err).Fatal("failed to start server")
	}
}
