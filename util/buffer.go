// Copyright 2020, The Go Authors. All rights reserved.
// Author: OnlyOneFace
// Date: 2020/7/24

package util

import (
	"bytes"
	"sync"
)

var bufPool = sync.Pool{New: func() interface{} {
	return new(bytes.Buffer)
}}

// BufGet 获取buf
func BufGet() *bytes.Buffer {
	for {
		b := bufPool.Get()
		buf, ok := b.(*bytes.Buffer)
		if !ok {
			continue
		}
		return buf
	}
}

// BufPut 放置buf
func BufPut(buf *bytes.Buffer) {
	buf.Reset()
	bufPool.Put(buf)
}
