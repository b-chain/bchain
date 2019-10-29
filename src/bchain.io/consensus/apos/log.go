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
// @File: log.go
// @Date: 2018/06/12 17:18:09
////////////////////////////////////////////////////////////////////////////////

package apos

import (
	"fmt"
	"bchain.io/log"
	"os"
)

var (
	logTag = "consensus.apos"
	logger log.Logger
)

func init() {
	logger = log.GetLogger(logTag)
	logger.SetLevel(log.LevelDebug)
	if logger == nil {
		fmt.Errorf("Can not get logger(%s)\n", logTag)
		os.Exit(1)
	}
}