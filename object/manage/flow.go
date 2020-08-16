package manage

import (
	"encoding/json"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/rikaaa0928/awsl/config"
	"github.com/rikaaa0928/awsl/tools"
)

// ServerFlowManager ServerFlowManager
var ServerFlowManager *SFlowManager

//var limit int

func init() {
	//limit = int(time.Hour / time.Second)
	ServerFlowManager = &SFlowManager{in: make(map[int]map[string]tools.Counter),
		out: make(map[int]map[string]tools.Counter),
		//inHistory:  make(map[int]map[string][]uint64, 0),
		//outHistory: make(map[int]map[string][]uint64, 0),
		inSum:     make(map[int]map[string]uint64),
		outSum:    make(map[int]map[string]uint64),
		inSecond:  make(map[int]map[string]uint64),
		outSecond: make(map[int]map[string]uint64),
		lock:      sync.RWMutex{}}
	config.GetConf()
	if config.Manage > 0 {
		go ServerFlowManager.Tick()
	}
}

// SFlowManager SFlowManager
type SFlowManager struct {
	in        map[int]map[string]tools.Counter
	out       map[int]map[string]tools.Counter
	inSecond  map[int]map[string]uint64
	outSecond map[int]map[string]uint64
	inSum     map[int]map[string]uint64
	outSum    map[int]map[string]uint64
	//inHistory  map[int]map[string][]uint64
	//outHistory map[int]map[string][]uint64
	allInSum  uint64
	allOutSum uint64
	lock      sync.RWMutex
}

func (fm *SFlowManager) add(id int, host string, count int64, m map[int]map[string]tools.Counter, sum map[int]map[string]uint64, second map[int]map[string]uint64) {
	fm.lock.RLock()
	hostMap, ok := m[id]
	fm.lock.RUnlock()
	if !ok {
		fm.lock.Lock()
		hostMap, ok = m[id]
		if !ok {
			m[id] = make(map[string]tools.Counter)
			sum[id] = make(map[string]uint64)
			second[id] = make(map[string]uint64)
			hostMap = m[id]
		}
		fm.lock.Unlock()
	}
	fm.lock.RLock()
	counter, ok := hostMap[host]
	fm.lock.RUnlock()
	if !ok {
		fm.lock.Lock()
		counter, ok = hostMap[host]
		if !ok {
			hostMap[host] = tools.NewCounter("atomic")
			counter = hostMap[host]
		}
		fm.lock.Unlock()
	}
	counter.Add(count)
}

// AddIn AddIn
func (fm *SFlowManager) AddIn(id int, host string, count int64) {
	fm.add(id, host, count, fm.in, fm.inSum, fm.inSecond)
	fm.allInSum += uint64(count)
}

// AddOut AddOut
func (fm *SFlowManager) AddOut(id int, host string, count int64) {
	fm.add(id, host, count, fm.out, fm.outSum, fm.outSecond)
	fm.allOutSum += uint64(count)
}

func (fm *SFlowManager) tickFor(m map[int]map[string]tools.Counter, sum map[int]map[string]uint64, second map[int]map[string]uint64) {
	for id := range m {
		for host := range m[id] {
			fm.lock.RLock()
			counter := m[id][host]
			fm.lock.RUnlock()
			num := counter.Set(0)
			fm.lock.Lock()
			second[id][host] = uint64(num)
			sum[id][host] += uint64(num)
			fm.lock.Unlock()
		}
	}
}

// Tick Tick
func (fm *SFlowManager) Tick() {
	t := time.Tick(time.Second)
	for {
		select {
		case <-t:
			go fm.tickFor(fm.in, fm.inSum, fm.inSecond)
			go fm.tickFor(fm.out, fm.outSum, fm.outSecond)
		}
	}
}

// GetRoot GetRoot
func (fm *SFlowManager) GetRoot() string {
	m := make(map[string]string)
	m["in sum"] = handleBytesNum(fm.allInSum)
	m["out sum"] = handleBytesNum(fm.allOutSum)
	bytes, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return err.Error()
	}
	return string(bytes)
}

type sIDMap struct {
	In     string
	Out    string
	SumIn  string
	SumOut string
}

// GetID GetID
func (fm *SFlowManager) GetID(id int) string {
	resultMap := make(map[string]sIDMap)
	fm.lock.RLock()
	//inMap := fm.inHistory[id]
	//outMap := fm.outHistory[id]
	inSecondMap := fm.inSecond[id]
	outSecondMap := fm.outSecond[id]
	idSumIn := uint64(0)
	idSumOut := uint64(0)
	secondInSum := uint64(0)
	secondOutSum := uint64(0)
	for k, inv := range inSecondMap {
		outv := outSecondMap[k]
		sumIn := fm.inSum[id][k]
		sumOut := fm.outSum[id][k]
		idSumIn += sumIn
		idSumOut += sumOut
		secondInSum += inv
		secondOutSum += outv
		resultMap[k] = sIDMap{In: handleBytesNum(inv), Out: handleBytesNum(outv), SumIn: handleBytesNum(sumIn), SumOut: handleBytesNum(sumOut)}
	}
	fm.lock.RUnlock()
	resultMap["0"] = sIDMap{In: handleBytesNum(secondInSum), Out: handleBytesNum(secondOutSum), SumIn: handleBytesNum(idSumIn), SumOut: handleBytesNum(idSumOut)}
	bytes, err := json.MarshalIndent(resultMap, "", "  ")
	if err != nil {
		return err.Error()
	}
	return string(bytes)
}

// GetIDHistory GetIDHistory
/*func (fm *SFlowManager) GetIDHistory(id int) string {
	fm.lock.RLock()
	defer fm.lock.RUnlock()
	inMap := fm.inHistory[id]
	outMap := fm.outHistory[id]
	res := "{\n"
	for k, inv := range inMap {
		outv := outMap[k]
		resM := make(map[string][]uint64)
		resM["ins"] = make([]uint64, len(inv))
		resM["outs"] = make([]uint64, len(outv))
		copy(resM["ins"], inv)
		copy(resM["outs"], outv)
		bytes, err := json.Marshal(resM)
		if err != nil {
			res += "\"" + k + "\":" + err.Error() + ",\n"
			continue
		}
		res += "\"" + k + "\":" + string(bytes) + ",\n"
	}
	li := strings.LastIndex(res, ",\n")
	if li >= 0 {
		res = res[:li] + res[li+1:]
	}
	return res + "}"
}*/

var mark = []string{"B", "KB", "MB", "GB", "TB", "PB", "EB"}

func handleBytesNum(num uint64) string {
	i := 0
	fNum := float64(num)
	for i < len(mark)-1 && fNum > 1024 {
		i++
		fNum /= 1024
	}
	str := strconv.FormatFloat(fNum, 'f', 2, 64)
	if strings.HasSuffix(str, ".00") {
		str = str[:len(str)-3]
	}
	return str + mark[i]
}

func handleBytesNumList(list []uint64) []string {
	res := make([]string, len(list))
	for i, v := range list {
		res[i] = handleBytesNum(v)
	}
	return res
}
