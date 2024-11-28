package client

import (
	"crypto/rand"
	"github.com/pkg/errors"
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
	read, err := rand.Read(buf)
	if err != nil {
		return errors.Wrap(err, "failed to populate the buffer for sending random size")
	}

	_, err = cl.conn.Write(buf[:read])
	if err != nil {
		return errors.Wrap(err, "failed to send file")
	}

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
