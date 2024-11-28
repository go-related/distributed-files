package server

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"net"
)

type FileStorage interface {
	Save(sc []byte) error
}

type FileServer struct {
	port    int
	storage FileStorage
}

func New(port int, fs FileStorage) *FileServer {
	return &FileServer{port: port, storage: fs}
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

	f.ReadMessageHandler(conn)
	return nil
}

func (f *FileServer) ReadMessageHandler(conn *net.UDPConn) {
	buffer := make([]byte, 1024)
	for {
		// read Message
		n, addr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			logrus.WithError(err).Error("error reading from UDP connection")
			continue
		}

		err = f.storage.Save(buffer[:n])
		if err != nil {
			logrus.WithError(err).Error("error storing file")
		}

		// respond to the client
		_, err = conn.WriteToUDP([]byte("Message received"), addr)
		if err != nil {
			logrus.WithError(err).Error("error responding to client")
		}
	}
}
