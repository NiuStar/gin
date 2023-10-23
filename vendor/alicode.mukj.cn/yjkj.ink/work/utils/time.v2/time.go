package time

import (
	"database/sql/driver"
	"fmt"
	"time"
)

var (
	Zero = Unix(0, 0) //Time{time.Time{}}
)

type Time struct {
	time.Time
}

func Now() Time {
	return Time{Time: time.Now()}
}

func Unix(sec, nsec int64) Time {
	return Time{Time: time.Unix(sec, nsec)}
}

// 重写 MarshaJSON 方法，在此方法中实现自定义格式的转换；程序中解析到JSON
func (t *Time) MarshalJSON() ([]byte, error) {
	formatted := fmt.Sprintf(`"%s"`, t.Format(LocalTimeFormat))
	return []byte(formatted), nil
}

// JSON中解析到程序中
func (t *Time) UnmarshalJSON(data []byte) (err error) {
	now, err := time.ParseInLocation(LocalTimeFormat, string(data), time.Local)
	*t = Time{Time: now}
	return
}

// 写入数据库时会调用该方法将自定义时间类型转换并写入数据库
func (t Time) Value() (driver.Value, error) {
	var zeroTime time.Time
	if t.Time.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return t.Time, nil
}

// 读取数据库时会调用该方法将时间数据转换成自定义时间类型
func (t *Time) Scan(v interface{}) error {
	v1, ok := v.(time.Time)
	if !ok {
		return fmt.Errorf("can not convert %v to timestamp", v)
	}
	t.Time = v1
	return nil
}

func (t *Time) String() string {
	if t == nil {
		return ""
	}
	return t.Time.Format(LocalTimeFormat)
}

func (t *Time) In(loc *time.Location) time.Time {
	if t == nil {
		return time.Now()
	}
	return t.Time.In(loc)
}
