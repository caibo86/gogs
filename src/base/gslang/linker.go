// -------------------------------------------
// @file      : linker.go
// @author    : 蔡波
// @contact   : caibo923@gmail.com
// @time      : 2023/12/16 下午9:57
// -------------------------------------------

package gslang

import (
	"bytes"
	"fmt"
	"gogs/base/gslang/ast"
	log "gogs/base/logger"
)

// link 编译器链接方法
func (compiler *Compiler) link(pkg *ast.Package) {
	// 新建连接器并访问包
	linker := &Linker{
		Compiler: compiler,
	}
	// 类型连接  连接后每一个TypeRef的Ref均不为空
	pkg.Accept(linker)

	// 新建属性连接器并访问包
	linker2 := &attrLinker{
		Compiler: compiler,
	}
	// 属性连接 确保每一个属性正确挂载在对应目标类型的节点
	pkg.Accept(linker2)

	// 新建协议连接器并访问包
	linker3 := &contractLinker{
		Compiler: compiler,
	}
	// 协议展开 每一个协议都包含自己所有父协议的所有函数 并按全局编号
	pkg.Accept(linker3)
}

// Linker 连接器 此连接器是将所有的类型引用连接到对应的类型
type Linker struct {
	*Compiler        // 所属编译器
	ast.EmptyVisitor // 空的访问者 用于实现访问者接口 部分访问方法自己实现 部分采用空访问者的方法
}

// VisitPackage 访问包
func (linker *Linker) VisitPackage(pkg *ast.Package) ast.Node {
	// 轮询访问包的属性
	for _, attr := range pkg.Attrs() {
		attr.Accept(linker)
	}
	// 轮询访问包的代码
	for _, script := range pkg.Scripts {
		script.Accept(linker)
	}
	return pkg
}

// VisitScript 访问代码
func (linker *Linker) VisitScript(script *ast.Script) ast.Node {
	// 轮询访问代码的属性
	for _, attr := range script.Attrs() {
		attr.Accept(linker)
	}
	// 轮询访问代码中的类型
	for _, expr := range script.Types {
		expr.Accept(linker)
	}
	return script
}

// VisitTable 访问结构体或者表
func (linker *Linker) VisitTable(table *ast.Table) ast.Node {
	// 轮询访问表或者结构体的属性
	for _, attr := range table.Attrs() {
		attr.Accept(linker)
	}
	// 论访问表或者结构体的域
	for _, field := range table.Fields {
		field.Accept(linker)
	}
	return table
}

// VisitField 访问字段
func (linker *Linker) VisitField(field *ast.Field) ast.Node {
	// 轮询访问字段的属性
	for _, attr := range field.Attrs() {
		attr.Accept(linker)
	}
	// 访问字段引用的类型
	field.Type.Accept(linker)
	return field
}

// VisitEnum 访问枚举
func (linker *Linker) VisitEnum(enum *ast.Enum) ast.Node {
	// 轮询访问枚举的属性
	for _, attr := range enum.Attrs() {
		attr.Accept(linker)
	}
	// 轮询访问枚举的单条枚举值
	for _, val := range enum.Values {
		val.Accept(linker)
	}
	return enum
}

// VisitEnumVal 访问单条枚举值
func (linker *Linker) VisitEnumVal(val *ast.EnumVal) ast.Node {
	// 轮询访问单条枚举值的属性
	for _, attr := range val.Attrs() {
		attr.Accept(linker)
	}
	return val
}

// VisitContract 访问服务
func (linker *Linker) VisitContract(contract *ast.Contract) ast.Node {
	// 轮询访问协议的属性
	for _, attr := range contract.Attrs() {
		attr.Accept(linker)
	}
	// 轮询访问协议的父协议
	for _, base := range contract.Bases {
		base.Accept(linker)
	}
	// 轮询访问协议的函数列表
	for _, method := range contract.Methods {
		method.Accept(linker)
	}
	return contract
}

// VisitMethod 访问方法
func (linker *Linker) VisitMethod(method *ast.Method) ast.Node {
	// 轮询访问函数的属性
	for _, attr := range method.Attrs() {
		attr.Accept(linker)
	}
	// 轮询访问函数的返回参数列表
	for _, expr := range method.Return {
		expr.Accept(linker)
	}
	// 轮询访问函数的输入参数列表
	for _, expr := range method.Params {
		expr.Accept(linker)
	}
	return method
}

