//
// main.go
// Copyright (C) 2017 weirdgiraffe <giraffe@cyberzoo.xyz>
//
// Distributed under terms of the MIT license.
//

package aws

import (
	"encoding/json"
	"net"
	"os"
	"strings"
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
	r *Record
}

func NewAWS() *AWS {
	return &AWS{
		r: loadDefaults(),
	}
}

func (a *AWS) IsAt(addr string) bool {
	ip := net.ParseIP(addr)
	if ip.To4() != nil {
		return a.isAtIPv4(ip.To4())
	}
	if ip.To16() != nil {
		return a.isAtIPv6(ip.To16())
	}
	panic("bad ip")
}

func (a *AWS) isAtIPv4(addr net.IP) bool {
	for i := range a.r.Prefix {
		ip, ipnet, err := net.ParseCIDR(a.r.Prefix[i].Prefix)
		if err != nil {
			panic(err)
		}
		if ip != nil && ip.Equal(addr) {
			return true
		}
		if ipnet != nil && ipnet.Contains(addr) {
			return true
		}
	}
	return false
}

func (a *AWS) isAtIPv6(addr net.IP) bool {
	for i := range a.r.Ipv6Prefix {
		ip, ipnet, err := net.ParseCIDR(a.r.Ipv6Prefix[i].Prefix)
		if err != nil {
			panic(err)
		}
		if ip != nil && ip.Equal(addr) {
			return true
		}
		if ipnet != nil && ipnet.Contains(addr) {
			return true
		}
	}
	return false
}

func loadDefaults() *Record {
	r := strings.NewReader(defaultIPRanges)
	var rec Record
	err := json.NewDecoder(r).Decode(&rec)
	if err != nil {
		panic(err)
	}
	return &rec
}

func load(name string) *Record {
	f, err := os.Open(name)
	if err != nil {
		panic(err)
	}
	var r Record
	err = json.NewDecoder(f).Decode(&r)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	return &r
}
