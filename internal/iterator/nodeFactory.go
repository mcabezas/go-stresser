package iterator

import "stresser/internal/endpoint"

type NodeFactory interface {
	build([]endpoint.EndPoint) Node
	canHandle([]endpoint.EndPoint) bool
}

func newNodeFactory(endpoints []endpoint.EndPoint) NodeFactory {
	builders := []NodeFactory{
		&cyclicEmptyFactory{},
		&cyclicSingleNodeFactory{},
		&cyclicMultiNodeFactory{},
	}
	for _, b := range builders {
		if b.canHandle(endpoints) {
			return b
		}
	}
	return &cyclicEmptyFactory{}
}

type cyclicEmptyFactory struct{}

func (enf *cyclicEmptyFactory) build([]endpoint.EndPoint) Node {
	node := &node{}
	node.next = node
	return node
}

func (enf *cyclicEmptyFactory) canHandle(endpoints []endpoint.EndPoint) bool {
	return len(endpoints) == 0
}

type cyclicSingleNodeFactory struct{}

func (snf *cyclicSingleNodeFactory) build(endpoints []endpoint.EndPoint) Node {
	node := &node{value: endpoints[0]}
	node.next = node
	return node
}

func (snf *cyclicSingleNodeFactory) canHandle(endpoints []endpoint.EndPoint) bool {
	return len(endpoints) == 1
}

type cyclicMultiNodeFactory struct{}

func (mnf *cyclicMultiNodeFactory) build(endpoints []endpoint.EndPoint) Node {
	size := len(endpoints)
	first := &node{value: endpoints[0]}
	lastNode := first
	for ii := 1; ii < size; ii++ {
		current := &node{value: endpoints[ii], next: lastNode}
		lastNode = current
	}
	first.next = lastNode
	return first
}

func (mnf *cyclicMultiNodeFactory) canHandle(endpoints []endpoint.EndPoint) bool {
	return len(endpoints) > 1
}
