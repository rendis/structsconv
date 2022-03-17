package structsconv

import (
	"log"
	"reflect"
	"sync"
)

// mapperRulesRegistry mapping rules register.
var mapperRulesRegistry = make(mapperRules)

// RegisterRulesDefinitions registers a mapper rules.
func RegisterRulesDefinitions(definitions ...RulesDefinition) {
	for _, d := range definitions {
		registerRules(d.Source, d.Target, d.Rules)
	}
}

// registerRules verifies and registers a mapper rules for specific mapping from structure to structure.
func registerRules(source interface{}, target interface{}, rules RulesSet) {
	key := buildKey(source, target)
	_, exists := mapperRulesRegistry[key]
	if exists {
		log.Fatalf("ERROR: Mapper with rulesKey (%s -> %s) already exists.", key.Source, key.Target)
	}
	checkMapperRules(key, rules)
	mapperRulesRegistry[key] = rules
}

// Map maps a source struct to a target struct.
func Map(source interface{}, target interface{}, args ...interface{}) {
	var wg sync.WaitGroup
	sourceV := reflect.ValueOf(source)
	targetV := reflect.ValueOf(target)
	if err := checkType(sourceV, targetV); err != nil { // check if the source and target are pointers
		log.Fatalf("ERROR: %s", err)
	}
	args = append([]interface{}{source}, args...)
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
	rules := mapperRulesRegistry[key]
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
			if pType := fieldToField(sourceValue, targetValue, args, wg); pType == IncompatibleTypes {
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
		targetValue.Set(mapperValue.Call([]reflect.Value{})[0])
	} else {
		params := getMethodParams(method, args, actualS)
		targetValue.Set(mapperValue.Call(params)[0])
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
	case StructsMapping:
		wg.Add(1)
		go cMappingStructLogic(sourceValue, targetValue, args, wg)
	case SlicesMapping:
		wg.Add(1)
		go cMappingSliceLogic(sourceValue, targetValue, args, wg)
	case MapsMapping:
		wg.Add(1)
		go cMappingMapLogic(sourceValue, targetValue, args, wg)
	case ArraysMapping:
		wg.Add(1)
		go cMappingArrayLogic(sourceValue, targetValue, args, wg)
	case DirectMapping:
		targetValue.Set(sourceValue)
	default:
		return IncompatibleTypes
	}
	return mappingType
}

// cMappingStructLogic is used to be called as a goroutine and map the source field to the destination field in a concurrent way
func cMappingStructLogic(source, target reflect.Value, args groupedArgs, wg *sync.WaitGroup) {
	defer wg.Done()
	structToStruct(source, target, source.Interface(), args, wg)
}

// cMappingMapLogic is used to be called as a goroutine and maps the structures of the source map to the destination map in a CONCURRENT way.
func cMappingMapLogic(sourceValue, targetValue reflect.Value, args groupedArgs, wg *sync.WaitGroup) {
	defer wg.Done()
	var mapWg sync.WaitGroup
	defer mapWg.Wait()
	itemType := targetValue.Type().Elem()
	var lock sync.Mutex
	targetValue.Set(reflect.MakeMap(targetValue.Type()))
	for _, key := range sourceValue.MapKeys() {
		mapWg.Add(1)
		go func(key reflect.Value) {
			defer mapWg.Done()
			var mapItemWg sync.WaitGroup
			item := reflect.New(itemType)
			sourceItem := sourceValue.MapIndex(key)
			structToStruct(sourceItem, item.Elem(), sourceItem.Interface(), args, &mapItemWg)
			mapItemWg.Wait()
			lock.Lock()
			targetValue.SetMapIndex(key, item.Elem())
			lock.Unlock()
		}(key)
	}
}

// cMappingArrayLogic is used to be called as a goroutine and maps the structures of the source array to the destination array in a CONCURRENT way.
func cMappingArrayLogic(sourceValue, targetValue reflect.Value, args groupedArgs, wg *sync.WaitGroup) {
	defer wg.Done()
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
			structToStruct(sourceItem, item.Elem(), sourceItem.Interface(), args, &arrayItemWg)
			arrayItemWg.Wait()
			lock.Lock()
			targetValue.Index(i).Set(item.Elem())
			lock.Unlock()
		}(i)
	}
}

// cMappingSliceLogic is used to be called as a goroutine and maps the structures of the source slice to the destination slice in a CONCURRENT way.
func cMappingSliceLogic(sourceValue, targetValue reflect.Value, args groupedArgs, wg *sync.WaitGroup) {
	defer wg.Done()
	var sliceWg sync.WaitGroup
	defer sliceWg.Wait()
	itemType := targetValue.Type().Elem()
	var lock sync.Mutex
	for i := 0; i < sourceValue.Len(); i++ {
		sliceWg.Add(1)
		go func(pos int) {
			defer sliceWg.Done()
			var sliceItemWg sync.WaitGroup
			item := reflect.New(itemType)
			sourceItem := sourceValue.Index(pos)
			structToStruct(sourceItem, item.Elem(), sourceItem.Interface(), args, &sliceItemWg)
			sliceItemWg.Wait()
			lock.Lock()
			targetValue.Set(reflect.Append(targetValue, item.Elem()))
			lock.Unlock()
		}(i)
	}
}
