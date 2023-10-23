package markdown

import (
	//"fmt"
	//"github.com/sirupsen/logrus"
	"reflect"
	"strings"
)

func ParseParms(_type reflect.Type, index int,resp bool) (params [][]string) {
	if index > 2 {
		return nil
	}
	prefix := ""
	for i := 0; i < index; i++ {
		prefix += "--"
	}
	for ;_type.Kind() == reflect.Ptr; {
		_type = _type.Elem()
	}
	if _type.Kind() == reflect.Struct {
		for i := 0; i < _type.NumField(); i++ {
			param := getParam(prefix, _type.Field(i),resp)
			if param != nil {
				params = append(params, param)
			}

			if _type.Field(i).Type.Kind() == reflect.Struct {
				params2 := ParseParms(_type.Field(i).Type, index+1,resp)
				params = append(params, params2...)
			} else if _type.Field(i).Type.Kind() == reflect.Ptr {
				params2 := ParseParms(_type.Field(i).Type.Elem(), index+1,resp)
				params = append(params, params2...)
			} else if _type.Field(i).Type.Kind() == reflect.Slice {
				//fmt.Println(_type.Field(i))
				_typeField := _type.Field(i).Type.Elem()
				for _typeField.Kind() == reflect.Ptr {
					_typeField = _typeField.Elem()
				}
				if _typeField.Kind() == reflect.Struct {
					params2 := ParseParms(_typeField, index+1,resp)
					params = append(params, params2...)
				}
				//fmt.Println("_type2", _type.Field(i).Type.Elem())
			}
		}
	}
	//fmt.Println("_type",_type,_type.String())

	return params
}

func getParam(prefix string, field reflect.StructField,resp bool) (param []string) {
	jsonName := field.Tag.Get("json")
	if jsonName == "" || jsonName == "-" {
		return
	}
	if strings.Contains(jsonName, ",") {
		jsonName = strings.Split(jsonName, ",")[0]
	}
	//fmt.Println("jsonName:",jsonName)
	param = append(param, prefix+jsonName)
	param = append(param, field.Type.String())
	param = append(param, field.Tag.Get("comment"))
	if !resp {

		optional := field.Tag.Get("optional")
		if optional == "false" {
			param = append(param, "是")
		} else {
			param = append(param, "否")
		}
	}

	return
}
func GetReflectType(value interface{}) reflect.Type {
	v := reflect.TypeOf(value)
	for reflect.Ptr == v.Kind() {
		//fmt.Println("v = v.Elem():",v.String())
		v = v.Elem()
	}
	return v
}

func GetReflectValue(value interface{}) reflect.Value {
	v := reflect.ValueOf(value)

	for reflect.Ptr == v.Kind() || reflect.Interface == v.Kind() {
		v = v.Elem()
	}
	return v
}

func writeResult(result reflect.Value, preStr string) (params [][]string) {

	if !result.IsValid() {
		return
	}
	_type := result.Type()
	_value := result

	if _type.Kind() != reflect.Struct {
		return
	}
	for i := 0; i < _type.NumField(); i++ {
		jsonName := _type.Field(i).Name
		if []byte(jsonName)[0] < 'A' || []byte(jsonName)[0] > 'Z' {
			continue
		}
		if _type.Field(i).Tag.Get("json") != "" {
			jsonName = _type.Field(i).Tag.Get("json")
		}

		if jsonName == _value.Field(i).Type().Name() {
			continue
		}

		var param []string
		param = append(param, preStr+jsonName)
		param = append(param, _value.Field(i).Type().Name())
		param = append(param, _type.Field(i).Tag.Get("comment"))
		params = append(params, param)

		parType := _type.Field(i).Type
		for parType.Kind() == reflect.Ptr {

			parType = parType.Elem()
		}

		if parType.Kind() == reflect.Struct {
			params = append(params, writeResult(_value.Field(i), preStr+"&nbsp;&nbsp;&nbsp;&nbsp;")...)
		} else if parType.Kind() == reflect.Slice {
			//fmt.Println("result type: reflect.Slice",_value.Field(i).CanSet())
			//_value.Field(i).SetLen(1)
			//v.SetCap(2)
			//fmt.Println("result type: reflect.Slice",_value.Field(i).Slice(0,0))
			/*for j := 0; j < _value.Field(i).NumField(); j++ {
				fmt.Printf("result type: Field %d: %v\n", i, _value.Field(i).Field(j).Type())
			}*/

			//	fmt.Println("result type: reflect.Slice")
			//params = append(params,writeResult(_value.Field(i),preStr + "&nbsp;&nbsp;&nbsp;&nbsp;")...)
		}
	}
	return params
}
