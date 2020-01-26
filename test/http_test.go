package test

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"testing"
	"time"
)

func transfer(destination io.WriteCloser, source io.ReadCloser) {
	defer destination.Close()
	defer source.Close()
	io.Copy(destination, source)
}
func TestHttp(t *testing.T) {
	fmt.Println("start")
	server := &http.Server{
		Addr: ":48880",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodConnect {
				fmt.Println(r.Host)
				dest_conn, err := net.DialTimeout("tcp", r.Host, 10*time.Second)
				if err != nil {
					http.Error(w, err.Error(), http.StatusServiceUnavailable)
					return
				}
				w.WriteHeader(http.StatusOK)
				hijacker, ok := w.(http.Hijacker)
				if !ok {
					http.Error(w, "Hijacking not supported", http.StatusInternalServerError)
					return
				}
				client_conn, _, err := hijacker.Hijack()
				if err != nil {
					http.Error(w, err.Error(), http.StatusServiceUnavailable)
				}
				go transfer(dest_conn, client_conn)
				go transfer(client_conn, dest_conn)
			} else {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}),
	}
	t.Log("listen")
	server.ListenAndServe()
}
