package apos

import (
	"container/heap"
	"fmt"
	"sort"
	"testing"
)

func TestBpWithPriorityHeap(t *testing.T) {
	hp := make(BpWithPriorityHeap, 0)
	heap.Init(&hp)
	for i := 0; i < 20; i++ {
		bp := new(BpWithPriority)
		bp.j = float64(i)
		heap.Push(&hp, bp)

	}
	//we just need the first one
	sort.Sort(hp)
	for index, v := range hp {
		fmt.Println("index :", index, " j :", v.j)
	}
}
