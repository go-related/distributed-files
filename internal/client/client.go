package client

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"io"
	"net"
)

type Client struct {
	conn net.Conn
}

func New(serverAddr string) (*Client, error) {
	conn, err := net.Dial("udp", serverAddr)

	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to server")
	}

	return &Client{
		conn: conn,
	}, nil
}

func (cl *Client) SendRandomSizeFile(size int) error {
	buf := make([]byte, size)
	_, err := rand.Read(buf)
	if err != nil {
		return errors.Wrap(err, "failed to populate the buffer for sending random size")
	}

	err = binary.Write(cl.conn, binary.LittleEndian, uint32(size))
	if err != nil {
		return errors.Wrap(err, "failed to send the size of the file to the server")
	}

	//
	//_, err = cl.conn.Write(buf[:read])
	// this will allow us more flexibility
	reader := bytes.NewReader(buf)
	_, err = io.Copy(cl.conn, reader)
	if err != nil {
		return errors.Wrap(err, "failed to send file")
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
