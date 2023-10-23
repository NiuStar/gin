package utils

import (
	"path"
	"time"
)

const LocalTimeFormat = `2006-01-02 15:04:05`
const FileTimeFormat = `20060102_150405`

func NowTimeSecond() int {
	return int(time.Now().Unix())
}

func NowTimeMSecond() int {
	return int(time.Now().UnixNano() / 1000000)
}

func NowTimeString() string {
	return time.Now().Local().Format(LocalTimeFormat)
}

func NowTimeForFileString() string {
	return time.Now().Local().Format(FileTimeFormat)
}

func NewFileNameFromTime(oldName string) string {
	ext := path.Ext(oldName)
	return oldName[:] + time.Now().Local().Format(FileTimeFormat) + ext
}

func ParseTime(t string) (time.Time, error) {
	return time.ParseInLocation(LocalTimeFormat, t, time.Local)
}
