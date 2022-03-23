package structsconv

import (
	"bytes"
	"log"
	"os"
	"reflect"
	"strings"
	"testing"
)

func Test_checkRootValuesTypes(t *testing.T) {
	type args struct {
		source reflect.Value
		target reflect.Value
	}
	strValue := "strValue"
	tests := []struct {
		name string
		args args
		want interface{}
	}{
		{
			name: "checkRootValuesTypes,(str,str)",
			args: args{
				source: reflect.ValueOf(strValue),
				target: reflect.ValueOf(strValue),
			},
			want: "MapRules error. Source must be a pointer",
		},
		{
			name: "checkRootValuesTypes,(*str,str)",
			args: args{
				source: reflect.ValueOf(&strValue),
				target: reflect.ValueOf(strValue),
			},
			want: "MapRules error. Target must be a pointer",
		},
		{
			name: "checkRootValuesTypes,(*str,*str)",
			args: args{
				source: reflect.ValueOf(&strValue),
				target: reflect.ValueOf(&strValue),
			},
			want: "MapRules error. Source must be a pointer to a struct",
		},
		{
			name: "checkRootValuesTypes,(*struct,*str)",
			args: args{
				source: reflect.ValueOf(&struct{}{}),
				target: reflect.ValueOf(&strValue),
			},
			want: "MapRules error. Target must be a pointer to a struct",
		},
		{
			name: "checkRootValuesTypes,(*struct,*struct)",
			args: args{
				source: reflect.ValueOf(&struct{}{}),
				target: reflect.ValueOf(&struct{}{}),
			},
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := checkRootValuesTypes(tt.args.source, tt.args.target); err != nil && err.Error() != tt.want {
				t.Errorf("logs expected to contain '%s', got '%s'", tt.want, err.Error())
			}
		})
	}
}

func Test_checkMapperRules_panics(t *testing.T) {
	key := buildKey(TestSource{}, TestTarget{})
	tests := []struct {
		name         string
		rules        RulesSet
		wantContains string
	}{
		{
			name:         "Non-existent target key name,panic expected",
			wantContains: "Field 'otherFieldName' is not present in target struct structsconv.TestTarget",
			rules: RulesSet{
				"otherFieldName": "fieldS1",
			},
		},
		{
			name:         "Non-existent source field name,panic expected",
			wantContains: "Field 'otherFieldName' is not present in source struct structsconv.TestSource",
			rules: RulesSet{
				"fieldT1": "otherFieldName",
			},
		},
		{
			name:         "Field kind is the different in origin and target struct,panic expected",
			wantContains: "Field 'fieldT1' has different type in source (fieldS2:int) and target (fieldT1:string) structs",
			rules: RulesSet{
				"fieldT1": "fieldS2",
			},
		},
		{
			name:         "Custom function returns different type than target,panic expected",
			wantContains: "Function 'fieldT1' must return type 'string', currently returns 'int'. Function = 'func() int'",
			rules: RulesSet{
				"fieldT1": func() int { return 314 },
			},
		},
		{
			name:         "Not valid rule, int value,panic expected",
			wantContains: "Rule 'fieldT1' is not valid",
			rules: RulesSet{
				"fieldT1": 123,
			},
		},
		{
			name:         "Not valid rule, int value,panic expected",
			wantContains: "Rule 'fieldT1' is not valid",
			rules: RulesSet{
				"fieldT1": struct{}{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var f = func() { checkMapperRules(key, tt.rules) }
			assertPanic(f, tt.wantContains, t)
		})
	}
}

func Test_checkMapperRules_ignorable(t *testing.T) {
	key := buildKey(TestSource{}, TestTarget{})
	tests := []struct {
		name         string
		rules        RulesSet
		wantContains string
	}{
		{
			name:         "Ignorable rule, no panic expected",
			wantContains: "Field 'ignorableTargetField' is marked as ignored",
			rules: RulesSet{
				"ignorableTargetField": nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var f = func() { checkMapperRules(key, tt.rules) }
			assertLog(f, tt.wantContains, t)
		})
	}
}

func assertLog(f func(), wantContains string, t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(os.Stderr)
	f()
	if !strings.Contains(buf.String(), wantContains) {
		t.Errorf("logs expected to contain '%s', got '%s'", wantContains, buf.String())
	}
}
