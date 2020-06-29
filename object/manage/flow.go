package manage

import (
	"sync"

	"github.com/Evi1/awsl/tools"
)

// FlowManager FlowManager
var FlowManager flowManager

func init() {
	FlowManager = flowManager{in: make(map[int]map[string]tools.Counter),
		out:        make(map[int]map[string]tools.Counter),
		inHistory:  make([]int64, 0),
		outHistory: make([]int64, 0),
		lock:       sync.Mutex{}}
}

type flowManager struct {
	in         map[int]map[string]tools.Counter
	out        map[int]map[string]tools.Counter
	inHistory  []int64
	outHistory []int64
	lock       sync.Mutex
}
