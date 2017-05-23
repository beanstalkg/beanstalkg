package backend

import (
	"fmt"
	"testing"
	"time"

	"github.com/vimukthi-git/beanstalkg/architecture"
)

const numberOfInserts = 3

var tt = []architecture.PriorityQueueItem{
	ownHeapItem{4, "1", time.Now().UnixNano()},
	ownHeapItem{9, "2", time.Now().UnixNano()},
	ownHeapItem{3, "3", time.Now().UnixNano()},

	ownHeapItem{1, "one", time.Now().UnixNano()},
	ownHeapItem{1, "two", time.Now().UnixNano()},
	ownHeapItem{1, "three", time.Now().UnixNano()},
	ownHeapItem{1, "four", time.Now().UnixNano()},
}

func TestMinHeap_Insert(t *testing.T) {
	m := MinHeap{}

	var expectedMins = []int64{4, 4, 3}
	var tests = tt[:numberOfInserts]

	// Test inserting items into heap.
	for i, testItem := range tests {
		m.Enqueue(testItem)

		if key := m.Min().Key(); key != expectedMins[i] {
			t.Errorf("Expected %d, got %d.", expectedMins[i], key)
		}
	}
}

func TestMinHeap_Dequeue(t *testing.T) {
	m := MinHeap{}

	var tt_expectedIds = []string{"one", "two", "three", "four"}

	for _, testItem := range tt[numberOfInserts:] {
		m.Enqueue(testItem)
	}

	for _, expectedId := range tt_expectedIds {
		item := m.Dequeue()

		if id := item.Id(); id != expectedId {
			t.Errorf("Expected %s, got %s.", expectedId, id)
		}
	}

	if item := m.Dequeue(); item != nil {
		t.Errorf("Expected nil, got %#v", item)
	}
}

func TestIntegration(t *testing.T) {
	tube := architecture.Tube{"test", &MinHeap{}, &MinHeap{}, &MinHeap{}, &MinHeap{}, &MinHeap{}, make(map[string]*architecture.AwaitingClient)}
	//m.Enqueue(ownHeapItem{4, string(1)})
	fmt.Println(tube)
	tube.Delayed.Enqueue(ownHeapItem{4, string(1), time.Now().UnixNano()})
	if tube.Delayed.Dequeue().Key() != 4 {
		t.Error("invalid key")
	}
	fmt.Println(tube.Delayed)
	if tube.Delayed.Find(string(1)) != nil {
		t.Error("delayed.find failed")
	}
	if tube.Delayed.Dequeue() != nil {
		t.Error("delayed.dequeue failed")
	}
}
