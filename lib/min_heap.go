package lib

import (
	"log"
	"math"
)

type HeapItem interface {
	Key() int
}

type MinHeap struct {
	Store [100000]int
	Size int
}

func (h *MinHeap) Init() {
	for i := 0; i < 100000; i++ {
		h.Store[i] = 1000000001
	}
}

func (h *MinHeap) Insert(key int) {
	h.Size = h.Size + 1
	h.Store[h.Size - 1] = 1000000001 // maximum
	h.DecreaseKey(h.Size - 1, key)
}

func (h *MinHeap) Delete(key int) {
	for i := 0; i < h.Size; i++ {
		if h.Store[i] == key {
			h.Store[i] = 1000000001
			h.MinHeapify(i)
			break
		}
	}
}

func (h *MinHeap) Find(key int) int {
	for i := 0; i < h.Size; i++ {
		if h.Store[i] == key {
			return i
		}
	}
	return 100001
}

func (h *MinHeap) DecreaseKey(i, key int) {
	if key > h.Store[i] {
		log.Fatal("new key can not be larger than the current")
	}
	h.Store[i] = key
	//fmt.Println(h.Size, key)
	for ;i > 0 && h.Store[h.Parent(i)] > h.Store[i]; {
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
	return 2 * i + 1
}

func (h *MinHeap) Right(i int) int {
	return 2 * i + 2
}

func (h *MinHeap) MinHeapify(i int) {
	l := h.Left(i)
	r := h.Right(i)
	smallest := 0
	if l <= h.Size && h.Store[l] < h.Store[i] {
		smallest = l
	} else {
		smallest = i
	}
	if r <= h.Size && h.Store[r] < h.Store[smallest] {
		smallest = r
	}
	if smallest != i {
		temp := h.Store[i]
		h.Store[i] = h.Store[smallest]
		h.Store[smallest] = temp
		h.MinHeapify(smallest)
	}
}

func (h *MinHeap) Min() int {
	return h.Store[0]
}