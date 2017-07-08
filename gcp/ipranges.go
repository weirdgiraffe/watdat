//
// ipranges.go
// Copyright (C) 2017 weirdgiraffe <giraffe@cyberzoo.xyz>
//
// Distributed under terms of the MIT license.
//

package gcp

import "net"

func mustParseCIDR(addr string) *net.IPNet {
	_, n, err := net.ParseCIDR(addr)
	if err != nil {
		panic(err)
	}
	return n
}

var defaultIpRanges = []*net.IPNet{
	mustParseCIDR("8.34.208.0/20"),
	mustParseCIDR("8.35.192.0/21"),
	mustParseCIDR("8.35.200.0/23"),
	mustParseCIDR("108.59.80.0/20"),
	mustParseCIDR("108.170.192.0/20"),
	mustParseCIDR("108.170.208.0/21"),
	mustParseCIDR("108.170.216.0/22"),
	mustParseCIDR("108.170.220.0/23"),
	mustParseCIDR("108.170.222.0/24"),
	mustParseCIDR("162.216.148.0/22"),
	mustParseCIDR("162.222.176.0/21"),
	mustParseCIDR("173.255.112.0/20"),
	mustParseCIDR("192.158.28.0/22"),
	mustParseCIDR("199.192.112.0/22"),
	mustParseCIDR("199.223.232.0/22"),
	mustParseCIDR("199.223.236.0/23"),
	mustParseCIDR("23.236.48.0/20"),
	mustParseCIDR("23.251.128.0/19"),
	mustParseCIDR("107.167.160.0/19"),
	mustParseCIDR("107.178.192.0/18"),
	mustParseCIDR("146.148.2.0/23"),
	mustParseCIDR("146.148.4.0/22"),
	mustParseCIDR("146.148.8.0/21"),
	mustParseCIDR("146.148.16.0/20"),
	mustParseCIDR("146.148.32.0/19"),
	mustParseCIDR("146.148.64.0/18"),
	mustParseCIDR("130.211.4.0/22"),
	mustParseCIDR("35.203.240.0/20"),
	mustParseCIDR("130.211.8.0/21"),
	mustParseCIDR("130.211.16.0/20"),
	mustParseCIDR("130.211.32.0/19"),
	mustParseCIDR("130.211.64.0/18"),
	mustParseCIDR("130.211.128.0/17"),
	mustParseCIDR("104.154.0.0/15"),
	mustParseCIDR("104.196.0.0/14"),
	mustParseCIDR("208.68.108.0/23"),
	mustParseCIDR("35.184.0.0/14"),
	mustParseCIDR("35.188.0.0/15"),
	mustParseCIDR("35.206.64.0/18"),
	mustParseCIDR("35.202.0.0/16"),
	mustParseCIDR("35.190.0.0/17"),
	mustParseCIDR("35.190.128.0/18"),
	mustParseCIDR("35.190.192.0/19"),
	mustParseCIDR("35.190.224.0/20"),
	mustParseCIDR("35.192.0.0/14"),
	mustParseCIDR("35.196.0.0/15"),
	mustParseCIDR("35.198.0.0/16"),
	mustParseCIDR("35.199.0.0/17"),
	mustParseCIDR("35.199.128.0/18"),
	mustParseCIDR("35.200.0.0/16"),
	mustParseCIDR("35.201.0.0/17"),
	mustParseCIDR("2600:1900::/35"),
}