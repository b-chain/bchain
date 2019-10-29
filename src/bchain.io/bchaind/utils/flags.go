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
// @File: flags.go
// @Date: 2018/05/08 16:38:08
////////////////////////////////////////////////////////////////////////////////

package utils

import (
	"bchain.io/bchaind/defaults"
	"gopkg.in/urfave/cli.v1"
	"os"
	"path/filepath"
	"bchain.io/utils/metrics"
)

var (
	Version = "1.1.1"
)

func init() {
	cli.AppHelpTemplate = `{{.HelpName}} {{if .VisibleFlags}}[global options]{{end}} {{if .Commands}}command [command options]{{end}} {{if .ArgsUsage}}{{.ArgsUsage}}{{else}}[arguments...]{{end}}

   {{if .Copyright}}{{.Copyright}}{{end}}

VERSION:
   {{.Version}}{{if .Commands}}

COMMANDS:
   {{range .Commands}}{{join .Names ", "}}{{ "\t" }}{{.Usage}}
   {{end}}{{end}}{{if .VisibleFlags}}
GLOBAL OPTIONS:
   {{range .Flags}}{{.}}
   {{end}}{{end}}
`

	cli.CommandHelpTemplate = `NAME:
   {{.HelpName}} - {{.Usage}}

USAGE:
   {{if .UsageText}}{{.UsageText}}{{else}}{{.HelpName}}{{if .VisibleFlags}} [command options]{{end}} {{if .ArgsUsage}}{{.ArgsUsage}}{{else}}[arguments...]{{end}}{{end}}{{if .Category}}

DESCRIPTION:
   {{.Description}}{{end}}{{if .VisibleFlags}}

OPTIONS:
   {{range .VisibleFlags}}{{.}}
   {{end}}{{end}}
`
	//cli.SubcommandHelpTemplate	// use default
}

// NewApp creates an app with sane defaults.
func NewApp(gitCommit, usage string) *cli.App {
	app := cli.NewApp()
	app.Name = filepath.Base(os.Args[0])
	app.Author = ""
	//app.Authors = nil
	app.Email = ""
	app.Version = Version
	if len(gitCommit) >= 8 {
		app.Version += "-" + gitCommit[:8]
	}
	app.Usage = usage
	return app
}

// These are all the command line flags we support.
// If you add to this list, please remember to include the
// flag in the appropriate command definition.
//
// The flags are defined here so their names and help texts
// are the same for all commands.
var (
	// General settings
	DataDirFlag = DirectoryFlag{
		Name:  "datadir",
		Usage: "Data directory for the databases and keystore",
		Value: DirectoryString{defaults.DefaultDataDir},
	}

	KeysotreFlag = cli.StringFlag{
		Name:  "keystore",
		Usage: "the keystore file directory",
		Value: defaults.DefaultKeystore,
	}

	ConfigFileFlag = cli.StringFlag{
		Name:  "config",
		Usage: "TOML configuration file",
		Value: defaults.DefaultTOMLConfigPath,
	}

	LogFileFlag = cli.StringFlag{
		Name:  "log",
		Usage: "The path of log file",
		Value: defaults.DefaultLogPath,
	}

	LogLevelFlag = cli.StringFlag{
		Name:  "loglevel",
		Usage: "The level of log [trace | debug | info | warn | error | critical | off]",
		Value: defaults.DefaultLogLevel,
	}

	NodeNameFlag = cli.StringFlag{
		Name:  "nodename",
		Usage: "The name of local node",
		Value: defaults.DefaultNodeName,
	}

	ListenPortFlag = cli.IntFlag{
		Name:  "port",
		Usage: "Network listening port",
		Value: defaults.DefaultNodePort,
	}

	BootNodeUrlFlag = cli.StringFlag{
		Name:  "bootnode",
		Usage: "The url of bootstrap node (bnode://id@ip:port)",
	}

	StaticNodeUrlFlag = cli.StringFlag{
		Name:  "staticnode",
		Usage: "The url of static node (bnode://id@ip:port)",
	}
	//RPC
	HttpModulesFlag = cli.StringFlag{
		Name: "httpmodules",
		Usage: "A list of API modules to expose via the HTTP RPC interface [bchain,personal,txpool,blockproducer]",
		Value: defaults.DefaultHttpModules,
	}

	HttpPortFlag = cli.IntFlag{
		Name:  "httpport",
		Usage: "The port Rpc listen",
		Value: defaults.DefaultHttpPort,
	}

	HttpHostFlag = cli.StringFlag{
		Name:  "httphost",
		Usage: "The host Rpc listen",
		Value: defaults.DefaultHttpHost,
	}

	StartBlockproducerFlag = cli.BoolFlag{
		Name:  "blockproducer",
		Usage: "The trigger of blockproducer",
	}

	MetricsEnabledFlag = cli.BoolFlag{
		Name:  metrics.MetricsEnabledFlag,
		Usage: "Enable metrics collection and reporting",
	}

	WorkingNetFlag = cli.StringFlag{
		Name:	"net",
		Usage:	"A working net(main/beta/alpha) to run. This option is not used in current version",
		Value:	defaults.DefaultWorkingNet,
	}

	ResyncBlockFlag = cli.BoolFlag{
		Name: 	"resync-block",
		Usage:	"Clear chain database and rebuild blockchain.",
	}
)
