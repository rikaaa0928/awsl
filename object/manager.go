package object

import (
	"github.com/rikaaa0928/awsl/config"
)

type defaultManager struct {
	stop     bool
	stopChan chan int8
	o        Object
}

func (m *defaultManager) RunObject() {
	for !m.stop {
		go func() {
			m.o = NewObject(*config.GetConf())
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
