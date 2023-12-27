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
	ID       uint16 // ID 从0开始的
	Type     Expr   // 类型表达式
}

// Table 表或者结构体,表达式
type Table struct {
	BaseExpr          // 内嵌基本表达式实现
	Fields   []*Field // 结构体的字段列表
}

// NewTable 在代码节点内新建结构体
func (script *Script) NewTable(name string) *Table {
	s := &Table{}
	// 设置结构体节点为给定的名字 设置所属代码节点
	s.Init(name, script)
	return s
}

// Field 在结构体内查找给定名字的字段,返回该字段和是否找到
func (table *Table) Field(name string) (*Field, bool) {
	for _, field := range table.Fields {
		if field.Name() == name {
			return field, true
		}
	}
	return nil, false
}

// FieldByID 在结构体内查找给定ID的字段,返回该字段和是否找到
func (table *Table) FieldByID(id uint16) (*Field, bool) {
	for _, field := range table.Fields {
		if field.ID == id {
			return field, true
		}
	}
	return nil, false
}

// NewField 在结构体内新建字段
func (table *Table) NewField(name string, id uint16, t Expr) (*Field, bool) {
	for _, field := range table.Fields {
		// 字段重名
		if field.Name() == name {
			return field, false
		}
		// ID重复
		if field.ID == id {
			return field, false
		}
	}
	// 新建字段 ID为结构体的当前字段列表长度
	field := &Field{
		ID:   id,
		Type: t,
	}
	// 设置名字 设置所属代码为 所属结构体的所属代码节点
	field.Init(name, table.Script())
	// 设置父节点为此结构体节点
	field.SetParent(table)
	// 将字段添加到结构体的字段列表
	table.Fields = append(table.Fields, field)
	return field, true
}
