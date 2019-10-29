////////////////////////////////////////////////////////////////////////////////
// Copyright (c) 2018 The bchain-go Authors.
//
// The bchain-go is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.
//
// @File: cmd.go
// @Date: 2018/05/08 15:18:08
////////////////////////////////////////////////////////////////////////////////

package main

import (
	"fmt"
	"bchain.io/bchaind/utils"
	"gopkg.in/urfave/cli.v1"
	"os"
	"runtime"
)

// TODO: To be complete

func version(ctx *cli.Context) error {
	fmt.Println("Version:", utils.Version)
	if gitCommit != "" {
		fmt.Println("Git Commit:", gitCommit)
	}
	fmt.Println("Architecture:", runtime.GOARCH)
	//fmt.Println("Protocol Versions:", bchain.ProtocolVersions)
	//fmt.Println("Network Id:", bchain.DefaultConfig.NetworkId)
	fmt.Println("Go Version:", runtime.Version())
	fmt.Println("Operating System:", runtime.GOOS)
	fmt.Printf("GOPATH=%s\n", os.Getenv("GOPATH"))
	fmt.Printf("GOROOT=%s\n", runtime.GOROOT())
	return nil
}

var (
	versionCommand = cli.Command{
		Action:      version,
		Name:        "version",
		Usage:       "Print version numbers",
		ArgsUsage:   " ",
		Category:    "MISCELLANEOUS COMMANDS",
		Description: `The output of this command is supposed to be machine-readable.`,
	}
)
