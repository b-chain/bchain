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
// @File: block_signing.go
// @Date: 2018/05/14 18:11:05
////////////////////////////////////////////////////////////////////////////////

package block

import (
	"crypto/ecdsa"
	"errors"
	"fmt"
	"math/big"
	"bchain.io/common"
	"bchain.io/common/types"
	"bchain.io/utils/crypto"
)

var (
	ErrInvalidSig     = errors.New("invalid block v, r, s values")
	ErrInvalidChainId = errors.New("invalid chain id for block signer")
)

// SignHeader signs the header using the given signer and private key
func SignHeader(h *Header, s Signer, prv *ecdsa.PrivateKey) (*Header, error) {
	hash := s.Hash(h)
	sig, err := crypto.Sign(hash[:], prv)
	if err != nil {
		return nil, err
	}
	return h.WithSignature(s, sig)
}

// SignHeaderInner signs the header(modify R S V) using the given signer and private key
func SignHeaderInner(h *Header, s Signer, prv *ecdsa.PrivateKey) error {
	hash := s.Hash(h)
	sig, err := crypto.Sign(hash[:], prv)
	if err != nil {
		return err
	}
	return h.AddSignature(s, sig)
}

// Signer encapsulates transaction signature handling. Note that this interface is not a
// stable API and may change at any time to accommodate new protocol rules.
type Signer interface {
	// Sender returns the sender address of the transaction.
	Sender(h *Header) (types.Address, error)
	// SignatureValues returns the raw R, S, V values corresponding to the
	// given signature.
	SignatureValues(h *Header, sig []byte) (r, s, v *big.Int, err error)

	// Hash returns the hash to be signed.
	Hash(h *Header) types.Hash

	// Equal returns true if the given signer is the same as the receiver.
	Equal(Signer) bool
}

type BlockSigner struct {
	chainId, chainIdMul *big.Int
}

// NewBlockSigner returns a Signer based on the given chain config
func NewBlockSigner(chainId *big.Int) BlockSigner {
	if chainId == nil {
		chainId = new(big.Int)
	}
	return BlockSigner{
		chainId:    chainId,
		chainIdMul: new(big.Int).Mul(chainId, common.Big2),
	}
}

func (s BlockSigner) Equal(signer Signer) bool {
	bSigner, ok := signer.(BlockSigner)
	return ok && bSigner.chainId.Cmp(s.chainId) == 0
}

func (s BlockSigner) Sender(h *Header) (types.Address, error) {
	if deriveChainId(&h.V.IntVal).Cmp(s.chainId) != 0 {
		return types.Address{}, ErrInvalidChainId
	}

	//empty := types.Address{}
	//if h.BlockProducer != empty {
	//	return h.BlockProducer, nil
	//}

	V := &big.Int{}
	if s.chainId.Sign() != 0 {
		V = V.Sub(&h.V.IntVal, s.chainIdMul)
		V.Sub(V, common.Big35)
	} else {
		V = V.Sub(&h.V.IntVal, common.Big27)
	}
	address, err := recoverPlain(h.HashNoSig(), &h.R.IntVal, &h.S.IntVal, V, true)
	h.Producer = address
	return address, err
}

func (s BlockSigner) VerifySignature(h *Header) (bool, error) {
	//chain id check
	if deriveChainId(&h.V.IntVal).Cmp(s.chainId) != 0 {
		return false, ErrInvalidChainId
	}

	//R S V check
	var V uint64
	if s.chainId.Sign() != 0 {
		V = h.V.IntVal.Uint64() - s.chainIdMul.Uint64() - 35
	} else {
		V = h.V.IntVal.Uint64() - 27
	}

	if !crypto.ValidateSignatureValues(byte(V), &h.R.IntVal, &h.S.IntVal, true) {
		return false, ErrInvalidSig
	}

	//encode the snature in uncompressed format
	R, S := h.R.IntVal.Bytes(), h.S.IntVal.Bytes()
	sig := make([]byte, 65)
	copy(sig[32-len(R):32], R)
	copy(sig[64-len(S):64], S)
	sig[64] = byte(V)

	// recover the public key from the signature
	hash := s.Hash(h).Bytes()
	pub, err := crypto.Ecrecover(hash, sig)
	if err != nil {
		return false, err
	}
	if len(pub) == 0 || pub[0] != 4 {
		return false, errors.New("invalid public key")
	}

	ret := crypto.VerifySignature(pub, hash, sig[:64])

	return ret, nil
}

// SignatureValues returns a header's R S V based given signature. This signature
// needs to be in the [R || S || V] format where V is 0 or 1.
func (s BlockSigner) SignatureValues(h *Header, sig []byte) (R, S, V *big.Int, err error) {
	if len(sig) != 65 {
		errStr := fmt.Sprintf("wrong size for signature: got %d, want 65", len(sig))
		err = errors.New(errStr)
		return nil, nil, nil, err
	} else {
		R = new(big.Int).SetBytes(sig[:32])
		S = new(big.Int).SetBytes(sig[32:64])

		if s.chainId.Sign() != 0 {
			V = big.NewInt(int64(sig[64] + 35))
			V.Add(V, s.chainIdMul)
		} else {
			V = new(big.Int).SetBytes([]byte{sig[64] + 27})
		}
	}

	return R, S, V, nil
}

func (s BlockSigner) Hash(h *Header) types.Hash {
	return h.HashNoSig()
}

func recoverPlain(sighash types.Hash, R, S, Vb *big.Int, homestead bool) (types.Address, error) {
	if Vb.BitLen() > 8 {
		return types.Address{}, ErrInvalidSig
	}
	V := byte(Vb.Uint64())
	if !crypto.ValidateSignatureValues(V, R, S, homestead) {
		return types.Address{}, ErrInvalidSig
	}
	// encode the snature in uncompressed format
	r, s := R.Bytes(), S.Bytes()
	sig := make([]byte, 65)
	copy(sig[32-len(r):32], r)
	copy(sig[64-len(s):64], s)
	sig[64] = V
	// recover the public key from the snature
	pub, err := crypto.Ecrecover(sighash[:], sig)
	if err != nil {
		return types.Address{}, err
	}
	if len(pub) == 0 || pub[0] != 4 {
		return types.Address{}, errors.New("invalid public key")
	}
	var addr types.Address
	copy(addr[:], crypto.Keccak256(pub[1:])[12:])
	return addr, nil
}

// deriveChainId derives the chain id from the given v parameter
func deriveChainId(v *big.Int) *big.Int {
	if v.BitLen() <= 64 {
		v := v.Uint64()
		if v == 27 || v == 28 {
			return new(big.Int)
		}
		return new(big.Int).SetUint64((v - 35) / 2)
	}
	v = new(big.Int).Sub(v, big.NewInt(35))
	return v.Div(v, big.NewInt(2))
}
