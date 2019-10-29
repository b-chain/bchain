package txprocessor

import (
	"sync"
	"bchain.io/core/transaction"
	"math/big"
	"math/rand"
)

type testInterpreter struct {
	lock sync.RWMutex
}

func (this *testInterpreter)GetPriorityFromTransaction(tx *transaction.Transaction)*big.Int{

	return big.NewInt(int64(rand.Intn(1000)))
}

func (this *testInterpreter)SetPriorityForTransaction(tx *transaction.Transaction){
	tx.SetPriority(big.NewInt(int64(rand.Intn(1000))))
}




