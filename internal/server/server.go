package server

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"net"
)

type FileServer struct {
	port int
}

func New(port int) *FileServer {
	return &FileServer{port: port}
}

func (f *FileServer) StartServer() error {

	udpAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", f.port))
	if err != nil {
		return errors.Wrap(err, "error resolving UDP address")
	}

	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		return errors.Wrap(err, "failed to start server on specified port")
	}

	defer func(conn *net.UDPConn) {
		err := conn.Close()
		if err != nil {
			logrus.Warnf("failed to close UDP connection: %v", err)
		}
	}(conn)

	logrus.Infof("starting server on port %d", f.port)

	//clientConns := new(sync.Map)
	// this is for the future

	buffer := make([]byte, 1024)
	for {
		//error reading from
		n, addr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			logrus.WithError(err).Error("error reading from UDP connection")
			continue
		}

		// Print the received message and client address
		message := string(buffer[:n])
		fmt.Printf("Received '%s' from %s\n", message, addr)

		// Respond to the client
		_, err = conn.WriteToUDP([]byte("Message received"), addr)
		if err != nil {
			fmt.Println("Error responding to client:", err)
		}
	}
}
