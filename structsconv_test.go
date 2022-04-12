package structsconv

import (
	"reflect"
	"testing"
)

func Test_RegisterRulesDefinitions_duplicate_error(t *testing.T) {
	type source struct {
		Name string
	}
	type target struct {
		Nick string
	}

	ruleDef1 := RulesDefinition{
		Source: source{Name: "name"},
		Target: target{Nick: "nick"},
		Rules:  RulesSet{},
	}
	ruleDef2 := RulesDefinition{
		Source: source{Name: "name"},
		Target: target{Nick: "nick"},
		Rules:  RulesSet{},
	}

	f := func() { RegisterRulesDefinitions(ruleDef1, ruleDef2) }
	assertPanic(f, "Mapper with rulesKey (structsconv.source -> structsconv.target) already exists", t)
}

func Test_Map_source_target_type_panics(t *testing.T) {
	type args struct {
		source interface{}
		target interface{}
	}
	var test = []struct {
		name         string
		args         args
		wantContains string
	}{
		{
			name: "Mapping with source!=pointer,panic expected",
			args: args{
				source: struct{}{},
				target: &struct{}{},
			},
			wantContains: "rules error: source must be a pointer",
		},
		{
			name: "Mapping with target!=pointer,panic expected",
			args: args{
				source: &struct{}{},
				target: struct{}{},
			},
			wantContains: "rules error: target must be a pointer",
		},
	}

	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			var f = func() { Map(tt.args.source, tt.args.target) }
			assertPanic(f, tt.wantContains, t)
		})
	}
}

// ptr struct -> ptr struct
func Test_Map_nested_ptr_struct_case1(t *testing.T) {
	type nestedSource struct{ Field string }
	type nestedTarget struct{ Field string }
	type source struct {
		Name   string
		Age    int
		Nested *nestedSource
	}
	type target struct {
		Name   string
		Age    int
		Nested *nestedTarget
	}

	o := &source{
		Name:   "name",
		Age:    10,
		Nested: &nestedSource{Field: "nested"},
	}

	d := &target{}

	want := &target{
		Name:   "name",
		Age:    10,
		Nested: &nestedTarget{Field: "nested"},
	}

	Map(o, d)

	if !reflect.DeepEqual(d, want) {
		t.Errorf("Map(%v, %v) = %v, want %v", o, d, d, want)
	}
}

// ptr struct -> struct
func Test_Map_nested_ptr_struct_case2(t *testing.T) {
	type nestedSource struct{ Field string }
	type nestedTarget struct{ Field string }
	type source struct {
		Name   string
		Age    int
		Nested *nestedSource
	}
	type target struct {
		Name   string
		Age    int
		Nested nestedTarget
	}

	o := &source{
		Name:   "name",
		Age:    10,
		Nested: &nestedSource{Field: "nested"},
	}

	d := &target{}

	want := &target{
		Name:   "name",
		Age:    10,
		Nested: nestedTarget{Field: "nested"},
	}

	Map(o, d)

	if !reflect.DeepEqual(d, want) {
		t.Errorf("Map(%v, %v) = %v, want %v", o, d, d, want)
	}
}

// struct -> ptr struct
func Test_Map_nested_ptr_struct_case3(t *testing.T) {
	type nestedSource struct{ Field string }
	type nestedTarget struct{ Field string }
	type source struct {
		Name   string
		Age    int
		Nested nestedSource
	}
	type target struct {
		Name   string
		Age    int
		Nested *nestedTarget
	}

	o := &source{
		Name:   "name",
		Age:    10,
		Nested: nestedSource{Field: "nested"},
	}

	d := &target{}

	want := &target{
		Name:   "name",
		Age:    10,
		Nested: &nestedTarget{Field: "nested"},
	}

	Map(o, d)

	if !reflect.DeepEqual(d, want) {
		t.Errorf("Map(%v, %v) = %v, want %v", o, d, d, want)
	}
}