// VisitParam 访问参数
func (linker *Linker) VisitParam(param *ast.Param) ast.Node {
	// 轮询访问参数的属性
	for _, attr := range param.Attrs() {
		attr.Accept(linker)
	}
	// 访问参数的类型
	param.Type.Accept(linker)
	return param
}

// VisitBinaryOp 访问二元操作
func (linker *Linker) VisitBinaryOp(op *ast.BinaryOp) ast.Node {
	// 访问左操作数
	op.Left.Accept(linker)
	// 访问右操作数
	op.Right.Accept(linker)
	return op
}

// VisitList 访问切片
func (linker *Linker) VisitList(list *ast.List) ast.Node {
	// 访问切片的元素类型
	list.Element.Accept(linker)
	return list
}

// VisitArray 访问数组
func (linker *Linker) VisitArray(array *ast.Array) ast.Node {
	// 访问数组的元素类型
	array.Element.Accept(linker)
	return array
}

// VisitMap 访问字典
func (linker *Linker) VisitMap(m *ast.Map) ast.Node {
	// 访问字典的键类型
	m.Key.Accept(linker)
	// 访问字典的值类型
	m.Value.Accept(linker)
	return m
}

// VisitAttr 访问属性
func (linker *Linker) VisitAttr(attr *ast.Attr) ast.Node {
	// 访问属性 的类型应用
	attr.Type.Accept(linker)
	// 访问属性的参数列表
	if attr.Args != nil {
		attr.Args.Accept(linker)
	}
	return attr
}

// VisitArgs 访问参数列表
func (linker *Linker) VisitArgs(args *ast.Args) ast.Node {
	// 轮询访问参数列表中单个参数
	for _, arg := range args.Items {
		arg.Accept(linker)
	}
	return args
}

// VisitNamedArgs 访问命名参数列表
func (linker *Linker) VisitNamedArgs(args *ast.NamedArgs) ast.Node {
	// 轮询访问命名参数列表中的单个参数
	for _, arg := range args.Items {
		arg.Accept(linker)
	}
	return args
}

// VisitTypeRef 访问类型引用
func (linker *Linker) VisitTypeRef(ref *ast.TypeRef) ast.Node {
	if ref.Ref == nil { // 引用表达式为空的时候 需要检查路径名字
		// 路径长度需要大于1
		nodes := len(ref.NamePath)
		if nodes == 0 {
			log.Panic("the NamePath,can not be nil")
		}
		log.Debug(ref.NamePath, ref.Ref, nodes)
		switch nodes { // 根据类型路径长度判断
		case 1: // 长度为1 则NamePath[0]就是类型名
			// 在代码节点引用的代码包中查找指定名字目标包
			// 引用的包不能跟类型重名 如果有同名包则报错
			if pkg, ok := ref.Script().Imports[ref.NamePath[0]]; !ok {
				pkg := ref.Package() // 类型引用所属包 必须不为空
				if pkg == nil {
					log.Panicf("ref(%s) must bind ast tree ", ref)
				}
				// 在包内类型列表中查找对应类型 添加引用
				if expr, ok := pkg.Types[ref.NamePath[0]]; ok {
					ref.Ref = expr
					return ref
				}
			} else {
				linker.errorf(Pos(ref),
					"type name conflict with import package name:\n\tsee: %s",
					Pos(pkg))
			}
		case 2: // 路径长度为2  eg: ast.Node
			// 在代码应用的包列表中查找NamePath[0],即目标类型所属的包
			log.Debug(ref.Script().Imports)
			if pkg, ok := ref.Script().Imports[ref.NamePath[0]]; ok {

				if pkg.Ref == nil {
					log.Panicf("(%s)first parse phase must link import package:%s",
						ref.Script(), pkg)
				}
				// 在引用的包的类型列表中查找对应名字的类型并引用
				if expr, ok := pkg.Ref.Types[ref.NamePath[1]]; ok {
					ref.Ref = expr
					return ref
				}
			} else { // 如果不是引用包中的类型 则判断是否是当前包中的枚举类型
				log.Debug("下面")
				if expr, ok := ref.Package().Types[ref.NamePath[0]]; ok {
					if enum, ok := expr.(*ast.Enum); ok {
						if val, ok := enum.Values[ref.NamePath[1]]; ok {
							ref.Ref = val
							return ref
						}
					}
				}
			}
		case 3: // 长度为3的情况 一定是引用包中的枚举类型
			if pkg, ok := ref.Script().Imports[ref.NamePath[0]]; ok {
				if expr, ok := pkg.Ref.Types[ref.NamePath[1]]; ok {
					if enum, ok := expr.(*ast.Enum); ok {
						if val, ok := enum.Values[ref.NamePath[2]]; ok {
							ref.Ref = val
							return ref
						}
					}
				}
			}
		}
	}
	// 以上情况均不符合则报错
	linker.errorf(Pos(ref), "unknown type(%s)", ref)
	return ref
}

