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
// @File: apos_signing.go
// @Date: 2018/06/13 11:12:13
////////////////////////////////////////////////////////////////////////////////

package apos

import (
	"crypto/ecdsa"
	"errors"
	"fmt"
	"math/big"
	"bchain.io/common"
	"bchain.io/common/types"
	"bchain.io/core/blockchain/block"
	"bchain.io/utils/crypto"
)

//go:generate msgp

var (
	ErrInvalidSig     = errors.New("invalid  v, r, s values")
	ErrInvalidChainId = errors.New("invalid chain id for signer")

	gCommonTools CommonTools
)

// Signer encapsulates apos signature handling. Note that this interface is not a
// stable API and may change at any time to accommodate new protocol rules.
type signer interface {
	// sign the obj
	sign(prv *ecdsa.PrivateKey) (R *big.Int, S *big.Int, V *big.Int, err error)

	// Sender returns the sender address of the Credential.
	sender() (types.Address, error)

	// hash
	hash() types.Hash
}

/*
ReadMe:
what the relationship between SeedDataSigForHash and SeedData?
if I say hash-SeedData is  Qr-1 and hash-SeedDataSigForHash is  Qr , it's not wrong.
although hash-SeedDataSigForHash just a hash not equal Qr and Qr-1,but it's used for producing Qr.
like this below:
h := hash-SeedDataSigForHash
SeedData.sign(h , privateKey),the sign will fill the R,S,V in SeedData
and the consensusData.para = SeedData.toBytes(),the function toBytes() equal SeedData.Signature.toBytes(),it's
meaning that the bytes equal R,S,V. But the consensusData not equal the Qr-1/r ,it's just a part of Qr-1/r

but now,how to get hash-SeedDataSigForHash? The defination of SeedDataSigForHash like below:
type SeedDataSigForHash struct {
	Round    uint64 // round
	Seed     []byte // quantity(seed, Qr-1)
}

here we know Round,but,what is Seed?
1.Seed-Round(Qr) is depend on Qr-1
2.Qr = H(Qr-1 , r)

so we know Seed is Qr-1(hash -> bytes)
sigBytes(r-1) = block(r-1).consensusData.Para
Signature(r-1) = Signature.get(sigBytes(r-1))
SeedData(Signature(r-1) , r-1).Hash()-->Qr-1

so ,here we know how to get Qr-1,now ,how to get Qr?
First,it's a bad question,the result we get is not Qr,just a part of Qr--we called Signature's bytes.
SeedDataSigForHash{Round:r , Seed:bytes(Qr-1)}-->hash,this hash just for sign,because we want get signature's bytes

h := hash-SeedDataSigForHash{Round:r , Seed:bytes(Qr-1)}
SeedData{Signature:init() , r}.sign(h , privateKey),the r is useless in sign,but we fill Signature in this fucntion.
so A PART OF Qr = SeedData{Signature:signed and filled , r}.toBytes().

Another tips:
Seed is continuous , but Credential is not
like:
Seed0->Seed1->Seed2->Seed3->Seed4
C1(Seed0)->C2(Seed0)->C3(Seed0)->C4(Seed1)->C5(Seed1)->C6(Seed1)......

*/
type SeedDataSigForHash struct {
	Round uint64 // round
	Seed  []byte // quantity(seed, Qr-1)
}

// Qr = H(SIGℓr (Qr−1), r),
type SeedData struct {
	Signature // SIGℓr (Qr−1) is the signature of leader bp
	Round     uint64
}

//use this , we get Qr-1
func (this *SeedData) Hash() types.Hash {
	h, err := common.MsgpHash(this)
	if err != nil {
		return types.Hash{}
	}
	return h
}

//use this,we get a hash for signature,what will be a part of Qr--ConsensusData.Para
//the different between Hash and hashForSig is that Hash get a Qr-1,but hashForSig get a param for signing
func (sd *SeedData) hashForSig() types.Hash {
	//get Qr-1:types.Hash{}
	seed_R_1, _, err := restoreSeed(sd.Round - 1)
	if err != nil {
		logger.Error("get Quantity fail")
		return types.Hash{}
	}
	qdforhash := &SeedDataSigForHash{
		sd.Round,
		seed_R_1.Bytes(),
	}
	hash, err := common.MsgpHash(qdforhash)
	if err != nil {
		return types.Hash{}
	}
	return hash
}

