package jsonv

import (
	"encoding/json"
	"reflect"
)

type Struct map[string]interface{}

func Unmarshal(data []byte, v any, x *Struct) error {
	*x = make(Struct)
	err := json.Unmarshal(data, x)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}

func Marshal(v any, x Struct) ([]byte, error) {
	t := reflect.TypeOf(v)
	thing := reflect.ValueOf(v)

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		fieldName, ok := f.Tag.Lookup("json")
		if ok {
			x[fieldName] = thing.Field(i).String()
		}
	}

	return json.Marshal(x)
}
