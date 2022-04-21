package structsconv

import (
	"log"
	"reflect"
	"sync"
	"unsafe"
)

// rulesRegistry mapping rules register.
var rulesRegistry = make(mapperRulesRegistry)

// RegisterRulesDefinitions it is used to register rule definitions.
func RegisterRulesDefinitions(definitions ...RulesDefinition) {
	for _, d := range definitions {
		registerRules(d.Source, d.Target, d.Rules)
	}
}

// registerRules verifies and registers a mapper rules for specific mapping from structure to structure.
func registerRules(source interface{}, target interface{}, rules RulesSet) {
	key := buildKey(source, target)
	_, exists := rulesRegistry[key]
	if exists {
		log.Panicf("ERROR: Mapper with rulesKey (%s -> %s) already exists.", key.source, key.target)
	}
	checkMapperRules(key, rules)
	rulesRegistry[key] = rules
}

// Map maps the source structure to a destination structure; source and target must be pointers to structs.
//
// The value will be mapped into the target structure.
func Map(source interface{}, target interface{}, args ...interface{}) {
	var wg sync.WaitGroup
	sourceV := reflect.ValueOf(source)
	targetV := reflect.ValueOf(target)
	if err := checkRootValuesTypes(sourceV, targetV); err != nil { // check if the source and target are pointers
		log.Panicf("ERROR: %s", err)
	}
	args = append([]interface{}{sourceV.Interface()}, args...)
	args = append(args, sourceV.Elem().Interface())
	structToStruct(sourceV.Elem(), targetV.Elem(), sourceV.Elem().Interface(), groupArgs(args), &wg)
	wg.Wait()
}

// groupArgs groups the arguments by their type.
func groupArgs(args []interface{}) groupedArgs {
	g := make(groupedArgs)
	for _, v := range args {
		t := reflect.TypeOf(v)
		g[t] = append(g[t], v)
	}
	return g
}

// structToStruct maps the source struct to the target struct
func structToStruct(source, target reflect.Value, actualS interface{}, args groupedArgs, wg *sync.WaitGroup) {
	key := rulesKey{source.Type(), target.Type()}
	rules := rulesRegistry[key]
	targetType := target.Type()

	for i := 0; i < target.NumField(); i++ {
		targetFieldName := targetType.Field(i).Name
		targetValue := target.Field(i)

		// if there is a rule for this field, use it
		if mapper, exists := rules[targetFieldName]; exists {
			// if the rule is not nil (is not ignorable) apply rule
			if mapper != nil {
				applyRule(source, targetValue, mapper, actualS, args, wg)
			}
			continue
		}

		// field-to-field mapping source field to target field by target field name
		sourceValue := source.FieldByName(targetFieldName)
		if sourceValue.IsValid() {
			if pType := fieldToField(sourceValue, targetValue, args, wg); pType == incompatibleTypes {
				logIgnoringMappingForIncompatibleTypes(key, targetFieldName, sourceValue, targetValue)
			}
			continue
		}

		// A target field without mapping value in source
		logTargetFieldWithoutMappingValueInSource(key, targetFieldName)
	}
}

// applyRule processes a rule for a target field.
func applyRule(source, targetValue reflect.Value, mapper, actualS interface{}, args groupedArgs, wg *sync.WaitGroup) {
	switch mapperValue := reflect.ValueOf(mapper); mapperValue.Kind() {
	case reflect.String: // mapper has the name of the source field
		fieldToField(source.FieldByName(mapper.(string)), targetValue, args, wg)
	default: // mapper is a function
		callFunc(targetValue, mapperValue, actualS, args)
	}
}

// callFunc calls a function with the given arguments
func callFunc(targetValue, mapperValue reflect.Value, actualS interface{}, args groupedArgs) {
	method := mapperValue.Type()
	if method.NumIn() == 0 {
		mappingDirectMapping(mapperValue.Call([]reflect.Value{})[0], targetValue)
	} else {
		params := getMethodParams(method, args, actualS)
		mappingDirectMapping(mapperValue.Call(params)[0], targetValue)
	}
}

