package main

import (
	"fmt"
	"github.com/juliant/distributed_file_server/internal/server"
	"github.com/sirupsen/logrus"
)

func main() {
	port := 8080
	runServer(port)

	fmt.Println("Server stopped")
}

func runServer(port int) {
	srv := server.New(port)
	err := srv.StartServer()
	if err != nil {
		logrus.WithError(err).Fatal("failed to start server")
	}
}