// 协议连接器
type contractLinker struct {
	*Compiler        // 所属连接器
	ast.EmptyVisitor // 内嵌空访问者
}

// VisitPackage 访问包
func (linker *contractLinker) VisitPackage(pkg *ast.Package) ast.Node {
	// 轮询访问包内代码列表
	for _, script := range pkg.Scripts {
		script.Accept(linker)
	}
	return pkg
}

// VisitScript 访问代码
func (linker *contractLinker) VisitScript(script *ast.Script) ast.Node {
	// 轮询访问代码内的类型
	for _, expr := range script.Types {
		expr.Accept(linker)
	}
	return script
}

// VisitContract 访问协议
func (linker *contractLinker) VisitContract(contract *ast.Contract) ast.Node {
	linker.unwind(contract, nil)
	return contract
}

// unwind 展开协议 协议展开后 每一个协议都持有所有父协议的函数 并重新分配 ID
func (linker *contractLinker) unwind(expr *ast.Contract, stack []*ast.Contract) []*ast.Contract {
	// 如果协议有已经展开的额外信息 则直接返回协议栈
	if _, ok := expr.Extra("unwind"); ok {
		return stack
	}
	var buff bytes.Buffer
	// 在协议栈中查找是否存在递归继承 如果有则报错
	for _, contract := range stack {
		if contract == expr || buff.Len() != 0 {
			buff.WriteString(fmt.Sprintf("\t%s inheri\n", contract))
		}
	}
	if buff.Len() != 0 {
		linker.errorf(Pos(expr), "circular inherit:\n%s\t%s", buff.String(), expr)
	}
	// 将该协议添加到栈尾
	stack = append(stack, expr)
	// 用于保存当前协议的父协议们的函数总数
	modify := uint32(0)
	// 检查协议的父协议 并展开
	for _, base := range expr.Bases {
		contract, ok := base.Ref.(*ast.Contract)
		if !ok { // 检查父协议的类型是否正确
			linker.errorf(Pos(base),
				"contract(%s) inherit type is not contract:\n\tsee: %s", expr,
				Pos(base.Ref))
		}
		// 将所有父协议压栈
		stack = linker.unwind(contract, stack)
		// 统计父协议函数总数
		modify = modify + uint32(len(contract.Methods))
	}
	// 处理协议的函数ID 加上父协议的函数总数
	for _, method := range expr.Methods {
		method.ID = method.ID + modify
	}
	// 将父协议的函数列表复制到当前协议
	modify = uint32(0)
	for _, base := range expr.Bases {
		contract := base.Ref.(*ast.Contract)
		for _, method := range contract.Methods {
			clone := &ast.Method{}
			*clone = *method
			clone.ID = clone.ID + modify
			if old, ok := expr.Methods[clone.Name()]; ok { // 不允许有重名函数
				linker.errorf(Pos(expr),
					"duplicate method name: %s\n\tsee: %s\n\t see: %s",
					clone,
					Pos(old),
					Pos(clone))
			}
			method.SetParent(expr)
			expr.Methods[clone.Name()] = clone
		}
		modify = modify + uint32(len(contract.Methods))
	}
	// 标记当前协议已经展开
	expr.NewExtra("unwind", true)
	// 丢弃栈尾元素
	stack = stack[:len(stack)-1]
	return stack
}

// attrLinker 属性连接器
type attrLinker struct {
	*Compiler                         // 所属编译器
	ast.EmptyVisitor                  // 内嵌空访问者
	attrTarget       map[string]int64 // 指定为gslang包中的AttrStruct枚举类型解析后的字典
	attrStruct       ast.Expr         // 指定为gslang包中的Struct类型
	attrError        ast.Expr         // 指定为gslang包中的Error类型
}

