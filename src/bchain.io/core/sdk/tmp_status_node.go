package sdk

import (
	"bchain.io/common/types"
	"bchain.io/utils/crypto"
)

type TmpKey struct {
	contractAddress types.Address
	key types.Address
	stateRoot types.Hash    //nothing or a last stateRoot
}

func (this *TmpKey)MakeHashKey()(types.Hash , error){
	keyHexLen := types.AddressLength + len(this.key) + types.HashLength
	keyHex := make([]byte , keyHexLen)
	keyHex = keyHex[:0]

	keyHex = append(keyHex , this.contractAddress[:]...)
	keyHex = append(keyHex , this.key[:]...)
	keyHex = append(keyHex , this.stateRoot[:]...)


	hash := crypto.Keccak256Hash(keyHex)

	return hash , nil
}




type TmpStatusNode struct {
	Modified map[TmpKey][]byte
}

func NewStatusNode()*TmpStatusNode{
	n := new(TmpStatusNode)
	n.Modified = make(map[TmpKey][]byte)
	return n
}


func (this *TmpStatusNode)ExistValue(tmpKey TmpKey)[]byte{
	if value , ok := this.Modified[tmpKey];ok{
		return value
	}
	return nil
}


func (this *TmpStatusNode)SetValue(tmpKey TmpKey , value []byte){
	this.Modified[tmpKey] = value
}

