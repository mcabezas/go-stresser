package executor

import (
	"io"

	"stresser/internal/endpoint"
	"stresser/internal/iterator"
	"stresser/internal/statistics"
)

type parallel struct {
	endpoints     iterator.Iterator
	Hits          int
	maxConcurrent int64
	formatter     statistics.Formatter
	writer        io.Writer
}

func NewParallel(hits int, endpoints []endpoint.EndPoint, formatter statistics.Formatter, writer io.Writer) Executor {
	return &parallel{
		endpoints: iterator.NewCyclicIterator(endpoints),
		Hits:      hits,
		formatter: formatter,
		writer:    writer,
	}
}

func (s *parallel) Start() []statistics.Statistics {
	stats := make([]statistics.Statistics, s.Hits)
	resultsChan := make(chan statistics.Statistics, s.Hits)
	for ii := 0; ii < s.Hits; ii++ {
		go func() {
			executable := s.endpoints.Next()
			stat := executable.Execute(s.formatter)
			_, _ = s.writer.Write(stat.Format())
			resultsChan <- stat
		}()
	}
	for ii := 0; ii < s.Hits; ii++ {
		stats[ii] = <-resultsChan
	}
	return stats
}
