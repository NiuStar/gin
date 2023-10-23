package markdown

import (
	"unicode"

	//"fmt"
	//"github.com/sirupsen/logrus"
	"reflect"
	"strings"
)

func makeData(_type reflect.Type) (params []string) {
	params = make([]string, 0)
	params = append(params, "data")

	if _type.Kind() == reflect.Struct && _type.NumField() == 1 {
		params = append(params, _type.Field(0).Name)
		params = append(params, _type.Field(0).Tag.Get("comment"))
	} else {
		params = append(params, _type.String())
		params = append(params, "")
	}
	return params
}

func ParseParms(_type reflect.Type, index int, resp bool) (params [][]string) {

	if _type == nil {
		return nil
	}
	kind := _type.Kind()
	if resp && index == 0 && kind == reflect.Interface {
		return nil
	}
	_typeString := strings.TrimPrefix(_type.String(), "*")

	if index > 4 {
		return
		//	fmt.Println("x")
	}
	prefix := ""
	for i := 0; i < index; i++ {
		prefix += "--"
	}
	for _type.Kind() == reflect.Ptr {
		_type = _type.Elem()
	}

	if index == 1 && resp {
		param := makeData(_type)
		params = append(params, param)
	}
	if _type.Kind() == reflect.Struct {
		exported, _ := countFields(_type)

		if index == 1 && resp && exported == 0 {
			return nil
		} else if index == 1 && resp && exported == 1 {
			params2 := getFieldParams(_type, index, resp, 0, _typeString)
			params = append(params, params2...)
			return params
		}
		for i := 0; i < _type.NumField(); i++ {
			f := _type.Field(i)
			if unicode.IsLower(rune(f.Name[0])) {
				continue
			}
			param := getParam(prefix, _type.Field(i), resp)
			if param != nil {
				params = append(params, param)
			}
			params2 := getFieldParams(_type, index, resp, i, _typeString)
			params = append(params, params2...)
		}
	}
	//fmt.Println("_type",_type,_type.String())

	return params
}

func countFields(t reflect.Type) (exported, unexported int) {
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if unicode.IsUpper(rune(f.Name[0])) {
			exported++
		} else if unicode.IsLower(rune(f.Name[0])) {
			unexported++
		}
	}
	return
}
func getFieldParams(_type reflect.Type, index int, resp bool, i int, _typeString string) (params [][]string) {
	if _type.Field(i).Type.Kind() == reflect.Struct {
		if strings.TrimPrefix(_type.Field(i).Type.String(), "*") == _typeString {
			return
		}
		params2 := ParseParms(_type.Field(i).Type, index+1, resp)
		params = append(params, params2...)
	} else if _type.Field(i).Type.Kind() == reflect.Ptr {
		if strings.TrimPrefix(_type.Field(i).Type.Elem().String(), "*") == _typeString {
			return
		}
		params2 := ParseParms(_type.Field(i).Type.Elem(), index+1, resp)
		params = append(params, params2...)
	} else if _type.Field(i).Type.Kind() == reflect.Slice {
		//fmt.Println(_type.Field(i))
		if strings.TrimPrefix(_type.Field(i).Type.Elem().String(), "*") == _typeString {
			return
		}
		_typeField := _type.Field(i).Type.Elem()
		for _typeField.Kind() == reflect.Ptr {
			_typeField = _typeField.Elem()
		}
		if _typeField.Kind() == reflect.Struct {
			params2 := ParseParms(_typeField, index+1, resp)
			params = append(params, params2...)
		}
		//fmt.Println("_type2", _type.Field(i).Type.Elem())
	}
	return
}
func getParam(prefix string, field reflect.StructField, resp bool) (param []string) {
	if prefix == "" && resp && field.Type.Kind() == reflect.Interface {
		return nil
	}
	jsonName := field.Tag.Get("json")
	if jsonName == "" {
		jsonName = field.Name
	}
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