// nil ptr -> ptr struct
func Test_Map_nested_ptr_struct_case4(t *testing.T) {
	type nestedSource struct{ Field string }
	type nestedTarget struct{ Field string }
	type source struct {
		Nested *nestedSource
	}
	type target struct {
		Nested *nestedTarget
	}

	o := &source{
		Nested: nil,
	}

	d := &target{}

	want := &target{
		Nested: nil,
	}

	Map(o, d)

	if !reflect.DeepEqual(d, want) {
		t.Errorf("Map(%v, %v) = %v, want %v", o, d, d, want)
	}
}

// struct -> ptr struct: with rule
func Test_Map_nested_ptr_struct_case5(t *testing.T) {
	type nestedSource struct{ Field string }
	type nestedTarget struct{ Field string }
	type source struct {
		Name   string
		Age    int
		Nested nestedSource
	}
	type target struct {
		Name    string
		Age     int
		Nested1 *nestedTarget
	}

	o := &source{
		Name:   "name",
		Age:    10,
		Nested: nestedSource{Field: "nested"},
	}

	d := &target{}

	want := &target{
		Name:    "name",
		Age:     10,
		Nested1: &nestedTarget{Field: "nested"},
	}

	var rules = RulesSet{"Nested1": "Nested"}
	RegisterRulesDefinitions(RulesDefinition{Source: source{}, Target: target{}, Rules: rules})

	Map(o, d)

	if !reflect.DeepEqual(d, want) {
		t.Errorf("Map(%v, %v) = %v, want %v", o, d, d, want)
	}
}

// ptr struct -> struct: with rule
func Test_Map_nested_ptr_struct_case6(t *testing.T) {
	type nestedSource struct{ Field string }
	type nestedTarget struct{ Field string }
	type source struct {
		Name   string
		Age    int
		Nested *nestedSource
	}
	type target struct {
		Name    string
		Age     int
		Nested1 nestedTarget
	}

	o := &source{
		Name:   "name",
		Age:    10,
		Nested: &nestedSource{Field: "nested"},
	}

	d := &target{}

	want := &target{
		Name:    "name",
		Age:     10,
		Nested1: nestedTarget{Field: "nested"},
	}

	var rules = RulesSet{"Nested1": "Nested"}
	RegisterRulesDefinitions(RulesDefinition{Source: source{}, Target: target{}, Rules: rules})

	Map(o, d)

	if !reflect.DeepEqual(d, want) {
		t.Errorf("Map(%v, %v) = %v, want %v", o, d, d, want)
	}
}

// ptr struct -> ptr struct: with rule
func Test_Map_nested_ptr_struct_case7(t *testing.T) {
	type nestedSource struct{ Field string }
	type nestedTarget struct{ Field string }
	type source struct {
		Name   string
		Age    int
		Nested *nestedSource
	}
	type target struct {
		Name    string
		Age     int
		Nested1 *nestedTarget
	}

	o := &source{
		Name:   "name",
		Age:    10,
		Nested: &nestedSource{Field: "nested"},
	}

	d := &target{}

	want := &target{
		Name:    "name",
		Age:     10,
		Nested1: &nestedTarget{Field: "nested"},
	}

	var rules = RulesSet{"Nested1": "Nested"}
	RegisterRulesDefinitions(RulesDefinition{Source: source{}, Target: target{}, Rules: rules})

	Map(o, d)

	if !reflect.DeepEqual(d, want) {
		t.Errorf("Map(%v, %v) = %v, want %v", o, d, d, want)
	}
}

// [] ptr struct -> [] ptr struct: with rule
func Test_Map_slice_ptr_struct_case1(t *testing.T) {
	type itemSource struct{ Field string }
	type source struct {
		Items []*itemSource
	}
	type itemTarget struct{ Field string }
	type target struct {
		Items []*itemTarget
	}

	o := &source{
		Items: []*itemSource{{Field: "item1"}, {Field: "item2"}},
	}

	d := &target{}

	want := &target{
		Items: []*itemTarget{{Field: "item1"}, {Field: "item2"}},
	}

	Map(o, d)

	if !reflect.DeepEqual(d, want) {
		t.Errorf("Map(%v, %v) = %#v, want %#v", o, d, d.Items, want.Items)
	}
}

