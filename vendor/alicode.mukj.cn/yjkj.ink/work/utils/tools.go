package utils

import (
	"bytes"
	"strings"
)
// Marshal 转换为大驼峰命名法则 首字母大写，“_” 忽略后大写
func Marshal(name string) string {
	if name == "" {
		return ""
	}

	temp := strings.Split(name, "_")
	buffer := bytes.Buffer{}
	var s string
	for index, v := range temp {
		if index == 0 {
			buffer.WriteString(v)
			continue
		}
		vv := []rune(v)
		if len(vv) > 0 {
			if bool(vv[0] >= 'a' && vv[0] <= 'z') { //首字母大写
				vv[0] -= 32
			}
			buffer.WriteString(string(vv))
		}
	}

	s = uncommonInitialismsReplacer.Replace(buffer.String())
	//smap.Set(name, s)
	return s
}

// UnMarshal 回退网络模式
func UnMarshal(name string) string {
	const (
		lower = false
		upper = true
	)

	if name == "" {
		return ""
	}

	var (
		value                                    = commonInitialismsReplacer.Replace(name)
		buf                                      = bytes.NewBufferString("")
		lastCase, currCase, nextCase, nextNumber bool
	)

	for i, v := range value[:len(value)-1] {
		nextCase = bool(value[i+1] >= 'A' && value[i+1] <= 'Z')
		nextNumber = bool(value[i+1] >= '0' && value[i+1] <= '9')

		if i > 0 {
			if currCase == upper {
				if lastCase == upper && (nextCase == upper || nextNumber == upper) {
					buf.WriteRune(v)
				} else {
					if value[i-1] != '_' && value[i+1] != '_' {
						buf.WriteRune('_')
					}
					buf.WriteRune(v)
				}
			} else {
				buf.WriteRune(v)
				if i == len(value)-2 && (nextCase == upper && nextNumber == lower) {
					buf.WriteRune('_')
				}
			}
		} else {
			currCase = upper
			buf.WriteRune(v)
		}
		lastCase = currCase
		currCase = nextCase
	}

	buf.WriteByte(value[len(value)-1])

	s := strings.ToLower(buf.String())
	return s
}


// Copied from golint
var commonInitialisms = []string{"ACL", "API", "ASCII", "CPU", "CSS", "DNS", "EOF", "GUID", "HTML", "HTTP", "HTTPS", "ID", "IP", "JSON", "LHS", "QPS", "RAM", "RHS", "RPC", "SLA", "SMTP", "SQL", "SSH", "TCP", "TLS", "TTL", "UDP", "UI", "UID", "UUID", "URI", "URL", "UTF8", "VM", "XML", "XMPP", "XSRF", "XSS"}
var commonInitialismsReplacer *strings.Replacer
var uncommonInitialismsReplacer *strings.Replacer

func init() {
	var commonInitialismsForReplacer []string
	var uncommonInitialismsForReplacer []string
	for _, initialism := range commonInitialisms {
		commonInitialismsForReplacer = append(commonInitialismsForReplacer, initialism, strings.Title(strings.ToLower(initialism)))
		uncommonInitialismsForReplacer = append(uncommonInitialismsForReplacer, strings.Title(strings.ToLower(initialism)), initialism)
	}
	commonInitialismsReplacer = strings.NewReplacer(commonInitialismsForReplacer...)
	uncommonInitialismsReplacer = strings.NewReplacer(uncommonInitialismsForReplacer...)
}
