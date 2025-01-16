// SPDX-License-Identifier: Apache-2.0
// Copyright 2021 Authors of KubeArmor

package version

import (
	"fmt"
	"runtime"
	"runtime/debug"

	"github.com/spf13/cobra"
)

var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display version information",
	Long:  `Display version information`,
	Run: func(cmd *cobra.Command, args []string) {
		printVersion()
	},
}

func printVersion() {
	info, _ := debug.ReadBuildInfo()
	vcsTime := ""
	for _, s := range info.Settings {
		if s.Key == "vcs.time" {
			vcsTime = s.Value
		}
	}
	fmt.Printf("sfctl version: %s %s/%s buildtime: %s\n",
		info.Main.Version,
		runtime.GOOS,
		runtime.GOARCH,
		vcsTime,
	)
}
