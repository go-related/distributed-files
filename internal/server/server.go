package server

import (
	"encoding/binary"
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
	for {
		// read Message
		var fileSize uint32
		err := binary.Read(conn, binary.LittleEndian, &fileSize)
		if err != nil {
			logrus.WithError(err).Error("pls send the size of the file before sending the actual file")
			continue
		}
		receivedBytes := uint32(0)
		chunkSize := 1024
		fileBuffer := make([]byte, 0, fileSize)
		var clientAddr *net.UDPAddr
		counter := 0
		for receivedBytes < fileSize {
			remainingBytes := int(fileSize - receivedBytes)
			currentChunkSize := chunkSize
			if remainingBytes < chunkSize {
				currentChunkSize = remainingBytes // Last chunk
			}

			chunk := make([]byte, currentChunkSize)
			n, addr, err := conn.ReadFromUDP(chunk)
			if err != nil {
				logrus.WithError(err).Error("error reading file chunk from UDP connection")
				break
			}
			if clientAddr == nil {
				clientAddr = addr
			}
			fileBuffer = append(fileBuffer, chunk[:n]...)
			receivedBytes += uint32(n)
			counter++
			logrus.WithField("chunk", counter).Debug("received file chunk")
		}

		if receivedBytes < fileSize {
			logrus.WithFields(logrus.Fields{
				"expected": fileSize,
				"received": receivedBytes,
			}).Error("Incomplete file received")
			_, err = conn.WriteTo([]byte(fmt.Sprintf("Failed to recieve received from: %s", clientAddr.String())), clientAddr)
			if err != nil {
				logrus.WithError(err).Error("error responding to client")
			}

			continue
		}
		//buf := make([]byte, size)
		//_, clientAddr, err := conn.ReadFromUDP(buf)
		//if err != nil {
		//	logrus.WithError(err).Error("error reading from UDP connection")
		//	continue
		//}

		err = f.storage.Save(fileBuffer)
		if err != nil {
			logrus.WithError(err).Error("error storing file")
		}
		// respond to the client
		_, err = conn.WriteTo([]byte(fmt.Sprintf("Message received from: %s", clientAddr.String())), clientAddr)
		if err != nil {
			logrus.WithError(err).Error("error responding to client")
		}
	}
}
