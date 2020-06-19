package object

// Object c+s
type Object interface {
	Run()
	Stop()
}

// Manager Manager
type Manager interface {
	RunObject()
	Restart()
	Stop()
}
