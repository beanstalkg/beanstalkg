package server

import "github.com/vimukthi-git/beanstalkg/lib"

type Beanstalkg struct {
	Data lib.MinHeap
	Coms chan string
}

func (b *Beanstalkg) Init() {
	// initialize the go routine to handle Heap ops and the comm channel
}

func (b *Beanstalkg) ExecCommand(c Command) string {
	return "USING " + c.Params["tube"]
}