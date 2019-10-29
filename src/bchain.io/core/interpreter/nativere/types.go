package nativere

import (
	"bchain.io/common/types"
)

//go:generate msgp
type NativePara struct {
	FuncName string   `json:"func_name" msg:"funcName"`
	Args     []string `json:"args"      msg:"args"`
}

type DumpData struct {
	Amount       uint64
	PledgeData   *PledgeData
	PoolData     *PoolData
	ProducerData *ProducerData
}

// doubly-linked unit
type doubly_linked struct {
}

type PledgeData struct {
	//part1 pd value
	Parent types.Address
	Prev   types.Address
	Next   types.Address
	Amount uint64
}

type PoolData struct {
	//part1 pd value
	Producer  types.Address
	Prev      types.Address
	Next      types.Address
	Amount    uint64
	ChildHead types.Address
}

type ProducerData struct {
	//part1 pd value
	ChildHead    types.Address
	Amount       uint64
	ProducerCert uint64
}
