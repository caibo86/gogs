// -------------------------------------------
// @file      : attr.go
// @author    : 蔡波
// @contact   : caibo923@gmail.com
// @time      : 2023/12/16 下午4:13
// -------------------------------------------

package ast

// Attr 属性节点
type Attr struct {
	BaseExpr
	Type *TypeRef
	Args Expr
}

// NewAttr 为代码节点创建属性
func (script *Script) NewAttr(attrType *TypeRef) *Attr {
	attr := &Attr{
		Type: attrType,
	}
	attr.Init(attrType.Name(), script)
	attr.Type.SetParent(attr)
	return attr
}
