// -------------------------------------------
// @file      : position.go
// @author    : 蔡波
// @contact   : caibo923@gmail.com
// @time      : 2023/12/16 下午8:54
// -------------------------------------------

package gslang

import (
	"fmt"
	"path/filepath"
)

// Position 代码位置
type Position struct {
	Line     int    // 行号 从1开始
	Column   int    // 列号 从1开始
	Filename string // 文件名
}

// BaseName 返回基础文件名
func (pos Position) BaseName() string {
	return filepath.Base(pos.Filename)
}

// String 返回代码位置字符串
func (pos Position) String() string {
	return fmt.Sprintf("%s(%d:%d)", pos.Filename, pos.Line, pos.Column)
}

// Valid 返回代码位置是否有效
func (pos Position) Valid() bool {
	return pos.Line != 0
}
