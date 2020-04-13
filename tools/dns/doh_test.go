package dns

import "testing"

func Test_doh(t *testing.T) {
	d := DoH{URL: "https://cloudflare-dns.com/dns-query"}
	body, err := d.Resolve("www.baidu.com")
	if err != nil {
		t.Error(err)
	}
	t.Log(body)
}
