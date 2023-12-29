// -------------------------------------------
// @file      : script.go
// @author    : 蔡波
// @contact   : caibo923@gmail.com
// @time      : 2023/12/16 下午4:35
// -------------------------------------------

package ast

import (
	"fmt"
	"gogs/base/gserrors"
)

// Script 代码节点,是一个Node,代表一个gs文件
type Script struct {
	BaseNode                        // 内嵌基本节点实现
	Imports  map[string]*PackageRef // 引用的代码包
	Types    []Expr                 // 包内类型声明,用于防止重复声明类型
	pkg      *Package
}

// NewScript 在包节点内新建一个代码节点
func (pkg *Package) NewScript(name string) (*Script, error) {
	if pkg == nil {
		gserrors.Panic("pkg can not be nil")
	}
	if old, ok := pkg.Scripts[name]; ok {
		err := fmt.Errorf("duplicate script named:%s", old.Name())
		return old, err
	}
	script := &Script{
		pkg:     pkg,
		Imports: make(map[string]*PackageRef),
	}
	// 初始化代码节点,设置包节点为代码节点的父节点
	script.Init(name, pkg)
	// 将代码节点加入到包节点的代码列表
	pkg.Scripts[name] = script
	return script, nil
}

// PackageRef 包引用节点,代表一个gs文件中引用的其他包,是一个Node
type PackageRef struct {
	BaseNode
	Ref *Package
}

// NewPackageRef 在代码节点中新建一个包引用节点
func (script *Script) NewPackageRef(name string, pkg *Package) (*PackageRef, bool) {
	// 检查已引用的包列表内是否有同名包,有的则返回该同名包,并设置新引用失败
	if ref, ok := script.Imports[name]; ok {
		return ref, false
	}
	// 新建包引用
	ref := &PackageRef{
		Ref: pkg,
	}
	// 设置包引用名字,设置父节点为此代码节点
	ref.Init(name, script)
	// 将包引用加入到代码节点的包引用列表
	script.Imports[name] = ref
	return ref, true
}

// NewType 在代码节点中新建一个类型节点,类型节点在代码节点所属包节点中唯一. 包和代码节点分别以字典和切片保存此类型节点的引用
func (script *Script) NewType(expr Expr) (Expr, bool) {
	old, ok := script.pkg.NewType(expr)
	if ok {
		script.Types = append(script.Types, expr)
		// 类型节点的父节点为此代码节点
		expr.SetParent(script)
	}
	return old, ok
}

// Package 获取代码节点所属的包
func (script *Script) Package() *Package {
	return script.pkg
}
