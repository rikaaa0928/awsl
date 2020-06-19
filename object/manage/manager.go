package manage

import (
	"github.com/Evi1/awsl/config"
	"github.com/Evi1/awsl/object"
)

type defaultManager struct {
	stop     bool
	stopChan chan int8
	o        object.Object
}

func (m *defaultManager) RunObject() {
	for !m.stop {
		go func() {
			m.o = object.NewObject(*config.GetConf())
			m.o.Run()
		}()
		<-m.stopChan
	}
}

func (m *defaultManager) Restart() {
	m.stop = false
	m.o.Stop()
	m.stopChan <- 1
}

func (m *defaultManager) Stop() {
	m.stop = true
	m.o.Stop()
	m.stopChan <- 1
}
