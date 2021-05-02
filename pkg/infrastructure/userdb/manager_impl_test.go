package userdb

import (
	"os"
	"strconv"
	"testing"
)

func skipOrDo(t *testing.T) {
	if ok, _ := strconv.ParseBool(os.Getenv("ENABLE_DB_TESTS")); !ok {
		t.SkipNow()
	}
}
