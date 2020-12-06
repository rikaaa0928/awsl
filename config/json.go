package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"sync"
)

func NewJsonConfig() *JsonConfig {
	return &JsonConfig{m: make(map[string]json.RawMessage)}
}

type JsonConfig struct {
	m     map[string]json.RawMessage
	value []byte
	sync.RWMutex
}

func (c *JsonConfig) Open(file string) error {
	f, err := os.OpenFile(file, os.O_RDONLY, 0)
	if err != nil {
		return err
	}
	defer f.Close()
	c.value, err = ioutil.ReadAll(f)
	if err != nil {
		return err
	}
	err = json.Unmarshal(c.value, &c.m)
	if err != nil {
		return err
	}
	return nil
}

func (c *JsonConfig) Get(path ...string) (interface{}, error) {
	if len(path) <= 0 {
		return json.RawMessage(c.value), nil
	}
	var result interface{}
	c.Lock()
	defer c.Unlock()
	m := c.m

	for i := 0; i < len(path); i++ {
		k := path[i]
		v, ok := m[k]
		if !ok {
			return nil, fmt.Errorf("Wrong path: %v", path)
		}
		if i != len(path)-1 {
			err := json.Unmarshal(v, &m)
			if err != nil {
				nextKey := path[i+1]
				nextIntKey, err2 := strconv.Atoi(nextKey)
				if err2 != nil {
					return nil, errors.New(err.Error() + " & " + err2.Error())
				}
				var sliceMap []map[string]json.RawMessage
				err2 = json.Unmarshal(v, &sliceMap)
				if err2 != nil {
					return nil, errors.New(err.Error() + " & " + err2.Error())
				}
				if nextIntKey >= len(sliceMap) {
					return nil, errors.New(err.Error() + " & key: " + nextKey + " out of len")
				}
				m = sliceMap[nextIntKey]
				result = m
				i++
			}
		} else {
			result = v
		}
	}
	switch result.(type) {
	case map[string]json.RawMessage:
		bytes, err := json.Marshal(result)
		if err != nil {
			return nil, err
		}
		result = json.RawMessage(bytes)
	}
	return result, nil
}
func (c *JsonConfig) GetString(path ...string) (string, error) {
	v, err := c.Get(path...)
	if err != nil {
		return "", err
	}
	str := string(v.(json.RawMessage))
	str = strings.Trim(str, "\"")
	return str, nil
}
func (c *JsonConfig) GetInt(path ...string) (int64, error) {
	v, err := c.Get(path...)
	if err != nil {
		return 0, err
	}
	reInt, err := strconv.ParseInt(strings.Trim(string(v.(json.RawMessage)), "\""), 10, 64)
	if err != nil {
		return 0, err
	}
	return reInt, nil
}
func (c *JsonConfig) GetBool(path ...string) (bool, error) {
	v, err := c.Get(path...)
	if err != nil {
		return false, err
	}
	reBool, err := strconv.ParseBool(strings.Trim(string(v.(json.RawMessage)), "\""))
	if err != nil {
		return false, err
	}
	return reBool, nil
}
func (c *JsonConfig) GetFloat(path ...string) (float64, error) {
	v, err := c.Get(path...)
	if err != nil {
		return 0, err
	}
	re, err := strconv.ParseFloat(strings.Trim(string(v.(json.RawMessage)), "\""), 64)
	if err != nil {
		return 0, err
	}
	return re, nil
}
func (c *JsonConfig) GetSlice(path ...string) ([]string, error) {
	v, err := c.Get(path...)
	if err != nil {
		return nil, err
	}
	var re []string
	err = json.Unmarshal(v.(json.RawMessage), &re)
	if err != nil {
		return nil, err
	}
	return re, nil
}
func (c *JsonConfig) GetStrMap(path ...string) (map[string]string, error) {
	v, err := c.Get(path...)
	if err != nil {
		return nil, err
	}
	var re map[string]string
	err = json.Unmarshal(v.(json.RawMessage), &re)
	if err != nil {
		return nil, err
	}
	return re, nil
}
func (c *JsonConfig) GetMap(path ...string) (map[string]interface{}, error) {
	v, err := c.Get(path...)
	if err != nil {
		return nil, err
	}
	var re map[string]interface{}
	err = json.Unmarshal(v.(json.RawMessage), &re)
	if err != nil {
		return nil, err
	}
	return re, nil
}
func (c *JsonConfig) Watch(path ...string) (Watcher, error) {
	return nil, errors.New("not suported")
}
