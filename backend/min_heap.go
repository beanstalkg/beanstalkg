package backend

import (
	"github.com/beanstalkg/beanstalkg/architecture"
	//"os"
	//"runtime/pprof"

	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("BEANSTALKG")

type MinHeap struct {
	Store []architecture.PriorityQueueItem

	tubeName string
}

// +++++++++++++ START - PriorityQueue Interface methods +++++++++++++++++

func (h *MinHeap) Init(tubeName string) {
	h.tubeName = tubeName
	//for i := 0; i < 100000; i++ {
	//	h.Store[i] = 1000000001
	//}
}

func (h *MinHeap) Enqueue(item architecture.PriorityQueueItem) {
	// h.Size = h.Size + 1
	h.DecreaseKey(item)
	item.Enqueued()
}

func (h *MinHeap) Peek() architecture.PriorityQueueItem {
	return h.Min()
}

func (h *MinHeap) Dequeue() architecture.PriorityQueueItem {
	if h.Size() == 1 {
		min := h.Store[0]
		h.Store = nil
		min.Dequeued()
		return min
	} else if h.Size() > 1 {
		min := h.Store[0]
		h.Store[0] = h.Store[h.Size()-1]
		h.Store = h.Store[:(h.Size() - 1)]
		h.MinHeapify(0)
		min.Dequeued()
		return min
	}
	return nil
}

func (h *MinHeap) Find(id string) architecture.PriorityQueueItem {
	for _, item := range h.Store {
		if item.Id() == id {
			return item
		}
	}
	return nil
}

func (h *MinHeap) Delete(id string) architecture.PriorityQueueItem {
	for i, item := range h.Store {
		if item.Id() == id {
			if len(h.Store) == 1 {
				h.Store = nil
			} else {
				// remove item in the middle
				h.Store = append(h.Store[:i], h.Store[i+1:]...)
				h.MinHeapify(i)
			}
			return item
		}
	}

	return nil
}

func (h *MinHeap) Size() int {
	return len(h.Store)
}

// +++++++++++++ END - PriorityQueue Interface methods +++++++++++++++++

func (h *MinHeap) DecreaseKey(item architecture.PriorityQueueItem) {
	// Index of next slot in slice.
	i := h.Size()

	h.Store = append(h.Store, item)

	// Re-sort slice to put the new item in the proper place.
	for i > 0 && h.Store[h.Parent(i)].Key() > h.Store[i].Key() {
		// Swap item locationss.
		h.Store[i], h.Store[h.Parent(i)] = h.Store[h.Parent(i)], h.Store[i]

		i = h.Parent(i)
	}
}

func (h *MinHeap) Parent(i int) int {
	return i >> 1
}

func (h *MinHeap) Left(i int) int {
	return 2*i + 1
}

func (h *MinHeap) Right(i int) int {
	return 2*i + 2
}

func (h *MinHeap) MinHeapify(i int) {
	// log.Println("i=", i, h.Store[i].Timestamp())
	l := h.Left(i)
	r := h.Right(i)
	// log.Println("l=", l)
	// log.Println("r=", r)
	smallest := i
	if l < h.Size() {
		if left, parent := h.Store[l], h.Store[i]; left.Key() < parent.Key() ||
			(left.Key() == parent.Key() &&
				left.Timestamp() < parent.Timestamp()) {
			// log.Println("l=", l, h.Store[l].Timestamp())
			smallest = l
		}
	}
	if r < h.Size() {
		if right, parent := h.Store[r], h.Store[smallest]; right.Key() < parent.Key() ||
			(right.Key() == parent.Key() &&
				right.Timestamp() < parent.Timestamp()) {
			// log.Println("r=", r, h.Store[r].Timestamp())
			smallest = r
		}
	}
	// log.Println("smallest=", smallest)
	if smallest != i {
		// log.Println("smallest=", smallest, ", i=", i)
		h.Store[i], h.Store[smallest] = h.Store[smallest], h.Store[i]

		h.MinHeapify(smallest)
	}
}

func (h *MinHeap) Min() architecture.PriorityQueueItem {
	if h.Size() > 0 {
		return h.Store[0]
	}

	return nil
}