// [] struct -> [] ptr struct: with rule
func Test_Map_slice_ptr_struct_case2(t *testing.T) {
	type itemSource struct{ Field string }
	type source struct {
		Items []itemSource
	}
	type itemTarget struct{ Field string }
	type target struct {
		Items []*itemTarget
	}

	o := &source{
		Items: []itemSource{{Field: "item1"}, {Field: "item2"}},
	}

	d := &target{}

	want := &target{
		Items: []*itemTarget{{Field: "item1"}, {Field: "item2"}},
	}

	Map(o, d)

	if !reflect.DeepEqual(d, want) {
		t.Errorf("Map(%v, %v) = %#v, want %#v", o, d, d.Items, want.Items)
	}
}

// [] ptr struct -> [] struct: with rule
func Test_Map_slice_ptr_struct_case3(t *testing.T) {
	type itemSource struct{ Field string }
	type source struct {
		Items []*itemSource
	}
	type itemTarget struct{ Field string }
	type target struct {
		Items []itemTarget
	}

	o := &source{
		Items: []*itemSource{{Field: "item1"}, {Field: "item2"}},
	}

	d := &target{}

	want := &target{
		Items: []itemTarget{{Field: "item1"}, {Field: "item2"}},
	}

	Map(o, d)

	if !reflect.DeepEqual(d, want) {
		t.Errorf("Map(%v, %v) = %#v, want %#v", o, d, d.Items, want.Items)
	}
}

func Test_Map_nested_struct(t *testing.T) {
	type nestedSource struct{ Field string }
	type nestedTarget struct{ Field string }
	type source struct {
		Name   string
		Age    int
		Nested nestedSource
	}
	type target struct {
		Name   string
		Age    int
		Nested nestedTarget
	}

	o := &source{
		Name:   "name",
		Age:    10,
		Nested: nestedSource{Field: "nested"},
	}

	d := &target{}

	want := &target{
		Name:   "name",
		Age:    10,
		Nested: nestedTarget{Field: "nested"},
	}

	Map(o, d)

	if !reflect.DeepEqual(d, want) {
		t.Errorf("Map(%v, %v) = %v, want %v", o, d, d, want)
	}
}

func Test_Map_pkg_visibility(t *testing.T) {
	type args struct {
		source *TestSource
		target *TestTarget
		rules  RulesSet
		args   []interface{}
	}
	var test = []struct {
		name string
		args args
		want TestTarget
	}{
		{
			name: "Mapping structs,unexportedField->unexportedField,exportedField->exportedField",
			args: args{
				source: &TestSource{
					fieldS1:              "valueS1",
					fieldS2:              314,
					fieldUnexportedEqual: 314,
					FieldExportedEqual:   598,
					singleListFieldS:     []string{"valueS1", "valueS2"},
				},
				target: &TestTarget{},
			},
			want: TestTarget{
				fieldT1:              "",
				fieldT2:              "",
				fieldT3:              0,
				ignorableTargetField: "",
				fieldUnexportedEqual: 314,
				FieldExportedEqual:   598,
				singleListFieldT:     []string(nil),
			},
		},
		{
			name: "Mapping structs,unexportedField->exportedField,exportedField->unexportedField",
			args: args{
				source: &TestSource{
					fieldS1:              "valueS1",
					fieldS2:              314,
					fieldUnexportedEqual: 314,
					FieldExportedEqual:   598,
					singleListFieldS:     []string{"valueS1", "valueS2"},
				},
				target: &TestTarget{},
				rules: RulesSet{
					"fieldUnexportedEqual": "FieldExportedEqual",
					"FieldExportedEqual":   "fieldUnexportedEqual",
					"singleListFieldT":     "singleListFieldS",
				},
			},
			want: TestTarget{
				fieldT1:              "",
				fieldT2:              "",
				fieldT3:              0,
				ignorableTargetField: "",
				fieldUnexportedEqual: 598,
				FieldExportedEqual:   314,
				singleListFieldT:     []string{"valueS1", "valueS2"},
			},
		},
	}

	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			rulesRegistry = make(mapperRulesRegistry)
			if tt.args.rules != nil {
				RegisterRulesDefinitions(RulesDefinition{
					*tt.args.source,
					*tt.args.target,
					tt.args.rules,
				})
			}
			if tt.args.args != nil {
				Map(tt.args.source, tt.args.target, tt.args.args...)
			} else {
				Map(tt.args.source, tt.args.target)
			}
			if !reflect.DeepEqual(tt.args.target, &tt.want) {
				t.Errorf("Map() = %#v, want %#v", tt.args.target, &tt.want)
			}
		})
	}
}

