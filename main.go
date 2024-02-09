// Copyright Josh Komoroske. All rights reserved.
// Use of this source code is governed by the MIT license,
// a copy of which can be found in the LICENSE.txt file.
// SPDX-License-Identifier: MIT

// Package main defines the entrypoint for the actions-docker-shim tool.
package main

import (
	"fmt"
	"os"

	"github.com/joshdk/actions-docker-shim/docker"
)

func main() {
	if err := mainCmd(); err != nil {
		fmt.Fprintln(os.Stderr, "actions-docker-shim:", err)
		os.Exit(1)
	}
}

//nolint:forbidigo,wsl
func mainCmd() error {
	image := fmt.Sprintf("ghcr.io/%s:%s", os.Getenv("GITHUB_REPOSITORY"), os.Getenv("GITHUB_ACTION_REF"))

	var token string
	if value := os.Getenv("GITHUB_TOKEN"); value != "" {
		// Environment variable named "GITHUB_TOKEN".
		token = value
	} else if value := os.Getenv("INPUT_GITHUB_TOKEN"); value != "" {
		// Input named "github-token".
		token = value
	} else if value := os.Getenv("INPUT_TOKEN"); value != "" {
		// Input named "token".
		token = value
	}

	fmt.Printf("::group::%s\n", "Docker login")
	err := docker.Login(os.Getenv("GITHUB_ACTOR"), token)
	fmt.Println("::endgroup::")
	if err != nil {
		return err
	}

	fmt.Printf("::group::%s\n", "Docker pull")
	err = docker.Pull(image)
	fmt.Println("::endgroup::")
	if err != nil {
		return err
	}

	return nil
}
