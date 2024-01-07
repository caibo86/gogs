// -------------------------------------------
// @file      : array.go
// @author    : 蔡波
// @contact   : caibo923@gmail.com
// @time      : 2023/12/16 下午7:46
// -------------------------------------------

package ast

import (
	"fmt"
)

// Array 数组表达式
type Array struct {
	BaseExpr        // 内嵌基本表达式实现
	Length   uint32 // 数组长度
	Element  Expr   // 数组元素类型
}

func (array *Array) OriginName() string {
	return fmt.Sprintf("[%d]%s", array.Length, array.Element.OriginName())
}

// NewArray 在代码节点内新建数组表达式 此数组表达式所属代码节点为此代码节点
func (script *Script) NewArray(length uint32, element Expr) *Array {
	array := &Array{
		Length:  length,
		Element: element,
	}
	// 初始化数组表达式
	array.Init(element.Name(), script)
	// 设置父节点为此代码节点
	array.Element.SetParent(array)
	return array
}

// Slice 切片
type Slice struct {
	BaseExpr      // 内嵌基本表达式实现
	Element  Expr // 切片元素类型
}

func (slice *Slice) OriginName() string {
	return "[]" + slice.Element.OriginName()
}

// NewSlice 在代码节点内新建切片表达式 此切片表达式所属代码节点为此代码节点
func (script *Script) NewSlice(element Expr) *Slice {
	slice := &Slice{
		Element: element,
	}
	// 初始化切片表达式
	slice.Init(element.Name(), script)
	// 设置父节点为此代码节点
	slice.Element.SetParent(slice)
	return slice
}
