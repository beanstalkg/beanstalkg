package backend

import (
	"github.com/vimukthi-git/beanstalkg/architecture"
	"log"
	"math"
	"runtime/debug"
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
	h.DecreaseKey(h.Size - 1, item)
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
			h.clean()
			return temp
		}
	}

	return nil
}

// +++++++++++++ END - PriorityQueue Interface methods +++++++++++++++++

func (h *MinHeap) DecreaseKey(i int, item architecture.PriorityQueueItem) {
	// log.Println("queue", h, i)
	if item.Key() > h.Store[i].Key() {
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
	if l < len(h.Store) && h.Store[l].Key() < h.Store[i].Key() && h.Store[l].Key() != math.MaxInt64 {
		//log.Println("l=", l)
		//log.Println("h.Store[l]", h.Store[l])
		smallest = l
	} else {
		smallest = i
	}
	//log.Println(r, h.Size)
	//log.Println("l=", l)
	//log.Println("i=", i)
	if r < len(h.Store) && h.Store[r].Key() < h.Store[smallest].Key() && h.Store[r].Key() != math.MaxInt64 {
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
	if h.Size > 0 {
		if h.Store[0].Key() == math.MaxInt64 {
			log.Println("heap error - corrupted size", h)
			debug.PrintStack()
			h.Size = 0
			h.clean()
			return nil
		}
		return h.Store[0]
	} else {
		return nil
	}

}

func (h *MinHeap) clean() {
	// cleanup so that we don't waste memory
	for j := len(h.Store) - 1; j > 0; j-- {
		if (h.Store[j].Key() == math.MaxInt64 && j > h.Size) {
			h.Store = h.Store[:len(h.Store)-1]
		} else {
			break
		}
	}
}
