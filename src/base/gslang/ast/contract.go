// -------------------------------------------
// @file      : contract.go
// @author    : 蔡波
// @contact   : caibo923@gmail.com
// @time      : 2023/12/16 下午6:41
// -------------------------------------------

package ast

import (
	"fmt"
)

// Param 函数参数表达式
type Param struct {
	BaseExpr      // 内嵌基本表达式实现
	ID       int  // 参数ID
	Type     Expr // 参数类型
}

// Method 函数表达式
type Method struct {
	BaseExpr          // 内嵌基本表达式实现
	ID       uint32   // 函数ID
	Return   []*Param // 返回参数列表
	Params   []*Param // 输入参数列表
}

// InputParams 函数输入参数个数
func (method *Method) InputParams() uint16 {
	return uint16(len(method.Params))
}

// ReturnParams 函数返回参数个数
func (method *Method) ReturnParams() uint16 {
	return uint16(len(method.Return))
}

// NewReturn 在函数节点上新建返回参数 并加入到此函数返回参数列表
func (method *Method) NewReturn(paramType Expr) *Param {
	// 用给定类型表达式做类型及当前函数返回参数列表长度做ID 进行初始化
	param := &Param{
		ID:   len(method.Return),
		Type: paramType,
	}
	// 设置类型节点的父节点为此参数节点
	paramType.SetParent(param)
	// 给参数命名 设定所属代码节点为此函数节点所属的代码节点
	param.Init(fmt.Sprintf("return_param(%d)", param.ID), method.Script())
	// 参数节点的父节点为此函数节点
	param.SetParent(method)
	// 加入到此函数返回参数列表
	method.Return = append(method.Return, param)
	return param
}

// NewParam 在函数节点上新建输入参数 并加入到此函数输入参数列表
func (method *Method) NewParam(paramType Expr) *Param {
	// 用给定类型表达式做类型及当前函数返回参数列表长度做ID 进行初始化
	param := &Param{
		ID:   len(method.Params),
		Type: paramType,
	}
	// 设置类型节点的父节点为此参数节点
	paramType.SetParent(param)
	// 给参数命名 设定所属代码节点为此函数节点所属的代码节点
	param.Init(fmt.Sprintf("param(%d)", param.ID), method.Script())
	// 参数节点的父节点为此函数节点
	param.SetParent(method)
	// 加入到此函数输入参数列表
	method.Params = append(method.Params, param)
	return param
}

// Service 服务表达式
type Service struct {
	BaseExpr                    // 内嵌基本表达式实现
	Methods  map[string]*Method // 函数列表
	Bases    []*TypeRef         // 基类列表,服务可以继承多个服务
}

// NewService 在代码节点内新建服务
func (script *Script) NewService(name string) *Service {
	service := &Service{
		Methods: make(map[string]*Method),
	}
	// 设置服务节点为给定的名字 设置所属代码节点
	service.Init(name, script)
	return service
}

// NewBase 为此服务添加一个基类
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

// NewMethod 在服务内新建一个方法
func (service *Service) NewMethod(name string) (*Method, bool) {
	// 检查是否已经存在此函数
	method, ok := service.Methods[name]
	if ok {
		return method, false
	}
	// 新建函数
	method = &Method{
		// TODO 使用hash计算函数ID
		ID: uint32(len(service.Methods)),
	}
	// 初始化函数
	method.Init(name, service.Script())
	// 设置函数的父节点为此服务节点
	method.SetParent(service)
	// 将函数加入到服务的函数列表
	service.Methods[name] = method
	return method, true
}
