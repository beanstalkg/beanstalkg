package backend

import (
	"github.com/vimukthi-git/beanstalkg/architecture"
	"log"
	"math"
	//"os"
	//"runtime/pprof"
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
}

// +++++++++++++ START - PriorityQueue Interface methods +++++++++++++++++

func (h *MinHeap) Init() {
	//for i := 0; i < 100000; i++ {
	//	h.Store[i] = 1000000001
	//}
}

func (h *MinHeap) Enqueue(item architecture.PriorityQueueItem) {
	// h.Size = h.Size + 1
	h.Store = append(h.Store, ownHeapItem{math.MaxInt64, "-2"})
	h.DecreaseKey(h.Size() - 1, item)
}

func (h *MinHeap) Peek() architecture.PriorityQueueItem {
	return h.Min()
}

func (h *MinHeap) Dequeue() architecture.PriorityQueueItem {
	if h.Size() == 1 {
		min := h.Store[0]
		h.Store = nil
		return min
	} else if h.Size() > 1 {
		min := h.Store[0]
		h.Store[0] = h.Store[h.Size() - 1]
		h.Store = h.Store[:(h.Size() - 1)]
		h.MinHeapify(0)
		return min
	}
	return nil
}

func (h *MinHeap) Find(id string) architecture.PriorityQueueItem {
	for i := 0; i < h.Size(); i++ {
		if h.Store[i].Id() == id {
			return h.Store[i]
		}
	}
	return nil
}

func (h *MinHeap) Delete(id string) architecture.PriorityQueueItem {
	for i := 0; i < h.Size(); i++ {
		if h.Store[i].Id() == id {
			temp := h.Store[i]
			if i == 0 {
				h.Store = nil
			} else {
				h.Store[i] = h.Store[h.Size() - 1]
				h.MinHeapify(i)
			}
			return temp
		}
	}

	return nil
}

func (h *MinHeap) Size() int {
	return len(h.Store)
}

// +++++++++++++ END - PriorityQueue Interface methods +++++++++++++++++

func (h *MinHeap) DecreaseKey(i int, item architecture.PriorityQueueItem) {
	// log.Println("queue", h, i)
	if item.Key() > h.Store[i].Key() {
		log.Println(h, h.Store[i], item)
		log.Fatal("new key can not be larger than the current")
	}
	h.Store[i] = item
	//log.Println(h.Size, key)
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
	// log.Println("i=", i)
	l := h.Left(i)
	r := h.Right(i)
	// log.Println("l=", l)
	// log.Println("r=", r)
	smallest := 0
	if l < h.Size() && h.Store[l].Key() < h.Store[i].Key() {
		//log.Println("l=", l)
		//log.Println("h.Store[l]", h.Store[l])
		smallest = l
	} else {
		smallest = i
	}
	//log.Println(r, h.Size)
	//log.Println("l=", l)
	//log.Println("i=", i)
	if r < h.Size() && h.Store[r].Key() < h.Store[smallest].Key() {
		//log.Println("r=", r)
		//log.Println("h.Store[r]", h.Store[r])
		smallest = r
	}
	// log.Println("smallest=", smallest)
	if smallest != i {
		temp := h.Store[i]
		h.Store[i] = h.Store[smallest]
		h.Store[smallest] = temp
		h.MinHeapify(smallest)
	}
}

func (h *MinHeap) Min() architecture.PriorityQueueItem {
	if h.Size() > 0 {
		return h.Store[0]
	} else {
		return nil
	}

}
