package object

import (
	"github.com/Evi1/awsl/clients"
	"github.com/Evi1/awsl/router"
	"strconv"
)

import "github.com/Evi1/awsl/servers"

import "net"

import "github.com/Evi1/awsl/tools"

import "log"

// DefaultObject default
type DefaultObject struct {
	C     []clients.Client
	S     []servers.Server
	R     router.Router
	Msg   []chan DefaultRemoteMsg
	Close chan int8
	stop  bool
}

type DefaultRemoteMsg struct {
	c net.Conn
	a servers.ANetAddr
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
			log.Printf("%+v\n", m)
			c, err := o.C[i].Dail(m.a.Host, strconv.Itoa(m.a.Port))
			if err != nil {
				log.Println(err)
				return
			}
			o.C[i].Verify(c)
			go tools.PipeThenClose(m.c, c)
			go tools.PipeThenClose(c, m.c)
		case <-o.Close:
			o.stop = true
		}
	}
}

func (o *DefaultObject) handelServer() {
	i := 0
	for i < len(o.S)-1 {
		go o.handelOneServer(i)
		i++
	}
	o.handelOneServer(i)
}

func (o *DefaultObject) handelOneServer(i int) {
	log.Println("server")
	l := o.S[i].Listen()
	for !o.stop {
		c, err := l.Accept()
		log.Println("income")
		if err != nil {
			log.Println(err)
			return
		}
		go func() {
			addr, e := o.S[i].ReadRemote(c)
			if e != nil {
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
}
