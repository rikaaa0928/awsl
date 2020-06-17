package servers

import (
	"context"

	"github.com/Evi1/awsl/model"
)

// NewServers new servers
func NewServers(ctx context.Context, conf []model.In) []Server {
	r := make([]Server, 0, len(conf))
	for _, v := range conf {
		switch v.Type {
		case "socks":
			r = append(r, SockeServer{IP: v.Host, Port: v.Port})
		case "socks5tcp":
			r = append(r, Socke5TCPServer{IP: v.Host, Port: v.Port})
		case "socks5":
			r = append(r, SockeServer{IP: v.Host, Port: v.Port})
		case "awsl":
			r = append(r, NewAWSL(ctx, v.Host, v.Port, v.Awsl.URI, v.Awsl.Auth, v.Awsl.Key, v.Awsl.Cert, v.Awsl.Chan))
		case "http":
			r = append(r, NewHTTP(ctx, v.Host, v.Port, v.HTTP.Chan))
		/*case "ahl":
		r = append(r, NewAHL(v.Host, v.Port, v.Awsl.URI, v.Awsl.Auth, v.Awsl.Key, v.Awsl.Cert, v.Awsl.Chan))*/
		case "h2c":
			r = append(r, NewH2C(ctx, v.Host, v.Port, v.Awsl.URI, v.Awsl.Auth, v.Awsl.Key, v.Awsl.Cert, v.Awsl.Chan))
		case "tcp":
			r = append(r, NewTCP(v.Host, v.Port, v.TCP.Auth))
		default:
			panic(v.Type)
		}
	}
	return r
}
