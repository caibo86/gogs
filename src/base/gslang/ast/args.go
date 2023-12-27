// -------------------------------------------
// @file      : args.go
// @author    : 蔡波
// @contact   : caibo923@gmail.com
// @time      : 2023/12/16 上午11:52
// -------------------------------------------

package ast

// Args 参数列表节点,是一个Expr
type Args struct {
	BaseExpr        // 内嵌基本表达式实现
	Items    []Expr // 参数列表
}

// NewArgs 在代码内新建参数列表节点 此参数列表所属代码节点为此代码节点
func (script *Script) NewArgs() *Args {
	args := &Args{}
	args.Init("args", script)
	return args
}

// NewArg 在参数列表节点内 保存对应表达式对应的参数 此参数的父节点为此参数列表节点
func (args *Args) NewArg(arg Expr) Expr {
	// 添加到参数列表
	args.Items = append(args.Items, arg)
	// 设置参数的父节点
	arg.SetParent(args)
	return arg
}

// NamedArgs 命名参数列表节点,是一个Expr
type NamedArgs struct {
	BaseExpr                 // 内嵌基本表达式实现
	Items    map[string]Expr // 用字典保存命名参数列表
}

// NewNamedArgs 在代码节点内新建命名参数列表 此命名参数列表名字args 所属代码节点为此代码节点
func (script *Script) NewNamedArgs() *NamedArgs {
	args := &NamedArgs{
		Items: make(map[string]Expr),
	}
	args.Init("args", script)
	return args
}

// NewArg 用指定的名字和表达式在命名参数列表内添加参数 并返回此参数表达式和添加结果
func (args *NamedArgs) NewArg(name string, arg Expr) (Expr, bool) {
	// 先检查是否有同名参数 有则返回此参数 及 新建失败标志
	if item, ok := args.Items[name]; ok {
		return item, false
	}
	args.Items[name] = arg
	// 设置此参数的父节点为此命名参数列表
	arg.SetParent(args)
	return arg, true
}
