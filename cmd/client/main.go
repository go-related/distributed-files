package main

import (
	"fmt"
	"github.com/juliant/distributed_file_server/internal/client"
	"github.com/sirupsen/logrus"
	"strconv"
	"strings"
)

func main() {
	addr := "127.0.0.1:8080"
	err := runClient(addr)
	if err != nil {
		logrus.Fatal(err)
	}
}

func runClient(addr string) error {
	cl, err := client.New(addr)
	if err != nil {
		return err
	}

	defer func(cl *client.Client) {
		err := cl.Close()
		if err != nil {
			logrus.WithError(err).Error("failed to close client")
		}
	}(cl)

	for {
		fmt.Println("Write the size to send, or press 'q' to exit:")

		var input string
		_, err := fmt.Scanln(&input)
		if err != nil {
			fmt.Println("Error reading input:", err)
			continue
		}

		input = strings.TrimSpace(input)
		if strings.ToLower(input) == "q" {
			fmt.Println("Exiting...")
			break
		}

		size, err := strconv.Atoi(input)
		if err != nil || size <= 0 {
			fmt.Println("Please enter a valid positive integer for the size.")
			continue
		}

		err = cl.SendRandomSizeFile(size)
		if err != nil {
			return err
		}

	}
	return nil
}
