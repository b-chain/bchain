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
// @File: transaction_signing.go
// @Date: 2018/05/08 15:18:08
////////////////////////////////////////////////////////////////////////////////

package transaction

import (
	"crypto/ecdsa"
	"errors"

	"math/big"
	"bchain.io/params"
	"bchain.io/utils/crypto"

	"fmt"
	"bchain.io/common"
	"bchain.io/common/types"
)

var (
	ErrInvalidChainId = errors.New("invalid chain id for signer")
)

// sigCache is used to cache the derived sender and contains
// the signer used to derive it.
type sigCache struct {
	signer Signer
	from   types.Address
}

// MakeSigner returns a Signer based on the given chain config and block number.
func MakeSigner(config *params.ChainConfig, blockNumber *big.Int) Signer {
	_ = blockNumber // blockNumber is reserve field
	var signer Signer
	//use latest Signer
	signer = NewMSigner(config.ChainId)
	return signer
}

// SignTx signs the transaction using the given signer and private key
func SignTx(tx *Transaction, s Signer, prv *ecdsa.PrivateKey) (*Transaction, error) {
	h := s.Hash(tx)
	sig, err := crypto.Sign(h[:], prv)
	if err != nil {
		return nil, err
	}
	return tx.WithSignature(s, sig)
}

// Sender returns the address derived from the signature (V, R, S) using secp256k1
// elliptic curve and an error if it failed deriving or upon an incorrect
// signature.
//
// Sender may cache the address, allowing it to be used regardless of
// signing method. The cache is invalidated if the cached signer does
// not match the signer used in the current call.
func Sender(signer Signer, tx *Transaction) (types.Address, error) {
	if sc := tx.from.Load(); sc != nil {
		sigCache := sc.(sigCache)
		// If the signer used to derive from in a previous
		// call is not the same as used current, invalidate
		// the cache.
		if sigCache.signer.Equal(signer) {
			return sigCache.from, nil
		}
	}

	addr, err := signer.Sender(tx)
	if err != nil {
		return types.Address{}, err
	}
	tx.from.Store(sigCache{signer: signer, from: addr})
	return addr, nil
}

// Signer encapsulates transaction signature handling. Note that this interface is not a
// stable API and may change at any time to accommodate new protocol rules.
type Signer interface {
	// Sender returns the sender address of the transaction.
	Sender(tx *Transaction) (types.Address, error)
	// SignatureValues returns the raw R, S, V values corresponding to the
	// given signature.
	SignatureValues(tx *Transaction, sig []byte) (r, s, v *big.Int, err error)
	// Hash returns the hash to be signed.
	Hash(tx *Transaction) types.Hash
	// Equal returns true if the given signer is the same as the receiver.
	Equal(Signer) bool
}

type MSigner struct {
	chainId, chainIdMul *big.Int
}

func NewMSigner(chainId *big.Int) MSigner {
	if chainId == nil {
		chainId = new(big.Int)
	}
	return MSigner{
		chainId:    chainId,
		chainIdMul: new(big.Int).Mul(chainId, big.NewInt(2)),
	}
}

func (s MSigner) Equal(s2 Signer) bool {
	eip155, ok := s2.(MSigner)
	return ok && eip155.chainId.Cmp(s.chainId) == 0
}

var big8 = big.NewInt(8)

func (s MSigner) Sender(tx *Transaction) (types.Address, error) {

	if tx.ChainId().Cmp(s.chainId) != 0 {
		return types.Address{}, ErrInvalidChainId
	}
	V := new(big.Int).Sub(&tx.Data.V.IntVal, s.chainIdMul)
	V.Sub(V, big8)
	return recoverPlain(s.Hash(tx), &tx.Data.R.IntVal, &tx.Data.S.IntVal, V, true)
}

// WithSignature returns a new transaction with the given signature. This signature
// needs to be in the [R || S || V] format where V is 0 or 1.
func (s MSigner) SignatureValues(tx *Transaction, sig []byte) (R, S, V *big.Int, err error) {
	//here use Frontier SignatureValues Function directly
	if len(sig) != 65 {
		errStr := fmt.Sprintf("wrong size for signature: got %d, want 65", len(sig))
		err = errors.New(errStr)
	} else {
		R = new(big.Int).SetBytes(sig[:32])
		S = new(big.Int).SetBytes(sig[32:64])
		V = new(big.Int).SetBytes([]byte{sig[64] + 27})
	}

	if err != nil {
		return nil, nil, nil, err
	}
	if s.chainId.Sign() != 0 {
		V = big.NewInt(int64(sig[64] + 35))
		V.Add(V, s.chainIdMul)
	}
	return R, S, V, nil
}

// Hash returns the hash to be signed by the sender.
// It does not uniquely identify the transaction.
func (s MSigner) Hash(tx *Transaction) types.Hash {
	h, err := common.MsgpHash([]interface{}{
		tx.Data.H,
		tx.Actions(),
		types.BigInt{*s.chainId},
		uint(0),
		uint(0),
	})

	if err != nil {
		panic(err)
	}

	return h
}

func recoverPlain(sighash types.Hash, R, S, Vb *big.Int, homestead bool) (types.Address, error) {
	if Vb.BitLen() > 8 {
		return types.Address{}, ErrInvalidSig
	}
	V := byte(Vb.Uint64() - 27)
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
