package server

import (
	"github.com/rikaaa0928/awsl/alistener"
)

type AServer interface {
	Listen() alistener.AListener
	Handler() AHandler
}
