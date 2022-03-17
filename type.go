package structsconv

import (
	"reflect"
)

// groupedArgs groups the arguments map by their type.
type groupedArgs map[reflect.Type][]interface{}

// RulesDefinition defines the source, target and rules for mapping.
type RulesDefinition struct {
	Source interface{}
	Target interface{}
	Rules  RulesSet
}

// rulesKey identifies the rules for specific mapping from structure to structure.
type rulesKey struct {
	Source reflect.Type
	Target reflect.Type
}

// RulesSet is a set of rules for mapping 2 specific structs, where the rulesKey is the name of the target field.
// and value is the rule for the mapping. Value can be:
// - a string, which is the name of the source field
// - a function, which will be called to get the target value
type RulesSet map[string]interface{}

// mapperRules contains all the rules registered for the mapping, where rulesKey identifies the rules for a specific mapping.
type mapperRules map[rulesKey]RulesSet

// processingResultType types of the processing result
// - StructsMapping 		 (0): the structs are processed using the mapping
// - SlicesMapping  		 (1): the slices are processed by mapping each element
// - ArraysMapping  		 (2): the arrays are processed by mapping each element
// - MapsMapping    		 (3): the maps are processed using the mapping
// - DirectMapping         	 (4): the values are processed directly
// - IncompatibleTypes 		 (5): the mapping is ignored
type processingResultType int

const (
	StructsMapping processingResultType = iota
	SlicesMapping
	ArraysMapping
	MapsMapping
	DirectMapping
	IncompatibleTypes
)

// buildKey builds a rulesKey from the source and target types.
func buildKey(source, target interface{}) rulesKey {
	return rulesKey{
		Source: reflect.TypeOf(source),
		Target: reflect.TypeOf(target),
	}
}

// getMappingType returns the mapping type for the given values.
func getMappingType(sourceValue, targetValue reflect.Value) processingResultType {
	switch {
	// *S* -> *S*
	case targetValue.Type().AssignableTo(sourceValue.Type()):
		return DirectMapping
	// {S} -> {N}
	case targetValue.Kind() == reflect.Struct && sourceValue.Kind() == reflect.Struct:
		return StructsMapping
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
		return IncompatibleTypes
	}
}

// getMapsMappingType returns the processing type for the given maps.
func getMapsMappingType(sourceValue, targetValue reflect.Value) processingResultType {
	// map(K)[{S}] -> map(K)[{N}]: Same keys types and the values are different structs
	if targetValue.Type().Key() == sourceValue.Type().Key() &&
		sourceValue.Type().Elem().Kind() == reflect.Struct &&
		targetValue.Type().Elem().Kind() == reflect.Struct {
		return MapsMapping
	}

	// map(KT)[]  ->  map(KW)[]: Different keys types
	// map(K)[{}] ->  map(K)[S]
	// map(K)[S]  ->  map(K)[{}]
	// map(K)[S]  ->  map(K)[N]
	return IncompatibleTypes
}

// getSlicesMappingType returns the processing type for the given slices.
func getSlicesMappingType(sourceValue, targetValue reflect.Value) processingResultType {
	// [{}] -> [{}]
	if sourceValue.Type().Elem().Kind() == reflect.Struct && targetValue.Type().Elem().Kind() == reflect.Struct {
		return SlicesMapping
	}
	// [s...] -> [n...]
	// [s...] -> [{}...]
	// [{}...] -> [s...]
	return IncompatibleTypes
}

// getArraysMappingType returns the processing type for the given arrays.
func getArraysMappingType(sourceValue, targetValue reflect.Value) processingResultType {
	// [{}] -> [{}]
	if sourceValue.Type().Elem().Kind() == reflect.Struct && targetValue.Type().Elem().Kind() == reflect.Struct {
		return ArraysMapping
	}
	// [s...] -> [n...]
	// [s...] -> [{}...]
	// [{}...] -> [s...]
	return IncompatibleTypes
}
