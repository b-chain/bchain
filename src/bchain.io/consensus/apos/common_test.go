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
// @File: common_test.go
// @Date: 2018/06/29 11:06:29
////////////////////////////////////////////////////////////////////////////////

package apos

import (
	"container/heap"
	"fmt"
	"math/big"
	"testing"
)

type testObj struct {
	aa int
	bb *big.Int
}

func TestPg(t *testing.T) {
	pq := make(priorityQueue, 0)
	heap.Init(&pq)

	for i := 0; i < 10; i++ {
		to0 := &testObj{i, big.NewInt(int64(i))}
		item0 := &pqItem{to0, to0.bb}
		heap.Push(&pq, item0)
	}

	for i := 0; i < 10; i++ {
		to1 := &testObj{100 - i, big.NewInt(int64(100 - i + 1))}
		item1 := &pqItem{to1, to1.bb}
		heap.Push(&pq, item1)
	}
	for pq.Len() > 0 {
		item := heap.Pop(&pq).(*pqItem)
		fmt.Println("output", item.priority, item.value.(*testObj).aa, item.value.(*testObj).bb)
	}
}

func TestPg1(t *testing.T) {
	pq := make(priorityQueue, 0)
	heap.Init(&pq)

	for i := 0; i < 10; i++ {
		to0 := &testObj{i, big.NewInt(int64(i))}
		item0 := &pqItem{to0, to0.bb}
		heap.Push(&pq, item0)
	}

	for i := 0; i < 10; i++ {
		to1 := &testObj{100 - i, big.NewInt(int64(100 - i + 1))}
		item1 := &pqItem{to1, to1.bb}
		heap.Push(&pq, item1)
	}
	for i := 0; i < pq.Len(); i++ {
		fmt.Println("output", pq[i].priority, pq[i].value.(*testObj).aa, pq[i].value.(*testObj).bb)
	}
}
