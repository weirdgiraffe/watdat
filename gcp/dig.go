//
// dig.go
// Copyright (C) 2017 weirdgiraffe <giraffe@cyberzoo.xyz>
//
// Distributed under terms of the MIT license.
//

package gcp

import "os/exec"

type dig interface {
	GetTXT(fqdn string) (string, error)
}

type execDig struct{}

func (d *execDig) GetTXT(fqdn string) (string, error) {
	out, err := exec.Command("dig", "@8.8.8.8", "txt", fqdn, "+short").Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}
