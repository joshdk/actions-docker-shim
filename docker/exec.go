// Copyright Josh Komoroske. All rights reserved.
// Use of this source code is governed by the MIT license,
// a copy of which can be found in the LICENSE.txt file.
// SPDX-License-Identifier: MIT

// Package docker exposes functions for executing specific docker commands.
package docker

import (
	"os"
	"os/exec"
	"strings"
)

// Login executes a docker login to ghcr.io with the given username and
// password.
func Login(username, password string) error {
	cmd := exec.Command("/usr/bin/docker", "login", "ghcr.io", "--username", username, "--password-stdin")
	cmd.Stdin = strings.NewReader(password + "\n")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// Pull executes a docker pull of the given image.
func Pull(image string) error {
	cmd := exec.Command("/usr/bin/docker", "pull", image)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
