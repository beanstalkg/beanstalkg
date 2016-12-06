package backend

import (
	"fmt"
	"github.com/vimukthi-git/beanstalkg/architecture"
	"testing"
	"time"
)

type testHeapItem struct {
	key       int64
	id        string
	timestamp int64
}

func (t testHeapItem) Key() int64 {
	return t.key
}

func (t testHeapItem) Id() string {
	return t.id
}

func (t testHeapItem) Timestamp() int64 {
	return t.timestamp
}

func (t testHeapItem) Enqueued() {
	t.timestamp = time.Now().UnixNano()
}

func (t testHeapItem) Dequeued() {
}

/**
5
INSERT 4
INSERT 9
DELETE 4
*/
func TestMinHeap_Insert(t *testing.T) {
	m := MinHeap{}
	m.Enqueue(testHeapItem{4, string(1), time.Now().UnixNano()})
	fmt.Println(m.Min().Key())
	if m.Min().Key() != 4 {
		t.Fail()
	}
	m.Enqueue(testHeapItem{9, string(2), time.Now().UnixNano()})
	fmt.Println(m.Min())
	if m.Min().Key() != 4 {
		t.Fail()
	}
	m.Delete(string(1))
	if m.Size() != 1 {
		t.Fail()
	}
	fmt.Println(m.Min().Key())
	if m.Min().Key() != 9 {
		t.Fail()
	}
	// m.Delete(string(2))
	fmt.Println(m.Dequeue().Key(), string(3))
}

func TestMinHeap_InsertCheckDelete(t *testing.T) {
	m := MinHeap{}
	m.Enqueue(testHeapItem{1, "one", time.Now().UnixNano()})
	m.Enqueue(testHeapItem{1, "two", time.Now().UnixNano()})
	m.Enqueue(testHeapItem{1, "three", time.Now().UnixNano()})
	m.Enqueue(testHeapItem{1, "four", time.Now().UnixNano()})
	fmt.Println(m)
	item := m.Dequeue().(testHeapItem)
	if item.Id() != "one" {
		t.Fail()
	}
	fmt.Println(item, m)
	item = m.Dequeue().(testHeapItem)
	if item.Id() != "two" {
		t.Fail()
	}
	fmt.Println(item, m)
	item = m.Dequeue().(testHeapItem)
	if item.Id() != "three" {
		t.Fail()
	}
	fmt.Println(item, m)
	item = m.Dequeue().(testHeapItem)
	if item.Id() != "four" {
		t.Fail()
	}
	fmt.Println(item, m)
	if m.Dequeue() != nil {
		t.Fail()
	}
	m.Enqueue(testHeapItem{1, "one", time.Now().UnixNano()})
	item = m.Dequeue().(testHeapItem)
	if item.Id() != "one" {
		t.Fail()
	}
	//item = m.Dequeue().(testHeapItem)
	//fmt.Println(item, m)
}

func TestIntegration(t *testing.T) {
	tube := architecture.Tube{"test", &MinHeap{}, &MinHeap{}, &MinHeap{}, &MinHeap{}, &MinHeap{}, make(map[string]*architecture.AwaitingClient)}
	//m.Enqueue(testHeapItem{4, string(1)})
	fmt.Println(tube)
	tube.Delayed.Enqueue(testHeapItem{4, string(1), time.Now().UnixNano()})
	if tube.Delayed.Dequeue().Key() != 4 {
		t.Fail()
	}
	fmt.Println(tube.Delayed)
	if tube.Delayed.Find(string(1)) != nil {
		t.Fail()
	}
	if tube.Delayed.Dequeue() != nil {
		t.Fail()
	}
}
