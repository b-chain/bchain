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
// @File: bchaind.go
// @Date: 2018/05/08 17:29:08
////////////////////////////////////////////////////////////////////////////////

package main

import (
	"fmt"
	"gopkg.in/urfave/cli.v1"
	"os"
	"runtime"
	"sort"

	"bchain.io/log"
	"bchain.io/bchaind/config"
	"bchain.io/bchaind/defaults"
	"bchain.io/bchaind/limits"
	"bchain.io/bchaind/utils"
	"bchain.io/node"
	"net/http"
	_ "net/http/pprof"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	clientIdentifier = "bchain" // Client identifier to advertise over the network
)

var (
	// Git SHA1 commit hash of the release (set via linker flags)
	gitCommit = ""
	// The app that holds all commands and flags.
	app = utils.NewApp(gitCommit, "the "+defaults.AppName+" command line interface")
	// basic flags
	basicFlags = []cli.Flag{
		utils.ConfigFileFlag,
		utils.LogFileFlag,
		utils.LogLevelFlag,
		utils.StartBlockproducerFlag,
		utils.MetricsEnabledFlag,
		utils.WorkingNetFlag,
		utils.ResyncBlockFlag,
	}

	logTag = "bchaind.main"
	logger log.Logger
)

func init() {
	// get a logger
	logger = log.GetLogger(logTag)
	if logger == nil {
		fmt.Errorf("Can not get logger(%s)\n", logTag)
		os.Exit(1)
	}

	// Initialize
	app.Action = bchaind
	app.HideVersion = true // we have a command to print the version
	app.Copyright = "Copyright 2018 The " + defaults.AppName + " Authors"

	// add commands
	app.Commands = []cli.Command{
		versionCommand,
	}
	sort.Sort(cli.CommandsByName(app.Commands))

	// add flags
	app.Flags = append(app.Flags, basicFlags...)
	// set before action
	app.Before = func(ctx *cli.Context) error {
		// TODO:

		// Use all processor cores.
		runtime.GOMAXPROCS(runtime.NumCPU())

		// Up some limits.
		if err := limits.SetLimits(); err != nil {
			fmt.Fprintf(os.Stderr, "failed to set limits: %v\n", err)
			return err
		}

		return nil
	}

	// set after action
	app.After = func(ctx *cli.Context) error {
		// TODO:
		return nil
	}
}

// remove block database based boot parameter --resync-block
func resyncBlockProc(ctx *cli.Context) error {
	resyncBlock := ctx.GlobalBool(utils.ResyncBlockFlag.Name)
	if resyncBlock {
		c := config.GetConfigInstance()
		c.SetPath(ctx.GlobalString(utils.ConfigFileFlag.Name))
		conf := &node.Config{}
		c.Register("node", conf)
		defer c.UnRegister("node")
		//need remove block data
		dataDir := conf.DataDir
		appName := conf.NameValue()
		path := filepath.Join(dataDir, appName, "chaindata")
		logger.Info("Now Remove the block database ", path)
		err := os.RemoveAll(path)
		if err != nil {
			logger.Error("remove block database error  ", err)
			return err
		}
		logger.Info("Remove block database completed !")
	}
	return nil
}

// bchaind is the real main entry pointclear
func bchaind(ctx *cli.Context) error {
	// log instance init
	err := log.InitInstance(ctx.GlobalString(utils.LogFileFlag.Name), ctx.GlobalString(utils.LogLevelFlag.Name))
	if err != nil {
		os.Exit(1)
	}
	defer log.CloseInstance()
	logger.Infof("")
	logger.Infof("===============================")
	logger.Infof("Hi, %s is starting ...", defaults.AppName)
	logger.Infof("===============================")

	if err := resyncBlockProc(ctx); err != nil {
		return err
	}

	node := createNode(ctx)
	if node == nil {
		logger.Critical("Create node failed.")
		os.Exit(1)
	}

	startNode(node)
	node.Wait()
	logger.Infof("%s is shutdown.", defaults.AppName)
	return nil
}

func main() {

	{

		//var err error
		for _, arg := range os.Args {

			if strings.Contains(arg, "pprof") {
				pprofArgs := strings.Split(arg, "-")
				pprofPort, err := strconv.Atoi(pprofArgs[1])

				if err == nil {
					go func() {
						addr := fmt.Sprintf("localhost:%d", pprofPort)
						fmt.Println("pprof addr:", addr)
						fmt.Println(http.ListenAndServe(addr, nil))
					}()
				} else {
					fmt.Println("[Error] Pprof params err:", err.Error())
				}
				break
			}

		}

	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
