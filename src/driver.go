package main

import (
	"context"
	"encoding/binary"
	"fmt"
	containerdFifo "github.com/containerd/fifo"
	dockerLogDto "github.com/docker/docker/api/types/plugins/logdriver"
	dockerDaemonLogger "github.com/docker/docker/daemon/logger"
	protoIo "github.com/gogo/protobuf/io"
	"io"
	"os"
	"path"
	"strings"
	"sync"
	"syscall"
)

type FileLoggerContext struct {
	file                   string
	dockerDaemonLoggerInfo dockerDaemonLogger.Info
	inputFile              io.ReadCloser
	stdOutLogFile          *os.File
	stdErrLogFile          *os.File
}

type FileLoggingDriver struct {
	mu      sync.Mutex
	loggers map[string]*FileLoggerContext
}

func (fileLoggingDriver *FileLoggingDriver) StartLogging(file string, dockerDaemonLoggerInfo dockerDaemonLogger.Info) error {
	fileLoggerContext, err := fileLoggingDriver.makeFileLoggerContext(file, dockerDaemonLoggerInfo)
	if err != nil {
		return err
	}
	go fileLoggerContext.consumeLogsFromFile()
	return nil
}

func (fileLoggingDriver *FileLoggingDriver) StopLogging(file string) error {
	fileLoggingDriver.mu.Lock()
	fileLogger, exists := fileLoggingDriver.loggers[path.Base(file)]
	if exists {
		_ = fileLogger.inputFile.Close()
		_ = fileLogger.stdOutLogFile.Close()
		_ = fileLogger.stdErrLogFile.Close()
		// fixing permission problem, somehow the OpenFile perm is not effective...
		_ = os.Chmod(fileLogger.stdOutLogFile.Name(), 0666)
		_ = os.Chmod(fileLogger.stdErrLogFile.Name(), 0666)
		delete(fileLoggingDriver.loggers, path.Base(file))
	}
	fileLoggingDriver.mu.Unlock()
	return nil
}

func newFileLoggingDriver() *FileLoggingDriver {
	return &FileLoggingDriver{
		loggers: make(map[string]*FileLoggerContext),
	}
}

func (fileLoggingDriver *FileLoggingDriver) makeFileLoggerContext(file string, dockerDaemonLoggerInfo dockerDaemonLogger.Info) (*FileLoggerContext, error) {
	fileLoggingDriver.mu.Lock()
	if _, exists := fileLoggingDriver.loggers[path.Base(file)]; exists {
		fileLoggingDriver.mu.Unlock()
		return nil, fmt.Errorf("logger for %q already exists", file)
	}
	fileLoggingDriver.mu.Unlock()

	inputFile, err := containerdFifo.OpenFifo(context.Background(), file, syscall.O_RDONLY, 0700)
	if err != nil {
		return nil, err
	}
	logFilePath := dockerDaemonLoggerInfo.ContainerName
	if logFileDir, ok := dockerDaemonLoggerInfo.Config["log-file-dir"]; ok {
		logFilePath = path.Join("/hostRoot", logFileDir, logFilePath)
	} else {
		logFilePath = path.Join("/hostRoot", "tmp", logFilePath)
	}
	stdOutLogFile, err := os.OpenFile(logFilePath+".out.log", os.O_TRUNC|os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		_ = inputFile.Close()
		return nil, err
	}
	stdErrLogFile, err := os.OpenFile(logFilePath+".err.log", os.O_TRUNC|os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		_ = inputFile.Close()
		_ = stdOutLogFile.Close()
		return nil, err
	}
	fileLoggerContext := &FileLoggerContext{
		file:                   file,
		dockerDaemonLoggerInfo: dockerDaemonLoggerInfo,
		inputFile:              inputFile,
		stdOutLogFile:          stdOutLogFile,
		stdErrLogFile:          stdErrLogFile,
	}
	fileLoggingDriver.mu.Lock()
	fileLoggingDriver.loggers[path.Base(file)] = fileLoggerContext
	fileLoggingDriver.mu.Unlock()
	return fileLoggerContext, nil
}

func (fileLoggerContext *FileLoggerContext) consumeLogsFromFile() {
	dec := protoIo.NewUint32DelimitedReader(fileLoggerContext.inputFile, binary.BigEndian, 1e6)
	defer dec.Close()
	var logEntry dockerLogDto.LogEntry
	for {
		if err := dec.ReadMsg(&logEntry); err != nil {
			if err == io.EOF || err == os.ErrClosed || strings.Contains(err.Error(), "file already closed") {
				_ = fileLoggerContext.inputFile.Close()
				return
			}
			dec = protoIo.NewUint32DelimitedReader(fileLoggerContext.inputFile, binary.BigEndian, 1e6)
		}
		if logEntry.Source == "stdout" {
			_, _ = fileLoggerContext.stdOutLogFile.Write(logEntry.Line)
			_, _ = fileLoggerContext.stdOutLogFile.Write([]byte{0x0a})
		} else {
			_, _ = fileLoggerContext.stdErrLogFile.Write(logEntry.Line)
			_, _ = fileLoggerContext.stdErrLogFile.Write([]byte{0x0a})
		}
		logEntry.Reset()
	}
}
