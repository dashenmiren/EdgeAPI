// Copyright 2022 GoEdge CDN goedge.cdn@gmail.com. All rights reserved. Official site: https://cdn.foyeseo.com .

package services

import (
	"github.com/dashenmiren/EdgeAPI/internal/utils"
	"github.com/iwind/TeaGo/assert"
	"testing"
)

func TestHTTPCacheTaskService_ParseDomain(t *testing.T) {
	var a = assert.NewAssertion(t)

	a.IsTrue(utils.ParseDomainFromKey("aaa") == "aaa")
	a.IsTrue(utils.ParseDomainFromKey("AAA") == "aaa")
	a.IsTrue(utils.ParseDomainFromKey("a.b-c.com") == "a.b-c.com")
	a.IsTrue(utils.ParseDomainFromKey("a.b-c.com/hello/world") == "a.b-c.com")
	a.IsTrue(utils.ParseDomainFromKey("https://a.b-c.com") == "a.b-c.com")
	a.IsTrue(utils.ParseDomainFromKey("http://a.b-c.com/hello/world") == "a.b-c.com")
	a.IsTrue(utils.ParseDomainFromKey("http://a.B-c.com/hello/world") == "a.b-c.com")
	a.IsTrue(utils.ParseDomainFromKey("http:/aaaa.com") == "http")
	a.IsTrue(utils.ParseDomainFromKey("北京") == "")
}