func Test_Map_pkg_func_and_args(t *testing.T) {
	type args struct {
		source *TestSource
		target *TestTarget
		rules  RulesSet
		args   []interface{}
	}
	var tests = []struct {
		name string
		args args
		want TestTarget
	}{
		{
			name: "Mapping structs,constant",
			args: args{
				source: &TestSource{
					fieldS1:              "",
					fieldS2:              0,
					fieldUnexportedEqual: 0,
					FieldExportedEqual:   0,
				},
				target: &TestTarget{},
				rules: RulesSet{
					"fieldT3": func() int { return 596 },
				},
			},
			want: TestTarget{
				fieldT1:              "",
				fieldT2:              "",
				fieldT3:              596,
				ignorableTargetField: "",
				fieldUnexportedEqual: 0,
				FieldExportedEqual:   0,
				singleListFieldT:     []string(nil),
			},
		},
		{
			name: "Mapping structs,1 simple arg",
			args: args{
				source: &TestSource{
					fieldS1:              "",
					fieldS2:              0,
					fieldUnexportedEqual: 0,
					FieldExportedEqual:   0,
				},
				target: &TestTarget{},
				rules: RulesSet{
					"fieldT3": func(i int) int { return i },
				},
				args: []interface{}{3454},
			},
			want: TestTarget{
				fieldT1:              "",
				fieldT2:              "",
				fieldT3:              3454,
				ignorableTargetField: "",
				fieldUnexportedEqual: 0,
				FieldExportedEqual:   0,
				singleListFieldT:     []string(nil),
			},
		},
		{
			name: "Mapping structs,2 simple args",
			args: args{
				source: &TestSource{
					fieldS1:              "",
					fieldS2:              0,
					fieldUnexportedEqual: 0,
					FieldExportedEqual:   0,
				},
				target: &TestTarget{},
				rules: RulesSet{
					"fieldT2": func(s string) string { return s },
					"fieldT3": func(i int) int { return i },
				},
				args: []interface{}{3454, "hello"},
			},
			want: TestTarget{
				fieldT1:              "",
				fieldT2:              "hello",
				fieldT3:              3454,
				ignorableTargetField: "",
				fieldUnexportedEqual: 0,
				FieldExportedEqual:   0,
				singleListFieldT:     []string(nil),
			},
		},
		{
			name: "Mapping structs,2 simple args,actual source struct",
			args: args{
				source: &TestSource{
					fieldS1:              "qwerty",
					fieldS2:              0,
					fieldUnexportedEqual: 0,
					FieldExportedEqual:   0,
				},
				target: &TestTarget{},
				rules: RulesSet{
					"fieldT1": func(s TestSource) string { return s.fieldS1 },
					"fieldT2": func(s string) string { return s },
					"fieldT3": func(i int) int { return i },
				},
				args: []interface{}{3454, "hello"},
			},
			want: TestTarget{
				fieldT1:              "qwerty",
				fieldT2:              "hello",
				fieldT3:              3454,
				ignorableTargetField: "",
				fieldUnexportedEqual: 0,
				FieldExportedEqual:   0,
				singleListFieldT:     []string(nil),
			},
		},
		{
			name: "Mapping structs,2 simple args,pointer to actual source struct",
			args: args{
				source: &TestSource{
					fieldS1:              "qwerty",
					fieldS2:              0,
					fieldUnexportedEqual: 0,
					FieldExportedEqual:   0,
				},
				target: &TestTarget{},
				rules: RulesSet{
					"fieldT1": func(s TestSource) string { return s.fieldS1 },
					"fieldT2": func(s string) string { return s },
					"fieldT3": func(i int) int { return i },
				},
				args: []interface{}{3454, "hello"},
			},
			want: TestTarget{
				fieldT1:              "qwerty",
				fieldT2:              "hello",
				fieldT3:              3454,
				ignorableTargetField: "",
				fieldUnexportedEqual: 0,
				FieldExportedEqual:   0,
				singleListFieldT:     []string(nil),
			},
		},
		{
			name: "Mapping structs,zero value",
			args: args{
				source: &TestSource{},
				target: &TestTarget{},
				rules: RulesSet{
					"fieldT2": func(s string) string { return s },
				},
				args: []interface{}{},
			},
			want: TestTarget{
				fieldT1:              "",
				fieldT2:              "",
				fieldT3:              0,
				ignorableTargetField: "",
				fieldUnexportedEqual: 0,
				FieldExportedEqual:   0,
				singleListFieldT:     []string(nil),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.rules != nil {
				rulesRegistry = make(mapperRulesRegistry)
				RegisterRulesDefinitions(RulesDefinition{
					*tt.args.source,
					*tt.args.target,
					tt.args.rules,
				})
			}
			if tt.args.args != nil {
				Map(tt.args.source, tt.args.target, tt.args.args...)
			} else {
				Map(tt.args.source, tt.args.target)
			}
			if !reflect.DeepEqual(tt.args.target, &tt.want) {
				t.Errorf("Map() = %#v, want %#v", tt.args.target, &tt.want)
			}
		})
	}
}