// VisitPackage	访问包
func (linker *attrLinker) VisitPackage(pkg *ast.Package) ast.Node {
	if len(pkg.Scripts) == 0 {
		return pkg
	}
	log.Debugf("我来看下包名:%s", pkg.Name())
	// 设置属性连接器的 属性目标为 gslang编译器内置的 指定名字的枚举值 解析成的字典
	if pkg.Name() == GSLangPackage {
		if expr, ok := pkg.Types[GSLangAttrTarget]; ok {
			if enum, ok := expr.(*ast.Enum); ok {
				linker.attrTarget = Enum(enum)
			}
		}
	} else {
		if pkg1, ok := linker.Loaded[GSLangPackage]; ok {
			if expr, ok := pkg1.Types[GSLangAttrTarget]; ok {
				if enum, ok := expr.(*ast.Enum); ok {
					linker.attrTarget = Enum(enum)
				}
			}
		}
	}
	if linker.attrTarget == nil {
		log.Panicf("inner error: can't found gslang.AttrTarget enum")
	}
	// 设置结构和枚举两种内置类型
	log.Debug(pkg.Name(), " ", GSLangPackage)
	if pkg.Name() == GSLangPackage {
		linker.attrStruct = pkg.Types[GSLangAttrStruct]
		if linker.attrStruct == nil {
			log.Panicf("inner error: can't found gslang.Table attribute type")
		}
		linker.attrError = pkg.Types[GSLangAttrError]
		if linker.attrError == nil {
			log.Panicf("inner error: can't found gslang.Error attribute type")
		}
	} else {
		attrStruct, err := linker.Type(GSLangPackage, GSLangAttrStruct)
		if err != nil {
			log.Panicf("inner error: can't found gslang.Table attribute type. err:%v", err)
		}
		linker.attrStruct = attrStruct
		attrError, err := linker.Type(GSLangPackage, GSLangAttrError)
		if err != nil {
			log.Panicf("inner error: can't found gslang.Error attribute type. err:%v", err)
		}
		linker.attrError = attrError
	}
	// 轮询访问包中代码
	for _, script := range pkg.Scripts {
		script.Accept(linker)
	}
	return pkg
}

// VisitScript 访问代码
func (linker *attrLinker) VisitScript(script *ast.Script) ast.Node {
	// 轮询代码的属性
	for _, attr := range script.Attrs() {
		target := linker.EvalAttrUsage(attr)
		// 如果属性目标不是 AttrTarget.Script
		if target&linker.attrTarget["Script"] == 0 {
			// 如果属性()中是 AttrTarget.GSLangPackage
			if target&linker.attrTarget["Package"] != 0 {
				// 将此属性从代码节点删除 并添加到代码节点所属的包节点下
				script.RemoveAttr(attr)
				script.Package().AddAttr(attr)
			} else {
				linker.errorf(Pos(attr), "attr(%s) can't be used to attribute script:\n\tsee: %s",
					attr,
					Pos(attr.Type.Ref))
			}
		}
	}
	// 轮询访问代码节点中类型
	for _, expr := range script.Types {
		expr.Accept(linker)
	}
	return script
}

// VisitTable 访问Table
func (linker *attrLinker) VisitTable(table *ast.Table) ast.Node {
	// 是否Struct
	var isStruct bool
	// 如果table的属性中有内置 标识符为Struct的Table类型
	if len(ast.GetAttrs(table, linker.attrStruct)) > 0 {
		// 则认为是Struct
		isStruct = true
		// 标记Struct额外信息
		markAsStruct(table)
	}
	// 轮询判断table的属性的目标是不是table的类型 不是则移动到对应的类型节点  代码节点或者包节点
	log.Debugf("现在看%s的attrs", table.Name())
	for _, attr := range table.Attrs() {
		log.Debug("参数:", attr.Name())
		target := linker.EvalAttrUsage(attr)
		log.Debug("target:", target)
		var toMove bool
		if isStruct {
			if target&linker.attrTarget["Struct"] == 0 {
				toMove = true
			}
		} else {
			if target&linker.attrTarget["Table"] == 0 {
				toMove = true
			}
		}
		if toMove {
			if target&linker.attrTarget["Script"] != 0 {
				table.RemoveAttr(attr)
				table.Script().AddAttr(attr)
				continue
			}
			if target&linker.attrTarget["Package"] != 0 {
				table.RemoveAttr(attr)
				table.Package().AddAttr(attr)
				continue
			}
			linker.errorf(Pos(attr),
				"attr(%s) can't be used to attribute table/struct:\n\tsee: %s",
				attr,
				Pos(attr.Type.Ref))
		}
	}
	for _, field := range table.Fields {
		field.Accept(linker)
	}
	return table
}

