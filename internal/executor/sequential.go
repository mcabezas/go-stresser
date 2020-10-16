package executor

import (
	"stresser/internal/endpoint"
	"stresser/internal/iterator"
	"stresser/internal/statistics"
)

type sequential struct {
	endpoints iterator.Iterator
	Hits      int
	formatter statistics.Formatter
}

func NewSequential(hits int, endpoints []endpoint.EndPoint) Executor {
	return &sequential{
		endpoints: iterator.NewCyclicIterator(endpoints),
		Hits:      hits,
	}
}

func (s *sequential) Start() []statistics.Statistics {
	stats := make([]statistics.Statistics, s.Hits)
	for ii := 0; ii < s.Hits; ii++ {
		executable := s.endpoints.Next()
		stats[ii] = executable.Execute(s.formatter)
	}
	return stats
}
