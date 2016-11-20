package db

import (
	"testing"
	"fmt"
)

type testHeapItem struct {
	key int64
	id string
}

func (t testHeapItem) Key() int64 {
	return t.key
}

func (t testHeapItem) Id() string {
	return t.id
}

/**
5
INSERT 4
INSERT 9
DELETE 4
 */
func TestMinHeap_Insert(t *testing.T) {
	m := MinHeap{}
	m.Insert(testHeapItem{4, string(1)})
	fmt.Println(m.Min().Key())
	if m.Min().Key() != 4 {
		t.Failed()
	}
	m.Insert(testHeapItem{9, string(2)})
	fmt.Println(m.Min().Key())
	if m.Min().Key() != 4 {
		t.Failed()
	}
	m.Delete(string(1))
	if m.Size != 1 {
		t.Failed()
	}
	fmt.Println(m.Min().Key())
	if m.Min().Key() != 9 {
		t.Failed()
	}
}