func (sd *SeedData) sign(prv *ecdsa.PrivateKey) (R *types.BigInt, S *types.BigInt, V *types.BigInt, err error) {
	if prv == nil {
		err := errors.New(fmt.Sprintf("private key is empty"))
		return nil, nil, nil, err
	}
	//Qr hash,{Qr-1 , r}-->Hash
	hash := sd.hashForSig()
	if (hash == types.Hash{}) {
		err := errors.New(fmt.Sprintf("the hash of QuantityData is empty"))
		return nil, nil, nil, err
	}

	sig, err := crypto.Sign(hash[:], prv)
	if err != nil {
		return nil, nil, nil, err
	}

	err = sd.get(sig)
	if err != nil {
		return nil, nil, nil, err
	}
	R = sd.R
	S = sd.S
	V = sd.V

	return R, S, V, nil
}

func (sd *SeedData) sender(parent *block.Header) (types.Address, error) {
	sd.checkObj()
	if Config().chainId != nil && deriveChainId(&sd.V.IntVal).Cmp(Config().chainId) != 0 {
		return types.Address{}, ErrInvalidChainId
	}
	if Config().chainId == nil {
		panic("Config().chainId == nil")
	}
	V := &big.Int{}
	if Config().chainId.Sign() != 0 {
		V = V.Sub(&sd.V.IntVal, Config().chainIdMul)
		V.Sub(V, common.Big35)
	} else {
		V = V.Sub(&sd.V.IntVal, common.Big27)
	}

	seed_R_1, _, err := generateSeedBySeedHeader(parent)
	if err != nil {
		return types.Address{}, err
	}

	qdforhash := &SeedDataSigForHash{
		sd.Round,
		seed_R_1.Bytes(),
	}
	hash, err := common.MsgpHash(qdforhash)
	if err != nil {
		return types.Address{}, err
	}

	address, err := recoverPlain(hash, &sd.R.IntVal, &sd.S.IntVal, V, true)
	return address, err
}

// empty block Qr = H(Qr−1, r)
type QuantityEmpty struct {
	LstQuantity types.Hash
	Round       uint64
}

func (this *QuantityEmpty) Hash() types.Hash {
	h, err := common.MsgpHash(this)
	if err != nil {
		return types.Hash{}
	}
	return h
}

// signature R, S, V
type Signature struct {
	R *types.BigInt
	S *types.BigInt
	V *types.BigInt
}

func MakeEmptySignature() *Signature {
	s := new(Signature)
	s.init()
	return s
}

func (s *Signature) init() {
	s.R = new(types.BigInt)
	s.S = new(types.BigInt)
	s.V = new(types.BigInt)
}

func (s *Signature) Init() {
	s.init()
}

type signValue interface {
	// check the signature obj is initialized, if not, throw painc
	checkObj()

	// get() computes R, S, V values corresponding to the
	// given signature.
	get(sig []byte) (err error)

	// convert to bytes
	toBytes() (sig []byte)
}

func (s *Signature) checkObj() {
	if s.R == nil || s.S == nil || s.V == nil {
		panic(fmt.Errorf("Signature obj is not initialized"))
	}
}

func (s *Signature) hashBytes() []byte {
	return s.Hash().Bytes()

}

func (s *Signature) Hash() types.Hash {

	h, err := common.MsgpHash(s)
	if err != nil {
		return types.Hash{}
	}
	return h

}
func (s *Signature) get(sig []byte) (err error) {
	s.checkObj()

	if len(sig) != 65 {
		return errors.New(fmt.Sprintf("wrong size for Signature: got %d, want 65", len(sig)))
	} else {
		s.R.IntVal.SetBytes(sig[:32])
		s.S.IntVal.SetBytes(sig[32:64])

		if Config().chainId != nil && Config().chainId.Sign() != 0 {
			s.V.IntVal.SetInt64(int64(sig[64] + 35))
			s.V.IntVal.Add(&s.V.IntVal, Config().chainIdMul)
		} else {
			s.V.IntVal.SetBytes([]byte{sig[64] + 27})
		}
	}
	return nil
}

