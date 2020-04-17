package object

import (
	"io"
	"log"
	"net"
	"net/http"
	"strconv"
	"sync"

	"github.com/Evi1/awsl/clients"
	"github.com/Evi1/awsl/config"
	"github.com/Evi1/awsl/model"
	"github.com/Evi1/awsl/router"
	"github.com/Evi1/awsl/servers"
	"github.com/Evi1/awsl/tools"
)

// NewDefault NewDefault
func NewDefault(cs []clients.Client, ss []servers.Server, r router.Router) *DefaultObject {
	m := make([]chan DefaultRemoteMsg, len(cs))
	for i := range m {
		m[i] = make(chan DefaultRemoteMsg, config.GetConf().BufSize)
	}
	return &DefaultObject{
		C:         cs,
		S:         ss,
		R:         r,
		Msg:       m,
		CloseChan: make(chan int8),
		stop:      false,
	}
}

// DefaultObject default
type DefaultObject struct {
	C         []clients.Client
	S         []servers.Server
	R         router.Router
	Msg       []chan DefaultRemoteMsg
	CloseChan chan int8
	stop      bool
}

// DefaultRemoteMsg DEFAULT
type DefaultRemoteMsg struct {
	c        net.Conn
	a        model.ANetAddr
	rs       []int
	reRouted bool
	src      int
}

// Run object
func (o *DefaultObject) Run() {
	o.handelClient()
	o.handelServer()
}

// Stop object
func (o *DefaultObject) Stop() {
	defer func() {
		recover()
	}()
	o.stop = true
	close(o.CloseChan)
}

func (o *DefaultObject) handelClient() {
	for i := range o.C {
		go o.handelOneClient(i)
	}
}

func (o *DefaultObject) handelOneClient(i int) {
	for !o.stop {
		select {
		case m := <-o.Msg[i]:
			go func() {
				c, err := o.C[i].Dial(m.a)
				if err != nil {
					if len(m.rs) > 0 {
						r := m.rs[0]
						m.rs = m.rs[1:]
						m.reRouted = true
						o.Msg[r] <- m
						if config.Debug {
							log.Printf("swith route to %d.\n", r)
						}
						return
					}
					m.c.Close()
					log.Println("client Dial error. client no.", i, " error = ", err)
					return
				}
				// temp route
				tr, ok := o.R.(router.TempRoute)
				if m.reRouted && ok {
					tr.TempRoute(m.src, m.a, i)
				}
				// client connection
				err = o.C[i].Verify(c)
				if err != nil {
					m.c.Close()
					log.Println("client Verify error. client no.", i, " error = ", err)
					c.Close()
					return
				}
				if hc, ok := m.c.(*servers.HTTPGetConn); ok {
					trans := http.Transport{Dial: func(network, addr string) (net.Conn, error) {
						return c, nil
					}}
					resp, err := trans.RoundTrip(hc.R)
					if err != nil {
						http.Error(hc.W, err.Error(), http.StatusServiceUnavailable)
						return
					}
					defer resp.Body.Close()
					tools.CopyHeader(hc.W.Header(), resp.Header)
					hc.W.WriteHeader(resp.StatusCode)
					io.Copy(hc.W, resp.Body)
					hc.Close()
				} else {
					go tools.PipeThenClose(m.c, c)
					tools.PipeThenClose(c, m.c)
				}

			}()
		case <-o.CloseChan:
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
			if err != servers.ErrUDP {
				log.Println("server Accept error. server no.", i, " error = ", err)
			}
			if c != nil {
				c.Close()
			}
			continue
		}
		go func() {
			addr, err := o.S[i].ReadRemote(c)
			if err != nil {
				log.Println("server ReadRemote error. server no.", i, " error = ", err)
				c.Close()
				return
			}
			rs := o.R.Route(i, addr)
			if len(rs) == 0 {
				log.Printf("Fatal error, no route for %d, %s.\n", i, addr.Host)
				return
			}
			r := rs[0]
			rs = rs[1:]
			if r > len(o.Msg)-1 {
				r = 0
			}
			o.Msg[r] <- DefaultRemoteMsg{c: c, a: addr, rs: rs, src: i}
		}()
	}
	l.Close()
	w.Done()
}
