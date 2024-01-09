// -------------------------------------------
// @file      : ast.go
// @author    : 蔡波
// @contact   : caibo923@gmail.com
// @time      : 2023/12/16 下午2:06
// -------------------------------------------

package ast

import (
	"bytes"
	"fmt"
	"gogs/base/cberrors"
	log "gogs/base/logger"
	"reflect"
)

// Node 抽象语法树的节点
type Node interface {
	fmt.Stringer
	Name() string                           // 节点名称
	Path() string                           // 获取节点路径字符串
	Parent() Node                           // 父节点
	SetParent(parent Node) Node             // 设置父节点
	Package() *Package                      // 包
	Attrs() []*Attr                         // 属性列表
	AddAttr(attr *Attr)                     // 添加属性
	AddAttrs(attr []*Attr)                  // 添加属性列表
	RemoveAttr(attr *Attr)                  // 移除属性
	NewExtra(name string, data interface{}) // 新建额外数据
	Extra(name string) (interface{}, bool)  // 获取额外数据
	DelExtra(name string)                   // 删除额外数据
	Accept(visitor Visitor) Node            // 接受访问者
}

// Path 节点路径 从祖先节点到当前节点的路径
func Path(node Node) []Node {
	var nodes []Node
	var result []Node
	current := node
	for current != nil {
		nodes = append(nodes, current)
		current = current.Parent()
	}
	for i := len(nodes) - 1; i >= 0; i-- {
		result = append(result, nodes[i])
	}
	return result
}

// GetAttrs 获取节点内给定属性类型的属性列表
func GetAttrs(node Node, attrType Expr) []*Attr {
	var attrs []*Attr
	for _, attr := range node.Attrs() {
		if attr.Type.Ref == attrType {
			attrs = append(attrs, attr)
		}
	}
	return attrs
}

// BaseNode 基本节点
type BaseNode struct {
	name   string         // 名字
	parent Node           // 父节点
	attrs  []*Attr        // 属性列表
	extras map[string]any // 额外数据
}

// Init 初始化 设置名字和父节点
func (node *BaseNode) Init(name string, parent Node) {
	node.name = name
	node.parent = parent
}

// Name 获取名字
func (node *BaseNode) Name() string {
	return node.name
}

// String 实现fmt.Stringer接口
func (node *BaseNode) String() string {
	return node.name
}

// Path 获取节点路径
func (node *BaseNode) Path() string {
	var writer bytes.Buffer
	pathNodes := Path(node)
	for i, n := range pathNodes {
		writer.WriteString(n.Name())
		if i < len(pathNodes)-2 {
			writer.WriteRune('/')
		} else if i == len(pathNodes)-2 {
			writer.WriteRune('#')
		}
	}
	ret := writer.String()
	log.Debugf("我是%s,我的路径是:%s", node.Name(), ret)
	return ret
}

// Package 获取包
func (node *BaseNode) Package() *Package {
	// 节点所属包是其祖先节点的包
	if node.Parent() == nil {
		return nil
	}
	return node.Parent().Package()
}

// Parent 获取父节点
func (node *BaseNode) Parent() Node {
	return node.parent
}

// SetParent 设置父节点并返回旧的父节点
func (node *BaseNode) SetParent(parent Node) Node {
	old := node.parent
	node.parent = parent
	return old
}

// getExtra 获取节点额外数据
func (node *BaseNode) getExtra() map[string]any {
	if node.extras == nil {
		node.extras = make(map[string]any)
	}
	return node.extras
}

// Attrs 获取属性列表
func (node *BaseNode) Attrs() []*Attr {
	return node.attrs
}

// AddAttr 添加属性
func (node *BaseNode) AddAttr(attr *Attr) {
	for _, old := range node.attrs {
		if old == attr {
			return
		}
	}
	// 设置属性的父节点为此节点
	attr.SetParent(node)
	node.attrs = append(node.attrs, attr)
}

// AddAttrs 添加属性列表
func (node *BaseNode) AddAttrs(attrs []*Attr) {
	for _, attr := range attrs {
		node.AddAttr(attr)
	}
}

// RemoveAttr 移除属性
func (node *BaseNode) RemoveAttr(attr *Attr) {
	var attrs []*Attr
	for _, old := range node.attrs {
		if old == attr {
			continue
		}
		attrs = append(attrs, old)
	}
	node.attrs = attrs
	attr.SetParent(nil)
}

// NewExtra 新建额外数据
func (node *BaseNode) NewExtra(name string, data interface{}) {
	node.getExtra()[name] = data
}

// Extra 获取额外数据
func (node *BaseNode) Extra(name string) (interface{}, bool) {
	data, ok := node.getExtra()[name]
	return data, ok
}

// DelExtra 删除额外数据
func (node *BaseNode) DelExtra(name string) {
	delete(node.getExtra(), name)
}

// Accept 接受访问者
func (node *BaseNode) Accept(visitor Visitor) Node {
	cberrors.Panic("type(%s) not implement Accept method", reflect.TypeOf(node))
	return nil
}

// Expr 抽象语法树中的表达式 比Node多了一个归属的代码节点
type Expr interface {
	Node
	Script() *Script
	OriginName() string
}

// BaseExpr 基本表达式
type BaseExpr struct {
	BaseNode
	script *Script
}

// Init 初始化
func (expr *BaseExpr) Init(name string, script *Script) {
	if script == nil {
		cberrors.Panic("the param script can not be nil")
	}
	expr.BaseNode.Init(name, nil)
	expr.script = script
}

// Script 获取基本表达式所属的代码节点
func (expr *BaseExpr) Script() *Script {
	return expr.script
}

// Package 获取基本表达式所属的包
func (expr *BaseExpr) Package() *Package {
	// 基本表达式所属包是其所属代码节点所属的包
	return expr.Script().Package()
}

func (expr *BaseExpr) OriginName() string {
	return expr.Name()
}
