package backend

import (
	"github.com/beanstalkg/beanstalkg/architecture"
)

const defaultBackend = "minheap"

var validBackends map[string]architecture.PriorityQueueCreator = map[string]architecture.PriorityQueueCreator{
	"minheap": func() architecture.PriorityQueue { return &MinHeap{} },
}

// QueueCreator retrieves the PriorityQueueCreator that returns a
// PriorityQueue for the specified backend.  Falls back to
// defaultBackend if the requested backend is invalid.
func QueueCreator(backend string) architecture.PriorityQueueCreator {
	if _, ok := validBackends[backend]; !ok {
		log.Debugf("%s backend not supported, falling back to %s.", backend, defaultBackend)
		backend = defaultBackend
	}

	return validBackends[backend]
}
