package object

import "github.com/Evi1/awsl/clients"

import "github.com/Evi1/awsl/servers"

import "net"

import "github.com/Evi1/awsl/tools"

import "log"

// DefaultObject default
type DefaultObject struct {
	C     clients.Client
	S     servers.Server
	Msg   chan DefaultRemoteMsg
	Close chan int8
	stop  bool
}

type DefaultRemoteMsg struct {
	c net.Conn
	h string
	p string
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
	for !o.stop {
		select {
		case m := <-o.Msg:
			log.Printf("%+v\n", m)
			c, err := o.C.Dail(m.h, m.p)
			if err != nil {
				log.Println(err)
				return
			}
			o.C.Verify(c)
			go tools.PipeThenClose(m.c, c)
			go tools.PipeThenClose(c, m.c)
		case <-o.Close:
			o.stop = true
		}
	}
}

func (o *DefaultObject) handelServer() {
	log.Println("server")
	l := o.S.Listen()
	for !o.stop {
		c, err := l.Accept()
		log.Println("income")
		if err != nil {
			// handle error
			log.Println(err)
			return
		}
		go func() {
			h, p, e := o.S.ReadRemote(c)
			if e != nil {
				return
			}
			o.Msg <- DefaultRemoteMsg{c: c, h: h, p: p}
		}()
	}
}
