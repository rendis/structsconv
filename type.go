package structsconv

import (
	"reflect"
)

// RulesDefinition defines the source, target and rules for mapping.
type RulesDefinition struct {
	Source interface{}
	Target interface{}
	Rules  RulesSet
}

// RulesSet is a set of rules for mapping 2 specific structs, where the rulesKey is the name of the target field.
// and value is the rule for the mapping. Value can be:
// - a string, which is the name of the source field
// - a function, which will be called to get the target value
type RulesSet map[string]interface{}

// groupedArgs groups the arguments map by their type.
type groupedArgs map[reflect.Type][]interface{}

// rulesKey identifies the rules for specific mapping from structure to structure.
type rulesKey struct {
	source reflect.Type
	target reflect.Type
}

// mapperRulesRegistry contains all the rules registered for the mapping, where rulesKey identifies the rules for a specific mapping.
type mapperRulesRegistry map[rulesKey]RulesSet

// processingResultType types of the processing result
// - structsMapping 		 (0): the structs are processed using the mapping
// - slicesMapping  		 (1): the slices are processed by mapping each element
// - arraysMapping  		 (2): the arrays are processed by mapping each element
// - mapsMapping    		 (3): the maps are processed using the mapping
// - directMapping         	 (4): the values are mapping directly
// - incompatibleTypes 		 (5): incompatible types, so the mapping will be ignored
type processingResultType int

const (
	structsMapping processingResultType = iota
	slicesMapping
	arraysMapping
	mapsMapping
	directMapping
	incompatibleTypes
)

// buildKey builds a rulesKey from the source and target types.
func buildKey(source, target interface{}) rulesKey {
	return rulesKey{
		source: reflect.TypeOf(source),
		target: reflect.TypeOf(target),
	}
}

// getMappingType returns the mapping type for the given values.
func getMappingType(sourceValue, targetValue reflect.Value) processingResultType {
	switch {
	// *S* -> *S*
	case targetValue.Type().AssignableTo(sourceValue.Type()):
		return directMapping
	// {S} -> {N}
	case targetValue.Kind() == reflect.Struct && sourceValue.Kind() == reflect.Struct:
		return structsMapping
	// [] -> []
	case targetValue.Kind() == reflect.Slice && sourceValue.Kind() == reflect.Slice:
		return getSlicesMappingType(sourceValue, targetValue)
	// array -> array
	case targetValue.Kind() == reflect.Array && sourceValue.Kind() == reflect.Array:
		return getArraysMappingType(sourceValue, targetValue)
	// map -> map
	case targetValue.Kind() == reflect.Map && sourceValue.Kind() == reflect.Map:
		return getMapsMappingType(sourceValue, targetValue)
	// S -> N
	default:
		return incompatibleTypes
	}
}

// getMapsMappingType returns the processing type for the given maps.
func getMapsMappingType(sourceValue, targetValue reflect.Value) processingResultType {
	// map(K)[{S}] -> map(K)[{N}]: Same keys types and the values are different structs
	if targetValue.Type().Key() == sourceValue.Type().Key() &&
		sourceValue.Type().Elem().Kind() == reflect.Struct &&
		targetValue.Type().Elem().Kind() == reflect.Struct {
		return mapsMapping
	}

	// map(KT)[]  ->  map(KW)[]: Different keys types
	// map(K)[{}] ->  map(K)[S]
	// map(K)[S]  ->  map(K)[{}]
	// map(K)[S]  ->  map(K)[N]
	return incompatibleTypes
}

// getSlicesMappingType returns the processing type for the given slices.
func getSlicesMappingType(sourceValue, targetValue reflect.Value) processingResultType {
	// [{}] -> [{}]
	if sourceValue.Type().Elem().Kind() == reflect.Struct && targetValue.Type().Elem().Kind() == reflect.Struct {
		return slicesMapping
	}
	// [s...] -> [n...]
	// [s...] -> [{}...]
	// [{}...] -> [s...]
	return incompatibleTypes
}

// getArraysMappingType returns the processing type for the given arrays.
func getArraysMappingType(sourceValue, targetValue reflect.Value) processingResultType {
	// [{}] -> [{}]
	if sourceValue.Type().Elem().Kind() == reflect.Struct && targetValue.Type().Elem().Kind() == reflect.Struct {
		return arraysMapping
	}
	// [s...] -> [n...]
	// [s...] -> [{}...]
	// [{}...] -> [s...]
	return incompatibleTypes
}
