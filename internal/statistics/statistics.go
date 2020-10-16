package statistics

import (
	"time"
)

type Statistics interface {
	Format() []byte
}

func NewStatistics(endpointName string, startAt, endAt time.Time, status int, body []byte, formatter Formatter) Statistics {
	return &statistics{
		baseStatistics: newBaseStatistics(endpointName, startAt, endAt),
		status:         status,
		body:           body,
		formatter:      formatter,
	}
}

type baseStatistics struct {
	endpointName string
	startAt      time.Time
	endAt        time.Time
}

func newBaseStatistics(endpointName string, startAt, endAt time.Time) *baseStatistics {
	return &baseStatistics{
		endpointName: endpointName,
		startAt:      startAt,
		endAt:        endAt,
	}
}

type statistics struct {
	*baseStatistics
	status    int
	body      []byte
	formatter Formatter
}

func (s *statistics) Format() []byte {
	return s.formatter.Format(s.endpointName, s.startAt, s.endAt, s.status, s.body)
}

type errStatistics struct {
	*baseStatistics
	err       error
	formatter Formatter
}

func NewErrStatistics(endpointName string, startAt, endAt time.Time, err error, formatter Formatter) Statistics {
	return &errStatistics{
		baseStatistics: newBaseStatistics(endpointName, startAt, endAt),
		err:            err,
		formatter:      formatter,
	}
}

func (es *errStatistics) Format() []byte {
	return es.formatter.Format(es.endpointName, es.startAt, es.endAt, 0, []byte(es.err.Error()))
}
