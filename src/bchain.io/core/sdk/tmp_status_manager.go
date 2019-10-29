package sdk

import (
	"bchain.io/common/types"
	"sync"
	"bchain.io/utils/database"
	"bchain.io/core/state"
	"bchain.io/utils/crypto"
)

/*
sdk is a last system status manager and modification manager,if you want get some last value esaily,
you can call a system function in sdk
*/




//TmpStatusManager,we can get a last value from TmpStatusManager ,or store a modified value (k-v) in TmpStatusManager
//by simple Api (f(stateRoot , contractAddress , key))

type TmpStatusManager struct {
	mu sync.RWMutex
	db database.IDatabaseGetter
	state *state.StateDB
	coinBase types.Address
	TmpConTracts map[types.Address]*TmpStatusNode
}

func NewTmpStatusManager(db database.IDatabaseGetter, state *state.StateDB , coinbase types.Address)*TmpStatusManager{
	t := new(TmpStatusManager)
	t.db = db
	t.state = state
	t.coinBase = coinbase
	t.TmpConTracts = make(map[types.Address]*TmpStatusNode)

	return t
}

//SetValue always set into memery
func (this *TmpStatusManager)SetValue(contractAddress types.Address , key []byte , value []byte)error{
	this.mu.Lock()
	defer this.mu.Unlock()

	//step 1: get a statusNode from manager
	statusNode := this.ExistContract(contractAddress)
	if statusNode == nil {
		//if not exist create One
		statusNode = this.CreateStatusNode(contractAddress)
	}

	//step 2: make TmpKey
	tmpKey := TmpKey{contractAddress:contractAddress , key:types.BytesToAddress(key)}

	//step 3:set value
	statusNode.SetValue(tmpKey , value)
	return nil
}


func (this *TmpStatusManager)GetValue(contractAddress types.Address , key []byte)[]byte{
	this.mu.RLock()
	defer this.mu.RUnlock()

	tmpKey := TmpKey{contractAddress:contractAddress , key:types.BytesToAddress(key)}

	tmpNode := this.ExistContract(contractAddress)
	if tmpNode != nil {
		dataExist :=  tmpNode.ExistValue(tmpKey)
		if dataExist != nil {
			return dataExist
		}
	}

	stateKey := crypto.Keccak256Hash(append(contractAddress.Bytes(), key...))
	LevelDbKey := this.state.GetState(contractAddress, stateKey)

	//if not find in memery,check in the LDB
	data , err := this.db.Get(LevelDbKey[:])
	if err != nil{
		logger.Error("db.Get(hashKey):" , err.Error(),"key:", LevelDbKey.String(),"stateKey:", stateKey.String())
		return nil
	}
	return data
}


//TmpStatusManager basic functions,should not control the mu(lock), the lock should hold by Upper caller


//check a Contract is exist in the tmpStatusManager
func (this *TmpStatusManager)ExistContract(contractAddress types.Address)*TmpStatusNode{
	if node , ok := this.TmpConTracts[contractAddress];ok{
		return node
	}
	return nil
}

func (this *TmpStatusManager)CreateStatusNode(contractAddress types.Address)*TmpStatusNode{
	node := NewStatusNode()
	this.TmpConTracts[contractAddress] = node
	return node
}
