package servers

import (
	"context"

	"github.com/Evi1/awsl/model"
)

// NewServers new servers
func NewServers(ctx context.Context, conf []model.In) []Server {
	r := make([]Server, 0, len(conf))
	for i, v := range conf {
		switch v.Type {
		case "socks":
			r = append(r, Socks{IP: v.Host, Port: v.Port, tag: v.Tag, id: i})
		case "socks5tcp":
			r = append(r, Socks5TCP{IP: v.Host, Port: v.Port, tag: v.Tag, id: i})
		case "socks5":
			r = append(r, Socks{IP: v.Host, Port: v.Port, tag: v.Tag, id: i})
		case "awsl":
			r = append(r, NewAWSL(ctx, v, i))
		case "http":
			r = append(r, NewHTTP(ctx, v, i))
		/*case "ahl":
		r = append(r, NewAHL(v.Host, v.Port, v.Awsl.URI, v.Awsl.Auth, v.Awsl.Key, v.Awsl.Cert, v.Awsl.Chan))*/
		case "h2c":
			r = append(r, NewH2C(ctx, v, i))
		case "tcp":
			r = append(r, NewTCP(v.Host, v.Port, v.TCP.Auth, v.Tag, i))
		default:
			panic(v.Type)
		}
	}
	return r
}
