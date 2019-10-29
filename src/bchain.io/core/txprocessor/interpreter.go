package txprocessor

import (
	"bchain.io/core/transaction"
	"math/big"
)

type Interpreter interface {
	/*
	if a transaction want join in txpool , a priority we should have,
	what is the priority?it should determined by miner's config file.
	*/
	GetPriorityFromTransaction(tx *transaction.Transaction)*big.Int
	SetPriorityForTransaction(tx *transaction.Transaction)

}
