//
// watdatcloud.go
// Copyright (C) 2017 weirdgiraffe <giraffe@cyberzoo.xyz>
//
// Distributed under terms of the MIT license.
//

package watdatcloud

import (
	"fmt"
	"net"
	"strings"
)

type Range struct {
	Provider string
	Info     string
	CIDR     string
}

type RangeLoader interface {
	LoadRanges() ([]Range, error)
}

type RangeNotFound struct {
	addr string
}

func (r RangeNotFound) Error() string {
	return fmt.Sprintf("No IP Range found for '%s'", r.addr)
}

type searchRange struct {
	Range
	net *net.IPNet
}

type RangeLookuper struct {
	ranges []searchRange
	loader []RangeLoader
}

func NewRangeLookuper(l ...RangeLoader) *RangeLookuper {
	return &RangeLookuper{loader: l}
}

func (r *RangeLookuper) UpdateRanges() error {
	nr := []searchRange{}
	for i := range r.loader {
		or, err := r.oneLoaderRanges(r.loader[i])
		if err != nil {
			return err
		}
		nr = append(nr, or...)
	}
	r.ranges = nr
	return nil
}

func (r *RangeLookuper) oneLoaderRanges(l RangeLoader) ([]searchRange, error) {
	rl, err := l.LoadRanges()
	if err != nil {
		return nil, err
	}
	ret := make([]searchRange, len(rl))
	for i := range rl {
		ret[i].Range = rl[i]
		_, ret[i].net, err = net.ParseCIDR(rl[i].CIDR)
		if err != nil {
			return nil, fmt.Errorf("Failed to parse CIDR '%s': %v", rl[i].CIDR, err)
		}
	}
	return ret, nil
}

func (r *RangeLookuper) Lookup(addr string) (Range, error) {
	ip := net.ParseIP(addr)
	if ip == nil {
		return Range{}, &RangeNotFound{addr}
	}
	for _, oneRange := range r.ranges {
		if oneRange.net.Contains(ip) {
			return oneRange.Range, nil
		}
	}
	return Range{}, &RangeNotFound{addr}
}

func AddrToCIDR(addr string) string {
	if strings.LastIndex(addr, "/") == -1 {
		// no subnet mask
		if strings.LastIndex(addr, ":") == -1 {
			// ip v4
			addr += "/32"
		} else {
			// ip v6
			addr += "/128"
		}
	}
	return addr
}
