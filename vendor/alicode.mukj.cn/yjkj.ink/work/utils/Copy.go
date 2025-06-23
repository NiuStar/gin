package utils

import (
	time2 "alicode.mukj.cn/yjkj.ink/work/utils/time.v2"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

/*
func CopyTo(src, dst interface{}) error {
	srcValue := reflect.Indirect(reflect.ValueOf(src))
	dstValue := reflect.ValueOf(dst)
	if dstValue.Kind() != reflect.Ptr {
		return errors.New("dst must on ptr Value")
	}
	dstValue = dstValue.Elem()
	srcType := srcValue.Type()
	for i := 0; i < srcValue.NumField(); i++ {
		field := getFiledByName(dstValue, srcValue.Field(i).Type(), srcType.Field(i).Name, srcType.Field(i).Tag.Get("json"))
		if field.IsValid() {
			field.Set(srcValue.Field(i))
		}
	}
	return nil
}*/

func getFiledByName(pValue reflect.Value, fieldType reflect.Type, name, jsonName string) reflect.Value {
	fieldTypeOld := fieldType
	if fieldType.Kind() == reflect.Ptr {
		fieldType = fieldType.Elem()
	}
	for pValue.Kind() == reflect.Ptr {
		pValue = reflect.Indirect(pValue)
	}
	for i := 0; i < pValue.NumField(); i++ {
		name1 := pValue.Type().Field(i).Name
		jsonName1 := pValue.Type().Field(i).Tag.Get("json")

		pFieldType := pValue.Field(i).Type()
		if pFieldType.Kind() == reflect.Ptr {
			pFieldType = pFieldType.Elem()
		}

		fkind := fieldType.Kind()
		pkind := pFieldType.Kind()

		if strings.ToLower(name1) == strings.ToLower(name) || strings.ToLower(jsonName) == strings.ToLower(jsonName1) {

			if fieldType == pFieldType || fkind == pkind || fkind == reflect.Interface {
				return pValue.Field(i)
			} else if pkind == reflect.String && fkind == reflect.Struct {
				if _, ok := fieldTypeOld.MethodByName("String"); ok {
					return pValue.Field(i)
				} else {
					pkgPath := fieldType.PkgPath()
					fmt.Println("pkgPath", pkgPath)
					pFieldTypepkgPath := pFieldType.PkgPath()
					fmt.Println("pFieldType pkgPath", pFieldTypepkgPath)
					if pkgPath == "" && pFieldTypepkgPath == "alicode.mukj.cn/yjkj.ink/work/utils/time.v2" {
						return pValue.Field(i)
					} else if pFieldTypepkgPath == "" && pkgPath == "alicode.mukj.cn/yjkj.ink/work/utils/time.v2" {
						return pValue.Field(i)
					}
				}
			} else {
				pkgPath := fieldType.PkgPath()
				fmt.Println("pkgPath", pkgPath)
				pFieldTypepkgPath := pFieldType.PkgPath()
				fmt.Println("pFieldType pkgPath", pFieldTypepkgPath)
				if pkgPath == "" && pFieldTypepkgPath == "alicode.mukj.cn/yjkj.ink/work/utils/time.v2" {
					return pValue.Field(i)
				} else if pFieldTypepkgPath == "" && pkgPath == "alicode.mukj.cn/yjkj.ink/work/utils/time.v2" {
					return pValue.Field(i)
				}
			}
		}

	}
	return reflect.ValueOf(nil)
}

// CopyObject 从map或数组数据中copy到对象中
func CopyObject(obj interface{}, value interface{}) {
	v := reflect.ValueOf(value)
	dstValue := v
	k1 := v.Type().Kind()
	if k1 != reflect.Ptr {
		return
	}
	k2 := v.Type().Elem().Kind()
	srcIndex := 0
	dstIndex := 1
	for k2 == reflect.Ptr {
		dstIndex++
		dstValue = dstValue.Elem()
		k2 = dstValue.Type().Elem().Kind()
	}
	fmt.Println(k1, k2)
	v1 := getValue(obj, v.Type().Elem())
	k3 := v1.Kind()
	srcValue := v1
	for k3 == reflect.Ptr {
		srcIndex++
		srcValue = srcValue.Elem()
		k3 = srcValue.Kind()
	}
	fmt.Println("v1", v1.Interface(), v1.Type())
	if v.Kind() == reflect.Ptr {
		if dstIndex == srcIndex {
			v.Elem().Set(v1.Elem())
		} else {
			v.Elem().Set(v1)
		}
		/*if v1.Kind() == reflect.Ptr {
			v.Elem().Set(v1.Elem())
		} else {
			v.Elem().Set(v1)
		}*/
	} else {
		if v.IsValid() {
			v.Set(v1)
		}
	}
}

func getMapValueByKeys(mapObj map[string]interface{}, field reflect.StructField) interface{} {
	jsonName := field.Tag.Get("json")
	name := field.Name
	lowerJsonName := strings.ToLower(jsonName)
	lowerName := strings.ToLower(name)

	names := []string{jsonName, name, lowerJsonName, lowerName}
	for _, n := range names {
		if len(n) > 0 {
			mapObj2 := mapObj[n]
			if mapObj2 != nil {
				return mapObj2
			}
		}
	}
	return nil
}
func getValue(obj interface{}, rType reflect.Type) reflect.Value {
	if obj == nil {
		return reflect.ValueOf(nil)
	}
	kind := rType.Kind()
	if reflect.Ptr == kind {
		v := reflect.New(rType.Elem())
		v1 := getValue(obj, rType.Elem())
		fmt.Println("v", v.Type())
		fmt.Println("v1", v1.Type())
		if v1.Kind() == reflect.Ptr {
			v.Elem().Set(v1.Elem())
		} else if v.Kind() == reflect.Ptr {
			if v.Elem().Kind() == reflect.Int {
				if v1.Kind() == reflect.Float64 {
					v.Elem().SetInt(int64(v1.Float()))
				} else {
					fmt.Println("类型不一致：v1", v1.Interface(), "v .kind", v.Elem().Kind())
				}
			} else if v.Elem().Kind() == reflect.Bool {
				if v1.Kind() == reflect.String {
					if strings.ToLower(v1.String()) == "true" {
						v.Elem().SetBool(true)
					} else {
						v.Elem().SetBool(false)
					}
				} else if v1.Kind() == reflect.Float64 {
					if v1.Float() == 1 {
						v.Elem().SetBool(true)
					} else {
						v.Elem().SetBool(false)
					}
				} else if v1.Kind() == reflect.Bool {
					v.Elem().SetBool(v1.Bool())
				} else {
					fmt.Println("类型不一致：v1", v1.Interface(), "v .kind", v.Elem().Kind())
				}
			} else if v.Elem().Kind() == v1.Kind() {
				v.Elem().Set(v1)
			} else {
				v.Elem().Set(v1.Convert(v.Elem().Type()))
			}
		} else {
			v.Elem().Set(v1.Convert(v.Elem().Type()))
		}
		return v
	} else if kind != reflect.Struct && kind != reflect.Slice {
		return reflect.ValueOf(obj)
	}

	rType2 := rType

	if rType2.Kind() == reflect.Struct {
		v := reflect.New(rType2)
		mapObj := obj.(map[string]interface{})
		for k, v := range mapObj {
			mapObj[strings.ToLower(k)] = v
		}
		for i := 0; i < rType2.NumField(); i++ {

			mapObj2 := getMapValueByKeys(mapObj, rType2.Field(i))
			if mapObj2 == nil {
				continue
			}
			v1 := getValue(mapObj2, rType2.Field(i).Type)
			if v1.Kind() == v.Elem().Field(i).Kind() {
				v.Elem().Field(i).Set(v1.Convert(v.Elem().Field(i).Type()))
			} else {
				copyObject(v.Elem().Field(i), v1)
			}
		}
		return v
	} else if rType2.Kind() == reflect.Slice {
		valueValueTmp := reflect.MakeSlice(rType, 0, 0)
		for _, resultObj := range obj.([]interface{}) {
			childType := rType2.Elem()
			for childType.Kind() == reflect.Ptr {
				childType = childType.Elem()
			}
			vValue := getValue(resultObj, childType)
			valueValueTmp = reflect.Append(valueValueTmp, vValue)
		}
		return valueValueTmp
	} else if rType2.Kind() == reflect.String {
		return reflect.ValueOf(obj.(string))
	} else {
		if rType2.Kind() == reflect.Int {
			switch obj.(type) {
			case float64:
				{
					obj = int(obj.(float64))
				}
			case string:
				{
					obj, _ = strconv.Atoi(obj.(string))
				}
			}
		} else if rType2.Kind() == reflect.Int8 {
			switch obj.(type) {
			case float64:
				{
					obj = int8(obj.(float64))
				}
			case string:
				{
					obj, _ = strconv.Atoi(obj.(string))
					obj = int8(obj.(int))
				}
			}
		} else if rType2.Kind() == reflect.Int64 {
			switch obj.(type) {
			case float64:
				{
					obj = int64(obj.(float64))
				}
			case string:
				{
					obj, _ = strconv.Atoi(obj.(string))
					obj = int64(obj.(int))
				}
			}
		} else if rType2.Kind() == reflect.Bool {
			switch obj.(type) {
			case float64:
				{
					if obj.(float64) == 1 {
						obj = true
					} else {
						obj = false
					}
				}
			case string:
				{
					objlower := strings.ToLower(obj.(string))
					if objlower == "1" || objlower == "true" {
						obj = true
					} else {
						obj = false
					}
				}
			}
		}
		return reflect.ValueOf(obj)
	}
	return reflect.ValueOf(nil)
}

func copyValue(dest reflect.Value, src reflect.Value) {
	if dest.Kind() == src.Kind() {
		dest.Set(src)
	} else if dest.CanFloat() && src.CanFloat() {
		dest.SetFloat(src.Float())
	} else if dest.CanInt() && src.CanInt() {
		dest.SetInt(src.Int())
	} else if dest.CanAddr() && src.CanAddr() {
		copyValue(dest.Addr(), src.Addr())
	} else if dest.CanSet() && src.CanSet() {
		dest.Set(src)
	} else if dest.CanFloat() && src.CanInt() {
		dest.SetFloat(float64(src.Int()))
	} else if dest.CanInt() && src.CanFloat() {
		dest.SetInt(int64(src.Float()))
	} else if dest.Kind() == reflect.Bool && src.Kind() == reflect.String {
		if strings.ToLower(src.String()) == "true" {
			dest.SetBool(true)
		} else {
			dest.SetBool(false)
		}
	} else {
		fmt.Println("copyValue: can not copy", dest.Kind(), src.Kind())
	}
}
func copyFromMap2Obj(srcValue, dstValue reflect.Value) error {
	if dstValue.Elem().Kind() == reflect.Ptr && dstValue.Elem().IsZero() {
		dstValue.Elem().Set(reflect.New(dstValue.Type().Elem().Elem()))
	}
	for _, key := range srcValue.MapKeys() {
		if key.Kind() != reflect.String {
			continue
		}
		field := getFiledByName(dstValue, srcValue.MapIndex(key).Type(), key.String(), key.String())
		if field.IsValid() {
			if field.Type() == srcValue.MapIndex(key).Type() {
				field.Set(srcValue.MapIndex(key))
			} else if field.Kind() == reflect.Ptr {
				field.Set(reflect.New(field.Type().Elem()))
				if field.Elem().Kind() == reflect.Struct {
					dstInterface := reflect.New(field.Type()).Interface()
					CopyTo(srcValue.MapIndex(key).Interface(), dstInterface)
					field.Set(reflect.ValueOf(dstInterface).Elem())
				} else {
					if srcValue.MapIndex(key).Kind() == reflect.Interface {
						copyValue(field.Elem(), srcValue.MapIndex(key).Elem())
						//field.Elem().Set(srcValue.MapIndex(key).Elem())
					} else {
						copyValue(field.Elem(), srcValue.MapIndex(key))
						//field.Elem().Set(srcValue.MapIndex(key))
					}
				}

			} else if field.Kind() == reflect.Slice {
				rSlice := reflect.MakeSlice(field.Type(), 0, 0)

				srcSlice := srcValue.MapIndex(key)
				if srcSlice.Kind() == reflect.Interface {
					srcSlice = srcValue.MapIndex(key).Elem()
				}

				for j := 0; j < srcSlice.Len(); j++ {
					isStruct := false
					dstFieldType := field.Type().Elem()
					if dstFieldType.Kind() == reflect.Ptr {
						if dstFieldType.Elem().Kind() == reflect.Struct {
							isStruct = true
						}
					} else if dstFieldType.Kind() == reflect.Struct {
						isStruct = true
					}
					dstFieldValue := reflect.New(dstFieldType)
					dstFieldInterface := dstFieldValue.Interface()

					if isStruct {
						CopyTo(srcSlice.Index(j).Interface(), dstFieldInterface)
					} else {
						/*kind := dstFieldType.Kind()
						if kind == reflect.Ptr {
							dstFieldValue.Elem().Set(srcValue.Field(i).Index(j))
						} else {
							dstFieldValue.Set(reflect.ValueOf(srcValue.Field(i).Index(j).Interface()))
						}*/
						dstFieldInterface = srcSlice.Index(j).Interface() //dstFieldValue.Interface()
					}
					if dstFieldType.Kind() == reflect.Ptr && dstFieldType.Elem().Kind() == reflect.Struct {
						rSlice = reflect.Append(rSlice, reflect.ValueOf(dstFieldInterface).Elem())
					} else {
						rSlice = reflect.Append(rSlice, reflect.ValueOf(dstFieldInterface))
					}

				}
				field.Set(rSlice)
			} else if field.Kind() == reflect.Map {

				if !srcValue.MapIndex(key).IsNil() {
					//s := srcValue.Field(i).Convert(reflect.TypeOf(field.Type()))
					//field.Set(s)
					fieldValue := reflect.MakeMapWithSize(field.Type(), srcValue.MapIndex(key).Len())
					for _, key2 := range srcValue.MapIndex(key).MapKeys() {
						fieldValue.SetMapIndex(key2, reflect.ValueOf(srcValue.MapIndex(key).MapIndex(key).Interface()))
					}
					field.Set(fieldValue)
				}
			} else {
				if !srcValue.MapIndex(key).IsNil() {
					srcKind := srcValue.MapIndex(key).Elem().Kind()
					fieldKind := field.Kind()
					if srcKind == fieldKind {
						field.Set(srcValue.MapIndex(key).Elem())
					} else {
						copyObject(field, srcValue.MapIndex(key).Elem())
					}
				}
			}
		} else {
			//kind := srcValue.Field(i).Kind()

		}
	}

	return nil
}

func CopyObjectValue(srcValue, dstValue reflect.Value) error {
	srcValue = reflect.Indirect(srcValue)
	srcType := srcValue.Type()
	for i := 0; i < srcValue.NumField(); i++ {
		fieldName := srcType.Field(i).Name
		field := getFiledByName(dstValue, srcValue.Field(i).Type(), fieldName, srcType.Field(i).Tag.Get("json"))
		if field.IsValid() {
			if field.Type() == srcValue.Field(i).Type() {
				field.Set(srcValue.Field(i))
			} else if field.Kind() == reflect.Ptr {
				field.Set(reflect.New(field.Type().Elem()))
				if field.Elem().Kind() == reflect.Struct {
					dstInterface := reflect.New(field.Type()).Interface()
					CopyTo(srcValue.Field(i).Interface(), dstInterface)
					field.Set(reflect.ValueOf(dstInterface).Elem())
				} else {
					field.Elem().Set(srcValue.Field(i))
				}

			} else if field.Kind() == reflect.Slice {
				rSlice := reflect.MakeSlice(field.Type(), 0, 0)
				for j := 0; j < srcValue.Field(i).Len(); j++ {
					isStruct := false
					dstFieldType := field.Type().Elem()
					if dstFieldType.Kind() == reflect.Ptr {
						if dstFieldType.Elem().Kind() == reflect.Struct {
							isStruct = true
						}
					} else if dstFieldType.Kind() == reflect.Struct {
						isStruct = true
					}
					dstFieldValue := reflect.New(dstFieldType)
					dstFieldInterface := dstFieldValue.Interface()

					if isStruct {
						CopyTo(srcValue.Field(i).Index(j).Interface(), dstFieldInterface)
					} else {
						/*kind := dstFieldType.Kind()
						if kind == reflect.Ptr {
							dstFieldValue.Elem().Set(srcValue.Field(i).Index(j))
						} else {
							dstFieldValue.Set(reflect.ValueOf(srcValue.Field(i).Index(j).Interface()))
						}*/
						dstFieldInterface = srcValue.Field(i).Index(j).Interface() //dstFieldValue.Interface()
					}
					if dstFieldType.Kind() == reflect.Ptr && dstFieldType.Elem().Kind() == reflect.Struct {
						rSlice = reflect.Append(rSlice, reflect.ValueOf(dstFieldInterface).Elem())
					} else {
						rSlice = reflect.Append(rSlice, reflect.ValueOf(dstFieldInterface))
					}

				}
				field.Set(rSlice)
			} else if field.Kind() == reflect.Map {

				if !srcValue.Field(i).IsNil() {
					//s := srcValue.Field(i).Convert(reflect.TypeOf(field.Type()))
					//field.Set(s)
					fieldValue := reflect.MakeMapWithSize(field.Type(), srcValue.Field(i).Len())
					for _, key := range srcValue.Field(i).MapKeys() {
						fieldValue.SetMapIndex(key, reflect.ValueOf(srcValue.Field(i).MapIndex(key).Interface()))
					}
					field.Set(fieldValue)
				}
			} else {
				if !srcValue.Field(i).IsNil() {
					srcFieldValue := srcValue.Field(i)
					if srcFieldValue.Kind() == reflect.Ptr {
						srcFieldValue = srcFieldValue.Elem()
					}
					if field.Kind() == reflect.String && srcFieldValue.Kind() == reflect.Struct {
						methodValue := srcValue.Field(i).MethodByName("String")
						if methodValue.IsValid() {
							inNumber := methodValue.Type().NumIn()
							if inNumber == 0 {
								rs := methodValue.Call([]reflect.Value{})
								if len(rs) > 0 {
									field.Set(rs[0])
								}
							}
						}

					} else {
						field.Set(srcValue.Field(i).Elem())

					}
				}
			}
		} else {
			//kind := srcValue.Field(i).Kind()

		}
	}
	return nil
}
func CopyTo(src, dst interface{}) error {

	srcValue := reflect.Indirect(reflect.ValueOf(src))
	dstValue := reflect.ValueOf(dst)
	dstKind := dstValue.Kind()
	if dstKind != reflect.Ptr {
		return errors.New("dst must on ptr Value")
	}
	if dstValue.Elem().Kind() == reflect.Interface {
		dstValue = dstValue.Elem().Elem()
		dstKind = dstValue.Kind()

	}
	if dstValue.IsZero() {
		return errors.New("dst IsZero")
	}
	if srcValue.Kind() == reflect.Map {
		return copyFromMap2Obj(srcValue, dstValue)
	}
	dstValue = dstValue.Elem()
	if dstValue.Kind() == reflect.Ptr && dstValue.IsNil() {
		dstValue.Set(reflect.New(dstValue.Type().Elem()))
		dstValue = dstValue.Elem()
	}
	if dstValue.Kind() == reflect.Ptr || dstValue.Kind() == reflect.Interface {
		dstValue = dstValue.Elem()
	}
	if srcValue.Kind() == reflect.Ptr || srcValue.Kind() == reflect.Interface {
		srcValue = srcValue.Elem()
	}
	srcType := srcValue.Type()
	if srcType.Kind() == reflect.String {
		if dstValue.Type().PkgPath() == "alicode.mukj.cn/yjkj.ink/work/utils/time.v2" {
			if srcValue.String() == "" {
				return errors.New("srcValue time is empty")
			}
			tm, err := time.Parse("2006-01-02 15:04:05", srcValue.String())
			if err != nil {
				tm, err = time.Parse("2006-01-02", srcValue.String())
			}
			dstValue.Set(reflect.ValueOf(time2.Time{tm}))

		}
		return nil
	}

	for i := 0; i < srcValue.NumField(); i++ {
		srcTypeField := srcType.Field(i)

		// 检查字段是否可以导出
		if srcTypeField.PkgPath != "" {
			fmt.Printf("Field %s is unexported\n", srcTypeField.Name)
			continue
		}
		fieldName := srcType.Field(i).Name
		if strings.HasSuffix(fieldName, "Time") {
			fmt.Println("123")
		}
		field := getFiledByName(dstValue, srcValue.Field(i).Type(), fieldName, srcType.Field(i).Tag.Get("json"))
		if field.IsValid() {
			if field.Type() == srcValue.Field(i).Type() {
				field.Set(srcValue.Field(i))
			} else if field.Kind() == reflect.Ptr {
				field.Set(reflect.New(field.Type().Elem()))
				if field.Elem().Kind() == reflect.Struct {
					dstInterface := reflect.New(field.Type()).Interface()
					if err := CopyTo(srcValue.Field(i).Interface(), dstInterface); err == nil {
						field.Set(reflect.ValueOf(dstInterface).Elem())
					} else {
						field.Set(reflect.Zero(field.Type()))
					}

				} else {
					field.Elem().Set(srcValue.Field(i))
				}

			} else if field.Kind() == reflect.Slice {
				rSlice := reflect.MakeSlice(field.Type(), 0, 0)
				for j := 0; j < srcValue.Field(i).Len(); j++ {
					isStruct := false
					isPtr := false
					var dstFieldInterface interface{}
					dstFieldType := field.Type().Elem()
					if dstFieldType.Kind() == reflect.Ptr {
						if dstFieldType.Elem().Kind() == reflect.Struct {
							dstFieldType = dstFieldType.Elem()
							isPtr = true
							isStruct = true
						}
					} else if dstFieldType.Kind() == reflect.Struct {

						isStruct = true
					}
					dstFieldValue := reflect.New(dstFieldType)
					dstFieldInterface = dstFieldValue.Interface()
					if isStruct {
						if isPtr {
							CopyTo(srcValue.Field(i).Index(j).Interface(), &dstFieldInterface)
						} else {
							CopyTo(srcValue.Field(i).Index(j).Interface(), dstFieldInterface)
						}
					} else {
						/*kind := dstFieldType.Kind()
						if kind == reflect.Ptr {
							dstFieldValue.Elem().Set(srcValue.Field(i).Index(j))
						} else {
							dstFieldValue.Set(reflect.ValueOf(srcValue.Field(i).Index(j).Interface()))
						}*/
						dstFieldInterface = srcValue.Field(i).Index(j).Interface() //dstFieldValue.Interface()
					}
					if dstFieldType.Kind() == reflect.Ptr && dstFieldType.Elem().Kind() == reflect.Struct {
						rSlice = reflect.Append(rSlice, reflect.ValueOf(dstFieldInterface).Elem())
					} else {
						rSlice = reflect.Append(rSlice, reflect.ValueOf(dstFieldInterface))
					}

				}
				field.Set(rSlice)
			} else if field.Kind() == reflect.Map {

				if !srcValue.Field(i).IsNil() {
					//s := srcValue.Field(i).Convert(reflect.TypeOf(field.Type()))
					//field.Set(s)
					fieldValue := reflect.MakeMapWithSize(field.Type(), srcValue.Field(i).Len())
					for _, key := range srcValue.Field(i).MapKeys() {
						fieldValue.SetMapIndex(key, reflect.ValueOf(srcValue.Field(i).MapIndex(key).Interface()))
					}
					field.Set(fieldValue)
				}
			} else {
				if !srcValue.Field(i).IsNil() {
					srcFieldValue := srcValue.Field(i)
					if srcFieldValue.Kind() == reflect.Ptr {
						srcFieldValue = srcFieldValue.Elem()
					}
					if field.Kind() == reflect.String && srcFieldValue.Kind() == reflect.Struct {
						methodValue := srcValue.Field(i).MethodByName("String")
						if methodValue.IsValid() {
							inNumber := methodValue.Type().NumIn()
							if inNumber == 0 {
								rs := methodValue.Call([]reflect.Value{})
								if len(rs) > 0 {
									field.Set(rs[0])
								}
							}
						}

					} else {
						field.Set(srcValue.Field(i).Elem())

					}

				}
			}
		} else {
			//kind := srcValue.Field(i).Kind()

		}
	}
	return nil
}

func isType(kind reflect.Kind, t string) bool {
	return strings.HasPrefix(kind.String(), t)
}

func copyObject(dst, src reflect.Value) {
	dstKind := dst.Kind()
	srcKind := src.Kind()
	if srcKind == dstKind {
		dst.Set(src)
	} else if dstKind == reflect.String {
		dst.SetString(fmt.Sprintf("%v", src.Interface()))
	} else if isType(dstKind, "int") {
		if srcKind == reflect.String {
			i, err := strconv.Atoi(src.String())
			if err == nil {
				dst.SetInt(int64(i))
			}
		} else if isType(srcKind, "float") {
			dst.SetInt(int64(src.Float()))
		} else if isType(srcKind, "int") {
			dst.SetInt(src.Int())
		} else if isType(srcKind, "uint") {
			dst.SetInt(int64(src.Uint()))
		}
	} else if isType(dstKind, "float") {
		if srcKind == reflect.String {
			f, err := strconv.ParseFloat(src.String(), 64)
			if err == nil {
				dst.SetFloat(f)
			}
		} else if isType(srcKind, "float") {
			dst.SetFloat(src.Float())
		} else if isType(srcKind, "int") {
			dst.SetFloat(float64(src.Int()))
		} else if isType(srcKind, "uint") {
			dst.SetFloat(float64(src.Uint()))
		}
	} else if isType(dstKind, "bool") {
		if srcKind == reflect.String {
			if strings.ToLower(src.String()) == "true" || src.String() == "1" {
				dst.SetBool(true)
			} else {
				dst.SetBool(false)
			}
		} else if isType(srcKind, "float") {
			dst.SetBool(src.Float() == 1.0)
		} else if isType(srcKind, "int") {
			dst.SetBool(src.Int() == 1)
		} else if isType(srcKind, "uint") {
			dst.SetBool(src.Uint() == 1)
		}
	}

}

func CopyReflectValue(srcValue, dstValue reflect.Value) reflect.Value {
	srcValue = reflect.Indirect(srcValue)

	var dstValue2 reflect.Value
	if dstValue.IsZero() {
		if dstValue.Type().Kind() == reflect.Ptr {
			dstValue2 = reflect.New(dstValue.Type().Elem()).Elem()

		} else {
			dstValue2 = reflect.New(dstValue.Type()).Elem()

		}
	} else {
		dstValue2 = dstValue
	}
	if srcValue.IsZero() {
		dstValue.Set(dstValue2)
		return dstValue2
	}
	if !dstValue2.IsValid() {
		fmt.Println("123456")
	}
	if dstValue2.IsValid() && dstValue2.IsZero() {
		fmt.Println("111")
	}
	if srcValue.Type() == dstValue2.Type() {
		if dstValue.Type().Kind() == reflect.Ptr {
			dstValue.Elem().Set(srcValue)
			return dstValue

		} else {
			dstValue.Set(srcValue)
		}

		return dstValue
	}
	if srcValue.Kind() == reflect.Struct && dstValue2.Kind() == reflect.Struct {
		return CopyReflectStruct(srcValue, dstValue2)
	} else if srcValue.Kind() == reflect.Slice && dstValue2.Kind() == reflect.Slice {
		return CopyReflectSlice(srcValue, dstValue2)
	} else if srcValue.Kind() == reflect.Slice && dstValue2.Kind() == reflect.Struct {
		if srcValue.Len() > 0 {
			return CopyReflectStruct(srcValue.Index(0), dstValue2)
		}
	} else if srcValue.Kind() == reflect.Struct && dstValue2.Kind() == reflect.Slice {
		dstValueItem := reflect.New(dstValue.Type().Elem()).Elem()
		dstValueItem = CopyReflectValue(srcValue, dstValueItem)
		fmt.Println("dstValueItem", dstValueItem.Interface())
		dstValue2 = reflect.Append(dstValue2, dstValueItem.Addr())

	} else if srcValue.Kind() == reflect.Struct && dstValue2.Kind() == reflect.Ptr {
		if dstValue2.Elem().Type() == srcValue.Type() {
			dstValue2.Elem().Set(srcValue)
		} else {
			dstValue2.Elem().Set(CopyReflectStruct(srcValue, dstValue2.Elem()))
		}

	} else if srcValue.Kind() == reflect.Slice && dstValue2.Kind() == reflect.Ptr {
		if srcValue.Len() > 0 {
			return CopyReflectStruct(srcValue.Index(0), dstValue2.Elem())
		}
	} else {
		fmt.Println("srcValue.Kind()", srcValue.Kind())
		fmt.Println("dstValue.Kind()", dstValue2.Kind())
	}
	/*	if dstValue.Type().Kind() == dstValue2.Type().Kind() {
			fmt.Println("dstValue.CanAddr()",dstValue.CanAddr())
			fmt.Println("dstValue2.CanAddr()",dstValue2.CanAddr())
			fmt.Println("srcValue.Kind()",srcValue.Kind())
			fmt.Println("dstValue.Kind()",dstValue2.Kind())
			dstValue.Set(dstValue2)

		} else {
			dstValue.Set(dstValue2.Addr())

		}*/
	return dstValue2
}
func CopyReflectSlice(srcValue, dstValue reflect.Value) reflect.Value {
	if srcValue.Kind() == reflect.Slice && dstValue.Kind() == reflect.Slice {
		for i := 0; i < srcValue.Len(); i++ {
			srcValueItem := srcValue.Index(i)
			dstValueItem := reflect.New(dstValue.Type().Elem()).Elem()
			dstValueItem = CopyReflectValue(srcValueItem, dstValueItem)
			fmt.Println("dstValueItem", dstValueItem.Interface())
			dstValue = reflect.Append(dstValue, dstValueItem.Addr())
			fmt.Println("dstValue", dstValue.Interface())

		}
	}
	return dstValue
}
func CopyReflectStruct(srcValue, dstValue reflect.Value) reflect.Value {
	for i := 0; i < srcValue.NumField(); i++ {
		if srcValue.Type().Field(i).IsExported() {
			dstField := dstValue.FieldByName(srcValue.Type().Field(i).Name)
			CopyReflectValue(srcValue.Field(i), dstField)
			fmt.Println(srcValue.Type().Field(i).Name, srcValue.Field(i).Interface(), dstField.Interface())
		}
	}
	fmt.Println("dstValue", dstValue.Interface())
	return dstValue
}
