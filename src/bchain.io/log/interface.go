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
// @File: interface.go
// @Date: 2018/05/08 18:02:08
////////////////////////////////////////////////////////////////////////////////

package log

// Logger is an interface which describes a level-based logger.  A default
// implementation of Logger is implemented by this package and can be created
// by calling (*Backend).Logger.
type Logger interface {
	// Tracef formats message according to format specifier and writes to
	// to log with LevelTrace.
	Tracef(format string, params ...interface{})

	// Debugf formats message according to format specifier and writes to
	// log with LevelDebug.
	Debugf(format string, params ...interface{})

	// Infof formats message according to format specifier and writes to
	// log with LevelInfo.
	Infof(format string, params ...interface{})

	// Warnf formats message according to format specifier and writes to
	// to log with LevelWarn.
	Warnf(format string, params ...interface{})

	// Errorf formats message according to format specifier and writes to
	// to log with LevelError.
	Errorf(format string, params ...interface{})

	// Criticalf formats message according to format specifier and writes to
	// log with LevelCritical.
	Criticalf(format string, params ...interface{})

	// Trace formats message using the default formats for its operands
	// and writes to log with LevelTrace.
	Trace(v ...interface{})

	// Debug formats message using the default formats for its operands
	// and writes to log with LevelDebug.
	Debug(v ...interface{})

	// Info formats message using the default formats for its operands
	// and writes to log with LevelInfo.
	Info(v ...interface{})

	// Warn formats message using the default formats for its operands
	// and writes to log with LevelWarn.
	Warn(v ...interface{})

	// Error formats message using the default formats for its operands
	// and writes to log with LevelError.
	Error(v ...interface{})

	// Critical formats message using the default formats for its operands
	// and writes to log with LevelCritical.
	Critical(v ...interface{})

	// Level returns the current logging level.
	Level() Level

	// SetLevel changes the logging level to the passed level.
	SetLevel(level Level)
}
