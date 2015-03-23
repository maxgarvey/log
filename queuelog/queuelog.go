// Package queuelog provides a buffered logging class that accumulates logs in
// memory until the flush threshold is reached which release stored logs, the
// buffered logger then act as a normal logger.
//
// Basic example:
//                            vvv stomp params vvv           vvv cutoff for buffering
//     logger := queuelog.New(host, port, user, pass, queue, log.Info,
//                            vvv maximum buffer depth vvv   vvv timeout in s for buffer send vvv
//                            maxBufferDepth,                bufferTimeout)
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

	// TODO: buffering stuff
	//bufferThreshold log.Level // <- cutoff for immediate queueing
	//maxBufferDepth  int
	//bufferTimeout   int

	//buffer chan []byte
}

// New creates a new Logger. Relevant parameters to the stomp queue
// being connecting to are specified
func New(host, port, username, password, queueName string, ignoreThreshold log.Level) { //, bufferThreshold log.Level, maxBufferDepth, bufferTimeout int) *Logger {

	queueWriter, err = stompWriter.New(host, port, username, password, queueName)
	if err != nil {
		fmt.Printf("error connecting to stomp broker: %v\n", err.Error())
	}

	logger := &Logger{
		Logger:      golog.New(queueWriter, ignoreThreshold),
		hostname:    host,
		port:        port,
		username:    username,
		password:    password,
		queueName:   queueName,
		queueWriter: queueWriter,
		//bufferThreshold: bufferThreshold,
		//maxBufferDepth:  maxBufferDepth,
		//bufferTimeout:   bufferTimeout,
		ignoreThreshold: ignoreThreshold,
	}
	//logger.buffer = make(chan []byte, maxBufferDepth)

	// go routine to check buffer

	return &logger
}

func (l *Logger) output(level log.Level, args ...interface{}) {
	/*
		// vvv desired logic
		if log.Level >= l.bufferThreshold {
			l.Logger.output(level, args)
		} else if log.Level >= l.ignoreThreshold {
			l.buffer <- l.Logger.Formatter(log.Level, args)
		}
	*/

	if level > l.ignoreThreshold {
		return
	}

	buffer := l.getBuffer()
	l.Formatter(&buffer.Buffer, level, args...)

	l.writeMutex.Lock()
	l.Writer(l.out, buffer.Bytes(), level)
	l.writeMutex.Unlock()

	l.putBuffer(buffer)
}