func (s *Signature) FillBySig(sig []byte) (R, S, V *big.Int, err error) {
	err = s.get(sig)
	if err != nil {
		return nil, nil, nil, err
	}

	R = new(big.Int)
	R.Set(&s.R.IntVal)

	S = new(big.Int)
	S.Set(&s.S.IntVal)

	V = new(big.Int)
	V.Set(&s.V.IntVal)

	return R, S, V, nil
}

func (s Signature) toBytes() (sig []byte) {
	s.checkObj()

	sV := s.V
	V := types.BigInt{}
	if Config().chainId.Sign() != 0 {
		V.IntVal.Sub(&sV.IntVal, Config().chainIdMul)
		V.IntVal.Sub(&V.IntVal, common.Big35)
	} else {
		V.IntVal.Sub(&sV.IntVal, common.Big27)
	}

	vb := byte(V.IntVal.Uint64())
	if !crypto.ValidateSignatureValues(vb, &s.R.IntVal, &s.S.IntVal, true) {
		logger.Debugf("invalid Signature\n")
		return nil
	}

	rb, sb := s.R.IntVal.Bytes(), s.S.IntVal.Bytes()
	sig = make([]byte, 65)
	copy(sig[32-len(rb):32], rb)
	copy(sig[64-len(sb):64], sb)
	sig[64] = vb

	return sig
}

// long-term key singer
type CredentialSign struct {
	Signature
	Round      uint64 // round
	Step       uint64 // step
	ParentHash types.Hash
	Time       *types.BigInt
	votes      float64
}

type CredentialSigForHash struct {
	Round    uint64 // round
	Step     uint64 // step
	Quantity []byte // quantity(seed, Qr-1)
}

func (cs *CredentialSign) CsHash() types.Hash {
	hash, err := common.MsgpHash(cs)
	if err != nil {
		return types.Hash{}
	}
	return hash
}

func (a *CredentialSign) sigHashBig() *big.Int {
	h := a.Signature.hashBytes()
	return new(big.Int).SetBytes(h)
}

func (a *CredentialSign) sigHashHashBig()*big.Int{
	h := a.Signature.hashBytes()
	hh := crypto.Keccak256Hash(h)
	return new(big.Int).SetBytes(hh.Bytes())
}

func (a *CredentialSign) Cmp(b *CredentialSign) int {
	h := a.Signature.hashBytes()
	aInt := new(big.Int).SetBytes(h)

	h = b.Signature.hashBytes()
	bInt := new(big.Int).SetBytes(h)

	return aInt.Cmp(bInt)
}

//CredentialSign.sign do everything except holding the privateKey
//called by commonTools,because commonTools hold the privateKey for signing

//the commontools has give us the privateKey
func (cret *CredentialSign) sign(prv *ecdsa.PrivateKey) (R *types.BigInt, S *types.BigInt, V *types.BigInt, err error) {
	if prv == nil {
		err := errors.New(fmt.Sprintf("private key is empty"))
		return nil, nil, nil, err
	}

	hash := cret.hash()
	if (hash == types.Hash{}) {
		err := errors.New(fmt.Sprintf("the hash of credential is empty"))
		return nil, nil, nil, err
	}

	sig, err := crypto.Sign(hash[:], prv)
	if err != nil {
		return nil, nil, nil, err
	}

	err = cret.get(sig)
	if err != nil {
		return nil, nil, nil, err
	}
	R = cret.R
	S = cret.S
	V = cret.V

	return R, S, V, nil
}

func (cret *CredentialSign) Sender() (types.Address, error) {
	return cret.sender()
}

func (cret *CredentialSign) sender() (types.Address, error) {
	cret.checkObj()
	if Config().chainId != nil && deriveChainId(&cret.V.IntVal).Cmp(Config().chainId) != 0 {
		return types.Address{}, ErrInvalidChainId
	}
	if Config().chainId == nil {
		panic("Config().chainId == nil")
	}
	V := &big.Int{}
	if Config().chainId.Sign() != 0 {
		V = V.Sub(&cret.V.IntVal, Config().chainIdMul)
		V.Sub(V, common.Big35)
	} else {
		V = V.Sub(&cret.V.IntVal, common.Big27)
	}
	address, err := recoverPlain(cret.hash(), &cret.R.IntVal, &cret.S.IntVal, V, true)
	return address, err
}

