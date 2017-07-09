//
// aws.go
// Copyright (C) 2017 weirdgiraffe <giraffe@cyberzoo.xyz>
//
// Distributed under terms of the MIT license.
//

package aws

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/weirdgiraffe/watdatcloud"
)

type Record struct {
	CreateDate string `json:"createDate"`
	Ipv6Prefix []IPv6 `json:"ipv6_prefixes"`
	Prefix     []IPv4 `json:"prefixes"`
	SyncToken  string `json:"syncToken"`
}

type IPv4 struct {
	Prefix  string `json:"ip_prefix"`
	Region  string `json:"region"`
	Service string `json:"service"`
}

type IPv6 struct {
	Prefix  string `json:"ipv6_prefix"`
	Region  string `json:"region"`
	Service string `json:"service"`
}

type AWS struct {
	loaded bool
}

func NewAWS() *AWS {
	return &AWS{}
}

func (a *AWS) LoadRanges() ([]watdatcloud.Range, error) {
	if !a.loaded {
		a.loaded = true
		return decode(strings.NewReader(defaultIpRanges))
	}
	res, err := http.Get("https://ip-ranges.amazonaws.com/ip-ranges.json")
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode == http.StatusOK {
		return decode(res.Body)
	}
	return nil, fmt.Errorf("Unexpected HTTP Status: %d", res.StatusCode)
}

func decode(r io.Reader) ([]watdatcloud.Range, error) {
	var rec Record
	err := json.NewDecoder(r).Decode(&rec)
	if err != nil {
		return nil, err
	}
	ret := []watdatcloud.Range{}
	for _, ip := range rec.Prefix {
		ret = append(ret, watdatcloud.Range{
			Provider: "AWS",
			Info:     ip.Region + " " + ip.Service,
			CIDR:     watdatcloud.AddrToCIDR(ip.Prefix),
		})
	}
	for _, ip := range rec.Ipv6Prefix {
		ret = append(ret, watdatcloud.Range{
			Provider: "AWS",
			Info:     ip.Region + " " + ip.Service,
			CIDR:     watdatcloud.AddrToCIDR(ip.Prefix),
		})
	}
	return ret, nil
}
