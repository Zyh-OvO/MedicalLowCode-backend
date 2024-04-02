package exportCode

import (
	"MedicalLowCode-backend/util"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

type Layer interface {
	IsLayer()
}

func RawData2Layer(layer Layer, rawData map[string]any) Layer {
	tmpJson, err := json.Marshal(rawData)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(tmpJson, layer)
	if err != nil {
		panic(err)
	}
	return layer
}

func Layer2Code(layer Layer) string {
	layerValue := reflect.ValueOf(layer).Elem()
	layerType := layerValue.Type()
	var code string
	for i := 0; i < layerValue.NumField(); i++ {
		field := layerValue.Field(i)
		switch field.Kind() {
		case reflect.Ptr:
			if field.IsNil() {
				continue
			}
			field = field.Elem()
		case reflect.Slice:
			if field.IsNil() {
				continue
			}
		case reflect.Interface:
			//对于any类型，不能用IsNil判断，只能用IsZero
			if field.IsZero() {
				continue
			}
			field = field.Elem()
		}
		//if field.Kind() == reflect.Ptr {
		//	field = field.Elem()
		//	value := field.Interface()
		//	code += util.CamelCaseToSnakeCase(layerType.Field(i).Name) + "=" + fmt.Sprintf("%v", value) + ", "
		//} else
		if field.Kind() == reflect.Slice {
			//转换为tuple
			tuple := "("
			for j := 0; j < field.Len(); j++ {
				tuple += fmt.Sprintf("%v", field.Index(j).Interface()) + ", "
			}
			tuple += ")"
			code += util.CamelCaseToSnakeCase(layerType.Field(i).Name) + "=" + tuple + ", "
		} else if field.Kind() == reflect.String {
			value := field.Interface()
			code += util.CamelCaseToSnakeCase(layerType.Field(i).Name) + "=" + "\"" + fmt.Sprintf("%v", value) + "\"" + ", "
		} else if field.Kind() == reflect.Bool {
			value := field.Interface()
			code += util.CamelCaseToSnakeCase(layerType.Field(i).Name) + "=" + strings.Title(fmt.Sprintf("%v", value)) + ", "
		} else {
			value := field.Interface()
			code += util.CamelCaseToSnakeCase(layerType.Field(i).Name) + "=" + fmt.Sprintf("%v", value) + ", "
		}
	}
	return code
}
