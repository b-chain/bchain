package aposgensis

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/tinylib/msgp/msgp"
	"io/ioutil"
	"bchain.io/common/types"
	"bchain.io/core/blockchain/block"
	"bchain.io/bchaind/defaults"
	"bchain.io/params"
	"bchain.io/utils/crypto"
	"os"
	"path/filepath"
)

//go:generate msgp

type WeightInfo struct {
	AddrStr types.Address `json:"addr"`
	Wt      uint64        `json:"wt"`
}

type WeightInfos []WeightInfo

var (
	BootCommitteeKey = []byte("apos_bootCommittee")
)

// Load boot committee weight config
func LoadBootCommittee() WeightInfos {
	path := filepath.Join(defaults.DefaultTOMLConfigPath, "bootCommittee.json")
	file, err := os.Open(path)
	if err != nil {
		//panic(fmt.Sprintf("boot committee config file(%s) is not exist...", path))
		logger.Critical("boot committee is not exist, please config...")
		logger.Critical("bchain will exit!")
		os.Exit(-1)
	}

	all, err := ioutil.ReadAll(file)
	if err != nil {
		panic(fmt.Sprintf("boot committee config file read err:%s", err.Error()))
	}

	wtConfig := WeightInfos{}
	err = json.Unmarshal(all, &wtConfig)
	if err != nil {
		panic(fmt.Sprintf("Unmarshal wtConfig err:%s", err.Error()))
	}

	return wtConfig
}

func MakeAposGenesisConsensusData() (bcd *block.ConsensusData, key, val []byte) {
	bcd = &block.ConsensusData{}
	bcd.Id = "apos"

	weight := LoadBootCommittee()
	//fmt.Println(weight)
	var encData bytes.Buffer
	err := msgp.Encode(&encData, &weight)
	if err != nil {
		panic(fmt.Sprintf("msgp encode err:%s", err.Error()))
	}
	seed := crypto.Keccak256(encData.Bytes())
	fmt.Printf("apos Genesis seed:0x%x \n", seed)

	sig, err := crypto.Sign(seed[:], params.RewordPrikey)
	if err != nil {
		panic(fmt.Sprintf("sign err:%s", err.Error()))
	}
	bcd.Para = sig

	key = BootCommitteeKey
	val = encData.Bytes()
	return bcd, key, val
}
