package test

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/Evi1/awsl/clients"
	"github.com/Evi1/awsl/config"
	"github.com/Evi1/awsl/model"
	"github.com/Evi1/awsl/servers"
)

func TestH2C(t *testing.T) {
	config.Debug = true
	conf := config.GetConf()
	conf.NoVerify = true
	server := servers.NewH2C("127.0.0.1", "1928", "h2c", "123", "server.key", "server.crt", 32)
	l := server.Listen()
	go func() {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println(err)
		}
		buf := make([]byte, 65536)
		for {
			n, err := conn.Read(buf)
			fmt.Println("server read", string(buf[:n]), err)
			if err != nil && err != io.EOF {
				fmt.Println(err)
			}
			conn.Write(append([]byte("rep "), buf[:n]...))
			fmt.Println("server write")
		}
	}()
	fmt.Println("start client")
	client := clients.NewH2C("lo.bilibili.network", "1928", "h2c", "123")
	conn, err := client.Dial(model.ANetAddr{Host: "www.bilibili.network", Port: 443, Typ: model.IPV4ADDR})
	fmt.Println(err)
	if err != nil {
		return
	}
	buf := make([]byte, 65536)
	fmt.Println("start client read")
	go func() {
		for {
			n, err := conn.Read(buf)
			fmt.Println("client read", string(buf[:n]), err)
			if err != nil && err != io.EOF {
				fmt.Println(err)
			}
		}
	}()
	var i int
	fmt.Println("start client write")
	for i < 10 {
		conn.Write([]byte("client : " + time.Now().Format("Mon Jan 2 15:04:05 -0700 MST 2006")))
		fmt.Println("client write")
		i++
		time.Sleep(time.Second)
	}
}

func TestClient(t *testing.T) {
	url := "https://http2.golang.org/ECHO"
	// Create a pipe - an object that implements `io.Reader` and `io.Writer`.
	// Whatever is written to the writer part will be read by the reader part.
	pr, pw := io.Pipe()

	// Create an `http.Request` and set its body as the reader part of the
	// pipe - after sending the request, whatever will be written to the pipe,
	// will be sent as the request body.
	// This makes the request content dynamic, so we don't need to define it
	// before sending the request.
	req, err := http.NewRequest(http.MethodPut, url, ioutil.NopCloser(pr))
	if err != nil {
		log.Fatal(err)
	}

	// Send the request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Got: %d", resp.StatusCode)

	// Run a loop which writes every second to the writer part of the pipe
	// the current time.
	go func() {
		for {
			time.Sleep(1 * time.Second)
			fmt.Fprintf(pw, "It is now %v\n", time.Now())
		}
	}()

	// Copy the server's response to stdout.
	_, err = io.Copy(os.Stdout, resp.Body)
	log.Fatal(err)
}
