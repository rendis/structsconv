package structsconv

import (
	"reflect"
	"testing"
)

type AStruct struct{}
type BStruct struct{}

func Test_buildKey(t *testing.T) {
	type args struct {
		source interface{}
		target interface{}
	}
	tests := []struct {
		name string
		args args
		want rulesKey
	}{
		{
			name: "buildKey",
			args: args{
				source: AStruct{},
				target: BStruct{},
			},
			want: rulesKey{
				source: reflect.TypeOf(AStruct{}),
				target: reflect.TypeOf(BStruct{}),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := buildKey(tt.args.source, tt.args.target); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("buildKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getArraysMappingType(t *testing.T) {
	type args struct {
		sourceValue reflect.Value
		targetValue reflect.Value
	}
	tests := []struct {
		name string
		args args
		want processingResultType
	}{
		{
			name: "getArraysMappingType",
			args: args{
				sourceValue: reflect.ValueOf([]AStruct{}),
				targetValue: reflect.ValueOf([]BStruct{}),
			},
			want: arraysMapping,
		},
		{
			name: "getArraysMappingType",
			args: args{
				sourceValue: reflect.ValueOf([]AStruct{}),
				targetValue: reflect.ValueOf([]int{}),
			},
			want: incompatibleTypes,
		},
		{
			name: "getArraysMappingType",
			args: args{
				sourceValue: reflect.ValueOf([]string{}),
				targetValue: reflect.ValueOf([]BStruct{}),
			},
			want: incompatibleTypes,
		},
		{
			name: "getArraysMappingType",
			args: args{
				sourceValue: reflect.ValueOf([]string{}),
				targetValue: reflect.ValueOf([]int{}),
			},
			want: incompatibleTypes,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getArraysMappingType(tt.args.sourceValue, tt.args.targetValue); got != tt.want {
				t.Errorf("getArraysMappingType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getSlicesMappingType(t *testing.T) {
	type args struct {
		sourceValue reflect.Value
		targetValue reflect.Value
	}
	tests := []struct {
		name string
		args args
		want processingResultType
	}{
		{
			name: "getSlicesMappingType",
			args: args{
				sourceValue: reflect.ValueOf([]AStruct{}),
				targetValue: reflect.ValueOf([]BStruct{}),
			},
			want: slicesMapping,
		},
		{
			name: "getSlicesMappingType",
			args: args{
				sourceValue: reflect.ValueOf([]AStruct{}),
				targetValue: reflect.ValueOf([]int{}),
			},
			want: incompatibleTypes,
		},
		{
			name: "getSlicesMappingType",
			args: args{
				sourceValue: reflect.ValueOf([]string{}),
				targetValue: reflect.ValueOf([]BStruct{}),
			},
			want: incompatibleTypes,
		},
		{
			name: "getSlicesMappingType",
			args: args{
				sourceValue: reflect.ValueOf([]string{}),
				targetValue: reflect.ValueOf([]int{}),
			},
			want: incompatibleTypes,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getSlicesMappingType(tt.args.sourceValue, tt.args.targetValue); got != tt.want {
				t.Errorf("getSlicesMappingType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getMapsMappingType(t *testing.T) {
	type args struct {
		sourceValue reflect.Value
		targetValue reflect.Value
	}
	tests := []struct {
		name string
		args args
		want processingResultType
	}{
		{
			name: "getMapsMappingType",
			args: args{
				sourceValue: reflect.ValueOf(map[string]AStruct{}),
				targetValue: reflect.ValueOf(map[string]BStruct{}),
			},
			want: mapsMapping,
		},
		{
			name: "getMapsMappingType",
			args: args{
				sourceValue: reflect.ValueOf(map[string]AStruct{}),
				targetValue: reflect.ValueOf(map[int]BStruct{}),
			},
			want: incompatibleTypes,
		},
		{
			name: "getMapsMappingType",
			args: args{
				sourceValue: reflect.ValueOf(map[string]AStruct{}),
				targetValue: reflect.ValueOf(map[string]int{}),
			},
			want: incompatibleTypes,
		},
		{
			name: "getMapsMappingType",
			args: args{
				sourceValue: reflect.ValueOf(map[string]string{}),
				targetValue: reflect.ValueOf(map[string]BStruct{}),
			},
			want: incompatibleTypes,
		},
		{
			name: "getMapsMappingType",
			args: args{
				sourceValue: reflect.ValueOf(map[string]string{}),
				targetValue: reflect.ValueOf(map[string]int{}),
			},
			want: incompatibleTypes,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getMapsMappingType(tt.args.sourceValue, tt.args.targetValue); got != tt.want {
				t.Errorf("getMapsMappingType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getMappingType(t *testing.T) {
	type args struct {
		sourceValue reflect.Value
		targetValue reflect.Value
	}
	tests := []struct {
		name string
		args args
		want processingResultType
	}{
		{
			name: "getMappingType,(str,str)->directly",
			args: args{
				sourceValue: reflect.ValueOf(""),
				targetValue: reflect.ValueOf(""),
			},
			want: directMapping,
		},
		{
			name: "getMappingType,([]int,[]int)->directly",
			args: args{
				sourceValue: reflect.ValueOf([]int{}),
				targetValue: reflect.ValueOf([]int{}),
			},
			want: directMapping,
		},
		{
			name: "getMappingType,([3]AStruct{},[3]AStruct{})->directly",
			args: args{
				sourceValue: reflect.ValueOf([3]AStruct{}),
				targetValue: reflect.ValueOf([3]AStruct{}),
			},
			want: directMapping,
		},
		{
			name: "getMappingType,(AStruct{},AStruct{})->directly",
			args: args{
				sourceValue: reflect.ValueOf(AStruct{}),
				targetValue: reflect.ValueOf(AStruct{}),
			},
			want: directMapping,
		},
		{
			name: "getMappingType,(map[str]int,map[str]int)->directly",
			args: args{
				sourceValue: reflect.ValueOf(map[string]int{}),
				targetValue: reflect.ValueOf(map[string]int{}),
			},
			want: directMapping,
		},
		{
			name: "getMappingType,(AStruct{},BStruct{})->structsMapping",
			args: args{
				sourceValue: reflect.ValueOf(AStruct{}),
				targetValue: reflect.ValueOf(BStruct{}),
			},
			want: structsMapping,
		},
		{
			name: "getMappingType,([]AStruct{},[]BStruct{})->slicesMapping",
			args: args{
				sourceValue: reflect.ValueOf([]AStruct{}),
				targetValue: reflect.ValueOf([]BStruct{}),
			},
			want: slicesMapping,
		},
		{
			name: "getMappingType,([]AStruct{},[]int{})->incompatibleTypes",
			args: args{
				sourceValue: reflect.ValueOf([]AStruct{}),
				targetValue: reflect.ValueOf([]int{}),
			},
			want: incompatibleTypes,
		},
		{
			name: "getMappingType,([5]AStruct{},[5]BStruct{})->arraysMapping",
			args: args{
				sourceValue: reflect.ValueOf([5]AStruct{}),
				targetValue: reflect.ValueOf([5]BStruct{}),
			},
			want: arraysMapping,
		},
		{
			name: "getMappingType,([5]AStruct{},[5]int{})->incompatibleTypes",
			args: args{
				sourceValue: reflect.ValueOf([5]AStruct{}),
				targetValue: reflect.ValueOf([5]int{}),
			},
			want: incompatibleTypes,
		},
		{
			name: "getMappingType,(map[string]AStruct{},map[string]BStruct{})->mapsMapping",
			args: args{
				sourceValue: reflect.ValueOf(map[string]AStruct{}),
				targetValue: reflect.ValueOf(map[string]BStruct{}),
			},
			want: mapsMapping,
		},
		{
			name: "getMappingType,(map[string]AStruct{},map[string]string{})->incompatibleTypes",
			args: args{
				sourceValue: reflect.ValueOf(map[string]AStruct{}),
				targetValue: reflect.ValueOf(map[string]string{}),
			},
			want: incompatibleTypes,
		},
		{
			name: "getMappingType,(string,float)->incompatibleTypes",
			args: args{
				sourceValue: reflect.ValueOf("str_value"),
				targetValue: reflect.ValueOf(3.14),
			},
			want: incompatibleTypes,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getMappingType(tt.args.sourceValue, tt.args.targetValue); got != tt.want {
				t.Errorf("getMappingType() = %v, want %v", got, tt.want)
			}
		})
	}
}
