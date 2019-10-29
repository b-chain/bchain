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
// @File: receipt.go
// @Date: 2018/05/08 15:18:08
////////////////////////////////////////////////////////////////////////////////

package transaction

import (
	"bytes"
	"fmt"
	"github.com/tinylib/msgp/msgp"
	"bchain.io/common/types"
	"bchain.io/utils/bloom"
)

//go:generate msgp

const (
	// ReceiptStatusFailed is the status code of a transaction if execution failed.
	ReceiptStatusFailed = uint(0)

	// ReceiptStatusSuccessful is the status code of a transaction if execution succeeded.
	ReceiptStatusSuccessful = uint(1)
)

// Receipt represents the results of a transaction.
type Receipt struct {
	// Consensus fields
	Status uint        `json:"status"`
	Bloom  types.Bloom `json:"logsBloom"         gencodec:"required"`
	Logs   []*Log      `json:"logs"              gencodec:"required"`

	// Implementation fields (don't reorder!)
	TxHash          types.Hash      `json:"transactionHash"   gencodec:"required"`
	ContractAddress []types.Address `json:"contractAddress"`
}

//ReceiptProtocol is the consensus encoding of a receipt
type ReceiptProtocol struct {
	Status uint
	Bloom  types.Bloom
	Logs   []*LogProtocol
}

// NewReceipt creates a barebone transaction receipt, copying the init fields.
func NewReceipt(failed bool) *Receipt {
	r := &Receipt{}
	if failed {
		r.Status = ReceiptStatusFailed
	} else {
		r.Status = ReceiptStatusSuccessful
	}
	return r
}

// String implements the Stringer interface.
func (r *Receipt) String() string {
	var tmpStr = fmt.Sprintf("receipt{status=%d   bloom=%x logs=%v}\n", r.Status, r.Bloom, r.Logs)
	tmpStr += fmt.Sprintf("receipt contains %d contract(s)\n", len(r.ContractAddress))
	for idx, s := range r.ContractAddress {
		tmpStr += fmt.Sprintf(" %d:%s", idx+1, s.Hex())
	}
	fmt.Println()
	return tmpStr
}

// Receipts is a wrapper around a Receipt array to implement DerivableList.
type Receipts []*Receipt

// Len returns the number of receipts in this list.
func (r Receipts) Len() int { return len(r) }

func (r Receipts) GetMsgp(i int) []byte {
	var buf bytes.Buffer
	logPs := []*LogProtocol{}
	for _, log := range r[i].Logs {
		logP := &LogProtocol{log.Address, log.Topics, log.Data}
		logPs = append(logPs, logP)
	}
	input := &ReceiptProtocol{r[i].Status, r[i].Bloom, logPs}
	err := msgp.Encode(&buf, input)
	if err != nil {
		return nil
	}
	return buf.Bytes()
}

func CreateBloom(receipts Receipts) types.Bloom {
	bloomIn := []bloom.BloomByte{}
	for _, receipt := range receipts {
		for _, log := range receipt.Logs {
			bloomIn = append(bloomIn, log.Address)
			for _, b := range log.Topics {
				bloomIn = append(bloomIn, b)
			}
		}
	}
	return bloom.CreateBloom(bloomIn)
}

//type Receipts_s [][]*Receipt
type ReceiptProtocols []*ReceiptProtocol
type Receipts_s struct {
	Receipts_s []ReceiptProtocols
}
