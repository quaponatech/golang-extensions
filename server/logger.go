package server

import (
	"fmt"
	"io"
	"log"
	"os"
	"sync"
)

const (
	chPrefix  = "[CHANNEL] "
	logPrefix = "[LOGGER]  "
)

// DebugLevel list
const (
	Debug int = iota
	Info
	Warning
	State
	Error
	Quiet
)

// Logger ...
type Logger struct {
	serverName string
	prefix     string

	useLogFile  bool
	logDir      string
	logFileName string
	logFile     *os.File

	//MultiWriter io.Writer

	// WaitGroup to wait in the callee until channels are stopped
	WaitGroup   sync.WaitGroup
	StatusChan  chan Status
	ErrorChan   chan error
	WarningChan chan string
	LogChan     chan string
	DebugChan   chan string
	DebugLevel  int
}

// New should not be used but returns a minimal valid instance of a ServerLogger
func New() *Logger {
	serverName := "Unnamed Server"

	return &Logger{serverName: serverName, prefix: serverName + " - ",
		useLogFile: false,
		StatusChan: make(chan Status), ErrorChan: make(chan error),
		WarningChan: make(chan string), LogChan: make(chan string),
		DebugChan: make(chan string), DebugLevel: Quiet}
}

// NewLogger returns a fully configured ServerLogger
func NewLogger(serverName string, logDirectory, logFileName string,
	statusChannel chan Status, errorChannel chan error,
	warningChannel chan string, logChannel chan string,
	debugChannel chan string, debugLevel int) *Logger {

	if "" == serverName {
		return nil
	}

	useLogFile := false
	if "" != logDirectory && "" != logFileName {
		useLogFile = true
	}

	var statusCh chan Status
	var errorCh chan error
	var warningCh chan string
	var logCh chan string
	var debugCh chan string
	if nil != statusChannel {
		statusCh = statusChannel
	} else {
		statusCh = make(chan Status)
	}
	if nil != errorChannel {
		errorCh = errorChannel
	} else {
		errorCh = make(chan error, 100)
	}
	if nil != warningChannel {
		warningCh = warningChannel
	} else {
		warningCh = make(chan string, 1000)
	}
	if nil != logChannel {
		logCh = logChannel
	} else {
		logCh = make(chan string, 10000)
	}
	if nil != debugChannel {
		debugCh = debugChannel
	} else {
		debugCh = make(chan string, 10000)
	}

	return &Logger{serverName: serverName, prefix: serverName + " - ",
		useLogFile: useLogFile, logDir: logDirectory, logFileName: logFileName,
		StatusChan: statusCh, ErrorChan: errorCh,
		WarningChan: warningCh, LogChan: logCh,
		DebugChan: debugCh, DebugLevel: debugLevel}
}

// StartLogger opens a predefined log file, sets a multiwriter with it and runs channels for logging
func (l *Logger) StartLogger() error {
	log.Println(l.prefix + logPrefix + "Starting Server Logger")
	var err error

	if l.useLogFile {
		l.logFile, err = OpenLogFile(l.prefix, l.logDir, l.logFileName)
		if nil != err {
			return err
		}

		//l.MultiWriter = io.MultiWriter(os.Stdout, l.logFile)
		//log.SetOutput(l.MultiWriter)
		log.SetOutput(io.MultiWriter(os.Stdout, l.logFile))
	}

	go l.listenToErrorChannel()
	go l.listenToStatusChannel()
	go l.listenToWarningChannel()
	go l.listenToLogChannel()
	go l.listenToDebugChannel()

	log.Println(l.prefix + logPrefix + "Started Server Logger")
	return nil
}

// StopLogger stops the channels and closes the log files
func (l *Logger) StopLogger() {
	log.Println(l.prefix + logPrefix + "Stopping Server Logger")
	l.closeChannels()
	l.closeLogFiles()
	log.Println(l.prefix + logPrefix + "Stopped Server Logger")
}

// Close channels
func (l *Logger) closeChannels() {
	log.Println(l.prefix + logPrefix + "Closing Channels")
	close(l.DebugChan)
	close(l.LogChan)
	close(l.WarningChan)
	close(l.ErrorChan)
	close(l.StatusChan)

	l.WaitGroup.Wait()

	log.Println(l.prefix + logPrefix + "Closed Channels")
}

