package iterator

import (
	"sync"

	"stresser/internal/endpoint"
)

func NewCyclicIterator(endpoints []endpoint.EndPoint) Iterator {
	nodeFactory := newNodeFactory(endpoints)
	return &cyclic{current: nodeFactory.build(endpoints), mu: &sync.Mutex{}}
}

type cyclic struct {
	current Node
	mu      *sync.Mutex
}

func (ci *cyclic) Next() endpoint.EndPoint {
	ci.mu.Lock()
	value := ci.current.Value()
	ci.current = ci.current.Next()
	ci.mu.Unlock()
	return value
}
