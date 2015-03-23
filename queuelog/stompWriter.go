package stompWriter

import (
	"fmt"

	"github.com/gmallard/stompngo"
)

// A wrapper for a stomp connection to facilitate making a io.Writer interface
type StompWriter struct {
	Connection        stompngo.Connection
	connectionHeaders stompngo.Headers
	netCon            net.Conn
	queueName         string
}

func New(hostname, port, username, password, queueName string) (*StompWriter, error) {
	newStompWriter := StompWriter{}

	newStompWriter.queueName = queueName
	newStompWriter.connectionHeaders = stompngo.Headers{
		"accept-version", "1.1",
		"login", username,
		"passcode", password,
		"host", hostname,
		"heart-beat", "5000,5000",
	}

	return &newStompWriter, nil
}

func (s *StompWriter) Connect() {
	// connect IO
	s.netCon, err = net.Dial("tcp", net.JoinHostPort(s.hostname, s.port))
	if err != nil {
		fmt.Printf("Couldn't connect to stomp host: %v\n", err.Error())
		return nil, err
	}
	s.Connection, err = stompngo.Connect(s.netCon, s.connectionHeaders)
	if err != nil {
		fmt.Printf("Couldn't Connect to stomp host: %v\n", err.Error())
		return nil, err
	}
}

func (s *StompWriter) Disconnect() {
	// disconnect IO
	s.Connection.Close()
	s.netCon.Close()
}

func (s *StompWriter) Write(payload []byte) (int, error) {
	// send message
	h := stompngo.Headers{
		"persistent", "true",
		"destination", "/queue/" + s.queueName,
		"content-type", "text/plain;charset=UTF-8",
	}
	err := s.Connection.Send(h, payload)
	if err != nil {
		fmt.Printf("error sending stomp message: %v\n", err.Error())
	}

	return len(payload), nil
}
