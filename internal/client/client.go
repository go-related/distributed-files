package client

import (
	"crypto/rand"
	"encoding/binary"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"net"
)

type Client struct {
	conn      net.Conn
	chunkSize int
}

func New(serverAddr string) (*Client, error) {
	conn, err := net.Dial("udp", serverAddr)

	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to server")
	}

	return &Client{
		conn:      conn,
		chunkSize: 1024,
	}, nil
}

func (cl *Client) SendRandomSizeFile(size int) error {
	buf := make([]byte, size)
	_, err := rand.Read(buf)
	if err != nil {
		return errors.Wrap(err, "failed to populate the buffer for sending random size")
	}
	logrus.Debug("finished writing randomly  to the file")

	err = binary.Write(cl.conn, binary.LittleEndian, uint32(size))
	if err != nil {
		return errors.Wrap(err, "failed to send the size of the file to the server")
	}

	// in udp conn max byts are  65,535 so we need to send data chunk by chunk
	counter := 0
	for start := 0; start < len(buf); start += cl.chunkSize {
		end := start + cl.chunkSize
		if end > len(buf) {
			end = len(buf)
		}

		// Write the chunk
		_, err = cl.conn.Write(buf[start:end])
		if err != nil {
			return errors.Wrap(err, "failed to send file chunk")
		}
		counter++
		logrus.WithField("chunk", counter).Debug("sent file chunk")
	}

	readBuffer := make([]byte, 1024)
	_, err = cl.conn.Read(readBuffer)
	if err != nil {
		logrus.WithError(err).Error("failed to read from server")
	}
	logrus.WithField("size", size).Infof("file sent to the server, recieved: %s ", string(readBuffer))
	return nil
}

func (cl *Client) Close() error {
	if cl.conn != nil {
		err := cl.conn.Close()
		if err != nil {
			return err
		}
	}
	return nil
}
