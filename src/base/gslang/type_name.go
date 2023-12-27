// -------------------------------------------
// @file      : type_name.go
// @author    : 蔡波
// @contact   : caibo923@gmail.com
// @time      : 2023/12/28 上午00:42
// -------------------------------------------

package gslang

import (
	"fmt"
	"gogs/base/gslang/ast"
	"strings"
)

// gslang内置类型对应的golang表示
var keyMapping = map[string]string{
	".gslang.Byte":    "byte",
	".gslang.Sbyte":   "int8",
	".gslang.Int16":   "int16",
	".gslang.Uint16":  "uint16",
	".gslang.Int32":   "int32",
	".gslang.Uint32":  "uint32",
	".gslang.Int64":   "int64",
	".gslang.Uint64":  "uint64",
	".gslang.Float32": "float32",
	".gslang.Float64": "float64",
	".gslang.Bool":    "bool",
	".gslang.String":  "string",
}

// TypeName 取ast中的类型对应的golang表示
func TypeName(expr ast.Expr) string {
	switch expr.(type) {
	case *ast.TypeRef:
		// 内置类型引用
		if key, ok := keyMapping[expr.Name()]; ok {
			return key
		}
		ref := expr.(*ast.TypeRef)
		// 枚举
		if _, ok := ref.Ref.(*ast.Enum); ok {
			return strings.TrimLeft(expr.Name(), ".")
		}
		// 自定义类型引用均用指针
		return "*" + strings.TrimLeft(expr.Name(), ".")
	case *ast.Array:
		// 数组
		array := expr.(*ast.Array)
		return fmt.Sprintf("[%d]%s", array.Length, TypeName(array.Element))
	case *ast.List:
		// 切片
		list := expr.(*ast.List)
		return fmt.Sprintf("[]%s", TypeName(list.Element))
	case *ast.Map:
		// 字典
		hash := expr.(*ast.Map)
		return fmt.Sprintf("map[%s]%s", TypeName(hash.Key), TypeName(hash.Value))
	}
	return "unknown"
}
