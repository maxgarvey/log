// Package queuelog provides a logging class that will send logs to a STOMP queueing server
// rather than to file.
//
// Basic example:
//                            vvv stomp params vvv           vvv cutoff threshold
//     logger := queuelog.New(host, port, user, pass, queue, log.Info)
//
//     logger.Debug("Connecting to the server...")   // Output goes to buffer
//     logger.Error("Connection failed")             // Output sends immediately

package queuelog

import (
	"bufio"
	"io"

	"github.com/alexcesaro/log"
	"github.com/alexcesaro/log/golog"
	"github.com/alexcesaro/log/queuelog/stompWriter"
)

// A Logger represents an active buffered logging object.
type Logger struct {
	*golog.Logger

	hostname  string
	port      string
	username  string
	password  string
	queueName string

	queueWriter stompWriter.StompWriter
}

// New creates a new Logger. Relevant parameters to the stomp queue
// being connecting to are specified
func New(host, port, username, password, queueName string, ignoreThreshold log.Level) *Logger {

	queueWriter, err = stompWriter.New(host, port, username, password, queueName)
	if err != nil {
		fmt.Printf("error connecting to stomp broker: %v\n", err.Error())
	}

	logger := &Logger{
		Logger:          golog.New(queueWriter, ignoreThreshold),
		hostname:        host,
		port:            port,
		username:        username,
		password:        password,
		queueName:       queueName,
		queueWriter:     queueWriter,
		ignoreThreshold: ignoreThreshold,
	}

	return &logger
}

func (l *Logger) output(level log.Level, args ...interface{}) {
	if level > l.ignoreThreshold {
		return
	}

	// additional logic to setup/teardown connections to queue
	l.queueWriter.Connect()
	defer l.queueWriter.Disconnect()

	l.Logger.output(level, args)
}
