// Copyright 2021 GoEdge CDN goedge.cdn@gmail.com. All rights reserved.

package utils_test

import (
	"github.com/dashenmiren/EdgeAPI/internal/utils"
	"github.com/dashenmiren/EdgeCommon/pkg/dnsconfigs"
	"testing"
)

func TestLookupCNAME(t *testing.T) {
	t.Log(utils.LookupCNAME("www.yun4s.cn"))
}

func TestLookupNS(t *testing.T) {
	t.Log(utils.LookupNS("goedge.cn", nil))
}

func TestLookupNSExtra(t *testing.T) {
	t.Log(utils.LookupNS("goedge.cn", []*dnsconfigs.DNSResolver{
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
	t.Log(utils.LookupTXT("yanzheng.goedge.cn", nil))
}
