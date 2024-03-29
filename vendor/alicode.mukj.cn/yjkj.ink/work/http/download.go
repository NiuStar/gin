package http

import (
	"alicode.mukj.cn/yjkj.ink/work/print"
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strconv"
)

func IsFileExist(filename string, filesize int64) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		fmt.Println(info)
		return false
	}
	if filesize == info.Size() {
		fmt.Println("安装包已存在！", info.Name(), info.Size(), info.ModTime())
		return true
	}
	del := os.Remove(filename)
	if del != nil {
		fmt.Println(del)
	}
	return false
}

func DownloadFile(url string, localPath string, fb func(length, downLen int64)) error {
	var (
		fsize   int64
		buf     = make([]byte, 32*1024)
		written int64
	)
	if fb == nil {
		print := print.NewPrint()
		out := fmt.Sprintf("%s %s %s [%s] %+v", localPath, "开始下载", 0, "/", 0)

		index2 := print.Add("0", "")
		fb = func(length, downLen int64) {
			out = fmt.Sprintf("%s %s %d%s%d  [%+v%s%d]", localPath, "下载进度", downLen, "/", length, downLen*100.0/length, "/", 100)
			print.Print(index2, out)
		}
	}
	tmpFilePath := localPath + ".download"
	fmt.Println(tmpFilePath)
	//创建一个http client
	client := new(http.Client)
	//client.Timeout = time.Second * 60 //设置超时时间
	//get方法获取资源
	resp, err := client.Get(url)
	if err != nil {
		return err
	}

	//读取服务器返回的文件大小
	fsize, err = strconv.ParseInt(resp.Header.Get("Content-Length"), 10, 32)
	if err != nil {
		fmt.Println(err)
	}
	if IsFileExist(localPath, fsize) {
		return err
	}
	fmt.Println("fsize", fsize)

	if resp.Body == nil {
		return errors.New("body is null")
	}
	defer resp.Body.Close()
	//下面是 io.copyBuffer() 的简化版本
	if resp.StatusCode != 200 {
		buffer := &bytes.Buffer{}
		for {
			//读取bytes
			nr, er := resp.Body.Read(buf)
			if nr > 0 {
				//写入bytes
				nw, ew := buffer.Write(buf[0:nr])
				//数据长度大于0
				if nw > 0 {
					written += int64(nw)
				}
				//写入出错
				if ew != nil {
					err = ew
					break
				}
				//读取是数据长度不等于写入的数据长度
				if nr != nw {
					err = io.ErrShortWrite
					break
				}
			}
			if er != nil {
				if er != io.EOF {
					err = er
				}
				break
			}
			//没有错误了快使用 callback
			if fb != nil {
				fb(fsize, written)
			}
		}
		return errors.New(buffer.String())
	} else {
		//创建文件
		os.MkdirAll(path.Dir(tmpFilePath), os.ModePerm)
		file, err := os.Create(tmpFilePath)
		if err != nil {
			return err
		}
		defer file.Close()
		for {
			//读取bytes
			nr, er := resp.Body.Read(buf)
			if nr > 0 {
				//写入bytes
				nw, ew := file.Write(buf[0:nr])
				//数据长度大于0
				if nw > 0 {
					written += int64(nw)
				}
				//写入出错
				if ew != nil {
					err = ew
					break
				}
				//读取是数据长度不等于写入的数据长度
				if nr != nw {
					err = io.ErrShortWrite
					break
				}
			}
			if er != nil {
				if er != io.EOF {
					err = er
				}
				break
			}
			//没有错误了快使用 callback
			if fb != nil {
				fb(fsize, written)
			}
		}
		fmt.Println(err)
		if err == nil {
			file.Close()
			err = os.Rename(tmpFilePath, localPath)
			fmt.Println(err)
		}
	}

	return err
}

func DownloadFileWithHeader(url string, localPath string, header *http.Header, fb func(length, downLen int64)) error {
	var (
		fsize   int64
		buf     = make([]byte, 32*1024)
		written int64
	)

	if fb == nil {
		print := print.NewPrint()
		out := fmt.Sprintf("%s %s %s [%s] %+v", localPath, "开始下载", 0, "/", 0)

		index2 := print.Add("0", "")
		fb = func(length, downLen int64) {
			out = fmt.Sprintf("%s %s %d%s%d  [%+v%s%d]", localPath, "下载进度", downLen, "/", length, downLen*100.0/length, "/", 100)
			print.Print(index2, out)
		}
	}
	tmpFilePath := localPath + ".download"
	fmt.Println(tmpFilePath)
	//创建一个http client
	client := new(http.Client)
	//client.Timeout = time.Second * 60 //设置超时时间
	//get方法获取资源
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	for _, c := range header.Values("Cookie") {
		req.Header.Add("Cookie", c)
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	//读取服务器返回的文件大小
	fsize, err = strconv.ParseInt(resp.Header.Get("Content-Length"), 10, 32)
	if err != nil {
		fmt.Println(err)
	}
	if IsFileExist(localPath, fsize) {
		return err
	}
	fmt.Println("fsize", fsize)

	if resp.Body == nil {
		return errors.New("body is null")
	}
	defer resp.Body.Close()
	//下面是 io.copyBuffer() 的简化版本
	if resp.StatusCode != 200 {
		buffer := &bytes.Buffer{}
		for {
			//读取bytes
			nr, er := resp.Body.Read(buf)
			if nr > 0 {
				//写入bytes
				nw, ew := buffer.Write(buf[0:nr])
				//数据长度大于0
				if nw > 0 {
					written += int64(nw)
				}
				//写入出错
				if ew != nil {
					err = ew
					break
				}
				//读取是数据长度不等于写入的数据长度
				if nr != nw {
					err = io.ErrShortWrite
					break
				}
			}
			if er != nil {
				if er != io.EOF {
					err = er
				}
				break
			}
			//没有错误了快使用 callback
			if fb != nil {
				fb(fsize, written)
			}
		}
		return errors.New(buffer.String())
	} else {
		//创建文件
		os.MkdirAll(path.Dir(tmpFilePath), os.ModePerm)
		file, err := os.Create(tmpFilePath)
		if err != nil {
			return err
		}
		defer file.Close()
		for {
			//读取bytes
			nr, er := resp.Body.Read(buf)
			if nr > 0 {
				//写入bytes
				nw, ew := file.Write(buf[0:nr])
				//数据长度大于0
				if nw > 0 {
					written += int64(nw)
				}
				//写入出错
				if ew != nil {
					err = ew
					break
				}
				//读取是数据长度不等于写入的数据长度
				if nr != nw {
					err = io.ErrShortWrite
					break
				}
			}
			if er != nil {
				if er != io.EOF {
					err = er
				}
				break
			}
			//没有错误了快使用 callback
			if fb != nil {
				fb(fsize, written)
			}
		}
		fmt.Println(err)
		if err == nil {
			file.Close()
			err = os.Rename(tmpFilePath, localPath)
			fmt.Println(err)
		}
	}

	return err
}
