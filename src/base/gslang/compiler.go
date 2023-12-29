// -------------------------------------------
// @file      : compiler.go
// @author    : 蔡波
// @contact   : caibo923@gmail.com
// @time      : 2023/12/16 下午7:52
// -------------------------------------------

package gslang

import (
	"bytes"
	"fmt"
	"gogs/base/gserrors"
	"gogs/base/gslang/ast"
	log "gogs/base/logger"
	"os"
	"path/filepath"
	"strings"
)

// setFilePath 设置代码节点的绝对文件名
func setFilePath(script *ast.Script, fullPath string) {
	script.NewExtra("FilePath", fullPath)
}

// FilePath 返回代码节点的绝对文件名
func FilePath(script *ast.Script) (string, bool) {
	path, ok := script.Extra("FilePath")
	if ok {
		return path.(string), true
	}
	return "", false
}

// Compiler 编译器
type Compiler struct {
	Loaded  map[string]*ast.Package // 已加载包节点字典
	loading []*ast.Package          // 正在加载的包节点列表
	goPath  []string                // 系统golang路径
}

// NewCompiler 新建一个编译器
func NewCompiler() *Compiler {
	goPath := os.Getenv("GOPATH")
	if goPath == "" {
		gserrors.Panicf(nil, "must set GOPATH first")
	}
	return &Compiler{
		Loaded: make(map[string]*ast.Package),
		goPath: strings.Split(goPath, string(os.PathListSeparator)),
	}
}

// searchPackage 从$GOPATH/src下查找指定名字的代码包,需要唯一
func (compiler *Compiler) searchPackage(pkgName string) string {
	var found []string
	for _, path := range compiler.goPath {
		fullPath := filepath.Join(path, "src", pkgName)
		fi, err := os.Stat(fullPath)
		if err == nil && fi.IsDir() {
			found = append(found, fullPath)
		}
	}
	// 多于1个或者少于1个包均报错
	if len(found) < 1 {
		compiler.errorf(Position{}, "found no package named:%s", pkgName)
	}
	if len(found) > 1 {
		var buff bytes.Buffer
		buff.WriteString(fmt.Sprintf("found more than one package named:%s", pkgName))
		for i, path := range found {
			buff.WriteString(fmt.Sprintf("\n\t%d:%s", i, path))
		}
		compiler.errorf(Position{}, buff.String())
	}
	// 返回唯一的包的绝对路径
	return found[0]
}

// circularRefCheck 检查循环引用,指定名字的包
func (compiler *Compiler) circularRefCheck(pkgName string) {
	var buff bytes.Buffer
	// 如果当前正在loading的包中包含对应的包名,则说明存在循环引用
	for _, pkg := range compiler.loading {
		if pkg.Name() == pkgName || buff.Len() != 0 {
			buff.WriteString(fmt.Sprintf("\t%s import\n", pkg.Name()))
		}
	}
	if buff.Len() != 0 {
		compiler.errorf(Position{}, "circular package import: %s %s", buff.String(), pkgName)
	}
}

// errorf 编译器报错
func (compiler *Compiler) errorf(position Position, template string, args ...any) {
	gserrors.Panicf(nil, fmt.Sprintf("compile: %s err: %s", position.String(), fmt.Sprintf(template, args...)))
}

// Accept 实现访问者模式,编译器访问入口
func (compiler *Compiler) Accept(visitor ast.Visitor) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
		}
	}()
	// 使用访问者对编译器已经加载的包进行遍历访问
	for _, pkg := range compiler.Loaded {
		log.Infof("Visit package: %s", pkg.Name())
		pkg.Accept(visitor)
	}
	return
}

// Compile 编译指定的代码包
func (compiler *Compiler) Compile(pkgName string) (pkg *ast.Package, err error) {
	defer func() {
		if e := recover(); e != nil {
			if _, ok := e.(gserrors.GSError); ok {
				err = e.(gserrors.GSError)
			} else {
				err = gserrors.New(e.(error))
			}
		}
	}()
	if loaded, ok := compiler.Loaded[pkgName]; ok {
		return loaded, nil
	}
	// 检查循环引用,在当前loading的包中已存在同名包,则报错
	compiler.circularRefCheck(pkgName)
	// 在系统中查找对应的包路径
	fullPath := compiler.searchPackage(pkgName)
	log.Infof("Found package: %s in: %s", pkgName, fullPath)
	// 生成一个包节点
	pkg = ast.NewPackage(pkgName)
	// 将包节点放入loading列表
	compiler.loading = append(compiler.loading, pkg)
	// 遍历目标包目录下的每一个文件
	err = filepath.Walk(fullPath, func(path string, info os.FileInfo, err error) error {
		// 系统遍历时报错则直接返回该错误
		if err != nil {
			return err
		}
		// 如果该文件是一个不同于fullPath的目录,则跳过
		if info.IsDir() && path != fullPath {
			return filepath.SkipDir
		}
		// 如果不是一个.gs文件,则跳过
		if filepath.Ext(path) != ".gs" {
			return nil
		}
		// 解析该gs文件,生成一个代码节点
		log.Info("Parsing file: ", path)
		script, err := compiler.parse(pkg, path)
		if err == nil {
			// 没有错误的话,把绝对路径保存为代码节点的额外信息
			setFilePath(script, path)
			log.Info("Done parse file: ", path)
		}
		return err
	})
	// 如果遍历过程中出现错误,则将该包从loading列表中移除
	if err != nil {
		compiler.loading = compiler.loading[:len(compiler.loading)-1]
		return
	}
	if pkg == nil {
		compiler.errorf(Position{}, "pkg should not be nil when err is nil")
	}
	compiler.link(pkg)
	// 加载完成后,将该包从loading列表中移除,并将其加入已加载列表
	compiler.loading = compiler.loading[:len(compiler.loading)-1]
	compiler.Loaded[pkgName] = pkg
	return
}

// Type 在当前编译器已加载的指定名字包中查找指定名字的类型表达式
func (compiler *Compiler) Type(pkgName, typeName string) (ast.Expr, error) {
	pkg, ok := compiler.Loaded[pkgName]
	if !ok {
		return nil, fmt.Errorf("can not find package(%s)", pkgName)
	}
	ret, ok := pkg.Types[typeName]
	if !ok {
		return nil, fmt.Errorf("can not find type(%s) in package(%s)", typeName, pkgName)
	}
	return ret, nil
}
