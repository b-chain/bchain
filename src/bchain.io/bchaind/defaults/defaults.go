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
// @File: defaults.go
// @Date: 2018/05/08 16:38:08
////////////////////////////////////////////////////////////////////////////////

package defaults

import (
	"os"
	"path/filepath"
	"strings"
	"fmt"
)

var (
	AppName 					= getAppName()

	DefaultDataDir        		= getAppCurrentDir()
	DefaultTOMLConfigPath 		= getAppCurrentDir() + "/"  + "config"
	DefaultLogPath        		= getAppCurrentDir() + "/log/" + AppName + ".log"
	DefaultLogLevel       		= "info"
	DefaultNodeName       		= AppName
	DefaultKeystore       		= getAppCurrentDir() + "/" + "keystore"
	DefaultNodePort       		= 36180

	//Rpc
	DefaultHttpModules    		= "bchain,personal,txpool"
	DefaultHttpHost       		= "localhost"
	DefaultHttpPort       		= 8989
	//Miner
	DefaultBlockproducerStart	= false
	//Net
	DefaultWorkingNet			="alpha"
)

func getAppName() string {
	name := os.Args[0]
	if strings.HasSuffix(name, ".exe") {
		name = strings.TrimSuffix(name, ".exe")
		if name == "" {
			panic("empty executable name")
		}
	}
	name = strings.Replace(name, "\\", "/", -1)
    v := strings.SplitAfterN(name, "/", -1)
	name = v[len(v) - 1]
	return name
}

func getAppCurrentDir() string {
	// discard error !!
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return ""
	}
	//fmt.Println("CurrentDir:",dir)

	return strings.Replace(dir, "\\", "/", -1)
}

func PrintAllDefalts(){
	fmt.Printf("DefaultDataDir:%s\n" , DefaultDataDir)
	fmt.Printf("DefaultKeystore:%s\n", DefaultKeystore)
}