package structsconv

import (
	"fmt"
	"log"
	"reflect"
)

func checkSetOfRules(set reflect.Value) {
	// set must be a struct or a pointer to a struct
	if set.Kind() != reflect.Ptr && (set.Kind() != reflect.Ptr && set.Kind() != reflect.Struct) {
		log.Panicf("ERROR: Set of Rules must be a pointer to a struct or a struct.\n")
	}
}

func checkSetDefinitionSupplier(supplier reflect.Value) {
	// supplier must be a function
	if supplier.Kind() != reflect.Func {
		log.Panicf("ERROR: Wrong type of RulesDefinition supplier.\n")
	}

	// supplier with no arguments
	if supplier.Type().NumIn() != 0 {
		log.Panicf("ERROR: Wrong number of arguments in RulesDefinition supplier.\n")
	}

	// supplier with one return value
	if supplier.Type().NumOut() != 1 {
		log.Panicf("ERROR: Wrong number of return values in RulesDefinition supplier.\n")
	}

	// supplier return value must be a pointer to a struct
	if supplier.Type().Out(0).Kind() != reflect.Struct || supplier.Type().Out(0).Kind() != reflect.Ptr {
		log.Panicf("ERROR: Wrong type of return value in RulesDefinition supplier. "+
			"Expected 'RulesDefinition' but got '%s'.\n", supplier.Type().Out(0).Kind().String(),
		)
	}

	// if supplier return value is a pointer, it must be a pointer to a struct
	if supplier.Type().Out(0).Kind() == reflect.Ptr && supplier.Type().Out(0).Elem().Kind() != reflect.Struct {
		log.Panicf("ERROR: Wrong type of return value in RulesDefinition supplier. "+
			"Expected 'RulesDefinition' but got '%s'.\n", supplier.Type().Out(0).Kind().String(),
		)
	}

	// if supplier return value is a struct, it must be a RulesDefinition struct
	if supplier.Type().Out(0).Kind() == reflect.Struct && supplier.Type().Out(0).Name() != "RulesDefinition" {
		log.Panicf("ERROR: Wrong type of return value in RulesDefinition supplier. "+
			"Expected 'RulesDefinition' but got '%s'.\n", supplier.Type().Out(0).Name(),
		)
	}

	// if supplier return value is a pointer, it must be a pointer to a RulesDefinition struct
	if supplier.Type().Out(0).Kind() == reflect.Ptr && supplier.Type().Out(0).Elem().Name() != "RulesDefinition" {
		log.Panicf("ERROR: Wrong type of return value in RulesDefinition supplier. "+
			"Expected 'RulesDefinition' but got '%s'.\n", supplier.Type().Out(0).Elem().Name(),
		)
	}
}

// checkRootValuesTypes checks if the ROOT source and target types are valid.
func checkRootValuesTypes(st, tt reflect.Value) error {
	if st.Kind() != reflect.Ptr {
		return fmt.Errorf("rules error: source must be a pointer")
	}
	if tt.Kind() != reflect.Ptr {
		return fmt.Errorf("rules error: target must be a pointer")
	}
	if st.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("rules error: source must be a pointer to a struct")
	}
	if tt.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("rules error: target must be a pointer to a struct")
	}
	return nil
}

// checkMapperRules checks if the mapper rules are valid
func checkMapperRules(key rulesKey, rules RulesSet) {
	log.Printf("Checking rules for mapping (%s -> %s).\n", key.source.String(), key.target.String())
	for k, r := range rules {
		if r == nil { // nil rule == ignore field
			log.Printf(
				"INFO: (%s -> %s) Field '%s' is marked as ignored.\n",
				key.source.String(), key.target.String(), k,
			)
			continue
		}
		checkTargetKeyName(k, key)

		t := reflect.TypeOf(r)
		switch t.Kind() {
		case reflect.String: // mapping source field name
			checkMappingName(r.(string), k, key)
		case reflect.Func: // mapping target field value from function
			checkFunc(t, k, key)
		default: // not valid rule
			log.Panicf(
				"ERROR: (%s -> %s) Rule '%s' is not valid. Rule = %v.\n",
				key.source.String(), key.target.String(), k, r,
			)
		}
	}
}

// checkTargetKeyName checks if field name (ruleKey) is present in target struct
func checkTargetKeyName(ruleKeyValue string, key rulesKey) {
	_, exist := key.target.FieldByName(ruleKeyValue)
	if !exist {
		log.Panicf(
			"ERROR: (%s -> %s) Field '%s' is not present in target struct %s.\n",
			key.source.String(), key.target.String(), ruleKeyValue, key.target.String(),
		)
	}
}

// checkMappingName checks field MappingName in source struct
//	- MappingName is present in source struct
// 	- field kind is the same in origin and target struct
func checkMappingName(mappingName, ruleKey string, key rulesKey) {
	sf, exist := key.source.FieldByName(mappingName)
	if !exist { // checks if MappingName is present in origin struct
		log.Panicf(
			"ERROR: (%s -> %s) Field '%s' is not present in source struct %s.\n",
			key.source.String(), key.target.String(), mappingName, key.source.String(),
		)
	}
	tf, _ := key.target.FieldByName(ruleKey)

	switch {
	case sf.Type.Kind() == reflect.Ptr && tf.Type.Kind() == reflect.Ptr && sf.Type.Elem().Kind() == tf.Type.Elem().Kind(), // both are pointers to the same type
		sf.Type.Kind() == reflect.Ptr && tf.Type.Kind() != reflect.Ptr && sf.Type.Elem().Kind() == tf.Type.Kind(), // source is pointer to the same type as target
		sf.Type.Kind() != reflect.Ptr && tf.Type.Kind() == reflect.Ptr && sf.Type.Kind() == tf.Type.Elem().Kind(), // target is pointer to the same type as source
		sf.Type.Kind() == tf.Type.Kind(): // both are the same type
		return
	default:
		log.Panicf(
			"ERROR: (%s -> %s) Field '%s' has different type in source (%s:%s) and target (%s:%s) structs.\n",
			key.source.String(), key.target.String(), ruleKey, mappingName, sf.Type.String(), ruleKey, tf.Type.String(),
		)
	}
}

// checkFunc checks if function is valid according to the following criteria:
//	- the function returns a value of the same type as the target
func checkFunc(f reflect.Type, ruleKey string, key rulesKey) {
	// checks if the function returns a value of the same type as the target
	if getFieldByName(ruleKey, key.target).Type != f.Out(0) {
		log.Panicf(
			"ERROR: (%s -> %s) Function '%s' must return type '%s', currently returns '%s'. Function = '%s'.\n",
			key.source.String(), key.target.String(), ruleKey, getFieldByName(ruleKey, key.target).Type.String(), f.Out(0).String(), f.String(),
		)
	}
}

// getFieldByName returns field by name
func getFieldByName(n string, t reflect.Type) reflect.StructField {
	f, _ := t.FieldByName(n)
	return f
}