// getMethodParams gets the arguments of the method based on its input parameters
func getMethodParams(method reflect.Type, args groupedArgs, current interface{}) []reflect.Value {
	params := make([]reflect.Value, method.NumIn())
	argsCounter := make(map[reflect.Type]int)
	cType := reflect.TypeOf(current)
	var cFlag bool
	for i := 0; i < method.NumIn(); i++ {
		mType := method.In(i)

		// Current struct
		if !cFlag && cType == mType {
			params[i] = reflect.ValueOf(current)
			cFlag = true
			continue
		}

		// Arguments
		if l, exists := args[mType]; exists && argsCounter[mType] < len(l) {
			c, _ := argsCounter[mType]
			params[i] = reflect.ValueOf(l[c])
			argsCounter[mType]++
			continue
		}

		// Zero value
		logPassingZeroValue(method, mType, i+1)
		params[i] = reflect.Zero(mType)
	}
	return params
}

// fieldToField field to field mapping orchestration
func fieldToField(sourceValue, targetValue reflect.Value, args groupedArgs, wg *sync.WaitGroup) processingResultType {
	mappingType := getMappingType(sourceValue, targetValue)
	switch mappingType {
	case structsMapping:
		wg.Add(1)
		go cMappingStructLogic(sourceValue, targetValue, args, wg)
	case slicesMapping:
		wg.Add(1)
		go cMappingSliceLogic(sourceValue, targetValue, args, wg)
	case mapsMapping:
		wg.Add(1)
		go cMappingMapLogic(sourceValue, targetValue, args, wg)
	case arraysMapping:
		wg.Add(1)
		go cMappingArrayLogic(sourceValue, targetValue, args, wg)
	case directMapping:
		mappingDirectMapping(sourceValue, targetValue)
	case ptrMapping:
		mappingPtrMapping(sourceValue, targetValue, args, wg)
	default:
		return mappingType
	}
	return mappingType
}

// cMappingStructLogic is used to be called as a goroutine and map the source field to the destination field in a concurrent way
func cMappingStructLogic(source, target reflect.Value, args groupedArgs, wg *sync.WaitGroup) {
	defer wg.Done()
	if !source.CanInterface() {
		source = getUnexportedField(source)
	}
	structToStruct(source, target, source.Interface(), args, wg)
}

// cMappingMapLogic is used to be called as a goroutine and maps the structures of the source map to the destination map in a concurrent way.
func cMappingMapLogic(sourceValue, targetValue reflect.Value, args groupedArgs, wg *sync.WaitGroup) {
	defer wg.Done()

	if !targetValue.CanInterface() {
		log.Printf(
			"WARNING: Operations on map type fields that are not exported are not supported. Operation ignored. Target = %s\n",
			targetValue.Type().String(),
		)
		return
	}

	var mapWg sync.WaitGroup
	defer mapWg.Wait()
	itemType := targetValue.Type().Elem()
	var lock sync.Mutex
	mappingDirectMapping(reflect.MakeMap(targetValue.Type()), targetValue)
	for _, key := range sourceValue.MapKeys() {
		mapWg.Add(1)
		go func(key reflect.Value) {
			defer mapWg.Done()
			var mapItemWg sync.WaitGroup
			item := reflect.New(itemType)
			sourceItem := sourceValue.MapIndex(key)
			if !sourceItem.CanInterface() {
				log.Printf(
					"WARNING: Operations on MAP type fields that are not exported are not supported. Operation ignored. Source Item = %s\n",
					sourceItem.Type().String(),
				)
				return
			}
			structToStruct(sourceItem, item.Elem(), sourceItem.Interface(), args, &mapItemWg)
			mapItemWg.Wait()
			lock.Lock()
			targetValue.SetMapIndex(key, item.Elem())
			lock.Unlock()
		}(key)
	}
}

