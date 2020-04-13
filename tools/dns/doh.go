package dns

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"
)

// DoH DoH
type DoH struct {
	URL string
}
type quest struct {
	Name string
	Type int
}
type answer struct {
	Name string
	Type int
	Data string
}

type dohResp struct {
	Status   int
	Question []quest
	Answer   []answer
}

// Resolve Resolve
func (d DoH) Resolve(host string) (Result, error) {
	r := Result{}
	wg := sync.WaitGroup{}
	wg.Add(2)
	var err error
	go func() {
		var e error
		r.V4, e = resolve(d.URL, host, "A")
		if e != nil {
			err = e
		}
		wg.Done()
	}()
	go func() {
		var e error
		r.V6, e = resolve(d.URL, host, "AAAA")
		if e != nil {
			err = e
		}
		wg.Done()
	}()
	wg.Wait()
	return r, err
}

func resolve(url, host, t string) (string, error) {
	req, err := http.NewRequest("GET", url+"?name="+host+"&type="+t, nil)
	if err != nil {
		return "", err
	}
	req.Header.Add("accept", "application/dns-json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", errors.New(string(bodyBytes))
	}
	return parse(bodyBytes)
}

func parse(str []byte) (string, error) {
	resp := dohResp{}
	json.Unmarshal(str, &resp)
	if resp.Status != 0 {
		return "", errors.New("dns error status=" + strconv.Itoa(resp.Status))
	}
	if len(resp.Answer) == 0 || len(resp.Question) == 0 {
		return "", errors.New("dns error len=0")
	}
	for _, v := range resp.Answer {
		if resp.Question[0].Type == v.Type {
			return v.Data, nil
		}
	}
	return "", errors.New("dns error " + strconv.Itoa(resp.Answer[0].Type))
}