func Test_Map_simple_slices(t *testing.T) {
	type source struct {
		vSliceV1 []int
		vSliceV2 []int
	}
	type target struct {
		vSliceV1 []string
		vSliceV2 []int
	}

	type args struct {
		source *source
		target *target
		rules  RulesSet
	}
	type test struct {
		name string
		args args
		want target
	}
	tests := []test{
		{
			name: "Source slice with different type than destination",
			args: args{
				source: &source{vSliceV1: []int{1, 2, 4}},
				target: &target{vSliceV1: []string(nil)},
			},
			want: target{
				vSliceV1: []string(nil),
				vSliceV2: []int(nil),
			},
		},
		{
			name: "Source slice with same type as the destination",
			args: args{
				source: &source{vSliceV2: []int{0, 1, 1, 2, 3, 5}},
				target: &target{vSliceV1: []string(nil)},
			},
			want: target{
				vSliceV1: []string(nil),
				vSliceV2: []int{0, 1, 1, 2, 3, 5},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.rules != nil {
				rulesRegistry = make(mapperRulesRegistry)
				RegisterRulesDefinitions(RulesDefinition{
					tt.args.source,
					tt.args.target,
					tt.args.rules,
				})
			}
			Map(tt.args.source, tt.args.target)
			if !reflect.DeepEqual(tt.args.target, &tt.want) {
				t.Errorf("Map() = %#v, want %#v", tt.args.target, &tt.want)
			}
		})
	}
}

