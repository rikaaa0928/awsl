package model

// HTTPIn HttpIn
type HTTPIn struct {
	Chan int
}

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
	HTTP *HTTPIn
	Type string
}

// Out out
type Out struct {
	Type string
	Awsl *AwslOut
}

// Object object
type Object struct {
	Ins      []In
	Outs     []Out
	BufSize  int
	NoVerify bool
}
