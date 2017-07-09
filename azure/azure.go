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
	"net/http"
	"regexp"
	"strings"

	"github.com/weirdgiraffe/watdatcloud"
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
	loaded bool
}

func NewAzure() *Azure {
	return &Azure{}
}

func (a *Azure) LoadRanges() ([]watdatcloud.Range, error) {
	if !a.loaded {
		a.loaded = true
		return decode(strings.NewReader(defaultIpRanges))
	}
	link, err := downloadLink()
	if err != nil {
		return nil, err
	}
	res, err := http.Get(link)
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
	err := xml.NewDecoder(r).Decode(&rec)
	if err != nil {
		return nil, err
	}
	ret := []watdatcloud.Range{}
	for _, reg := range rec.Region {
		for i := range reg.IpRange {
			ret = append(ret, watdatcloud.Range{
				Provider: "Azure",
				Info:     reg.Name,
				CIDR:     watdatcloud.AddrToCIDR(reg.IpRange[i].Subnet),
			})
		}
	}
	return ret, nil
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
