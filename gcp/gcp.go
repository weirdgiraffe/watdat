//
// gcp.go
// Copyright (C) 2017 weirdgiraffe <giraffe@cyberzoo.xyz>
//
// Distributed under terms of the MIT license.
//

package gcp

import (
	"fmt"
	"regexp"

	"github.com/weirdgiraffe/watdatcloud"
)

type GCP struct {
	dig    dig
	loaded bool
}

func NewGCP() *GCP {
	return &GCP{
		dig: &execDig{},
	}
}

func (g *GCP) LoadRanges() ([]watdatcloud.Range, error) {
	if !g.loaded {
		g.loaded = true
		r := make([]watdatcloud.Range, len(defaultIpRanges))
		for i := range defaultIpRanges {
			r[i] = watdatcloud.Range{
				Provider: "Google Cloud Platform",
				CIDR:     watdatcloud.AddrToCIDR(defaultIpRanges[i]),
			}
		}
		return r, nil
	}
	// check https://www.reddit.com/r/starcitizen/comments/3lce2k/list_of_google_cloud_ip_addresses_for_firewall/
	// or    https://gist.github.com/n0531m/f3714f6ad6ef738a3b0a
	out, err := g.dig.GetTXT("_cloud-netblocks.googleusercontent.com")
	if err != nil {
		return nil, err
	}
	r := []watdatcloud.Range{}
	for _, fqdn := range parseFQDN(out) {
		blockAddr, err := g.dig.GetTXT(fqdn)
		if err != nil {
			return nil, err
		}
		blockRanges, err := parseIP(blockAddr)
		if err != nil {
			return nil, err
		}
		r = append(r, blockRanges...)
	}
	return r, nil
}

func parseFQDN(txt string) (name []string) {
	re := regexp.MustCompile("include:([^ ]+?) ")
	l := re.FindAllStringSubmatch(txt, -1)
	if l != nil {
		name = make([]string, len(l))
		for i := range l {
			name[i] = l[i][1]
		}
	}
	return
}

func parseIP(txt string) (r []watdatcloud.Range, err error) {
	re := regexp.MustCompile("ip[46]:([^ ]+?) ")
	l := re.FindAllStringSubmatch(txt, -1)
	if l == nil {
		err = fmt.Errorf("TXT record has no IP addresses in")
		return
	}
	r = make([]watdatcloud.Range, len(l))
	for i := range l {
		r[i] = watdatcloud.Range{
			Provider: "Google Cloud Platform",
			CIDR:     watdatcloud.AddrToCIDR(l[i][1]),
		}
	}
	return r, nil
}
