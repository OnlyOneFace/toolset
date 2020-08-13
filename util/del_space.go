// Copyright 2020, The Go Authors. All rights reserved.
// Author: OnlyOneFace
// Date: 2020/6/27

package util

import "strings"


func DelSpace(args string) []string {
	buf := BufGet()
	defer BufPut(buf)
	buf.Grow(len(args))
	var flag bool
	for _, arg := range args {
		if arg != 32 { //32 空格
			flag = false
			buf.WriteByte(byte(arg))
			continue
		}
		if !flag {
			buf.WriteByte(byte(arg))
		}
		flag = true
	}
	return strings.Split(buf.String(), " ")
}
