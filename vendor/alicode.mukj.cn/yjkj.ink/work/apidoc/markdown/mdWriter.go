package markdown

//MarkDown文件内容类
import (
	"fmt"
	"os"
)

type MarkDown struct {
	name string
	path string

	content string
}

func NewMarkDown(name, path string) *MarkDown {

	return &MarkDown{name: name, path: path}
}

func (md *MarkDown) WriteTitle(level int, content string) {

	for i := 0; i < level; i++ {
		md.content += "#"
	}
	md.content += " " + content
}

func (md *MarkDown) WriteContent(content string) {
	md.content += content
}

func (md *MarkDown) WriteImportantContent(content string) {
	md.content += "**" + content + "**"
}

func (md *MarkDown) WriteCode(contents string, language string) {

	md.content += "```" + language + `
` + contents + `
` + "```" + `
`
}

func (md *MarkDown) WriteForm(contents [][]string) {

	if len(contents) == 1 {
		var param []string
		param = append(param, "")
		param = append(param, "")
		param = append(param, "")
		param = append(param, "")
		contents = append(contents, param)
	}
	for index, list := range contents {
		md.content += "\r\n|"
		c := ""
		for _, content := range list {
			c += content
			c += " |"
		}
		if len(c) == 0 {
			md.content += "\r\n|"
		} else {
			md.content += c
		}
		if index == 0 {
			md.content += "\r\n|"
			c = ""
			for _, _ = range list {
				c += " ---------- |"
			}
			if len(c) == 0 {
				md.content += " - |"
			} else {
				md.content += c
			}
		}
	}
	md.content += "|\r\n"

}

func (md *MarkDown) Save() {
	if err := os.Mkdir(md.path, 0777); err != nil {
		fmt.Println("创建文件夹失败", md.path, err)
		return
	}
	data := "# " + md.name + `
` + md.content
	WriteWithFileWrite(md.path+"/"+md.name+".md", data)
}

func (md *MarkDown) Content() string {
	return md.content
}

//使用os.OpenFile()相关函数打开文件对象，并使用文件对象的相关方法进行文件写入操作
//清空一次文件
func WriteWithFileWrite(name, content string) {
	fileObj, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		fmt.Println("Failed to open the file", err.Error())
		os.Exit(2)
	}
	defer fileObj.Close()
	if _, err := fileObj.WriteString(content); err == nil {
		fmt.Println("Successful writing to the file with os.OpenFile and *File.WriteString method.", content)
	}
	contents := []byte(content)
	if _, err := fileObj.Write(contents); err == nil {
		fmt.Println("Successful writing to thr file with os.OpenFile and *File.Write method.", content)
	}
}