// VisitField 访问域 确认域的属性的目标为AttrUsage.Field
func (linker *attrLinker) VisitField(field *ast.Field) ast.Node {
	for _, attr := range field.Attrs() {
		target := linker.EvalAttrUsage(attr)
		if target&linker.attrTarget["Field"] == 0 {
			linker.errorf(Pos(attr),
				"attr(%s) can't be used to attribute filed:\n\tsee: %s",
				attr,
				Pos(attr.Type.Ref),
			)
		}
	}
	return field
}

// VisitEnum 访问枚举
func (linker *attrLinker) VisitEnum(enum *ast.Enum) ast.Node {
	// 如果enum的属性中有 类型引用为内置 gslang.Error类型  则认为此枚举是一个错误枚举 并标记
	if len(ast.GetAttrs(enum, linker.attrError)) > 0 {
		markAsError(enum)
	}
	// 确认属性目标类型相符
	for _, attr := range enum.Attrs() {
		target := linker.EvalAttrUsage(attr)
		if target&linker.attrTarget["Enum"] == 0 {
			linker.errorf(
				Pos(attr),
				"attr(%s) can't be used to attribute enum:\n\tsee: %s", attr,
				Pos(attr.Type.Ref))
		}
	}
	// 轮询访问单挑枚举值
	for _, val := range enum.Values {
		val.Accept(linker)
	}
	return enum
}

// VisitEnumVal 访问单挑枚举值
func (linker *attrLinker) VisitEnumVal(val *ast.EnumVal) ast.Node {
	// 确认属性目标是AttrTarget.EnumVal
	for _, attr := range val.Attrs() {
		target := linker.EvalAttrUsage(attr)
		if target&linker.attrTarget["EnumVal"] == 0 {
			linker.errorf(
				Pos(attr),
				"attr(%s) can't be used to attribute enum value:\n\tsee: %s",
				attr, Pos(attr.Type.Ref))
		}
	}
	return val
}

// VisitContract 访问协议
func (linker *attrLinker) VisitContract(contract *ast.Contract) ast.Node {
	for _, attr := range contract.Attrs() {
		target := linker.EvalAttrUsage(attr)
		if target&linker.attrTarget["Script"] != 0 {
			contract.RemoveAttr(attr)
			contract.Script().AddAttr(attr)
			continue
		}
		if target&linker.attrTarget["Package"] != 0 {
			contract.RemoveAttr(attr)
			contract.Package().AddAttr(attr)
			continue
		}
		if target&linker.attrTarget["Contract"] != 0 {
			continue
		}
		linker.errorf(Pos(attr), "attr(%s) can't be used to contract:\n\tsee: %s",
			attr, Pos(attr.Type.Ref))
	}
	for _, method := range contract.Methods {
		method.Accept(linker)
	}
	return contract
}

// VisitMethod 访问函数
func (linker *attrLinker) VisitMethod(method *ast.Method) ast.Node {
	// 确保各属性的目标与 挂载的目标节点类型相符
	for _, attr := range method.Attrs() {
		target := linker.EvalAttrUsage(attr)
		if target&linker.attrTarget["Method"] == 0 {
			linker.errorf(Pos(attr),
				"attr(%s) can't be used to attribute method:\n\tsee: %s", attr,
				Pos(attr.Type.Ref))
		}
	}
	for _, expr := range method.Return {
		for _, attr := range expr.Attrs() {
			target := linker.EvalAttrUsage(attr)
			if target&linker.attrTarget["Return"] == 0 {
				linker.errorf(Pos(attr),
					"attr(%s) can't be used to attribute return param:\n\tsee: %s", attr,
					Pos(attr.Type.Ref))
			}
		}
	}
	for _, expr := range method.Params {
		for _, attr := range expr.Attrs() {
			target := linker.EvalAttrUsage(attr)
			if target&linker.attrTarget["Param"] == 0 {
				linker.errorf(Pos(attr),
					"attr(%s) can't be used to attribute method param:\n\tsee: %s", attr,
					Pos(attr.Type.Ref))
			}
		}
	}
	return method
}
