package clients

import "github.com/rikaaa0928/awsl/model"

// NewClients new clients
func NewClients(conf []model.Out) []Client {
	r := make([]Client, 0, len(conf))
	for i, v := range conf {
		switch v.Type {
		case "direct":
			r = append(r, DirectOut{id: i, tag: v.Tag})
		case "awsl":
			r = append(r, NewAWSL(i, v))
		/*case "ahl":
		r = append(r, NewAHL(v.Awsl.Host, v.Awsl.Port, v.Awsl.URI, v.Awsl.Auth))*/
		case "h2c":
			r = append(r, NewH2C(i, v))
		case "tcp":
			r = append(r, NewTCP(i, v))
		default:
			panic(v.Type)
		}
	}
	return r
}
