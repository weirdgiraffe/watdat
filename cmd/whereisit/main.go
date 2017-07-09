//
// main.go
// Copyright (C) 2017 weirdgiraffe <giraffe@cyberzoo.xyz>
//
// Distributed under terms of the MIT license.
//

package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"

	"github.com/weirdgiraffe/watdatcloud"
	"github.com/weirdgiraffe/watdatcloud/aws"
	"github.com/weirdgiraffe/watdatcloud/azure"
	"github.com/weirdgiraffe/watdatcloud/gcp"
)

func main() {
	addr, err := net.LookupHost(os.Args[1])
	if err != nil {
		panic(err)
	}
	l := watdatcloud.NewRangeLookuper(
		aws.NewAWS(),
		azure.NewAzure(),
		gcp.NewGCP(),
	)
	err = l.UpdateRanges()
	if err != nil {
		panic(err)
	}
	for i := range addr {
		res, err := l.Lookup(addr[i])
		if err != nil {
			if err.(*watdatcloud.RangeNotFound) != nil {
				continue
			}
			panic(err)
		}
		b, err := json.MarshalIndent(res, "", "  ")
		if err != nil {
			panic(err)
		}
		fmt.Println(string(b))
	}
}
