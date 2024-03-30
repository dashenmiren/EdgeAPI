package utils_test

import (
	"testing"

	"github.com/TeaOSLab/EdgeAPI/internal/utils"
	"github.com/TeaOSLab/EdgeCommon/pkg/dnsconfigs"
)

func TestLookupCNAME(t *testing.T) {
	t.Log(utils.LookupCNAME("google.com"))
}

func TestLookupNS(t *testing.T) {
	t.Log(utils.LookupNS("google.com", nil))
}

func TestLookupNSExtra(t *testing.T) {
	t.Log(utils.LookupNS("google.com", []*dnsconfigs.DNSResolver{
		{
			Host: "192.168.2.2",
		},
		{
			Host: "192.168.2.2",
			Port: 58,
		},
		{
			Host: "8.8.8.8",
			Port: 53,
		},
	}))
}

func TestLookupTXT(t *testing.T) {
	t.Log(utils.LookupTXT("google.com", nil))
}
