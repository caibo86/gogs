// -------------------------------------------
// @file      : reflect.go
// @author    : 蔡波
// @contact   : caibo923@gmail.com
// @time      : 2023/12/14 下午9:15
// -------------------------------------------

package gslang

import (
	"gogs/base/gslang/ast"
	log "gogs/base/logger"
)

// Enum 将Enum表达式内的EnumVal字典解析为map
// key为EnumVal的名称，value为EnumVal的值
func Enum(enum *ast.Enum) map[string]int32 {
	ret := make(map[string]int32)
	for _, val := range enum.Values {
		ret[val.Name()] = val.Value
	}
	return ret
}

// EvalFieldInitArg 在参数列表中找到与指定字段名匹配的参数表达式
func EvalFieldInitArg(field *ast.Field, expr ast.Expr) (ast.Expr, bool) {
	eval := &evalArg{
		field: field,
	}
	expr.Accept(eval)
	if eval.expr != nil {
		return eval.expr, true
	}
	return nil, false
}

// EvalEnumVal 访问枚举值节点的val值
func EvalEnumVal(expr ast.Expr) int32 {
	visitor := &evalEnumVal{}
	expr.Accept(visitor)
	return visitor.val
}

// IsAttrUsage 判断是不是内置AttrUsage结构
func IsAttrUsage(s *ast.Table) bool {
	if s.Name() == "AttrUsage" && s.Package().Name() == GSLangPackage {
		return true
	}
	return false
}

// IsStruct 判断是不是一个结构体
func IsStruct(s *ast.Table) bool {
	_, ok := s.Extra("isStruct")
	return ok
}

// markAsStruct 将结构体标记为结构体
func markAsStruct(s *ast.Table) {
	s.NewExtra("isStruct", true)
}

// IsError 检查枚举是不是表示错误声明
func IsError(enum *ast.Enum) bool {
	_, ok := enum.Extra("isError")
	return ok
}

// markAsError 将枚举标记为错误枚举
func markAsError(enum *ast.Enum) {
	enum.NewExtra("isError", true)
}

func markAsFlower(enum *ast.Enum) {
	enum.NewExtra("isFlower", true)
}

// EvalAttrUsage 评价属性是否是AttrUsage
func (compiler *Compiler) EvalAttrUsage(attr *ast.Attr) int32 {
	// 属性的类型引用必须先连接到对应类型
	if attr.Type.Ref == nil {
		log.Panicf("attr(%s) must linked first:\n\t%s", attr, Pos(attr).String())
	}

	// 对属性求值
	ea := &evalAttr{}
	attr.Accept(ea)

	// 只有Table才能被作为属性的类型引用
	s, ok := attr.Type.Ref.(*ast.Table)
	if !ok {
		log.Panicf("only table can be used as attr type:\n\tattr def:%s\n\ttype def:%s",
			Pos(attr), Pos(attr.Type.Ref))
	}
	// 轮询属性的类型引用的属性列表
	for _, metaAttr := range s.Attrs() {
		// 属性的类型引用必须是Struct
		usage, ok := metaAttr.Type.Ref.(*ast.Table)
		if !ok {
			log.Panicf("attr(%s) must linked first:\n\t%s", metaAttr, Pos(metaAttr))
		}
		if IsAttrUsage(usage) {
			field, ok := usage.Field("Target")
			if !ok {
				log.Panicf("inner gslang AttrUsage must declare field Target:\n\t%s", Pos(usage))
			}
			if target, ok := EvalFieldInitArg(field, metaAttr.Args); ok {
				return EvalEnumVal(target)
			}
			log.Panicf("AttrUsage attribute init list expect target val \n\tattr def:%s", Pos(metaAttr))
		}
	}
	// 能作为属性的Table必须有一个属性@AttrUsage
	log.Panicf("target table can not be used as attribute type:\n\tattr def:%s\n\ttype def:%s", Pos(attr), Pos(attr.Type.Ref))
	return 0
}
