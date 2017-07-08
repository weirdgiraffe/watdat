//
// gcp_test.go
// Copyright (C) 2017 weirdgiraffe <giraffe@cyberzoo.xyz>
//
// Distributed under terms of the MIT license.
//

package gcp

import (
	"strings"
	"testing"
)

func TestUpdate(t *testing.T) {
	g := NewGCP()
	err := g.UpdateRanges()
	if err != nil {
		t.Fatal(err)
	}
}

func TestBlocksResultParsing(t *testing.T) {
	txt := "v=spf1 include:_cloud-netblocks1.googleusercontent.com include:_cloud-netblocks2.googleusercontent.com include:_cloud-netblocks3.googleusercontent.com include:_cloud-netblocks4.googleusercontent.com include:_cloud-netblocks5.googleusercontent.com ?all"
	expected := []string{
		"_cloud-netblocks1.googleusercontent.com",
		"_cloud-netblocks2.googleusercontent.com",
		"_cloud-netblocks3.googleusercontent.com",
		"_cloud-netblocks4.googleusercontent.com",
		"_cloud-netblocks5.googleusercontent.com",
	}
	res := parseFQDN(txt)
	for i := range expected {
		found := false
		for j := range res {
			if strings.EqualFold(res[j], expected[i]) {
				found = true
				break
			}
		}
		if !found {
			t.Log(res)
			t.Fatalf("Not found expected block %s", expected[i])
		}
	}
}

func TestObeBlockResultParsing(t *testing.T) {
	txt := "v=spf1 ip4:35.190.0.0/17 ip4:35.190.128.0/18 ip4:35.190.192.0/19 ip4:35.190.224.0/20 ip4:35.192.0.0/14 ip4:35.196.0.0/15 ip4:35.198.0.0/16 ip4:35.199.0.0/17 ip4:35.199.128.0/18 ip4:35.200.0.0/16 ip4:35.201.0.0/17 ip6:2600:1900::/35 ?all"
	expected := []string{
		"35.190.0.0/17",
		"35.190.128.0/18",
		"35.190.192.0/19",
		"35.190.224.0/20",
		"35.192.0.0/14",
		"35.196.0.0/15",
		"35.198.0.0/16",
		"35.199.0.0/17",
		"35.199.128.0/18",
		"35.200.0.0/16",
		"35.201.0.0/17",
		"2600:1900::/35",
	}
	res, err := parseIP(txt)
	if err != nil {
		t.Fatal(err)
	}
	for i := range expected {
		found := false
		for j := range res {
			if strings.EqualFold(res[j], expected[i]) {
				found = true
				break
			}
		}
		if !found {
			t.Log(res)
			t.Fatalf("Not found expected block %s", expected[i])
		}
	}
}
