// -------------------------------------------
// @file      : package.go
// @author    : 蔡波
// @contact   : caibo923@gmail.com
// @time      : 2023/12/16 下午2:23
// -------------------------------------------

package ast

// Package 包节点 包是一个Node
type Package struct {
	BaseNode
	Scripts map[string]*Script // 脚本列表
	Types   map[string]Expr
}

// NewPackage 新建一个包节点
func NewPackage(name string) *Package {
	pkg := &Package{
		Scripts: make(map[string]*Script),
		Types:   make(map[string]Expr),
	}
	pkg.Init(name, nil)
	return pkg
}

// NewType 在包中添加类型,不能添加同名类型
func (pkg *Package) NewType(expr Expr) (Expr, bool) {
	if old, ok := pkg.Types[expr.Name()]; ok {
		return old, false
	}
	pkg.Types[expr.Name()] = expr
	return expr, true
}

// Package 包节点所属包节点为自身
func (pkg *Package) Package() *Package {
	return pkg
}
