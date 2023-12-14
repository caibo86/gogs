// -------------------------------------------
// @file      : args.go
// @author    : 蔡波
// @contact   : caibo923@gmail.com
// @time      : 2023/12/14 上午11:52
// -------------------------------------------

package ast

// Args 参数列表节点
type Args struct {
	BaseExpr        // 内嵌基本表达式实现
	Items    []Expr // 参数列表
}
