package models

import (
	"crypto/md5"
	"fmt"
	"time"
)

func UnixToDate(timestamp int) string {
	t := time.Unix(int64(timestamp), 0)
	return t.Format("2006-01-02 15:04:05")
}

func GetDay() string {
	template := "20060102"
	return time.Now().Format(template)
}

func GetUnix() int64 {
	return time.Now().Unix()
}

func Md5(str string) string {
	data := []byte(str)
	return fmt.Sprintf("%x\n", md5.Sum(data))
}
