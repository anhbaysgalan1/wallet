package job

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	ret := m.Run()
	os.Exit(ret)
}
