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
// @File: net.go
// @Date: 2018/03/19 09:35:19
////////////////////////////////////////////////////////////////////////////////

package types

import (
	"net"
)

//go:generate msgp
//msgp:shim net.IP as:[]byte using:toBytes/fromBytes

func toBytes(ip net.IP) []byte {
	return []byte(ip)
}

func fromBytes(b []byte) net.IP {
	return net.IP(b)
}

var (
	ipType int8
)

type IP struct {
	Ip net.IP `msg:"ip"`
}

func (ip IP) Get() net.IP {
	return ip.Ip
}

func (ip *IP) Put(in net.IP) *IP {
	ip.Ip = in
	return ip
}

func NewIP(in net.IP) *IP {
	ip := new(IP)
	ip.Ip = in
	return ip
}

// Here, we'll pick an arbitrary number between
// 0 and 127 that isn't already in use
func (*IP) ExtensionType() int8 {
	return ipType
}

// We'll always use 16 bytes to encode the data
func (ip *IP) Len() int {
	return len(ip.Ip)
}

// MarshalBinaryTo simply copies the value
// of the bytes into 'b'
func (ip *IP) MarshalBinaryTo(b []byte) error {
	if len(ip.Ip) <= net.IPv6len {
		copy(b, ip.Ip)
		return nil
	}
	return ErrBytesTooLong
}

// UnmarshalBinary copies the value of 'b'
// into the Hash object. (We might want to add
// a sanity check here later that len(b) <= HashLength.)
func (ip *IP) UnmarshalBinary(b []byte) error {
	// TODO: check b, only hex, len <= HashLength
	if len(b) <= net.IPv6len {
		if ipv4 := ip.Ip.To4(); ipv4 != nil {
			ip.Ip = ipv4
		} else {
			ip.Ip = make(net.IP, len(b))
			copy(ip.Ip, b)
		}
		return nil
	}

	return ErrBytesTooLong
}
