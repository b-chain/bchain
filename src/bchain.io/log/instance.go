////////////////////////////////////////////////////////////////////////////////
// Copyright (c) 2018 The bchain-go Authors.
//
// The bchain-go is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.
//
// @File: instance.go
// @Date: 2018/05/08 18:02:08
////////////////////////////////////////////////////////////////////////////////

package log

import (
	"fmt"
	"github.com/jrick/logrotate/rotator"
	"os"
	"path/filepath"
	"sync"
)

type BchainlogError struct {
	error string
}

func (e BchainlogError) Error() string {
	return e.error
}

// error
var (
	ErrNewBackend        = BchainlogError{error: "New Backend failed"}
	ErrNoInstance        = BchainlogError{error: "Not create instance"}
	ErrLogRotatorFile    = BchainlogError{error: "Failed to create file rotator"}
	ErrLogRotatorCreated = BchainlogError{error: "Log has been Created"}
	ErrLogRotatorDir     = BchainlogError{error: "Failed to create log directory"}
)

// logWriter implements an io.Writer that outputs to both standard output and
// the write-end pipe of an initialized log rotator.
type logWriter struct{}

func (logWriter) Write(p []byte) (n int, err error) {
	os.Stdout.Write(p)
	if logRotator != nil {
		logRotator.Write(p)
	}
	return len(p), nil
}

var (
	// backendLog is the logging backend used to create all subsystem loggers.
	// The backend must not be used before the log rotator has been initialized,
	// or data races and/or nil pointer dereferences will occur.
	backendLog *Backend

	// logRotator is one of the logging outputs.  It should be closed on
	// application shutdown.
	logRotator *rotator.Rotator

	// a map of logger
	subsystemLoggers map[string]Logger
	muLoggerMap      sync.Mutex

	// options
	lvl  Level = LevelOff
	opts BackendOption

	//xxxlog = backendLog.Logger("xxx")
)

// initLogRotator initializes the logging rotater to write logs to logFile and
// create roll files in the same directory.  It must be called before the
// package-global log rotater variables are used.
func initLogRotator(logFile string) error {
	if logRotator != nil {
		return ErrLogRotatorCreated
	}

	logDir, _ := filepath.Split(logFile)
	err := os.MkdirAll(logDir, 0700)
	if err != nil {
		//fmt.Fprintf(os.Stderr, "failed to create log directory: %v\n", err)
		return ErrLogRotatorDir
	}
	r, err := rotator.New(logFile, 10*1024, false, 3)
	if err != nil {
		//fmt.Fprintf(os.Stderr, "failed to create file rotator: %v\n", err)
		return ErrLogRotatorFile
	}

	logRotator = r
	return nil
}

// create a log instance
func createInstance() error {
	backendLog = NewBackend(logWriter{})
	if backendLog == nil {
		return ErrNewBackend
	}

	subsystemLoggers = make(map[string]Logger)

	return nil
}

// destroy a log instance
func CloseInstance() {
	if logRotator != nil {
		logRotator.Close()
	}
}

// init a instance
func InitInstance(logFile string, logLevel string) error {
	if backendLog == nil {
		return ErrNoInstance
	}

	err := initLogRotator(logFile)
	if err != nil {
		if err == ErrLogRotatorCreated {
			fmt.Fprintf(os.Stdout, err.Error())
		} else {
			return err
		}
	}

	lvl, _ = LevelFromString(logLevel)
	for _, logger := range subsystemLoggers {
		logger.SetLevel(lvl)
	}

	return nil
}

// get a logger obj, must be called after CreateInstance()
func GetLogger(tag string) Logger {
	// subsystemLoggers not init
	if backendLog == nil {
		return nil
	}

	muLoggerMap.Lock()
	defer muLoggerMap.Unlock()

	logger := subsystemLoggers[tag]
	if logger == nil {
		logger = backendLog.Logger(tag)
		subsystemLoggers[tag] = logger
	}
	logger.SetLevel(lvl)

	return logger
}
