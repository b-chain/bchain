package aposgensis

import (
	"bchain.io/log"
	"fmt"
	"os"
)

var (
	logTag = "consensus.apos.aposgensis"
	logger log.Logger
)

func init() {
	logger = log.GetLogger(logTag)
	logger.SetLevel(log.LevelDebug)
	if logger == nil {
		fmt.Errorf("Can not get logger(%s)\n", logTag)
		os.Exit(1)
	}
}
