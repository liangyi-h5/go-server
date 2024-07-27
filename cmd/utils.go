package main

import (
	"reflect"
	"regexp"
)

func isEmail(email string) bool {
	// 定义邮箱地址正则表达式模式
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`

	// 编译正则表达式
	re := regexp.MustCompile(emailRegex)
	e := re.MatchString(email)
	return e
}

func isString(str any) bool {
	return reflect.TypeOf(str).Kind() == reflect.String
}

func isInt(i any) bool {
	return reflect.TypeOf(i).Kind() == reflect.Int
}
