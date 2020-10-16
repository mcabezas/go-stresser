package executor

import (
	"stresser/internal/statistics"
)

type Executor interface {
	Start() []statistics.Statistics
}
