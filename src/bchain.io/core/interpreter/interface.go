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
// @Date: 2018/07/30 16:19:30
////////////////////////////////////////////////////////////////////////////////

package interpreter

import "time"

type Interpreter interface {
	Validator
	Executor
	Formatter
}

type Validator interface {

}

type Executor interface {
	//about retValue:0 everything is ok, != 0 , something is not correct,this transaction should be discard
	Exec(code []byte, param []byte, ttl time.Duration) (duration time.Duration)
}

type Formatter interface {

}