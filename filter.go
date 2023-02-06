package null

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type (
	filterOpts struct {
		tag string
	}

	filterOpt func(f *filterOpts)
)

var (
	defaultFilterOpts = filterOpts{
		tag: "json",
	}
)

// UseTag
func UseTag(tag string) filterOpt {
	if tag == "" {
		return func(f *filterOpts) {}
	}

	return func(f *filterOpts) {
		f.tag = tag
	}
}

// FilterStruct filters the given structure from unset nullable fields
func FilterStruct(s any, opts ...filterOpt) (map[string]any, error) {
	if s == nil {
		return nil, errors.New("input cannot be nil")
	}

	rt := reflect.TypeOf(s)
	if rt.Kind() != reflect.Struct {
		return nil, fmt.Errorf("invalid type %T. input must be a struct", s)
	}

	// set options
	fOpts := defaultFilterOpts
	for _, opt := range opts {
		opt(&fOpts)
	}

	retMap := filterStruct(fOpts.tag, s)
	return retMap, nil
}

// FilterMap filters the given map from unset nullable fields
func FilterMap(m map[string]any) (map[string]any, error) {
	if m == nil {
		return nil, errors.New("input cannot be nil")
	}

	return filterMap(m), nil
}

// filterMap filters a map from unset nullable variables.
// If keepOtherFields is true, then every other field that is not a nullable type will keep intact
func filterMap(m map[string]any) map[string]any {
	retMap := make(map[string]any)

	for k, v := range m {
		switch val := v.(type) {
		case nullVar:
			if !val.isSet() {
				continue
			}
			retMap[k] = val.getVal()
		case map[string]any:
			mm := filterMap(val)
			if len(mm) == 0 {
				continue
			}
			retMap[k] = mm
		default:
			retMap[k] = v
		}
	}

	return retMap
}

// 1. loop through struct fields
// 2. check each field
// a. unexported -> continue
// b. doesn't have the necessary tag -> continue
// c. struct and implements Filterable -> filterStruct
// d. struct and doesn't implement filterable -> use as is
// e. map[string]any -> filterMap
// f. anonymous
// 	i. struct -> filterStruct
// 	ii. map[string]any -> filterMap

// structFieldsToMap creates a map from the given struct via the assigned tags.
func filterStruct(tag string, s any) map[string]any {
	retMap := make(map[string]any)

	if s == nil {
		return retMap
	}

	val := reflect.ValueOf(s)
	for i := 0; i < val.NumField(); i++ {

		// skip unexported fields
		if !val.Field(i).CanInterface() {
			continue
		}

		structField := val.Type().Field(i)       // the struct field itself
		fieldKind := structField.Type.Kind()     // its kind
		fieldValue := val.Field(i).Interface()   // its value as an interface
		fieldIsEmbedded := structField.Anonymous // if its embedded
		fieldName := ""                          // default name

		// skip field if it doesn't have the necessary tag
		// but only if it's not embedded/promoted field
		fTag, tagOk := structField.Tag.Lookup(tag)
		if !tagOk && !fieldIsEmbedded {
			continue
		}
		tagOpts := strings.Split(fTag, ",")
		if len(tagOpts) > 0 {
			fieldName = tagOpts[0]
		}

		// skip the field if:
		// 	- has "-" as field name
		// 	- has no fieldname and is not embedded
		// 	- has no fieldname, is embedded and not a struct
		if fieldName == "-" ||
			(fieldName == "" && !fieldIsEmbedded) ||
			(fieldName == "" && fieldIsEmbedded && fieldKind != reflect.Struct) {
			continue
		}

		switch fieldKind {
		case reflect.Struct:
			// check if implements Filterable
			_, iOk := fieldValue.(Filterable)
			if iOk {
				fs := filterStruct(tag, fieldValue)

				// if embedded then the fields need to be on the same level as others
				if fieldIsEmbedded && fieldName == "" {
					for k, v := range fs {
						if _, rOk := retMap[k]; !rOk {
							retMap[k] = v
						}
					}
				} else {
					// else put it on the given key
					if len(fs) == 0 {
						continue
					}
					retMap[fieldName] = fs
				}
				continue
			}

			nv, iOk := fieldValue.(nullVar)
			if iOk {
				if nv.isSet() {
					retMap[fieldName] = nv.getVal()
				}
				continue
			}

			// if it doesn't implement Filterable just use it's value,
			// but only if it has a valid tag name
			if fieldName != "" {
				retMap[fieldName] = fieldValue
			}

		case reflect.Map:

			switch t := fieldValue.(type) {
			case map[string]any:
				fm := filterMap(t)
				if len(fm) == 0 {
					continue
				}
				retMap[fieldName] = fm
			default:
				retMap[fieldName] = fieldValue
			}
		default:
			retMap[fieldName] = fieldValue
		}
	}

	return retMap
}
