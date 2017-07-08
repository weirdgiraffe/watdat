//
// main.go
// Copyright (C) 2017 weirdgiraffe <giraffe@cyberzoo.xyz>
//
// Distributed under terms of the MIT license.
//

package main

import (
	"fmt"
	"net"
	"os"

	"github.com/weirdgiraffe/isat"
	"github.com/weirdgiraffe/isat/aws"
	"github.com/weirdgiraffe/isat/azure"
)

func main() {
	addr, err := net.LookupHost(os.Args[1])
	if err != nil {
		panic(err)
	}
	providers := []isat.Provider{
		aws.NewAWS(),
		azure.NewAzure(),
	}
	for _, p := range providers {
		for i := range addr {
			if p.IsAt(addr[i]) {
				fmt.Printf("%-15s IS AT %s\n", addr[i], p.Name())
			}
		}
	}
}
