package backend

import (
	"github.com/vimukthi-git/beanstalkg/architecture"
	"log"
	"math"
	//"fmt"
)

/**
+++++ MIN HEAP BACKEND ++++++
Dont want to use built in Heap for now. Easy to do optimizations
*/

type ownHeapItem struct {
	key int64
	id  string
}

func (t ownHeapItem) Key() int64 {
	return t.key
}

func (t ownHeapItem) Id() string {
	return t.id
}

type MinHeap struct {
	Store []architecture.PriorityQueueItem
	Size  int
}

// +++++++++++++ START - PriorityQueue Interface methods +++++++++++++++++

func (h *MinHeap) Init() {
	//for i := 0; i < 100000; i++ {
	//	h.Store[i] = 1000000001
	//}
}

func (h *MinHeap) Enqueue(item architecture.PriorityQueueItem) {
	h.Size = h.Size + 1
	h.Store = append(h.Store, ownHeapItem{math.MaxInt64, "-2"})
	h.DecreaseKey(h.Size-1, item)
}

func (h *MinHeap) Peek() architecture.PriorityQueueItem {
	return h.Min()
}

func (h *MinHeap) Dequeue() architecture.PriorityQueueItem {
	if h.Size > 0 {
		min := h.Min()
		h.Delete(min.Id())
		return min
	}
	return nil
}

func (h *MinHeap) Find(id string) architecture.PriorityQueueItem {
	for i := 0; i < h.Size; i++ {
		if h.Store[i].Id() == id {
			return h.Store[i]
		}
	}
	return nil
}

func (h *MinHeap) Delete(id string) architecture.PriorityQueueItem {
	for i := 0; i < h.Size; i++ {
		if h.Store[i].Id() == id {
			temp := h.Store[i]
			h.Store[i] = ownHeapItem{math.MaxInt64, "-2"}
			h.MinHeapify(i)
			h.Size = h.Size - 1
			return temp
		}
	}
	return nil
}

// +++++++++++++ END - PriorityQueue Interface methods +++++++++++++++++

func (h *MinHeap) DecreaseKey(i int, item architecture.PriorityQueueItem) {
	if item.Key() > h.Store[i].Key() {
		log.Fatal("new key can not be larger than the current")
	}
	h.Store[i] = item
	//fmt.Println(h.Size, key)
	for i > 0 && h.Store[h.Parent(i)].Key() > h.Store[i].Key() {
		temp := h.Store[i]
		h.Store[i] = h.Store[h.Parent(i)]
		h.Store[h.Parent(i)] = temp
		i = h.Parent(i)
	}
}

func (h *MinHeap) Parent(i int) int {
	return int(math.Floor(float64(i / 2)))
}

func (h *MinHeap) Left(i int) int {
	return 2*i + 1
}

func (h *MinHeap) Right(i int) int {
	return 2*i + 2
}

func (h *MinHeap) MinHeapify(i int) {
	l := h.Left(i)
	r := h.Right(i)
	smallest := 0
	if l < h.Size && h.Store[l].Key() < h.Store[i].Key() {
		smallest = l
	} else {
		smallest = i
	}
	//fmt.Println(r, h.Size)
	//fmt.Println("l=", l)
	//fmt.Println("i=", i)
	//fmt.Println("smallest=", smallest)
	if r < h.Size && h.Store[r].Key() < h.Store[smallest].Key() {
		smallest = r
	}
	if smallest != i {
		temp := h.Store[i]
		h.Store[i] = h.Store[smallest]
		h.Store[smallest] = temp
		h.MinHeapify(smallest)
	}
}

func (h *MinHeap) Min() architecture.PriorityQueueItem {
	if h.Size > 0 {
		return h.Store[0]
	} else {
		return nil
	}

}
