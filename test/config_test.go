package test

import "testing"

import "github.com/Evi1/awsl/config"

import "github.com/Evi1/awsl/object"

func TestConfig(t *testing.T) {
	o := object.NewObject(config.Conf)
	o.Run()
}
