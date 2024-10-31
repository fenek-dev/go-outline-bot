package utils

import "reflect"

func StructToMap(data interface{}) map[string]interface{} {
	// 1. Create an empty map named result to store the fields and their values.
	result := make(map[string]interface{})

	val := reflect.ValueOf(data)

	// 3. If the input object is a pointer, dereference it to get the underlying value.
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	typ := val.Type()

	// 4. Iterate through the fields of the struct using a for loop.
	for i := 0; i < val.NumField(); i++ {
		// 5. For each field, get its name and kind (e.g., string, int, struct).
		fieldName := typ.Field(i).Name
		fieldValueKind := val.Field(i).Kind()
		var fieldValue interface{}

		// 6. If the field is a struct, recursively call structToMap to get the map representation of the nested struct.
		// Otherwise, get the field value directly.
		if fieldValueKind == reflect.Struct {
			fieldValue = StructToMap(val.Field(i).Interface())
		} else {
			fieldValue = val.Field(i).Interface()
		}

		// 7. Add the field name and value to the result map.
		result[fieldName] = fieldValue
	}

	return result
}