// -------------------------------------------
// @file      : visitor.go
// @author    : 蔡波
// @contact   : caibo923@gmail.com
// @time      : 2023/12/16 下午5:50
// -------------------------------------------

package ast

// Visitor 访问者接口
type Visitor interface {
	VisitPackage(*Package) Node     // 访问包节点
	VisitScript(*Script) Node       // 访问代码节点
	VisitEnum(*Enum) Node           // 访问枚举节点
	VisitEnumVal(*EnumVal) Node     // 访问枚举值节点
	VisitTable(*Table) Node         // 访问结构体节点
	VisitField(*Field) Node         // 访问结构体字段节点
	VisitContract(*Contract) Node   // 访问服务节点
	VisitMethod(*Method) Node       // 访问函数节点
	VisitTypeRef(*TypeRef) Node     // 访问类型引用节点
	VisitAttr(*Attr) Node           // 访问属性节点
	VisitArray(*Array) Node         // 访问数组节点
	VisitList(*List) Node           // 访问链表节点
	VisitArgs(*Args) Node           // 访问参数列表节点
	VisitNamedArgs(*NamedArgs) Node // 访问命名参数列表节点
	VisitString(*String) Node       // 访问字符串节点
	VisitFloat(*Float) Node         // 访问浮点数节点
	VisitInt(*Int) Node             // 访问整数节点
	VisitBool(*Bool) Node           // 访问布尔值节点
	VisitBinaryOp(*BinaryOp) Node   // 访问二元表达式节点
	VisitMap(*Map) Node             // 访问Map节点
}

// 访问者模式
// 为每一种节点类型构造一个Accept方法 使其能够实现Node接口
// 每一种节点的接受一个访问者参数 然后以自身为参数调用访问者对应自己类型的访问方法

// Accept 为包节点实现Node接口
func (pkg *Package) Accept(visitor Visitor) Node {
	return visitor.VisitPackage(pkg)
}

// Accept 为代码节点实现Node接口
func (script *Script) Accept(visitor Visitor) Node {
	return visitor.VisitScript(script)
}

// Accept 为枚举节点实现Node接口
func (enum *Enum) Accept(visitor Visitor) Node {
	return visitor.VisitEnum(enum)
}

// Accept 为枚举值节点实现Node接口
func (enumVal *EnumVal) Accept(visitor Visitor) Node {
	return visitor.VisitEnumVal(enumVal)
}

// Accept 为结构体节点实现Node接口
func (table *Table) Accept(visitor Visitor) Node {
	return visitor.VisitTable(table)
}

// Accept 为结构体字段节点实现Node接口
func (field *Field) Accept(visitor Visitor) Node {
	return visitor.VisitField(field)
}

// Accept 为服务节点实现Node接口
func (contract *Contract) Accept(visitor Visitor) Node {
	return visitor.VisitContract(contract)
}

// Accept 为函数节点实现Node接口
func (method *Method) Accept(visitor Visitor) Node {
	return visitor.VisitMethod(method)
}

// Accept 为类型引用节点实现Node接口
func (ref *TypeRef) Accept(visitor Visitor) Node {
	return visitor.VisitTypeRef(ref)
}

// Accept 为属性节点实现Node接口
func (attr *Attr) Accept(visitor Visitor) Node {
	return visitor.VisitAttr(attr)
}

// Accept 为数组节点实现Node接口
func (array *Array) Accept(visitor Visitor) Node {
	return visitor.VisitArray(array)
}

// Accept 为链表节点实现Node接口
func (list *List) Accept(visitor Visitor) Node {
	return visitor.VisitList(list)
}

// Accept 为参数列表节点实现Node接口
func (args *Args) Accept(visitor Visitor) Node {
	return visitor.VisitArgs(args)
}

// Accept 为命名参数列表节点实现Node接口
func (args *NamedArgs) Accept(visitor Visitor) Node {
	return visitor.VisitNamedArgs(args)
}