func (cret *CredentialSign) senderBySeedBlock(seedBlock *block.Header) (types.Address, error) {
	if seedBlock == nil {
		return cret.sender()
	}
	cret.checkObj()
	if Config().chainId != nil && deriveChainId(&cret.V.IntVal).Cmp(Config().chainId) != 0 {
		return types.Address{}, ErrInvalidChainId
	}
	if Config().chainId == nil {
		panic("Config().chainId == nil")
	}
	V := &big.Int{}
	if Config().chainId.Sign() != 0 {
		V = V.Sub(&cret.V.IntVal, Config().chainIdMul)
		V.Sub(V, common.Big35)
	} else {
		V = V.Sub(&cret.V.IntVal, common.Big27)
	}

	seed_r, _, err := generateSeedBySeedHeader(seedBlock)
	if err != nil {
		return types.Address{}, err
	}

	cretforhash := &CredentialSigForHash{
		cret.Round,
		cret.Step,
		seed_r.Bytes(),
	}
	hash, err := common.MsgpHash(cretforhash)
	if err != nil {
		return types.Address{}, err
	}

	address, err := recoverPlain(hash, &cret.R.IntVal, &cret.S.IntVal, V, true)
	return address, err
}

func (cret *CredentialSign) hash() types.Hash {
	R := uint64(Config().R)
	currentRound := cret.Round
	seedRound := cret.Round
	if currentRound < R {
		seedRound = 0
	} else {
		seedRound = currentRound - 1 - (currentRound % R)
	}
	seed_r, _, err := restoreSeed(seedRound)
	if err != nil {
		logger.Error("get Quantity fail")
		return types.Hash{}
	}
	cretforhash := &CredentialSigForHash{
		cret.Round,
		cret.Step,
		seed_r.Bytes(), // TODO: to get Quantity !!!!!!!!!!!!!!! need to implement a global function(round)
	}
	hash, err := common.MsgpHash(cretforhash)
	if err != nil {
		return types.Hash{}
	}
	return hash
}

// TODO: In current, EphemeralSig is the same as the Credential, need to be modified in the next version
// ephemeral key singer
type EphemeralSign struct {
	Signature
	round uint64 // round
	step  uint64 // step
	val   []byte // Val = Hash(B), or Val = 0, or Val = 1
}

type EphemeralSigForHash struct {
	Round uint64 // round
	Step  uint64 // step
	Val   []byte // Val = Hash(B), or Val = 0, or Val = 1
}

func (esig *EphemeralSign) GetStep() uint64 {
	return esig.step
}

func (esig *EphemeralSign) sign(prv *ecdsa.PrivateKey) (R *types.BigInt, S *types.BigInt, V *types.BigInt, err error) {
	if prv == nil {
		err := errors.New(fmt.Sprintf("private key is empty"))
		return nil, nil, nil, err
	}

	hash := esig.hash()
	if (hash == types.Hash{}) {
		err := errors.New(fmt.Sprintf("the hash of credential is empty"))
		return nil, nil, nil, err
	}

	sig, err := crypto.Sign(hash[:], prv)
	if err != nil {
		return nil, nil, nil, err
	}

	err = esig.get(sig)
	if err != nil {
		return nil, nil, nil, err
	}
	R = esig.R
	S = esig.S
	V = esig.V

	return R, S, V, nil
}

func (esig *EphemeralSign) sender() (types.Address, error) {
	esig.checkObj()

	if Config().chainId != nil && deriveChainId(&esig.V.IntVal).Cmp(Config().chainId) != 0 {
		return types.Address{}, ErrInvalidChainId
	}

	V := &big.Int{}
	if Config().chainId.Sign() != 0 {
		V = V.Sub(&esig.V.IntVal, Config().chainIdMul)
		V.Sub(V, common.Big35)
	} else {
		V = V.Sub(&esig.V.IntVal, common.Big27)
	}
	address, err := recoverPlain(esig.hash(), &esig.R.IntVal, &esig.S.IntVal, V, true)
	return address, err
}

func (esig *EphemeralSign) hash() types.Hash {
	if esig.val == nil {
		panic(fmt.Errorf("EphemeralSign obj is not initialized"))
	}
	eisgforhash := &EphemeralSigForHash{
		esig.round,
		esig.step,
		esig.val,
	}
	hash, err := common.MsgpHash(eisgforhash)
	if err != nil {
		return types.Hash{}
	}
	return hash
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
	// recover the public key from the signature
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
