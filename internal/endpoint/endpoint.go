package endpoint

import "stresser/internal/statistics"

type EndPoint interface {
	Execute(formatter statistics.Formatter) statistics.Statistics
}
