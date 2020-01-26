package object

import (
	"log"
	"net"
	"strconv"
	"sync"

	"github.com/Evi1/awsl/clients"
	"github.com/Evi1/awsl/model"
	"github.com/Evi1/awsl/router"
	"github.com/Evi1/awsl/servers"
	"github.com/Evi1/awsl/tools"
)

// NewDefault NewDefault
func NewDefault(cs []clients.Client, ss []servers.Server) *DefaultObject {
	m := make([]chan DefaultRemoteMsg, len(cs))
	for i := range m {
		m[i] = make(chan DefaultRemoteMsg, 10)
	}
	return &DefaultObject{
		C:     cs,
		S:     ss,
		R:     router.ARouter{},
		Msg:   m,
		Close: make(chan int8),
		stop:  false,
	}
}

// DefaultObject default
type DefaultObject struct {
	C     []clients.Client
	S     []servers.Server
	R     router.Router
	Msg   []chan DefaultRemoteMsg
	Close chan int8
	stop  bool
}

// DefaultRemoteMsg DEFAULT
type DefaultRemoteMsg struct {
	c net.Conn
	a model.ANetAddr
	r int
}

// Run object
func (o *DefaultObject) Run() {
	go o.handelClient()
	o.handelServer()
}

// Stop object
func (o *DefaultObject) Stop() {
	o.stop = true
	o.Close <- 0
}

func (o *DefaultObject) handelClient() {
	for i := range o.C {
		o.handelOneClient(i)
	}
}

func (o *DefaultObject) handelOneClient(i int) {
	for !o.stop {
		select {
		case m := <-o.Msg[i]:
			go func() {
				c, err := o.C[i].Dial(m.a)
				if err != nil {
					log.Println(err)
					return
				}
				err = o.C[i].Verify(c)
				if err != nil {
					log.Println(err)
					return
				}
				go tools.PipeThenClose(m.c, c)
				tools.PipeThenClose(c, m.c)
			}()
		case <-o.Close:
			o.stop = true
		}
	}
}

func (o *DefaultObject) handelServer() {
	var w sync.WaitGroup
	for i := range o.S {
		w.Add(1)
		go o.handelOneServer(i, &w)
	}
	w.Wait()
}

func (o *DefaultObject) handelOneServer(i int, w *sync.WaitGroup) {
	log.Println("start server: " + strconv.Itoa(i))
	l := o.S[i].Listen()
	for !o.stop {
		c, err := l.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go func() {
			addr, err := o.S[i].ReadRemote(c)
			if err != nil {
				log.Println(err)
				return
			}
			r := o.R.Route(addr)
			if r > len(o.Msg)-1 {
				r = 0
			}
			o.Msg[r] <- DefaultRemoteMsg{c: c, a: addr, r: r}
		}()
	}
	w.Done()
}
