// -------------------------------------------
// @file      : array.go
// @author    : 蔡波
// @contact   : caibo923@gmail.com
// @time      : 2023/12/16 下午7:46
// -------------------------------------------

package ast

// Array 数组表达式
type Array struct {
	BaseExpr        // 内嵌基本表达式实现
	Length   uint32 // 数组长度
	Element  Expr   // 数组元素类型
}

func (array *Array) OriginName() string {
	return "[]" + array.Element.OriginName()
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

// List 链表表达式
type List struct {
	BaseExpr      // 内嵌基本表达式实现
	Element  Expr // 链表元素类型
}

func (list *List) OriginName() string {
	return "[]" + list.Element.OriginName()
}

// NewList 在代码节点内新建链表表达式 此链表表达式所属代码节点为此代码节点
func (script *Script) NewList(element Expr) *List {
	list := &List{
		Element: element,
	}
	// 初始化链表表达式
	list.Init(element.Name(), script)
	// 设置父节点为此代码节点
	list.Element.SetParent(list)
	return list
}
