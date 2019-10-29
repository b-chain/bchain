package sdk

import (
	"bchain.io/log"
	"fmt"
	"os"
)

var (
	logTag = "core.sdk"
	logger log.Logger
)



func init() {
	logger = log.GetLogger(logTag)
	if logger == nil {
		fmt.Errorf("Can not get logger(%s)\n", logTag)
		os.Exit(1)
	}
}
