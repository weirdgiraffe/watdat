//
// isat.go
// Copyright (C) 2017 weirdgiraffe <giraffe@cyberzoo.xyz>
//
// Distributed under terms of the MIT license.
//

package isat

type Provider interface {
	Name() string
	UpdateRanges() error
	IsAt(addr string) bool
}
