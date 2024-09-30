// SPDX-License-Identifier: Apache-2.0
// Copyright 2024 Authors of SentryFlow

package cmd

import (
	"runtime"
	"runtime/debug"
)

func logBuildInfo() {
	info, _ := debug.ReadBuildInfo()
	vcsRev := ""
	vcsTime := ""
	for _, s := range info.Settings {
		if s.Key == "vcs.revision" {
			vcsRev = s.Value
		} else if s.Key == "vcs.time" {
			vcsTime = s.Value
		}
	}
	logger.Infof("Git revision: %s, build time: %s, build version: %s, go os/arch: %s/%s", vcsRev, vcsTime, info.Main.Version, runtime.GOOS, runtime.GOARCH)
}
