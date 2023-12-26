package utils

import (
	"context"
	"fmt"
	"net"
	"time"
)

const (
	TCP_PORT = 6722
	TIMEOUT  = 30 * time.Second
)

func ResetWithDelay(relayIP string, channel uint8, delay int) (status string, err error) {
	if channel > 8 {
		return "", fmt.Errorf("channel must be 0 (all) or between 1 and 8")
	}
	if delay < 0 {
		return "", fmt.Errorf("delay must be greater than 0")
	}

	return Raw(relayIP, fmt.Sprintf("1%d:%d", channel, delay))
}

func Reset(relayIP string, channel uint8, delay int) (status string, err error) {
	// allow only > 0 and < 2
	if channel > 2 {
		err = fmt.Errorf("channel must be 0 (all) or between 1 and 2")
		return
	}

	ch := "X"
	if channel > 0 {
		ch = fmt.Sprintf("%d", channel)
	}

	status, err = Raw(relayIP, "1"+ch)
	if err != nil {
		return "", err
	}
	// wait for the delay
	fmt.Println("Status", status)
	fmt.Println("Waiting for delay")
	time.Sleep(time.Duration(delay) * time.Second)
	fmt.Println("Done waiting")

	status, err = Raw(relayIP, "2"+ch)
	if err != nil {
		return "", err
	}

	return status, nil
}

func GetStatus(relayIP string) (status string, err error) {
	return Raw(relayIP, "00")
}

func Raw(relayIP string, command string) (status string, err error) {
	var d net.Dialer
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()

	address := fmt.Sprintf("%s:%d", relayIP, TCP_PORT)
	socket, err := d.DialContext(ctx, "tcp", address)
	if err != nil {
		return "", err
	}
	defer socket.Close()

	_, err = socket.Write([]byte(command))
	if err != nil {
		return "", err
	}

	// read from the socket
	buf := make([]byte, 1024)
	n, err := socket.Read(buf)
	if err != nil {
		return "", err
	}

	return string(buf[:n]), nil
}
