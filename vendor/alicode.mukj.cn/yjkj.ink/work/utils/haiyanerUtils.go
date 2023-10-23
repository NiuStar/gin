package utils

import (
	"strings"
	"fmt"
	"strconv"
)
func FloatDeleteString(str string) string{
	strArr := []string{"a","b","c","d","e","f","g","h","i","j","k","l","m","n","o","p","q","r","s","t","u","v","w","x","y","z","A","B","C","D","E","F","H","I",
		"J","K","L","M","N","O","P","Q","R","S","T","U","V","W","X","Y","Z","+"}
	for _,v := range strArr{
		str = strings.Replace(str,v,"",-1)
	}
	return str
}
func FormatFloat_3(value float64) float64 {
	var zero float64 = 0
	value_ := fmt.Sprintf("%.3f", value)
	f, err := strconv.ParseFloat(value_,64)
	if err!= nil{
		fmt.Println("err:",err)
		return zero
	}
	return f
}
func FormatFloat_2(value float64) float64 {
	var zero float64 = 0
	value_ := fmt.Sprintf("%.2f", value)
	f, err := strconv.ParseFloat(value_,64)
	if err!= nil{
		fmt.Println("err:",err)
		return zero
	}
	return f
}
//float64取绝对值
func CalcAbs_float(x float64) float64{
	if x < 0 {
		return -x
	}
	if x == 0 {
		return 0
	}
	return x
}
func RemoveDuplicate_string(list []string) []string{
	var x []string
	for _, i := range list {
		if len(x) == 0 {
			x = append(x, i)
		} else {
			for k, v := range x {
				if strings.EqualFold(v,i) {
					break
				}
				if k == len(x)-1 {
					x = append(x, i)
				}
			}
		}
	}
	return x
}

func RemoveDuplicate_int(list []int) []int{
	var x []int
	for _, i := range list {
		if len(x) == 0 {
			x = append(x, i)
		} else {
			for k, v := range x {
				if v == i {
					break
				}
				if k == len(x)-1 {
					x = append(x, i)
				}
			}
		}
	}
	return x
}
//移除重复项--不分数据类型
func RemoveDuplicate(list []interface{}) []interface{}{
	var x []interface{}
	for _, i := range list {
		if len(x) == 0 {
			x = append(x, i)
		} else {
			switch typeof(i){
			case "string":
				for k, v := range x {
					if strings.EqualFold(v.(string),i.(string)) {
						break
					}
					if k == len(x)-1 {
						x = append(x, i)
					}
				}
				break
			default:
				for k, v := range x {
					if v == i {
						break
					}
					if k == len(x)-1 {
						x = append(x, i)
					}
				}
				break
			}
		}
	}
	return x
}
func typeof(v interface{}) string {
	return fmt.Sprintf("%T", v)
}


func RemoveDuplicate_map(list []map[string]interface{},key string) []map[string]interface{}{

	var x []map[string]interface{}
	for _, i := range list {
		if len(x) == 0 {
			x = append(x, i)
		} else {
			for k, v := range x {
				if v[key] == i[key] {
					break
				}
				if k == len(x)-1 {
					x = append(x, i)
				}
			}
		}
	}
	return x
}