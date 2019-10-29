package sdk

import (
	"fmt"
	"bchain.io/common/types"
	"bchain.io/core/state"
	"bchain.io/utils/database"
	"testing"
)

func TestDataStore(t *testing.T) {
	//create a database
	db, err := database.OpenMemDB()
	if err != nil {
		panic(err)
	}
	addr := types.Address{}
	statedb, _ := state.New(types.Hash{}, state.NewDatabase(db))
	sdkHandler := NewTmpStatusManager(db, statedb, addr)
	contractAddr := types.Address{}
	contractAddr[0] = 1

	accountAddr := types.Address{}
	Sys_SetValue(sdkHandler, contractAddr, accountAddr[:], []byte{1, 2, 3, 4, 5})
	r := Sys_GetValue(sdkHandler, contractAddr, accountAddr[:])
	fmt.Println("r:", r)
}
