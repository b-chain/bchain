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
// @Date: 2018/05/08 17:22:08
////////////////////////////////////////////////////////////////////////////////

package blockproducer

import (
	"fmt"
	"os"
	"bchain.io/log"
)

var (
	LogTag = "blockproducer"
	logger log.Logger
)

func init() {
	logger = log.GetLogger(LogTag)
	if logger == nil {
		fmt.Errorf("Can not get logger(%s)\n", LogTag)
		os.Exit(1)
	}
}