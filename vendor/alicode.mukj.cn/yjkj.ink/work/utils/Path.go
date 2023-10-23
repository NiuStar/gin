package utils

import (
	//"runtime"
	//"os/exec"
	//"path/filepath"
	//"os"
	"path"
	//"strings"
	"time"
)

func GetCurrPath() string {
	/*if strings.HasSuffix(os.Args[0],"___go_build_main_go") {
		return "./"
	}
	if runtime.GOOS == "windows" {
		/*file, _ := exec.LookPath(os.Args[0])
		path, _ := filepath.Abs(file)
		path = path[0:strings.LastIndex(path,"\\")] + "\\"
		return path*/
	/*	return "./"
	} else{
		dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			panic(err)
		}
		return strings.Replace(dir, "\\", "/", -1) + "/"


	}*/
	return "./"
}

func GetUploadNetPath(domain, hardPath, name string) string {

	time.Now().UnixNano()

	// domain http://cninct.com:10001
	//hardPath upload
	//name 1.png
	return domain + "\\" + hardPath + "\\" + name
}

func lastChar(str string) uint8 {
	size := len(str)
	if size == 0 {
		panic("The length of the string can't be 0")
	}
	return str[size-1]
}

func JoinPaths(absolutePath, relativePath string) string {
	if len(relativePath) == 0 {
		return absolutePath
	}

	finalPath := path.Join(absolutePath, relativePath)
	appendSlash := lastChar(relativePath) == '/' && lastChar(finalPath) != '/'
	if appendSlash {
		return finalPath + "/"
	}
	return finalPath
}
