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
// @File: bloom.go
// @Date: 2018/03/23 18:32:23
////////////////////////////////////////////////////////////////////////////////

package bloom

import(
	"math/big"
	"bchain.io/common/types"
	"bchain.io/utils/crypto"
)

type BloomByte interface {
	Bytes() []byte
}

func calculateBloom(b []byte) *big.Int {
	b = crypto.Keccak256(b[:])

	r := new(big.Int)

	for i := 0; i < 6; i += 2 {
		t := big.NewInt(1)
		b := (uint(b[i+1]) + (uint(b[i]) << 8)) & 2047
		r.Or(r, t.Lsh(t, b))
	}

	return r
}

func CreateBloom(topics []BloomByte) types.Bloom {
	bin := new(big.Int)
	for _, topic := range topics {
		bin.Or(bin, calculateBloom(topic.Bytes()[:]))
	}

	return types.BytesToBloom(bin.Bytes())
}

func BloomLookup(bin types.Bloom, topic BloomByte) bool {
	bloom := bin.Big()
	cmp := calculateBloom(topic.Bytes()[:])

	return bloom.And(bloom, cmp).Cmp(cmp) == 0
}
