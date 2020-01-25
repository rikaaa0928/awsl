package servers

import "github.com/Evi1/awsl/model"

// NewServers new servers
func NewServers(conf []model.In) []Server {
	r := make([]Server, 0, len(conf))
	for _, v := range conf {
		switch v.Type {
		case "socks5":
			r = append(r, Socke5Server{IP: v.Host, Port: v.Port})
		case "awsl":
			r = append(r, NewAWSL(v.Host, v.Port, v.Awsl.URI, v.Awsl.Key, v.Awsl.Cert))
		default:
			panic(v.Type)
		}
	}
	return r
}
