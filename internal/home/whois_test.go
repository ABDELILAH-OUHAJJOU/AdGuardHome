package home

import (
	"context"
	"testing"

	"github.com/AdguardTeam/AdGuardHome/internal/dnsforward"
	"github.com/stretchr/testify/assert"
)

func prepareTestDNSServer() error {
	config.DNS.Port = 1234
	Context.dnsServer = dnsforward.NewServer(dnsforward.DNSCreateParams{})
	conf := &dnsforward.ServerConfig{}
	conf.UpstreamDNS = []string{"8.8.8.8"}

	return Context.dnsServer.Prepare(conf)
}

// TODO(e.burkov): It's kind of complicated to get rid of network access in this
// test.  The thing is that *Whois creates new *net.Dialer each time it requests
// the server, so it becomes hard to simulate handling of request from test even
// with substituted upstream.  However, it must be done.
func TestWhois(t *testing.T) {
	assert.Nil(t, prepareTestDNSServer())

	w := Whois{timeoutMsec: 5000}
	resp, err := w.queryAll(context.Background(), "8.8.8.8")
	assert.Nil(t, err)
	m := whoisParse(resp)
	assert.Equal(t, "Google LLC", m["orgname"])
	assert.Equal(t, "US", m["country"])
	assert.Equal(t, "Mountain View", m["city"])
}