// cMappingArrayLogic is used to be called as a goroutine and maps the structures of the source array to the destination array in a concurrent way.
func cMappingArrayLogic(sourceValue, targetValue reflect.Value, args groupedArgs, wg *sync.WaitGroup) {
	defer wg.Done()
	if !targetValue.CanInterface() {
		targetValue = getUnexportedField(targetValue)
	}
	var arrayWg sync.WaitGroup
	defer arrayWg.Wait()
	itemType := targetValue.Type().Elem()
	var lock sync.Mutex
	for i := 0; i < targetValue.Cap(); i++ {
		arrayWg.Add(1)
		go func(i int) {
			defer arrayWg.Done()
			var arrayItemWg sync.WaitGroup
			item := reflect.New(itemType)
			sourceItem := sourceValue.Index(i)
			if !sourceItem.CanInterface() {
				sourceItem = getUnexportedField(sourceItem)
			}
			structToStruct(sourceItem, item.Elem(), sourceItem.Interface(), args, &arrayItemWg)
			arrayItemWg.Wait()
			lock.Lock()
			targetValue.Index(i).Set(item.Elem())
			lock.Unlock()
		}(i)
	}
}

// cMappingSliceLogic is used to be called as a goroutine and maps the structures of the source slice to the destination slice in a concurrent way.
func cMappingSliceLogic(sourceValue, targetValue reflect.Value, args groupedArgs, wg *sync.WaitGroup) {
	defer wg.Done()
	if !targetValue.CanInterface() {
		targetValue = getUnexportedField(targetValue)
	}
	itemType := targetValue.Type().Elem()
	for i := 0; i < sourceValue.Len(); i++ {
		item := reflect.New(itemType)
		sourceItem := sourceValue.Index(i)
		if !sourceItem.CanInterface() {
			sourceItem = getUnexportedField(sourceItem)
		}

		if itemType.Kind() == reflect.Ptr {
			mappingPtrMapping(sourceItem, item.Elem(), args, wg)
		} else if sourceItem.Kind() == reflect.Ptr {
			structToStruct(sourceItem.Elem(), item.Elem(), sourceItem.Interface(), args, wg)
		} else {
			structToStruct(sourceItem, item.Elem(), sourceItem.Interface(), args, wg)
		}
		mappingDirectMapping(reflect.Append(targetValue, item.Elem()), targetValue)
	}
}

// mappingPtrMapping is used to map ptr types
func mappingPtrMapping(sourceValue, targetValue reflect.Value, args groupedArgs, wg *sync.WaitGroup) {
	switch {
	case sourceValue.Kind() == reflect.Ptr && targetValue.Kind() != reflect.Ptr: // source is a pointer and target is not a pointer
		fieldToField(sourceValue.Elem(), targetValue, args, wg)
	case sourceValue.Kind() != reflect.Ptr && targetValue.Kind() == reflect.Ptr: // source is not a pointer and target is a pointer
		nv := reflect.New(targetValue.Type().Elem())
		targetValue.Set(nv)
		fieldToField(sourceValue, targetValue.Elem(), args, wg)
	default: // both are pointers
		nv := reflect.New(targetValue.Type().Elem())
		targetValue.Set(nv)
		fieldToField(sourceValue.Elem(), targetValue.Elem(), args, wg)
	}
}

// mappingDirectMapping is used to map direct types
func mappingDirectMapping(s, t reflect.Value) {
	switch {
	case s.CanInterface() && t.CanInterface():
		t.Set(s)
	case !s.CanInterface() && s.CanAddr() && !t.CanInterface() && t.CanAddr():
		s := getUnexportedField(s)
		setUnexportedField(t, s)
	case s.CanInterface() && !t.CanInterface() && t.CanAddr():
		setUnexportedField(t, s)
	case !s.CanInterface() && s.CanAddr() && t.CanInterface():
		s := getUnexportedField(s)
		t.Set(s)
	default:
		log.Printf(
			"WARNING: Operations on MAP type fields that are not exported are not supported. Operation ignored. Source = %s, Target = %s\n",
			s, t,
		)
	}
}

func getUnexportedField(field reflect.Value) reflect.Value {
	return reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).Elem()
}

func setUnexportedField(field, value reflect.Value) {
	reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).
		Elem().
		Set(value)
}
