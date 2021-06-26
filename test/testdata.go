package test

import (
	"strconv"
	"time"
)

func StringPointer(s string) *string {
	return &s
}

func BoolPointer(b bool) *bool {
	return &b
}

func GetTimeNow() string {
	year := strconv.Itoa(time.Now().Year())
	mouth := time.Now().Month().String()
	day := strconv.Itoa(time.Now().Day())

	return year + mouth + day
}
