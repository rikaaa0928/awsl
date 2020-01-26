package model

// AwslIn awsl
type AwslIn struct {
	Key  string
	Cert string
	URI  string
	Auth string
	Chan int
}

// AwslOut awsl
type AwslOut struct {
	Host string
	Port string
	URI  string
	Auth string
}

// In in
type In struct {
	Host string
	Port string
	Awsl *AwslIn
	Type string
}

// Out out
type Out struct {
	Type string
	Awsl *AwslOut
}

// Object object
type Object struct {
	Ins     []In
	Outs    []Out
	BufSize int
}
