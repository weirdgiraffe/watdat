//
// aws_test.go
// Copyright (C) 2017 weirdgiraffe <giraffe@cyberzoo.xyz>
//
// Distributed under terms of the MIT license.
//

package aws

import "testing"

func TestUpdate(t *testing.T) {
	a := NewAWS()
	_, err := a.LoadRanges()
	if err != nil {
		t.Fatal(err)
	}
}
