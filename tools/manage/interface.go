package manage

import (
	"github.com/Evi1/awsl/object"
)

// ObjectManager ObjectManager
type ObjectManager interface {
	RunObject(object.Object)
}
