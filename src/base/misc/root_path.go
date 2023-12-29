// -------------------------------------------
// @file      : root_path.go
// @author    : 蔡波
// @contact   : caibo923@gmail.com
// @time      : 2023/12/29 下午10:05
// -------------------------------------------

package misc

import (
	"os"
	"path/filepath"
	"strings"
)

const (
	DefaultRootPath = "./"
)

var RootLocate = ".root_locate"

// IsFileExist 文件是否存在
func IsFileExist(filename string) bool {
	_, err := os.Lstat(filename)
	return !os.IsNotExist(err)
}

// IsRootDir 判断是否根目录
func IsRootDir(path string) bool {
	locatePath := filepath.Join(path, RootLocate)
	if IsFileExist(locatePath) {
		return true
	}
	return false
}

// GetRootPath 获取根目录
func GetRootPath() string {
	cwd, err := os.Getwd()
	if err != nil {
		return DefaultRootPath
	}
	path, err := filepath.Abs(cwd)
	if err != nil {
		return DefaultRootPath
	}
	for depth := 32; depth > 0 && !IsRootDir(path); depth-- {
		idx := strings.LastIndexByte(path, filepath.Separator)
		if idx >= 0 {
			path = path[:idx]
		} else {
			return DefaultRootPath
		}
	}
	return path
}

func GetBinPath() string {
	rootDir := GetRootPath()
	return filepath.Join(rootDir, "bin")
}

func GetPathInRootDir(filename string) string {
	rootDir := GetRootPath()
	return filepath.Join(rootDir, filename)
}
