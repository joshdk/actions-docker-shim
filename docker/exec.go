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
	"syscall"
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

// Run executes a docker run replicating the GitHub Actions Docker runtime. If
// execution is successful, the current process is replaced and this function
// will not return.
func Run(image string, cmdargs []string) error {
	args := []string{
		"docker", "run",
		"--rm",
		"--workdir", "/github/workspace",
		"-v", "/var/run/docker.sock:/var/run/docker.sock",
		"-v", "/home/runner/work/_temp/_github_home:/github/home",
		"-v", "/home/runner/work/_temp/_github_workflow:/github/workflow",
		"-v", "/home/runner/work/_temp/_runner_file_commands:/github/file_commands",
		"-v", "/home/runner/work/go-action/go-action:/github/workspace",
	}

	environ := os.Environ()
	envs := environ[:0]

	// Loop over the current set of environment variables and either pass them
	// through to the docker run command arguments, or drop them (certain
	// environment variables are not explicitly passed).
	for _, env := range environ {
		key := strings.SplitN(env, "=", 2)[0] //nolint:gomnd
		switch key {
		case "HOSTNAME", "PATH":
			// Drop these environment variables.
			continue
		default:
			// Keep these environment variables and add them to the docker run
			// command line.
			envs = append(envs, env)
			args = append(args, "-e", key)
		}
	}

	args = append(args, image)
	args = append(args, cmdargs...)

	return syscall.Exec("/usr/bin/docker", args, envs) //nolint:gosec
}
