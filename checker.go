package structsconv

import (
	"fmt"
	"log"
	"reflect"
)

// checkType checks if the ROOT source and target types are valid.
func checkType(st, tt reflect.Value) error {
	if st.Kind() != reflect.Ptr {
		return fmt.Errorf("MapRules error. Source must be a pointer")
	}
	if tt.Kind() != reflect.Ptr {
		return fmt.Errorf("MapRules error. Target must be a pointer")
	}
	if st.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("MapRules error. Source must be a pointer to a struct")
	}
	if tt.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("MapRules error. Target must be a pointer to a struct")
	}
	return nil
}

// checkMapperRules checks if the mapper rules are valid
func checkMapperRules(key rulesKey, rules RulesSet) {
	log.Printf("Checking rules for mapping (%s -> %s).\n", key.source.String(), key.target.String())
	for k, r := range rules {
		if r == nil { // nil rule == ignore field
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
			log.Fatalf(
				"ERROR: (%s -> %s) Rule '%s' is not valid. Rule = '%s'.\n",
				key.source.String(), key.target.String(), k, reflect.TypeOf(t).String(),
			)
		}
	}
}

// checkTargetKeyName checks if field name (ruleKey) is present in target struct
func checkTargetKeyName(ruleKeyValue string, key rulesKey) {
	_, exist := key.target.FieldByName(ruleKeyValue)
	if !exist {
		log.Fatalf(
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
		log.Fatalf(
			"ERROR: (%s -> %s) Field '%s' is not present in source struct %s. Value = '%s'.\n",
			key.source.String(), key.target.String(), ruleKey, key.source.String(), mappingName,
		)
	}
	tf, _ := key.target.FieldByName(ruleKey)
	if sf.Type.Kind() != tf.Type.Kind() { // checks if field type is the same in origin and target struct
		log.Fatalf(
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
		log.Fatalf(
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
