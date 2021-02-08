package arouter

import (
	"context"
	"testing"

	"github.com/rikaaa0928/awsl/aconn"
	"github.com/rikaaa0928/awsl/config"
	"github.com/rikaaa0928/awsl/global"
)

func TestNewRouter(t *testing.T) {
	conf := config.NewJsonConfig()
	err := conf.Open("./test.json")
	if err != nil {
		panic(err)
	}
	ctx := context.WithValue(context.Background(), global.CTXInTag, "in1")
	NewRouter(conf)(ctx, aconn.NewAddr("www.baidu.com", 1234, "tcp"))
}
