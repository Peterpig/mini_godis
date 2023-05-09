package files

import (
	"fmt"
	"os"
)

func CheckNotPermission(src string) bool {
	_, err := os.Stat(src)
	return os.IsPermission(err)
}

func ChecNotExist(src string) bool {
	_, err := os.Stat(src)
	return os.IsNotExist(err)
}

func mkDir(src string) error {
	return os.MkdirAll(src, os.ModePerm)
}

func IsNotExistMkDir(src string) error {
	if ChecNotExist(src) {
		return mkDir(src)
	}
	return nil
}

func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	return err == nil && !info.IsDir()
}

func MustOpen(fileName, src string) (*os.File, error) {
	if CheckNotPermission(src) {
		return nil, fmt.Errorf("permission denied dir: %s", src)
	}

	if err := IsNotExistMkDir(src); err != nil {
		return nil, fmt.Errorf("error during make dir %s, err: %v", src, err)
	}

	f, err := os.OpenFile(src+string(os.PathSeparator)+fileName, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, fmt.Errorf("fail to open file, err %v", err)
	}

	return f, nil
}
