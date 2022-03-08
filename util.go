package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func Union(arr1 []string, arr2 []string) []string {
	m := make(map[string]bool)
	for _, str := range arr1 {
		m[str] = true
	}
	for _, str := range arr2 {
		m[str] = true
	}
	res := make([]string, 0, len(m))
	for k := range m {
		res = append(res, k)
	}
	return res
}

func Intersect(arr1 []string, arr2 []string) []string {
	m := make(map[string]bool)
	for _, str := range arr1 {
		m[str] = true
	}
	res := make([]string, 0, len(m))
	for _, str := range arr2 {
		if m[str] {
			res = append(res, str)
		}
	}
	return res
}

func Difference(arr1 []string, arr2 []string) []string {
	m := make(map[string]bool)
	for _, str := range arr1 {
		m[str] = true
	}
	for _, str := range arr2 {
		delete(m, str)
	}
	res := make([]string, 0, len(m))
	for k := range m {
		res = append(res, k)
	}
	return res
}

func DirList(dirpath string) ([]string, error) {
	var dir_list []string
	dir_err := filepath.Walk(dirpath,
		func(p string, f os.FileInfo, err error) error {
			if f == nil {
				return err
			}
			if f.IsDir() && p != dirpath {
				dir_list = append(dir_list, filepath.Base(p))
				return nil
			}

			return nil
		})
	return dir_list, dir_err
}

func FileList(dirpath string) ([]string, error) {
	var fileList []string
	err := filepath.Walk(dirpath,
		func(p string, f os.FileInfo, err error) error {
			if f == nil {
				return err
			}
			if !f.IsDir() {
				fileList = append(fileList, p)
				return nil
			}
			return nil
		})
	return fileList, err
}
func Md5sum(file string) string {
	f, err := os.Open(file)
	if err != nil {
		return ""
	}
	defer f.Close()

	h := md5.New()

	_, err = io.Copy(h, f)
	if err != nil {
		return ""
	}

	return fmt.Sprintf("%x", h.Sum(nil))

}
