//
// azure_test.go
// Copyright (C) 2017 weirdgiraffe <giraffe@cyberzoo.xyz>
//
// Distributed under terms of the MIT license.
//

package azure

import "testing"

func TestUpdate(t *testing.T) {
	a := NewAzure()
	err := a.UpdateRanges()
	if err != nil {
		t.Fatal(err)
	}
}
