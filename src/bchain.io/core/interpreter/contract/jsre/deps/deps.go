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
// @File: deps.go
// @Date: 2018/01/21 19:39:21
////////////////////////////////////////////////////////////////////////////////

// Package deps contains the console JavaScript dependencies Go embedded.
package deps

//go:generate jsmarshal consensus.js consensus.md jsre.JSRE

//go:generate jsmarshal bchain.js bchain.md jsre.JSRE

//go:generate jsmarshal system.js system.md jsre.JSRE

//go:generate jsmarshal pledge.js pledge.md jsre.JSRE

//go:generate go-bindata -nometadata -pkg deps -o bindata.go consensus.md bchain.md system.md pledge.md
