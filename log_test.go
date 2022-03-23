package structsconv

import (
	"bytes"
	"log"
	"os"
	"reflect"
	"strings"
	"testing"
)

func Test_Log_on(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(os.Stderr)

	SetLogWarning(true)

	rk := rulesKey{
		source: reflect.TypeOf(struct {
			A int
		}{}),
		target: reflect.TypeOf(struct {
			B int
		}{}),
	}
	wantContains := "(struct { A int } -> struct { B int }) No mapping found for name 'A'"
	logTargetFieldWithoutMappingValueInSource(rk, "A")
	if !strings.Contains(buf.String(), wantContains) {
		t.Errorf("log.Printf() = %q; want %q", buf.String(), wantContains)
	}
}

func Test_logTargetFieldWithoutMappingValueInSource_off(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(os.Stderr)

	logWarning = false
	defer func() { logWarning = true }()

	rk := rulesKey{
		source: reflect.TypeOf(struct {
			A int
		}{}),
		target: reflect.TypeOf(struct {
			B int
		}{}),
	}
	wantContains := ""
	logTargetFieldWithoutMappingValueInSource(rk, "A")
	if !strings.Contains(buf.String(), wantContains) {
		t.Errorf("log.Printf() = %q; want %q", buf.String(), wantContains)
	}
}

func Test_logIgnoringMappingForIncompatibleTypes_off(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(os.Stderr)

	logWarning = false
	defer func() { logWarning = true }()

	source := struct {
		A int
	}{}

	target := struct {
		B int
	}{}

	rk := rulesKey{
		source: reflect.TypeOf(source),
		target: reflect.TypeOf(target),
	}
	wantContains := ""
	logIgnoringMappingForIncompatibleTypes(rk, "A", reflect.ValueOf(source), reflect.ValueOf(target))
	if !strings.Contains(buf.String(), wantContains) {
		t.Errorf("log.Printf() = %q; want %q", buf.String(), wantContains)
	}
}

func Test_logPassingZeroValue_off(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(os.Stderr)

	logWarning = false
	defer func() { logWarning = true }()

	wantContains := ""

	ft := reflect.TypeOf(func(i int) {})
	at := reflect.TypeOf(1)

	logPassingZeroValue(ft, at, 0)
	if !strings.Contains(buf.String(), wantContains) {
		t.Errorf("log.Printf() = %q; want %q", buf.String(), wantContains)
	}
}