// Accept 为字符串节点实现Node接口
func (s *String) Accept(visitor Visitor) Node {
	return visitor.VisitString(s)
}

// Accept 为浮点数节点实现Node接口
func (f *Float) Accept(visitor Visitor) Node {
	return visitor.VisitFloat(f)
}

// Accept 为整数节点实现Node接口
func (i *Int) Accept(visitor Visitor) Node {
	return visitor.VisitInt(i)
}

// Accept 为布尔值节点实现Node接口
func (b *Bool) Accept(visitor Visitor) Node {
	return visitor.VisitBool(b)
}

// Accept 为二元表达式节点实现Node接口
func (op *BinaryOp) Accept(visitor Visitor) Node {
	return visitor.VisitBinaryOp(op)
}

// Accept 为Map节点实现Node接口
func (m *Map) Accept(visitor Visitor) Node {
	return visitor.VisitMap(m)
}

// EmptyVisitor 一个空的什么都不做的访问者
type EmptyVisitor struct{}

// VisitString 实现访问者接口
func (visitor *EmptyVisitor) VisitString(*String) Node {
	return nil
}

// VisitFloat 实现访问者接口
func (visitor *EmptyVisitor) VisitFloat(*Float) Node {
	return nil
}

// VisitInt 实现访问者接口
func (visitor *EmptyVisitor) VisitInt(*Int) Node {
	return nil
}

// VisitBool 实现访问者接口
func (visitor *EmptyVisitor) VisitBool(*Bool) Node {
	return nil
}

// VisitPackage 实现访问者接口
func (visitor *EmptyVisitor) VisitPackage(*Package) Node {
	return nil
}

// VisitScript 实现访问者接口
func (visitor *EmptyVisitor) VisitScript(*Script) Node {
	return nil
}

// VisitEnum 实现访问者接口
func (visitor *EmptyVisitor) VisitEnum(*Enum) Node {
	return nil
}

// VisitEnumVal 实现访问者接口
func (visitor *EmptyVisitor) VisitEnumVal(*EnumVal) Node {
	return nil
}

// VisitTable 实现访问者接口
func (visitor *EmptyVisitor) VisitTable(*Table) Node {
	return nil
}

// VisitField 实现访问者接口
func (visitor *EmptyVisitor) VisitField(*Field) Node {
	return nil
}

// VisitContract 实现访问者接口
func (visitor *EmptyVisitor) VisitContract(*Contract) Node {
	return nil
}

// VisitMethod 实现访问者接口
func (visitor *EmptyVisitor) VisitMethod(*Method) Node {
	return nil
}

// VisitTypeRef 实现访问者接口
func (visitor *EmptyVisitor) VisitTypeRef(*TypeRef) Node {
	return nil
}

// VisitAttr 实现访问者接口
func (visitor *EmptyVisitor) VisitAttr(*Attr) Node {
	return nil
}

// VisitArray 实现访问者接口
func (visitor *EmptyVisitor) VisitArray(*Array) Node {
	return nil
}

// VisitList 实现访问者接口
func (visitor *EmptyVisitor) VisitList(*List) Node {
	return nil
}

// VisitMap 实现访问者接口
func (visitor *EmptyVisitor) VisitMap(*Map) Node {
	return nil
}

// VisitArgs 实现访问者接口
func (visitor *EmptyVisitor) VisitArgs(*Args) Node {
	return nil
}

// VisitNamedArgs 实现访问者接口
func (visitor *EmptyVisitor) VisitNamedArgs(*NamedArgs) Node {
	return nil
}

// VisitParam 实现访问者接口
func (visitor *EmptyVisitor) VisitParam(*Param) Node {
	return nil
}

// VisitBinaryOp 实现访问者接口
func (visitor *EmptyVisitor) VisitBinaryOp(*BinaryOp) Node {
	return nil
}
