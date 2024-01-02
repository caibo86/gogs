// -------------------------------------------
// @file      : service.go
// @author    : 蔡波
// @contact   : caibo923@gmail.com
// @time      : 2023/12/16 下午6:41
// -------------------------------------------

package ast

import (
	"fmt"
	"gogs/base/misc"
)

// Param 参数表达式
type Param struct {
	BaseExpr      // 内嵌基本表达式实现
	ID       int  // 参数ID
	Type     Expr // 参数类型
}

// Method 方法表达式
type Method struct {
	BaseExpr          // 内嵌基本表达式实现
	ID       uint32   // 方法ID
	Return   []*Param // 返回参数列表
	Params   []*Param // 输入参数列表
}

// InputParams 方法输入参数个数
func (method *Method) InputParams() uint16 {
	return uint16(len(method.Params))
}

// ReturnParams 方法返回参数个数
func (method *Method) ReturnParams() uint16 {
	return uint16(len(method.Return))
}

// NewReturn 在方法节点上新建返回参数 并加入到此方法返回参数列表
func (method *Method) NewReturn(paramType Expr) *Param {
	// 用给定类型表达式做类型及当前方法返回参数列表长度做ID 进行初始化
	param := &Param{
		ID:   len(method.Return),
		Type: paramType,
	}
	// 设置类型节点的父节点为此参数节点
	paramType.SetParent(param)
	// 给参数命名 设定所属代码节点为此方法节点所属的代码节点
	param.Init(fmt.Sprintf("return_param(%d)", param.ID), method.Script())
	// 参数节点的父节点为此方法节点
	param.SetParent(method)
	// 加入到此方法返回参数列表
	method.Return = append(method.Return, param)
	return param
}

// NewParam 在方法节点上新建输入参数 并加入到此方法输入参数列表
func (method *Method) NewParam(paramType Expr) *Param {
	// 用给定类型表达式做类型及当前方法返回参数列表长度做ID 进行初始化
	param := &Param{
		ID:   len(method.Params),
		Type: paramType,
	}
	// 设置类型节点的父节点为此参数节点
	paramType.SetParent(param)
	// 给参数命名 设定所属代码节点为此方法节点所属的代码节点
	param.Init(fmt.Sprintf("param(%d)", param.ID), method.Script())
	// 参数节点的父节点为此方法节点
	param.SetParent(method)
	// 加入到此方法输入参数列表
	method.Params = append(method.Params, param)
	return param
}

// Service 协议,表达式
type Service struct {
	BaseExpr                      // 内嵌基本表达式实现
	Methods   map[string]*Method  // 方法列表
	MethodIDs map[uint32]struct{} // 已使用的方法ID列表
	Bases     []*TypeRef          // 基类列表,协议可以继承自多个协议
}

// NewService 在代码节点内新建协议
func (script *Script) NewService(name string) *Service {
	service := &Service{
		Methods:   make(map[string]*Method),
		MethodIDs: make(map[uint32]struct{}),
	}
	// 设置协议节点为给定的名字 设置所属代码节点
	service.Init(name, script)
	return service
}

// NewBase 为此协议添加一个基类
func (service *Service) NewBase(base *TypeRef) (*TypeRef, bool) {
	// 检查是否已经存在此基类
	for _, old := range service.Bases {
		if old == base {
			return old, false
		}
	}
	// 添加基类
	service.Bases = append(service.Bases, base)
	return base, true
}

// NewMethod 在协议内新建一个方法
func (service *Service) NewMethod(name string) (*Method, bool) {
	// 检查是否已经存在此方法
	method, ok := service.Methods[name]
	if ok {
		return method, false
	}
	methodID := misc.BKDRHash(name)
	if _, ok = service.MethodIDs[methodID]; ok {
		return nil, false
	}
	// 新建协议
	method = &Method{
		ID: methodID,
	}
	// 初始化协议
	method.Init(name, service.Script())
	// 设置方法的父节点为此协议节点
	method.SetParent(service)
	// 将方法加入到协议的方法列表
	service.Methods[name] = method
	service.MethodIDs[methodID] = struct{}{}
	return method, true
}

// CopyMethod 复制一个方法到本协议下
func (service *Service) CopyMethod(src *Method) (*Method, bool) {
	// 检查是否已经存在此方法
	_, ok := service.Methods[src.Name()]
	if ok {
		return nil, false
	}
	if _, ok = service.MethodIDs[src.ID]; ok {
		return nil, false
	}
	// 新建协议
	method := &Method{
		ID: src.ID,
	}
	// 初始化协议
	method.Init(src.name, service.Script())
	// 设置方法的父节点为此协议节点
	method.SetParent(service)
	// 将方法加入到协议的方法列表
	service.Methods[method.Name()] = method
	service.MethodIDs[method.ID] = struct{}{}
	// 检查是否把对应的包引用添加进来
	for name, pkgRef := range src.script.Imports {
		if _, ok = service.Script().Imports[name]; !ok {
			service.Script().NewPackageRef(name, pkgRef.Ref)
		}
	}

	for _, param := range src.Params {
		r := param.Type.(*TypeRef)
		namePath := r.NamePath
		if len(r.NamePath) == 1 && src.Package() != service.Package() {
			namePath = append([]string{src.Package().Name()}, r.NamePath[0])
		}
		ref := service.Script().NewTypeRef(namePath, "")
		method.NewParam(ref)
	}
	for _, param := range src.Return {
		r := param.Type.(*TypeRef)
		namePath := r.NamePath
		if len(r.NamePath) == 1 && src.Package() != service.Package() {
			namePath = append([]string{src.Package().Name()}, r.NamePath[0])
		}
		ref := service.Script().NewTypeRef(namePath, "")
		method.NewReturn(ref)
	}
	return method, true
}
