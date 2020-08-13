// Copyright 2020, The Go Authors. All rights reserved.
// Author: OnlyOneFace
// Date: 2020/6/6

package util

import (
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"
)

func GetCurPath() (string, error) {
	p, err := os.Executable()
	if err != nil {
		return "", err
	}
	fp := filepath.Dir(p)
	return filepath.EvalSymlinks(fp)
}

func SaveFile(file *multipart.FileHeader, dst string) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()
	var out *os.File
	if out, err = os.Create(dst); err != nil {
		return "", err
	}
	defer out.Close()
	_, err = io.Copy(out, src)
	return dst, err
}

// CreateLogPath Build log file path
func CreateLogPath(root, appName string) string {
	buf := BufGet()
	defer BufPut(buf)
	buf.WriteString(root)
	buf.WriteByte('/')
	buf.WriteString(time.Now().Local().Format("20060102"))
	_, err := os.Stat(buf.String())
	if err != nil {
		if !os.IsNotExist(err) {
			return err.Error()
		}
		if err = os.MkdirAll(buf.String(), os.ModeDir); err != nil {
			return err.Error()
		}
	}
	buf.WriteByte('/')
	buf.WriteString(appName)
	buf.WriteString(".log")
	return buf.String()
}
