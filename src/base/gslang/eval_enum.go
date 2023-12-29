// -------------------------------------------
// @file      : eval_enum.go
// @author    : 蔡波
// @contact   : caibo923@gmail.com
// @time      : 2023/12/14 下午9:11
// -------------------------------------------

package gslang

import (
	"gogs/base/gserrors"
	"gogs/base/gslang/ast"
)

// evalEnumVal 访问枚举值
type evalEnumVal struct {
	val int32
}

// VisitBinaryOp 访问二元运算
func (visitor *evalEnumVal) VisitBinaryOp(node *ast.BinaryOp) ast.Node {
	visitor.val = EvalEnumVal(node.Left) | EvalEnumVal(node.Right)
	return nil
}

// VisitTypeRef 访问类型引用
func (visitor *evalEnumVal) VisitTypeRef(node *ast.TypeRef) ast.Node {
	node.Ref.Accept(visitor)
	return node
}

// VisitEnumVal 访问枚举值
func (visitor *evalEnumVal) VisitEnumVal(node *ast.EnumVal) ast.Node {
	visitor.val = node.Value
	return node
}

// VisitString 仅为实现访问者接口
func (visitor *evalEnumVal) VisitString(node *ast.String) ast.Node {
	gserrors.Panicf(nil, "node is not an enum expr: %s", Pos(node))
	return nil
}

// VisitFloat 仅为实现访问者接口
func (visitor *evalEnumVal) VisitFloat(node *ast.Float) ast.Node {
	gserrors.Panicf(nil, "node is not an enum expr: %s", Pos(node))
	return nil
}

// VisitInt 仅为实现访问者接口
func (visitor *evalEnumVal) VisitInt(node *ast.Int) ast.Node {
	gserrors.Panicf(nil, "node is not an enum expr: %s", Pos(node))
	return nil
}

// VisitBool 仅为实现访问者接口
func (visitor *evalEnumVal) VisitBool(node *ast.Bool) ast.Node {
	gserrors.Panicf(nil, "node is not an enum expr: %s", Pos(node))
	return nil
}

// VisitPackage 仅为实现访问者接口
func (visitor *evalEnumVal) VisitPackage(node *ast.Package) ast.Node {
	gserrors.Panicf(nil, "node is not an enum expr: %s", Pos(node))
	return nil
}

// VisitScript 仅为实现访问者接口
func (visitor *evalEnumVal) VisitScript(node *ast.Script) ast.Node {
	gserrors.Panicf(nil, "node is not an enum expr: %s", Pos(node))
	return nil
}

// VisitEnum 仅为实现访问者接口
func (visitor *evalEnumVal) VisitEnum(node *ast.Enum) ast.Node {
	gserrors.Panicf(nil, "node is not an enum expr: %s", Pos(node))
	return nil
}

// VisitTable 仅为实现访问者接口
func (visitor *evalEnumVal) VisitTable(node *ast.Table) ast.Node {
	gserrors.Panicf(nil, "node is not an enum expr: %s", Pos(node))
	return nil
}

// VisitField 仅为实现访问者接口
func (visitor *evalEnumVal) VisitField(node *ast.Field) ast.Node {
	gserrors.Panicf(nil, "node is not an enum expr: %s", Pos(node))
	return nil
}

// VisitContract 仅为实现访问者接口
func (visitor *evalEnumVal) VisitContract(node *ast.Contract) ast.Node {
	gserrors.Panicf(nil, "node is not an enum expr: %s", Pos(node))
	return nil
}

// VisitMethod 仅为实现访问者接口
func (visitor *evalEnumVal) VisitMethod(node *ast.Method) ast.Node {
	gserrors.Panicf(nil, "node is not an enum expr: %s", Pos(node))
	return nil
}

// VisitAttr 仅为实现访问者接口
func (visitor *evalEnumVal) VisitAttr(node *ast.Attr) ast.Node {
	gserrors.Panicf(nil, "node is not an enum expr: %s", Pos(node))
	return nil
}

// VisitArray 仅为实现访问者接口
func (visitor *evalEnumVal) VisitArray(node *ast.Array) ast.Node {
	gserrors.Panicf(nil, "node is not an enum expr: %s", Pos(node))
	return nil
}

// VisitSlice 仅为实现访问者接口
func (visitor *evalEnumVal) VisitSlice(node *ast.Slice) ast.Node {
	gserrors.Panicf(nil, "node is not an enum expr: %s", Pos(node))
	return nil
}

// VisitMap 仅为实现访问者接口
func (visitor *evalEnumVal) VisitMap(node *ast.Map) ast.Node {
	gserrors.Panicf(nil, "node is not an enum expr: %s", Pos(node))
	return nil
}

// VisitArgs 仅为实现访问者接口
func (visitor *evalEnumVal) VisitArgs(node *ast.Args) ast.Node {
	gserrors.Panicf(nil, "node is not an enum expr: %s", Pos(node))
	return nil
}

// VisitNamedArgs 仅为实现访问者接口
func (visitor *evalEnumVal) VisitNamedArgs(node *ast.NamedArgs) ast.Node {
	gserrors.Panicf(nil, "node is not an enum expr: %s", Pos(node))
	return nil
}

// VisitParam 仅为实现访问者接口
func (visitor *evalEnumVal) VisitParam(node *ast.Param) ast.Node {
	gserrors.Panicf(nil, "node is not an enum expr: %s", Pos(node))
	return nil
}
