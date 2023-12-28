// -------------------------------------------
// @file      : eval_arg.go
// @author    : 蔡波
// @contact   : caibo923@gmail.com
// @time      : 2023/12/16 下午8:50
// -------------------------------------------

package gslang

import (
	"gogs/base/gserrors"
	"gogs/base/gslang/ast"
)

// evalArg 参数处理器
type evalArg struct {
	field *ast.Field
	expr  ast.Expr
}

// VisitArgs 实现访问者  访问参数列表节点 将参数列表中与field的id相同的参数 保存在expr 中
func (visitor *evalArg) VisitArgs(node *ast.Args) ast.Node {
	for idx, arg := range node.Items {
		if uint16(idx) == visitor.field.ID {
			visitor.expr = arg
		}
	}
	return nil
}

// VisitNamedArgs 实现访问者 访问命名参数列表 将命名参数列表中与field名字相同的参数保存在expr中
func (visitor *evalArg) VisitNamedArgs(node *ast.NamedArgs) ast.Node {
	for idx, arg := range node.Items {
		if idx == visitor.field.Name() {
			visitor.expr = arg
		}
	}
	return nil
}

// VisitString 仅仅为实现访问者
func (visitor *evalArg) VisitString(node *ast.String) ast.Node {
	panic(gserrors.Newf(nil, "node is not args or named args: %s", Pos(node)))
	return nil
}

// VisitFloat 仅仅为实现访问者
func (visitor *evalArg) VisitFloat(node *ast.Float) ast.Node {
	panic(gserrors.Newf(nil, "node is not args or named args: %s", Pos(node)))
	return nil
}

// VisitInt 仅仅为实现访问者
func (visitor *evalArg) VisitInt(node *ast.Int) ast.Node {
	panic(gserrors.Newf(nil, "node is not args or named args: %s", Pos(node)))
	return nil
}

// VisitBool 仅仅为实现访问者
func (visitor *evalArg) VisitBool(node *ast.Bool) ast.Node {
	panic(gserrors.Newf(nil, "node is not args or named args: %s", Pos(node)))
	return nil
}

// VisitPackage 仅仅为实现访问者
func (visitor *evalArg) VisitPackage(node *ast.Package) ast.Node {
	panic(gserrors.Newf(nil, "node is not args or named args: %s", Pos(node)))
	return nil
}

// VisitScript 仅仅为实现访问者
func (visitor *evalArg) VisitScript(node *ast.Script) ast.Node {
	panic(gserrors.Newf(nil, "node is not args or named args: %s", Pos(node)))
	return nil
}

// VisitEnum 仅仅为实现访问者
func (visitor *evalArg) VisitEnum(node *ast.Enum) ast.Node {
	panic(gserrors.Newf(nil, "node is not args or named args: %s", Pos(node)))
	return nil
}

// VisitEnumVal 仅仅为实现访问者
func (visitor *evalArg) VisitEnumVal(node *ast.EnumVal) ast.Node {
	panic(gserrors.Newf(nil, "node is not args or named args: %s", Pos(node)))
	return nil
}

// VisitTable 仅仅为实现访问者
func (visitor *evalArg) VisitTable(node *ast.Table) ast.Node {
	panic(gserrors.Newf(nil, "node is not args or named args: %s", Pos(node)))
	return nil
}

// VisitField 仅仅为实现访问者
func (visitor *evalArg) VisitField(node *ast.Field) ast.Node {
	panic(gserrors.Newf(nil, "node is not args or named args: %s", Pos(node)))
	return nil
}

// VisitContract 仅仅为实现访问者
func (visitor *evalArg) VisitContract(node *ast.Contract) ast.Node {
	panic(gserrors.Newf(nil, "node is not args or named args: %s", Pos(node)))
	return nil
}

// VisitMethod 仅仅为实现访问者
func (visitor *evalArg) VisitMethod(node *ast.Method) ast.Node {
	panic(gserrors.Newf(nil, "node is not args or named args: %s", Pos(node)))
	return nil
}

// VisitTypeRef 仅仅为实现访问者
func (visitor *evalArg) VisitTypeRef(node *ast.TypeRef) ast.Node {
	panic(gserrors.Newf(nil, "node is not args or named args: %s", Pos(node)))
	return nil
}

// VisitAttr 仅仅为实现访问者
func (visitor *evalArg) VisitAttr(node *ast.Attr) ast.Node {
	panic(gserrors.Newf(nil, "node is not args or named args: %s", Pos(node)))
	return nil
}

// VisitArray 仅仅为实现访问者
func (visitor *evalArg) VisitArray(node *ast.Array) ast.Node {
	panic(gserrors.Newf(nil, "node is not args or named args: %s", Pos(node)))
	return nil
}

// VisitSlice 仅仅为实现访问者
func (visitor *evalArg) VisitSlice(node *ast.Slice) ast.Node {
	panic(gserrors.Newf(nil, "node is not args or named args: %s", Pos(node)))
	return nil
}

// VisitMap 仅仅为实现访问者
func (visitor *evalArg) VisitMap(node *ast.Map) ast.Node {
	panic(gserrors.Newf(nil, "node is not args or named args: %s", Pos(node)))
	return nil
}

// VisitParam 仅仅为实现访问者
func (visitor *evalArg) VisitParam(node *ast.Param) ast.Node {
	panic(gserrors.Newf(nil, "node is not args or named args: %s", Pos(node)))
	return nil
}

// VisitBinaryOp 仅仅为实现访问者
func (visitor *evalArg) VisitBinaryOp(node *ast.BinaryOp) ast.Node {
	panic(gserrors.Newf(nil, "node is not args or named args: %s", Pos(node)))
	return nil
}
