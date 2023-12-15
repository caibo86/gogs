// -------------------------------------------
// @file      : map.go
// @author    : 蔡波
// @contact   : caibo923@gmail.com
// @time      : 2023/12/16 下午6:32
// -------------------------------------------

package ast

import (
	"fmt"
)

// Map 字典表达式
type Map struct {
	BaseExpr      // 内嵌基本表达式实现
	Key      Expr // 字典key类型
	Value    Expr // 字典value类型
}

// NewMap 在代码节点内新建字典表达式
func (script *Script) NewMap(key, value Expr) *Map {
	m := &Map{
		Key:   key,
		Value: value,
	}
	name := fmt.Sprintf("map[%s]%s", key.Name(), value.Name())
	m.Init(name, script)
	key.SetParent(m)
	value.SetParent(m)
	return m
}
