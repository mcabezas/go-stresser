package iterator

import (
	"stresser/internal/endpoint"
)

type Iterator interface {
	Next() endpoint.EndPoint
}

type Node interface {
	Next() Node
	Value() endpoint.EndPoint
}

type node struct {
	value endpoint.EndPoint
	next  Node
}

func (n *node) Next() Node {
	return n.next
}

func (n *node) Value() endpoint.EndPoint {
	return n.value
}
