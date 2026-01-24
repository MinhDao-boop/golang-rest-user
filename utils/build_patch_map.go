package utils

import "reflect"

func BuildPatchMap(req any) map[string]interface{} {
	updates := make(map[string]interface{})

	v := reflect.ValueOf(req)
	t := reflect.TypeOf(req)

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		tag := t.Field(i).Tag.Get("patch")

		if tag != "true" || field.IsNil() {
			continue
		}

		jsonKey := t.Field(i).Tag.Get("json")
		updates[jsonKey] = field.Elem().Interface()
	}

	return updates
}
