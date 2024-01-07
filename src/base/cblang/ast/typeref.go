// -------------------------------------------
// @file      : typeref.go
// @author    : 蔡波
// @contact   : caibo923@gmail.com
// @time      : 2023/12/16 下午4:18
// -------------------------------------------

package ast

import (
	"bytes"
	"gogs/base/cberrors"
)

// TypeRef 类型引用 是一个Expr
type TypeRef struct {
	BaseExpr
	Ref      Expr
	NamePath []string
	Origin   string // 原始代码
}

func (ref *TypeRef) OriginName() string {
	return ref.Origin
}

// NewTypeRef 在代码节点内新建类型引用
func (script *Script) NewTypeRef(namePath []string, origin string) *TypeRef {
	if len(namePath) == 0 {
		cberrors.Panic("namePath can not be nil")
	}
	typeRef := &TypeRef{
		NamePath: namePath,
		Origin:   origin,
	}
	var buff bytes.Buffer
	for _, nodeName := range namePath {
		buff.WriteRune('.')
		buff.WriteString(nodeName)
	}
	// 类型引用的名字为namePath按.连接,所属代码节点为此代码节点
	// 如引用ast.TypeRef namePath为[]string{"ast", "TypeRef"} 则名字为.ast.TypeRef
	typeRef.Init(buff.String(), script)
	return typeRef
}
