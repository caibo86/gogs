// -------------------------------------------
// @file      : op.go
// @author    : 蔡波
// @contact   : caibo923@gmail.com
// @time      : 2023/12/16 下午6:36
// -------------------------------------------

package ast

// BinaryOp 二元运算表达式
type BinaryOp struct {
	BaseExpr      // 内嵌基本表达式实现
	Left     Expr // 左操作数
	Right    Expr // 右操作数
}

// NewBinaryOp 在代码节点内新建二元运算表达式
func (script *Script) NewBinaryOp(name string, left, right Expr) *BinaryOp {
	op := &BinaryOp{
		Left:  left,
		Right: right,
	}
	op.Init(name, script)
	return op
}

// UnaryOp 一元运算表达式
type UnaryOp struct {
	BaseExpr      // 内嵌基本表达式实现
	Right    Expr // 右操作数
}

// NewUnaryOp 在代码节点内新建一元运算表达式
func (script *Script) NewUnaryOp(name string) *UnaryOp {
	op := &UnaryOp{}
	op.Init(name, script)
	return op
}
