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
// @File: protocol_params.go
// @Date: 2018/05/07 14:35:07
////////////////////////////////////////////////////////////////////////////////

package params

const (
	MaximumExtraDataSize  uint64 = 32    // Maximum size extra data may be after Genesis.
	EpochDuration    uint64 = 30000 	 // Duration between proof-of-stack epochs
	MaxCodeSize 	= 32768 			 // Maximum bytecode to permit for a contract
)
