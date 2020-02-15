package log

import (
	"fmt"
	"os"
	"testing"
)

func TestLog(t *testing.T) {
	fmt.Println(os.Getenv("MICRO_LOG_LEVEL"))
	Infof("Transport [%d] Listening", 1)
}
