package structsconv

import (
	"fmt"
	"strings"
	"testing"
)

type TestTarget struct {
	fieldT1              string
	fieldT2              string
	fieldT3              int
	fieldDiffType4       int
	fieldUnexportedEqual int
	FieldExportedEqual   int
	ignorableTargetField string
	singleListFieldT     []string
}

type TestSource struct {
	fieldS1              string
	fieldDiffType4       string
	fieldUnexportedEqual int
	FieldExportedEqual   int
	fieldS2              int
	singleListFieldS     []string
}

func assertPanic(f func(), contain string, t *testing.T) {
	defer func() {
		r := recover()
		if !strings.Contains(fmt.Sprintf("%v", r), contain) {
			t.Errorf("logs expected to contain '%s', got '%s'", contain, fmt.Sprintf("%v", r))
		}
	}()
	f()
}
