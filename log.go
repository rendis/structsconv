package structsconv

import (
	"log"
	"reflect"
	"sync"
)

var (
	logWarning = true
	setLogOnce sync.Once
)

// SetLogWarning turns on (true) or off (false) the warning log messages. Default is on (true).
//
// If you want to turn on the warning log messages, you can call this function before you call any other functions.
func SetLogWarning(b bool) {
	setLogOnce.Do(func() {
		logWarning = b
	})
}

func logTargetFieldWithoutMappingValueInSource(key rulesKey, targetFieldName string) {
	if logWarning {
		log.Printf("WARNING: (%s -> %s) No mapping found for name '%s'.\n",
			key.source, key.target, targetFieldName,
		)
	}
}

func logIgnoringMappingForIncompatibleTypes(key rulesKey, targetFieldName string, sourceValue, targetValue reflect.Value) {
	if logWarning {
		log.Printf(
			"WARNING: (%s -> %s) Ignoring mapping for name '%s' (%s) to (%s), cause: Incompatible types.\n",
			key.source, key.target, targetFieldName, sourceValue.Type(), targetValue.Type(),
		)
	}
}

func logPassingZeroValue(method, argType reflect.Type, argPosition int) {
	if logWarning {
		log.Printf(
			"WARNING: Passing 'ZeroValue' in custom function (%s) for argument of type '%s' in position %d.\n",
			method, argType, argPosition,
		)
	}
}
