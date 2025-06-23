package utils

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"reflect"
)

// Json2Excel json转excel []map[string]interface{} -> excel 或者 []struct -> excel
func Json2Excel(list interface{}, fileName string) {

	columnsMap := map[string]string{}
	columns := []string{}
	maps := []map[string]interface{}{}
	// todo 判断是否为[]map[string]interface{} 或者 []struct
	if list2, ok := list.([]interface{}); ok {
		if len(list2) == 0 {
			return
		}
		v1 := reflect.ValueOf(list2[0])
		v1 = reflect.Indirect(v1)
		if v1.Kind() == reflect.Map {
			// todo []map[string]interface{}
			for _, v := range list2 {
				if v2, ok2 := v.(map[string]interface{}); ok2 {
					// todo map[string]interface{}
					for k, _ := range v2 {
						// todo 获取字段名
						if _, ok3 := columnsMap[k]; !ok3 {
							columnsMap[k] = k
							columns = append(columns, k)
						}
					}
					maps = append(maps, v2)
				} else {
					// todo 不是map[string]interface{}而是struct
				}
			}
		} else if v1.Kind() == reflect.Struct {
			maps, columns = list2Map(list2[0])
		}
	} else {
		maps, columns = list2Map(list)
	}
	f := excelize.NewFile()
	sheetName := "Sheet1"
	// 创建新的Sheet
	f.NewSheet(sheetName)
	for index, v := range columns {
		f.SetCellValue(sheetName, fmt.Sprintf("%s%d", AAZZ(int64(index)), 1), v)
	}
	for index2, mapv := range maps {
		for index, v := range columns {
			cellName := fmt.Sprintf("%s%d", AAZZ(int64(index)), index2+2)
			f.SetCellValue(sheetName, cellName, mapv[v])
		}
	}
	if err := f.SaveAs(fmt.Sprintf("%s.xlsx", fileName)); err != nil {
		fmt.Println(err)
	}
}

func AAZZ(index int64) string {
	index2 := index
	if index2 < 26 {
		return string(rune(65 + index2))
	}
	var indexs []int64

	lastIndex := index2 % 26
	index2 -= lastIndex
	indexs = append(indexs, lastIndex)
	index2 /= 26
	for ; index2 >= 26; index2 /= 26 {
		lastIndex = index2 % 26
		indexs = append(indexs, lastIndex)
	}
	indexs = append(indexs, index2-1)

	var res string
	for i := len(indexs) - 1; i >= 0; i-- {
		res += string(rune(65 + indexs[i]))
	}
	return res
}

func list2Map(obj interface{}) (res []map[string]interface{}, columns []string) {
	columnsMap := map[string]interface{}{}
	vs := reflect.ValueOf(obj)
	if vs.Kind() == reflect.Slice {
		l := vs.Cap()
		l = vs.Len()
		for i := 0; i < l; i++ {
			v := vs.Index(i)
			for v.Kind() == reflect.Ptr {
				v = reflect.Indirect(v)
			}
			if v.Kind() == reflect.Struct {
				m, c := obj2Map(v)
				res = append(res, m)
				for _, vc := range c {
					if vcm := columnsMap[vc]; vcm == nil {
						columns = append(columns, vc)
						columnsMap[vc] = vc
					}
				}
			}
		}
	}
	return
}

func obj2Map(obj reflect.Value) (res map[string]interface{}, columns []string) {
	res = make(map[string]interface{})
	for obj.Kind() == reflect.Ptr {
		obj = reflect.Indirect(obj)
	}

	if obj.Kind() == reflect.Struct {
		for i := 0; i < obj.NumField(); i++ {
			field := obj.Type().Field(i)

			if field.IsExported() {
				excelName := field.Name
				excelName2 := field.Tag.Get("excel")
				if excelName2 == "" {
					commentName := field.Tag.Get("comment")
					if commentName != "" {
						excelName = commentName
					}
				} else {
					excelName = excelName2
				}

				columns = append(columns, excelName)
				res[excelName] = obj.Field(i).Interface()
			}

		}
	}
	return
}
