package convert

import (
	"fmt"
	"reflect"
	"strconv"
)

// InterfaceToStringMap converts an incomming struct or map[string]interface{} type
// to a map[string]string containing the fields of struct or map as keys
// and values of those as string converted values of those.
func InterfaceToStringMap(in interface{}) (map[string]string, error) {
	//interpret and test given (struct) type with reflect
	//fmt.Printf("Row: %v\n", row)
	unnested := reflect.ValueOf(in)
	var out map[string]string

	switch unnested.Kind() {
	case reflect.Struct:
		// allocate memory slice for contents of given struct
		out = make(map[string]string)
		var numberOfStructFields = unnested.NumField()
		var value reflect.Value
		var fieldName string

		// get names of every field in struct and their values alternating
		// and sample them for redis query in slice
		for i := 0; i < numberOfStructFields; i++ {
			fieldName = unnested.Type().Field(i).Name
			value = unnested.Field(i)

			//fmt.Printf("STRUCT unnested %d.:\tField: %s\n\t\t\tValue: %v\n",
			//	i, fieldName, value)

			// get value of field, write in string format (type dependent)
			// and organize match with field in string map
			switch value.Kind() {
			case reflect.Float32, reflect.Float64:
				out[fieldName] = strconv.FormatFloat(value.Float(), 'f', 10, 64)

			case reflect.Int, reflect.Int8, reflect.Int16,
				reflect.Int32, reflect.Int64:
				out[fieldName] = strconv.FormatInt(value.Int(), 10)

			case reflect.Uint, reflect.Uint8, reflect.Uint16,
				reflect.Uint32, reflect.Uint64:

				out[fieldName] = strconv.FormatUint(value.Uint(), 10)

			case reflect.String:
				out[fieldName] = value.String()

			default:
				return nil, fmt.Errorf("Wrong type in struct: %v", value.Kind())
			}
		}

	case reflect.Map:
		// allocate memory slice map for contents of given interface map
		helpMap := in.(map[string]interface{})
		out = make(map[string]string)

		for key, value := range helpMap {
			//fmt.Printf("MAP unnested:\tField: %s\n\t\t\tValue: %v\n", key, value)
			out[key] = fmt.Sprintf("%v", value)
		}

	default:
		return nil, fmt.Errorf(
			"Wrong input type for conversion to map[string]string. "+
				"Only structures and map[string]interface{} are accepted. "+
				"Given type: %v\n", unnested.Kind())
	}

	return out, nil
}
