// -------------------------------------------
// @file      : parse.go
// @author    : 蔡波
// @contact   : caibo923@gmail.com
// @time      : 2023/12/14 下午8:19
// -------------------------------------------

package gslang

import (
	"gogs/base/gslang/ast"
)

const (
	//  位置额外信息的key
	posExtra = "gslang_parser_pos"
	// 注释额外信息的key
	commentExtra = "gslang_parser_comment"
)

// Pos 节点的位置信息
func Pos(node ast.Node) Position {
	if val, ok := node.Extra(posExtra); ok {
		return val.(Position)
	}
	return Position{
		Line:     0,
		Column:   0,
		Filename: "unknown",
	}
}
