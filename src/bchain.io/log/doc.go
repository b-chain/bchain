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
// @File: doc.go
// @Date: 2018/05/08 18:02:08
////////////////////////////////////////////////////////////////////////////////

/*
Package btclog defines an interface and default implementation for subsystem
logging.

Log level verbosity may be modified at runtime for each individual subsystem
logger.

The default implementation in this package must be created by the Backend type.
Backends can write to any io.Writer, including multi-writers created by
io.MultiWriter.  Multi-writers allow log output to be written to many writers,
including standard output and log files.

Optional logging behavior can be specified by using the LOGFLAGS environment
variable and overridden per-Backend by using the WithFlags call option. Multiple
LOGFLAGS options can be specified, separated by commas.  The following options
are recognized:

  longfile: Include the full filepath and line number in all log messages

  shortfile: Include the filename and line number in all log messages.
  Overrides longfile.
*/
package log
