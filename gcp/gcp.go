//
// gcp.go
// Copyright (C) 2017 weirdgiraffe <giraffe@cyberzoo.xyz>
//
// Distributed under terms of the MIT license.
//

package gcp

import (
	"fmt"
	"net"
	"os/exec"
	"regexp"
	"strings"
)

type dig interface {
	GetTXT(fqdn string) (string, error)
}

type execDig struct{}

func (d *execDig) GetTXT(fqdn string) (string, error) {
	out, err := exec.Command("dig", "@8.8.8.8", "txt", fqdn, "+short").Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}

type GCP struct {
	dig dig
	r   []*net.IPNet
}

func NewGCP() *GCP {
	return &GCP{
		dig: &execDig{},
		r:   defaultIpRanges,
	}
}

func (g *GCP) Name() string {
	return "Google Cloud Platform"
}

func (g *GCP) UpdateRanges() error {
	// check https://www.reddit.com/r/starcitizen/comments/3lce2k/list_of_google_cloud_ip_addresses_for_firewall/
	// or    https://gist.github.com/n0531m/f3714f6ad6ef738a3b0a
	out, err := g.dig.GetTXT("_cloud-netblocks.googleusercontent.com")
	if err != nil {
		return err
	}
	iprange := []string{}
	for _, fqdn := range parseFQDN(out) {
		addr, err := g.dig.GetTXT(fqdn)
		if err != nil {
			return err
		}
		ip, err := parseIP(addr)
		if err != nil {
			return err
		}
		iprange = append(iprange, ip...)
	}
	g.r = make([]*net.IPNet, len(iprange))
	for i := range iprange {
		if strings.LastIndex(iprange[i], "/") == -1 {
			// no subnet mask
			if strings.LastIndex(iprange[i], ":") == -1 {
				// ip v4
				iprange[i] += "/32"
			} else {
				// ip v6
				iprange[i] += "/128"
			}
		}
		g.r[i] = mustParseCIDR(iprange[i])
	}
	return nil
}

func (g *GCP) IsAt(addr string) bool {
	ip := net.ParseIP(addr)
	if ip == nil {
		panic("bad ip address")
	}
	for _, r := range g.r {
		if r.Contains(ip) {
			return true
		}
	}
	return false
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

func parseIP(txt string) (ip []string, err error) {
	re := regexp.MustCompile("ip[46]:([^ ]+?) ")
	l := re.FindAllStringSubmatch(txt, -1)
	if l == nil {
		err = fmt.Errorf("TXT record has no IP addresses in")
		return
	}
	ip = make([]string, len(l))
	for i := range l {
		ip[i] = l[i][1]
	}
	return ip, nil
}
