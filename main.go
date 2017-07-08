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

	"github.com/weirdgiraffe/isitat/aws"
)

func main() {
	addr, err := net.LookupHost(os.Args[1])
	if err != nil {
		panic(err)
	}
	a := aws.NewAWS()
	if err != nil {
		panic(err)
	}
	for i := range addr {
		if a.IsAt(addr[i]) {
			fmt.Printf("%-15s IS AT AWS\n", addr[i])
		} else {
			fmt.Printf("%-15s NOT AT AWS\n", addr[i])
		}
	}
}
