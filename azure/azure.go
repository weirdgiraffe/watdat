//
// azure.go
// Copyright (C) 2017 weirdgiraffe <giraffe@cyberzoo.xyz>
//
// Distributed under terms of the MIT license.
//

package azure

import (
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"regexp"
	"strings"
)

type Record struct {
	Xsi    string   `xml:"xsi,attr"`
	Xsd    string   `xml:"xsd,attr"`
	Region []Region `xml:"Region"`
}

type Region struct {
	Name    string    `xml:"Name,attr"`
	IpRange []IpRange `xml:"IpRange"`
}

type IpRange struct {
	Subnet string `xml:"Subnet,attr"`
}

type Azure struct {
	r *Record
}

func NewAzure() *Azure {
	return &Azure{
		r: loadDefaults(),
	}
}

func (a *Azure) Name() string {
	return "Azure"

}

func (a *Azure) UpdateRanges() error {
	link, err := downloadLink()
	if err != nil {
		return err
	}
	res, err := http.Get(link)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode == http.StatusOK {
		rec, err := load(res.Body)
		if err != nil {
			return err
		}
		a.r = rec
		return nil
	}
	return fmt.Errorf("Unexpected HTTP Status: %d", res.StatusCode)
}

func downloadLink() (link string, err error) {
	var res *http.Response
	res, err = http.Get("https://www.microsoft.com/en-us/download/confirmation.aspx?id=41653")
	if err != nil {
		return
	}
	defer res.Body.Close()
	if res.StatusCode == http.StatusOK {
		var text []byte
		text, err = ioutil.ReadAll(res.Body)
		if err != nil {
			return
		}
		return findLinkInText(text)
	}
	return "", fmt.Errorf("Unexpected HTTP Status: %d", res.StatusCode)
}

func findLinkInText(text []byte) (link string, err error) {
	re := regexp.MustCompile(`.*(https://download\.microsoft\.com/.*PublicIPs.*\.xml).*`)
	matches := re.FindStringSubmatch(string(text))
	if len(matches) < 1 {
		return "", fmt.Errorf("Download link not found")
	}
	return matches[1], nil
}

func (a *Azure) IsAt(addr string) bool {
	ip := net.ParseIP(addr)
	for _, region := range a.r.Region {
		for i := range region.IpRange {
			if a.isAtIP(ip, region.IpRange[i].Subnet) {
				return true
			}
		}
	}
	return false
}

func (a *Azure) isAtIP(addr net.IP, subnet string) bool {
	ip, ipnet, err := net.ParseCIDR(subnet)
	if err != nil {
		panic(err)
	}
	if ip != nil && ip.Equal(addr) {
		return true
	}
	if ipnet != nil && ipnet.Contains(addr) {
		return true
	}
	return false
}

func loadDefaults() *Record {
	r := strings.NewReader(defaultIpRanges)
	rec, err := load(r)
	if err != nil {
		panic(err)
	}
	return rec
}

func load(r io.Reader) (*Record, error) {
	var rec Record
	err := xml.NewDecoder(r).Decode(&rec)
	if err != nil {
		return nil, err
	}
	return &rec, nil
}