// Close logfiles
func (l *Logger) closeLogFiles() {
	log.Println(l.prefix + logPrefix + "Closing Log Files")

	mw := io.MultiWriter(os.Stdout)
	log.SetOutput(mw)

	l.logFile.Close()

	log.Println(l.prefix + logPrefix + "Closed Log Files")
}

func (l *Logger) listenToStatusChannel() {
	l.WaitGroup.Add(1)
	defer l.WaitGroup.Done()

	log.Println(l.prefix + chPrefix + "Listening to status channel")
	var msg Status
	var ok bool
	for {
		msg, ok = <-l.StatusChan
		if !ok {
			break
		}
		if l.DebugLevel > State {
			continue
		}
		message := fmt.Sprintf("[STATUS]  %v", msg.String())
		log.Println(l.prefix + message)
	}
	log.Println(l.prefix + chPrefix + "Stopped status channel")
}

func (l *Logger) listenToErrorChannel() {
	l.WaitGroup.Add(1)
	defer l.WaitGroup.Done()

	log.Println(l.prefix + chPrefix + "Listening to error channel")
	var msg error
	var ok bool
	for {
		msg, ok = <-l.ErrorChan
		if !ok {
			break
		}
		if l.DebugLevel > Error {
			continue
		}
		message := fmt.Sprintf("[ERROR]   %v", msg)
		log.Println(l.prefix + message)
	}
	log.Println(l.prefix + chPrefix + "Stopped error channel")
}

func (l *Logger) listenToWarningChannel() {
	l.WaitGroup.Add(1)
	defer l.WaitGroup.Done()

	log.Println(l.prefix + chPrefix + "Listening to warning channel")
	var msg string
	var ok bool
	for {
		msg, ok = <-l.WarningChan
		if !ok {
			break
		}
		if l.DebugLevel > Warning {
			continue
		}
		message := fmt.Sprintf("[WARNING] %v", msg)
		log.Println(l.prefix + message)
	}
	log.Println(l.prefix + chPrefix + "Stopped warning channel")
}

func (l *Logger) listenToLogChannel() {
	l.WaitGroup.Add(1)
	defer l.WaitGroup.Done()

	log.Println(l.prefix + chPrefix + "Listening to log channel")
	var msg string
	var ok bool
	for {
		msg, ok = <-l.LogChan
		if !ok {
			break
		}
		if l.DebugLevel > Info {
			continue
		}
		message := fmt.Sprintf("[INFO]    %v", msg)
		log.Println(l.prefix + message)
	}
	log.Println(l.prefix + chPrefix + "Stopped info channel")
}

func (l *Logger) listenToDebugChannel() {
	l.WaitGroup.Add(1)
	defer l.WaitGroup.Done()

	log.Println(l.prefix + chPrefix + "Listening to debug channel")
	var msg string
	var ok bool
	for {
		msg, ok = <-l.DebugChan
		if !ok {
			break
		}
		if l.DebugLevel > Debug {
			continue
		}
		message := fmt.Sprintf("[DEBUG]   %v", msg)
		log.Println(l.prefix + message)
	}
	log.Println(l.prefix + chPrefix + "Stopped debug channel")
}

// OpenLogFile ...
func OpenLogFile(prefix, logDir, logFileName string) (*os.File, error) {
	var logPrefix = prefix + "[LOGFILE] " + logFileName + ": "
	log.Println(logPrefix + "Opening")

	filePath := logDir + "/" + logFileName

	var err error
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_APPEND, 0600)
	if nil != err {
		err = os.MkdirAll(logDir, 0700)
		if nil != err {
			return nil, fmt.Errorf("Error: Creating log file path: %v", err)
		}
		file, err = os.OpenFile(filePath, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0600)
		if nil != err {
			return nil, fmt.Errorf("Error: Opening log file: %v", err)
		}
		log.Printf(logPrefix+"Success: Created NEW log file at %v", filePath)
	}

	log.Printf(logPrefix+"Success: Opened %v", file.Name())

	return file, nil
}
