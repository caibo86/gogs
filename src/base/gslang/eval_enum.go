// -------------------------------------------
// @file      : eval_enum.go
// @author    : 蔡波
// @contact   : caibo923@gmail.com
// @time      : 2023/12/14 下午9:11
// -------------------------------------------

package gslang

import (
	"gogs/base/gslang/ast"
	log "gogs/base/logger"
)

// evalEnumVal 访问枚举值
type evalEnumVal struct {
	val int64
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
	log.Panicf("stmt is not enum expr :%s", Pos(node))
	return nil
}

// VisitFloat 仅为实现访问者接口
func (visitor *evalEnumVal) VisitFloat(node *ast.Float) ast.Node {
	log.Panicf("stmt is not enum expr :%s", Pos(node))
	return nil
}

// VisitInt 仅为实现访问者接口
func (visitor *evalEnumVal) VisitInt(node *ast.Int) ast.Node {
	log.Panicf("stmt is not enum expr :%s", Pos(node))
	return nil
}

// VisitBool 仅为实现访问者接口
func (visitor *evalEnumVal) VisitBool(node *ast.Bool) ast.Node {
	log.Panicf("stmt is not enum expr :%s", Pos(node))
	return nil
}

// VisitPackage 仅为实现访问者接口
func (visitor *evalEnumVal) VisitPackage(node *ast.Package) ast.Node {
	log.Panicf("stmt is not enum expr :%s", Pos(node))
	return nil
}

// VisitScript 仅为实现访问者接口
func (visitor *evalEnumVal) VisitScript(node *ast.Script) ast.Node {
	log.Panicf("stmt is not enum expr :%s", Pos(node))
	return nil
}

// VisitEnum 仅为实现访问者接口
func (visitor *evalEnumVal) VisitEnum(node *ast.Enum) ast.Node {
	log.Panicf("stmt is not enum expr :%s", Pos(node))
	return nil
}

// VisitStruct 仅为实现访问者接口
func (visitor *evalEnumVal) VisitStruct(node *ast.Struct) ast.Node {
	log.Panicf("stmt is not enum expr :%s", Pos(node))
	return nil
}

// VisitField 仅为实现访问者接口
func (visitor *evalEnumVal) VisitField(node *ast.Field) ast.Node {
	log.Panicf("stmt is not enum expr :%s", Pos(node))
	return nil
}

// VisitService 仅为实现访问者接口
func (visitor *evalEnumVal) VisitService(node *ast.Service) ast.Node {
	log.Panicf("stmt is not enum expr :%s", Pos(node))
	return nil
}

// VisitFunction 仅为实现访问者接口
func (visitor *evalEnumVal) VisitFunction(node *ast.Method) ast.Node {
	log.Panicf("stmt is not enum expr :%s", Pos(node))
	return nil
}

// VisitAttr 仅为实现访问者接口
func (visitor *evalEnumVal) VisitAttr(node *ast.Attr) ast.Node {
	log.Panicf("stmt is not enum expr :%s", Pos(node))
	return nil
}

// VisitArray 仅为实现访问者接口
func (visitor *evalEnumVal) VisitArray(node *ast.Array) ast.Node {
	log.Panicf("stmt is not enum expr :%s", Pos(node))
	return nil
}

// VisitList 仅为实现访问者接口
func (visitor *evalEnumVal) VisitList(node *ast.List) ast.Node {
	log.Panicf("stmt is not enum expr :%s", Pos(node))
	return nil
}

// VisitMap 仅为实现访问者接口
func (visitor *evalEnumVal) VisitMap(node *ast.Map) ast.Node {
	log.Panicf("stmt is not enum expr :%s", Pos(node))
	return nil
}

// VisitArgs 仅为实现访问者接口
func (visitor *evalEnumVal) VisitArgs(node *ast.Args) ast.Node {
	log.Panicf("stmt is not enum expr :%s", Pos(node))
	return nil
}

// VisitNamedArgs 仅为实现访问者接口
func (visitor *evalEnumVal) VisitNamedArgs(node *ast.NamedArgs) ast.Node {
	log.Panicf("stmt is not enum expr :%s", Pos(node))
	return nil
}

// VisitParam 仅为实现访问者接口
func (visitor *evalEnumVal) VisitParam(node *ast.Param) ast.Node {
	log.Panicf("stmt is not enum expr :%s", Pos(node))
	return nil
}
