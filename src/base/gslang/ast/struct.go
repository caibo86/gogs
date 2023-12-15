// -------------------------------------------
// @file      : struct.go
// @author    : 蔡波
// @contact   : caibo923@gmail.com
// @time      : 2023/12/16 下午6:20
// -------------------------------------------

package ast

// Field 结构体的字段,表达式
type Field struct {
	BaseExpr        // 内嵌基本表达式实现
	ID       uint16 // ID
	Type     Expr   // 类型表达式
}

// Struct 结构体,表达式
type Struct struct {
	BaseExpr          // 内嵌基本表达式实现
	Fields   []*Field // 结构体的字段列表
}

// NewStruct 在代码节点内新建结构体
func (script *Script) NewStruct(name string) *Struct {
	s := &Struct{}
	// 设置结构体节点为给定的名字 设置所属代码节点
	s.Init(name, script)
	return s
}

// Field 在结构体内查找给定名字的字段,返回该字段和是否找到
func (s *Struct) Field(name string) (*Field, bool) {
	for _, field := range s.Fields {
		if field.Name() == name {
			return field, true
		}
	}
	return nil, false
}

// NewField 在结构体内新建字段
func (s *Struct) NewField(name string) (*Field, bool) {
	// 如果已存在同名字段则直接返回
	for _, field := range s.Fields {
		if field.Name() == name {
			return field, false
		}
	}
	// 新建字段 ID为结构体的当前字段列表长度
	// TODO 改成idl指定id,不再按顺序自增
	field := &Field{
		ID: uint16(len(s.Fields)),
	}
	// 设置名字 设置所属代码为 所属结构体的所属代码节点
	field.Init(name, s.Script())
	// 设置父节点为此结构体节点
	field.SetParent(s)
	// 将字段添加到结构体的字段列表
	s.Fields = append(s.Fields, field)
	return field, true
}
