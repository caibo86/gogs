// -------------------------------------------
// @file      : linker.go
// @author    : 蔡波
// @contact   : caibo923@gmail.com
// @time      : 2023/12/16 下午9:57
// -------------------------------------------

package gslang

import (
	"gogs/base/gslang/ast"
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

// VisitStruct 访问结构体
func (linker *Linker) VisitStruct(s *ast.Struct) ast.Node {
	// 轮询访问表或者结构体的属性
	for _, attr := range s.Attrs() {
		attr.Accept(linker)
	}
	// 论访问表或者结构体的域
	for _, field := range s.Fields {
		field.Accept(linker)
	}
	return s
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

// VisitService 访问服务
func (linker *Linker) VisitService(service *ast.Service) ast.Node {
	// 轮询访问协议的属性
	for _, attr := range service.Attrs() {
		attr.Accept(linker)
	}
	// 轮询访问协议的父协议
	for _, base := range service.Bases {
		base.Accept(linker)
	}
	// 轮询访问协议的函数列表
	for _, method := range service.Methods {
		method.Accept(linker)
	}
	return service
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
		gserrors.Assert(nodes > 0, "the NamePath,can not be nil")
		switch nodes { // 根据类型路径长度判断
		case 1: // 长度为1 则NamePath[0]就是类型名
			// 在代码节点引用的代码包中查找指定名字目标包
			// 引用的包不能跟类型重名 如果有同名包则报错
			if pkg, ok := ref.Script().Imports[ref.NamePath[0]]; !ok {
				pkg := ref.Package() // 类型引用所属包 必须不为空
				gserrors.Assert(pkg != nil, "ref(%s) must bind ast tree ", ref)
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
			if pkg, ok := ref.Script().Imports[ref.NamePath[0]]; ok {
				gserrors.Assert(pkg.Ref != nil,
					"(%s)first parse phase must link import package:%s",
					ref.Script(), pkg)
				// 在引用的包的类型列表中查找对应名字的类型并引用
				if expr, ok := pkg.Ref.Types[ref.NamePath[1]]; ok {
					ref.Ref = expr
					return ref
				}
			} else { // 如果不是引用包中的类型 则判断是否是当前包中的枚举类型
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
