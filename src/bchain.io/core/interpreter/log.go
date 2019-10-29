package interpreter

import (
	"bchain.io/log"
	"fmt"
	"os"
)

var (
	logTag = "core.interpreter"
	logger log.Logger
)



func init() {
	logger = log.GetLogger(logTag)
	if logger == nil {
		fmt.Errorf("Can not get logger(%s)\n", logTag)
		os.Exit(1)
	}
	logger.SetLevel(log.LevelDebug)
}