package dns

// Result Result
type Result struct {
	V4 string
	V6 string
}

// DNS DNS
type DNS interface {
	Resolve(host string) (Result, error)
}