func Test_Map_complex_slices(t *testing.T) {
	type source struct {
		field1 string
		field2 int
	}
	type target struct {
		field1 string
		field2 int
	}

	type complexSourceSlice struct{ complexSlice []source }
	type complexTargetSlice struct{ complexSlice []target }

	type args struct {
		source *complexSourceSlice
		target *complexTargetSlice
		rules  RulesSet
	}
	type test struct {
		name string
		args args
		want complexTargetSlice
	}
	tests := []test{
		{
			name: "Source slice with different type than destination",
			args: args{
				source: &complexSourceSlice{
					complexSlice: []source{
						{
							field1: "i1",
							field2: 1,
						},
						{
							field1: "i2",
							field2: 2,
						},
					},
				},
				target: &complexTargetSlice{},
			},
			want: complexTargetSlice{
				complexSlice: []target{
					{
						field1: "i1",
						field2: 1,
					},
					{
						field1: "i2",
						field2: 2,
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.rules != nil {
				rulesRegistry = make(mapperRulesRegistry)
				RegisterRulesDefinitions(RulesDefinition{
					tt.args.source,
					tt.args.target,
					tt.args.rules,
				})
			}
			Map(tt.args.source, tt.args.target)

			for _, w := range tt.want.complexSlice {
				var exists bool
				for _, m := range tt.args.target.complexSlice {
					if reflect.DeepEqual(w, m) {
						exists = true
						break
					}
				}
				if !exists {
					t.Errorf("%#v is not present in %#v", w, tt.args.target.complexSlice)
				}
			}
		})
	}
}

func Test_Map_simple_array(t *testing.T) {
	type source struct {
		vSliceV1 [5]int
		vSliceV2 [5]int
	}
	type target struct {
		vSliceV1 [5]string
		vSliceV2 [5]int
	}

	type args struct {
		source *source
		target *target
		rules  RulesSet
	}
	type test struct {
		name string
		args args
		want target
	}
	tests := []test{
		{
			name: "Source array with different type than destination",
			args: args{
				source: &source{vSliceV1: [5]int{1, 2, 4}},
				target: &target{},
			},
			want: target{
				vSliceV1: [5]string{},
				vSliceV2: [5]int{},
			},
		},
		{
			name: "Source array with same type as the destination",
			args: args{
				source: &source{vSliceV2: [5]int{0, 1, 1, 2, 3}},
				target: &target{},
			},
			want: target{
				vSliceV1: [5]string{},
				vSliceV2: [5]int{0, 1, 1, 2, 3},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.rules != nil {
				rulesRegistry = make(mapperRulesRegistry)
				RegisterRulesDefinitions(RulesDefinition{
					tt.args.source,
					tt.args.target,
					tt.args.rules,
				})
			}
			Map(tt.args.source, tt.args.target)
			if !reflect.DeepEqual(tt.args.target, &tt.want) {
				t.Errorf("Map() = %#v, want %#v", tt.args.target, &tt.want)
			}
		})
	}
}

func Test_Map_complex_array(t *testing.T) {
	type source struct {
		field1 string
		field2 int
	}
	type target struct {
		field1 string
		field2 int
	}

	type complexSourceSlice struct{ complexSlice [5]source }
	type complexTargetSlice struct{ complexSlice [5]target }

	type args struct {
		source *complexSourceSlice
		target *complexTargetSlice
		rules  RulesSet
	}
	type test struct {
		name string
		args args
		want complexTargetSlice
	}
	tests := []test{
		{
			name: "Source slice with different type than destination",
			args: args{
				source: &complexSourceSlice{
					complexSlice: [5]source{
						{
							field1: "i1",
							field2: 1,
						},
						{
							field1: "i2",
							field2: 2,
						},
					},
				},
				target: &complexTargetSlice{},
			},
			want: complexTargetSlice{
				complexSlice: [5]target{
					{
						field1: "i1",
						field2: 1,
					},
					{
						field1: "i2",
						field2: 2,
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.rules != nil {
				rulesRegistry = make(mapperRulesRegistry)
				RegisterRulesDefinitions(RulesDefinition{
					tt.args.source,
					tt.args.target,
					tt.args.rules,
				})
			}
			Map(tt.args.source, tt.args.target)

			for _, w := range tt.want.complexSlice {
				var exists bool
				for _, m := range tt.args.target.complexSlice {
					if reflect.DeepEqual(w, m) {
						exists = true
						break
					}
				}
				if !exists {
					t.Errorf("%#v is not present in %#v", w, tt.args.target.complexSlice)
				}
			}
		})
	}
}

func Test_Map_simple_map(t *testing.T) {
	type source struct {
		vSliceV1 map[string]int
		vSliceV2 map[string]int
	}
	type target struct {
		vSliceV1 map[int]int
		vSliceV2 map[string]int
	}

	type args struct {
		source *source
		target *target
		rules  RulesSet
	}
	type test struct {
		name string
		args args
		want target
	}
	tests := []test{
		{
			name: "Source map with different type than destination",
			args: args{
				source: &source{
					vSliceV1: map[string]int{
						"one": 1,
						"two": 2,
					},
				},
				target: &target{},
			},
			want: target{
				vSliceV1: map[int]int(nil),
				vSliceV2: map[string]int(nil),
			},
		},
		{
			name: "Source array with same type as the destination",
			args: args{
				source: &source{
					vSliceV1: map[string]int{
						"one": 1,
						"two": 2,
					},
					vSliceV2: map[string]int{
						"one": 1,
						"two": 2,
					},
				},
				target: &target{},
			},
			want: target{
				vSliceV1: map[int]int(nil),
				vSliceV2: map[string]int{
					"one": 1,
					"two": 2,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.rules != nil {
				rulesRegistry = make(mapperRulesRegistry)
				RegisterRulesDefinitions(RulesDefinition{
					tt.args.source,
					tt.args.target,
					tt.args.rules,
				})
			}
			Map(tt.args.source, tt.args.target)
			if !reflect.DeepEqual(tt.args.target, &tt.want) {
				t.Errorf("Map() = %#v, want %#v", tt.args.target, &tt.want)
			}
		})
	}
}

func Test_Map_complex_map(t *testing.T) {
	type source struct {
		Field1 string
		Field2 int
	}
	type target struct {
		Field1 string
		Field2 int
	}
	type complexSourceMap struct{ ComplexMap map[string]source }
	type complexTargetMap struct{ ComplexMap map[string]target }

	type args struct {
		source *complexSourceMap
		target *complexTargetMap
		rules  RulesSet
	}
	type test struct {
		name string
		args args
		want complexTargetMap
	}
	tests := []test{
		{
			name: "Source map with same type as the destination",
			args: args{
				source: &complexSourceMap{
					ComplexMap: map[string]source{
						"one": {
							Field1: "one val",
							Field2: 1,
						},
						"two": {
							Field1: "two val",
							Field2: 2,
						},
					},
				},
				target: &complexTargetMap{},
			},
			want: complexTargetMap{
				ComplexMap: map[string]target{
					"one": {
						Field1: "one val",
						Field2: 1,
					},
					"two": {
						Field1: "two val",
						Field2: 2,
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.rules != nil {
				rulesRegistry = make(mapperRulesRegistry)
				RegisterRulesDefinitions(RulesDefinition{
					tt.args.source,
					tt.args.target,
					tt.args.rules,
				})
			}
			Map(tt.args.source, tt.args.target)
			if !reflect.DeepEqual(tt.args.target, &tt.want) {
				t.Errorf("Map() = %#v, want %#v", tt.args.target, &tt.want)
			}
		})
	}
}

func Test_Map_complex_map_unexported_case_1(t *testing.T) {
	type source struct {
		field1 string
		field2 int
	}
	type target struct {
		field1 string
		field2 int
	}
	type complexSourceMap struct{ ComplexMap map[string]source }
	type complexTargetMap struct{ ComplexMap map[string]target }

	type args struct {
		source *complexSourceMap
		target *complexTargetMap
		rules  RulesSet
	}
	type test struct {
		name string
		args args
		want complexTargetMap
	}
	tests := []test{
		{
			name: "Source map with same type as the destination",
			args: args{
				source: &complexSourceMap{
					ComplexMap: map[string]source{
						"one": {
							field1: "one val",
							field2: 1,
						},
						"two": {
							field1: "two val",
							field2: 2,
						},
					},
				},
				target: &complexTargetMap{},
			},
			want: complexTargetMap{
				ComplexMap: map[string]target{
					"one": {
						field1: "",
						field2: 0,
					},
					"two": {
						field1: "",
						field2: 0,
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.rules != nil {
				rulesRegistry = make(mapperRulesRegistry)
				RegisterRulesDefinitions(RulesDefinition{
					tt.args.source,
					tt.args.target,
					tt.args.rules,
				})
			}
			Map(tt.args.source, tt.args.target)
			if !reflect.DeepEqual(tt.args.target, &tt.want) {
				t.Errorf("Map() = %#v, want %#v", tt.args.target, &tt.want)
			}
		})
	}
}

func Test_Map_complex_map_unexported_case_2(t *testing.T) {
	type source struct {
		Field1 string
		Field2 int
	}
	type target struct {
		Field1 string
		Field2 int
	}
	type complexSourceMap struct{ complexMap map[string]source }
	type complexTargetMap struct{ complexMap map[string]target }

	type args struct {
		source *complexSourceMap
		target *complexTargetMap
		rules  RulesSet
	}
	type test struct {
		name string
		args args
		want complexTargetMap
	}
	tests := []test{
		{
			name: "Source map with same type as the destination",
			args: args{
				source: &complexSourceMap{
					complexMap: map[string]source{
						"one": {
							Field1: "one val",
							Field2: 1,
						},
						"two": {
							Field1: "two val",
							Field2: 2,
						},
					},
				},
				target: &complexTargetMap{},
			},
			want: complexTargetMap{
				complexMap: map[string]target(nil),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.rules != nil {
				rulesRegistry = make(mapperRulesRegistry)
				RegisterRulesDefinitions(RulesDefinition{
					tt.args.source,
					tt.args.target,
					tt.args.rules,
				})
			}
			Map(tt.args.source, tt.args.target)
			if !reflect.DeepEqual(tt.args.target, &tt.want) {
				t.Errorf("Map() = %#v, want %#v", tt.args.target, &tt.want)
			}
		})
	}
}

func Test_Map_complex_map_unexported_case_3(t *testing.T) {
	type source struct {
		field1 string
		field2 int
	}
	type target struct {
		field1 string
		field2 int
	}
	type complexSourceMap struct{ unExComplexMap map[string]source }
	type complexTargetMap struct{ ExComplexMap map[string]target }

	type args struct {
		source *complexSourceMap
		target *complexTargetMap
		rules  RulesSet
	}
	type test struct {
		name string
		args args
		want complexTargetMap
	}
	tests := []test{
		{
			name: "Source map with same type as the destination",
			args: args{
				source: &complexSourceMap{
					unExComplexMap: map[string]source{
						"one": {
							field1: "one val",
							field2: 1,
						},
						"two": {
							field1: "two val",
							field2: 2,
						},
					},
				},
				target: &complexTargetMap{},
				rules: RulesSet{
					"ExComplexMap": "unExComplexMap",
				},
			},
			want: complexTargetMap{
				ExComplexMap: map[string]target{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.rules != nil {
				rulesRegistry = make(mapperRulesRegistry)
				RegisterRulesDefinitions(RulesDefinition{
					*tt.args.source,
					*tt.args.target,
					tt.args.rules,
				})
			}
			Map(tt.args.source, tt.args.target)
			if !reflect.DeepEqual(tt.args.target, &tt.want) {
				t.Errorf("Map() = %#v, want %#v", tt.args.target, &tt.want)
			}
		})
	}
}
